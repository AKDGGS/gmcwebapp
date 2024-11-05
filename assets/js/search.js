let search_control = new SearchControl({ moretools: true });
let drawbox_control = new DrawBoxControl({ callback: e => doSearch() });
let result_source = new ol.source.Vector();
let fmt = new ol.format.GeoJSON({
	dataProjection: 'EPSG:4326',
	featureProjection: 'EPSG:3857'
});

// Extends URLSearchParams to apply custom behavior when appending values
URLSearchParams.prototype.apply = function(k,v){
	if(typeof v === 'string' && v.trim() === '') return this;
	switch(k){
		case 'size':
			if(Number(v) !== 25) this.append(k,v);
		break;
		case 'from':
			if(Number(v) > 0) this.append(k,v);
		break;
		case 'sort':
			if(v !== '_score') this.append(k,v);
		break;
		case 'dir':
			if(v.toLowerCase() !== 'asc') this.append(k,'desc');
		break;
		default: this.append(k,v.trim());
	}
	return this;
};

// Convenience function: makes a select and label for search tools
function createToolSelect(label, name, url){
	const pro = fetch(url).then(r => {
		if(!r.ok) throw `${label} response not ok`;
		return r.json();
	}).then(vals => {
		let lbl = document.createElement('label');
		lbl.htmlFor = name;
		lbl.textContent = label;

		let sel = document.createElement('select');
		sel.name = sel.id = name;
		sel.autocomplete = 'off';
		sel.multiple = true;
		sel.size = 5;
		sel.addEventListener('change', e => doSearch());
		vals.forEach(v => {
			let opt = document.createElement('option');
			switch(typeof v){
				case 'string':
					opt.textContent = v;
				break;
				case 'object':
					opt.textContent = v['name'];
					opt.value = v['id'];
				break;
			}
			sel.appendChild(opt);
		});

		let div = document.createElement('div');
		div.appendChild(lbl);
		div.appendChild(sel);
		search_control.getSearchTools().appendChild(div);
	}).catch(err => {
		if(window.console) console.log(err);
	});
	return pro;
}

// Convenience function: empties element of all child nodes
function elementEmpty(el){
	if(typeof el === 'string') el = document.getElementById(el);
	while(el.lastChild) el.removeChild(el.lastChild);
	return el
}

// Convenience function: alters CSS display property
function elementDisplay(el, disp){
	if(typeof el === 'string') el = document.getElementById(el);
	el.style.display = disp;
	return el;
}

let search_active = false;
function doSearch(dir){
	if(search_active) return;
	search_active = true;

	let new_sp = new URLSearchParams();
	document.querySelectorAll(
		'#result-control select, .ol-search-tools select'
	).forEach(e => Array.from(e.options).forEach(o => {
		if(!o.selected) return;
		new_sp.apply(e.name, (o.value !== '' ? o.value : o.textContent));
	}));
	new_sp.apply('q', search_control.getSearchBox().value).sort();
	/*
	let feat = drawbox_control.getFeature();
	if(feat !== null){
		let geojson = fmt.writeGeometry(feat.getGeometry());
		url += `${url?'&':''}geojson=${encodeURIComponent(geojson)}`;
	}
	*/

	let old_sp = new URLSearchParams(window.location.search);
	old_sp.sort();
	let old_qs = old_sp.toString();
	let nfrom = Number(old_sp.get('from'));
	old_sp.delete('from');

	if(new_sp.toString() !== old_sp.toString()) nfrom = 0;
	else if(dir === 1) nfrom += Math.max(Number(new_sp.get('size')), 25);
	else if(dir === -1) {
		nfrom = Math.max(nfrom - Math.max(Number(new_sp.get('size')), 25), 0);
	}
	new_sp.apply('from', nfrom).sort();

	url = `search.json?${new_sp.toString()}`;
	fetch(url).then(response => {
		if(!response.ok){ throw 'response not ok'; }
		return response.json();
	}).then(response => {
		if(new_sp.toString() !== old_qs){
			window.history.pushState(null, '', `search?${new_sp.toString()}`);
		}
		let result = elementEmpty('result');
		result_source.clear();

		// If there's no results
		if (response?.hits === undefined){
			let div = document.createElement('div');
			div.textContent = 'No results found.';
			div.className = 'noresults';
			result.appendChild(div);
			elementDisplay('result-control', 'none');
			search_active = false;
			return
		}

		elementEmpty('result-from').textContent = (response.from + 1);
		elementEmpty('result-to').textContent = (
			response.from + response.hits.length
		);
		elementEmpty('result-total').textContent = response.total;
		document.getElementById('result-prev').disabled = (response.from === 0);
		document.getElementById('result-next').disabled = (
			(response.from + response.hits.length) >= response.total
		);

		response.hits.forEach(hit => {
			if ('geometries' in hit){
				hit['geometries'].forEach(g => {
					let feat = fmt.readFeature(g)
					feat.setProperties(hit, true)
					result_source.addFeature(feat);
				});
			}
		});

		result.innerHTML = mustache.render(
			document.getElementById('tmpl-search').innerHTML,
			response, {}, ['[[', ']]']
		);

		elementDisplay('result-control', 'block');
		search_active = false;
	}).catch(err => {
		if(window.console) console.log(err);
		alert('An error occurred while talking to server, please try again.');
		search_active = false;
	});
}

function updateFromURL(){
	let params = new URLSearchParams(window.location.search);
	params.delete('from');

	// Handle the query separately
	let q = params.get('q');
	if(q !== null) search_control.getSearchBox().value = q;
	params.delete('q');

	// Iterate over the parameter keys and extract any duplicates
	let keys = params.keys().toArray().filter((v,i,a) => {
		return a.indexOf(v) === i;
	});
	// Iterate over the unique keys and push values
	// into same-named elements
	let show = false;
	keys.forEach(k => {
		let vals = params.getAll(k);
		let els = document.querySelectorAll(`[name=${k}]`)
		if(vals.length === 0 || els.length === 0) return;

		let eli = 0;
		vals.forEach(v => {
			if(!show && search_control.inSearchTools(els[eli])) show = true;
			switch(els[eli].tagName){
				case 'SELECT':
					Array.from(els[eli].options).every(o => {
						if(o.value === v || o.textContent === v){
							return o.selected = true;
						}
						return true;
					});
					break
				case 'INPUT': els[eli].value = v;
			}
			// Move the element index forward if possible,
			// to handle multiple elements with the same name
			eli = Math.min((els.length - 1),(eli + 1));
		});
	});
	if(show) search_control.showSearchTools();
	doSearch();
}

Promise.allSettled([
	createToolSelect('Keywords', 'keyword', '../keywords.json'),
	createToolSelect('Collections', 'collection_id', '../collections.json'),
	createToolSelect('Prospects', 'prospect_id', '../prospects.json')
]).then(() => {
	let map = new ol.Map({
		target: 'map',
		controls: MAP_DEFAULTS.Controls.extend([
			search_control, drawbox_control
		]),
		interactions: MAP_DEFAULTS.Interactions,
		view: MAP_DEFAULTS.View,
		layers: [
			MAP_DEFAULTS.BaseLayers,
			MAP_DEFAULTS.OverlayLayers,
			new ol.layer.Vector({
				style: MAP_DEFAULTS.DynamicStyle,
				source: result_source
			})
		]
	});

	map.on('pointermove', function(e){
		e.map.getTargetElement().style.cursor = (
			e.map.hasFeatureAtPixel(e.pixel) ? 'pointer' : ''
		);
	});

	search_control.getSearchBox().addEventListener('keydown', e => {
		if (e.keyCode == 13){
			doSearch();
			e.preventDefault();
			return false;
		}
	});
	search_control.getSearchButton().addEventListener('click', e => {
		doSearch();
		e.preventDefault();
		return false;
	});

	document.getElementById('size').addEventListener(
		'change', e => doSearch()
	);
	document.getElementById('result-prev').addEventListener(
		'click', e => doSearch(-1)
	);
	document.getElementById('result-next').addEventListener(
		'click', e => doSearch(1)
	);
	document.getElementById('result-reset').addEventListener('click', e => {
		elementEmpty('result');
		result_source.clear();
		elementDisplay('result-control', 'none');
		document.querySelectorAll(
			'#result-control select, .ol-search-tools select'
		).forEach(e => Array.from(e.options).forEach(o => {
			if(o.dataset.default) o.selected = true;
			else o.selected = false;
		}));
		window.history.pushState(null, '', 'search');
		search_control.hideSearchTools();
		search_control.getSearchBox().value = '';
		search_control.getSearchBox().focus();
	});
	document.querySelectorAll('select[name="sort"]').forEach(e => {
		e.addEventListener('change', x => doSearch());
	});
	document.querySelectorAll('select[name="dir"]').forEach(e => {
		e.addEventListener('change', x => doSearch());
	});

	// If there's a query string, use that to rebuild the search form
	if(window.location.search) updateFromURL();

	// Refresh search if user moves backwards or forwards in history
	window.addEventListener('popstate', updateFromURL);

	search_control.getSearchBox().focus();
});
