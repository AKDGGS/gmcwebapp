let search_control = new SearchControl({ moretools: true });
let drawbox_control = new DrawBoxControl({ callback: e => doSearch() });
let result_source = new ol.source.Vector();
let fmt = new ol.format.GeoJSON({
	dataProjection: 'EPSG:4326',
	featureProjection: 'EPSG:3857'
});
let map;

// Convenience function: makes a select and label for search tools
function createToolSelect(label, name, url){
	let div = document.createElement('div');
	search_control.getSearchTools().appendChild(div);
	return fetch(url).then(r => {
		if(!r.ok) throw `${label} response not ok`;
		return r.json();
	}).then(vals => {
		div.innerHTML = `<label for="${name}">${label}</label>` +
			`<select name="${name}" autocomplete="off" size="5" multiple>` +
			(vals.reduce((a,v) => {
				switch(typeof v){
					case 'string': return a += `<option>${v}</option>`;
					case 'object':
						return a += `<option value="${v.id}">${v.name}</option>`;
				}
			}, '')) + '</select>';
		div.querySelector('select').addEventListener('change', e => doSearch());
	}).catch(err => {
		if(window.console) console.log(err);
	});
}

let search_active = false;
function doSearch(dir){
	if(search_active) return;
	search_active = true;

	let new_sp = new URLSearchParams();
	document.querySelectorAll(
		'#result-control select, .ol-search-tools select, .ol-search-tools input'
	).forEach(e => {
		switch(e.tagName){
			case 'SELECT':
				Array.from(e.options).forEach(o => {
					if(!o.selected || o.dataset.default) return;
					new_sp.append(e.name, (o.value !== '' ? o.value : o.textContent));
				});
			break;
			case 'INPUT':
				if(e.value.trim() !== '') new_sp.append(e.name, e.value.trim());
			break;
		}
	});
	let q = search_control.getSearchBox().value.trim();
	if(q !== '') new_sp.append('q', q);
	new_sp.sort();
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
	if(nfrom > 0) new_sp.append('from', nfrom);
	new_sp.sort();

	url = `search.json?${new_sp.toString()}`;
	fetch(url).then(response => {
		if(!response.ok){ throw 'response not ok'; }
		return response.json();
	}).then(response => {
		if(new_sp.toString() !== old_qs || !window.location.href.includes('?')){
			window.history.pushState(null, '', `search?${new_sp.toString()}`);
		}
		result_source.clear();

		// If there's no results
		if (response?.hits === undefined){
			document.querySelector('#result').innerHTML =
				'<div class="noresults">No results found.</div>';
			document.querySelector('#result-control').style.display = 'none';
			search_active = false;
			map.updateSize();
			return
		}

		document.querySelector('#result-from').textContent = (response.from + 1);
		document.querySelector('#result-to').textContext = (
			response.from + response.hits.length
		);
		document.querySelector('#result-total').textContent = response.total;
		document.querySelector('#result-prev').disabled = (response.from === 0);
		document.querySelector('#result-next').disabled = (
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
		document.querySelector('#result').innerHTML = mustache.render(
			document.getElementById('tmpl-search').innerHTML,
			response, {}, ['[[', ']]']
		);

		document.querySelector('#result-control').style.display = 'block';
		map.updateSize();
		search_active = false;
	}).catch(err => {
		if(window.console) console.log(err);
		alert('An error occurred while talking to server, please try again.');
		search_active = false;
	});
}

function updateFromURL(){
	let pr = new URLSearchParams(window.location.search);
	pr.delete('from');

	// Handle the query separately
	let q = pr.get('q');
	search_control.getSearchBox().value = (q == null ? '' : q);
	pr.delete('q');

	// Fetch distinct keys
	let keys = pr.keys().toArray().filter((v,i,a) => {
		return a.indexOf(v) === i;
	});

	// Determine if the advanced controls need to be shown
	let shown = Array.from(document.querySelectorAll(
		'.ol-search-tools input, .ol-search-tools select'
	)).map(e => e.name).some(v => keys.includes(v));

	document.querySelectorAll(
		'#result-control select, .ol-search-tools select, .ol-search-tools input'
	).forEach(e => {
		switch(e.tagName){
			case 'SELECT':
				let vals = pr.getAll(e.name);
				let vlen = vals.length;

				Array.from(e.options).forEach(o => {
					if(vlen === 0){
						return o.selected = (o.dataset.default ? true : false);
					}

					let i = vals.indexOf((o.value !== '' ? o.value : o.textContent));
					o.selected = (i >= 0);
					if(i >= 0){
						pr.delete(e.name, vals[i]);
						if(!e.multiple) vals = [];
						else vals.splice(vals, i);
					}
				});
			break;

			case 'INPUT':
				let v = pr.get(e.name);
				e.value = (v === null ? '' : v);
				pr.delete(e.name, v);
			break;
		}
	});

	if(shown) search_control.showSearchTools();
	else search_control.hideSearchTools();
	if(window.location.href.includes('?')) doSearch();
	else {
		document.querySelector('#result').innerHTML = '';
		result_source.clear();
		document.querySelector('#result-control').style.display = 'none';
	}
}

Promise.allSettled([
	createToolSelect('Keywords', 'keyword', '../keywords.json'),
	createToolSelect('Collections', 'collection_id', '../collections.json'),
	createToolSelect('Prospects', 'prospect_id', '../prospects.json')
]).then(() => {
	map = new ol.Map({
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
		if (e.keyCode == 13){ doSearch(); e.preventDefault(); }
	});
	search_control.getSearchButton().addEventListener('click', e => {
		doSearch();
		e.preventDefault();
		return false;
	});

	let div = document.createElement('div');
	div.innerHTML = '<label for="top">Interval</label>' +
		'Top: <input type="text" name="top" autocomplete="off" size="5">' +
		'Bottom: <input type="text" name="bottom" autocomplete="off" size="5">';
	div.querySelectorAll('input').forEach(el => {
		el.addEventListener('keydown', e => {
			if (e.keyCode == 13){ doSearch(); e.preventDefault(); }
		});
	});
	search_control.getSearchTools().appendChild(div);

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
		window.history.pushState(null, '', 'search');
		updateFromURL();
		search_control.getSearchBox().focus();
	});
	document.querySelectorAll('select[name="sort"]').forEach(e => {
		e.addEventListener('change', x => doSearch());
	});
	document.querySelectorAll('select[name="dir"]').forEach(e => {
		e.addEventListener('change', x => doSearch());
	});

	// If there's a query string, use that to rebuild the search form
	if(window.location.href.includes('?')) updateFromURL();
	// Refresh search if user moves backwards or forwards in history
	window.addEventListener('popstate', updateFromURL);

	search_control.getSearchBox().focus();
});
