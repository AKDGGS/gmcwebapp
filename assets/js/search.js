let search_control = new SearchControl({ moretools: true });
let drawbox_control = new DrawBoxControl({ callback: doSearch });
let result_source = new ol.source.Vector();
let fmt = new ol.format.GeoJSON({
	dataProjection: 'EPSG:4326',
	featureProjection: 'EPSG:3857'
});

URLSearchParams.prototype.seapp = function(k,v){
	if(typeof v === 'string' && v.trim() === '') return;
	switch(k){
		case 'size':
			if(Number(v) !== 25) this.append(k,v);
			return;
		case 'from':
			if(Number(v) > 0) this.append(k,v);
			return
		case 'sort':
			if(v !== '_score') this.append(k,v);
			return
		case 'dir':
			if(v !== 'asc') this.append(k,v);
			return
		default:
			this.append(k,v.trim());
	}
};

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
// Convenience function: returns value of all selected options
function elementSelected(el){
	if(typeof el === 'string') el = document.getElementById(el);
	return Array.from(el.querySelectorAll('option'))
		.filter(o => o.selected).map(o => o.value);
}

let map = new ol.Map({
	target: 'map',
	controls: ol.control.defaults.defaults({ attribution: false }).extend([
		search_control, drawbox_control,
		new ol.control.ScaleLine({ units: "us" }),
		new ol.control.LayerSwitcher({
			tipLabel: 'Legend',
			groupSelectStyle: 'none'
		}),
		new ol.control.MousePosition({
			projection: 'EPSG:4326',
			placeholder: '',
			coordinateFormat: ol.coordinate.createStringXY(3)
		})
	]),
	interactions: ol.interaction.defaults.defaults({ mouseWheelZoom: false }).extend([
		new ol.interaction.MouseWheelZoom({
			condition: ol.events.condition.platformModifierKeyOnly
		})
	]),
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

let search_active = false;
function doSearch(dir){
	if(search_active) return;
	search_active = true;

	let cte = (p, e, v) => {
		let el = document.createElement(e);
		el.appendChild(document.createTextNode(v === undefined ? '' : v));
		p.appendChild(el);
		return el;
	};

	let oldparams = new URL(window.location.toString()).searchParams;
	let nfrom = Number(oldparams.get('from'));
	oldparams.delete('from');
	oldparams.sort();

	let params = new URLSearchParams();
	document.querySelectorAll('#result-control select').forEach(e => {
		e.querySelectorAll('option:checked').forEach(o => {
			params.seapp(e.getAttribute('name'), o.value);
		});
	});
	params.seapp('q', search_control.getSearchBox().value);
	/*
	let feat = drawbox_control.getFeature();
	if(feat !== null){
		let geojson = fmt.writeGeometry(feat.getGeometry());
		url += `${url?'&':''}geojson=${encodeURIComponent(geojson)}`;
	}
	*/
	params.sort();

	// Search is dirty. Restart from beginning
	if(params.toString() !== oldparams.toString()){
		nfrom = 0;
	} else if(dir === 1){
		nfrom += Math.max(Number(params.get('size')), 25);
	} else if(dir === -1){
		nfrom = Math.max(nfrom - Math.max(Number(params.get('size')), 25), 0);
	}
	params.seapp('from', nfrom);

	url = `search.json?${params.toString()}`;
	fetch(url).then(response => {
		if(!response.ok){ throw 'response not ok'; }
		return response.json();
	}).then(response => {
		window.history.replaceState(null, '', `search?${params.toString()}`);
		let result = elementEmpty('result');
		result_source.clear();

		// If there's no results
		if (response?.hits === undefined){
			let div = cte(result, 'div', 'No results found.');
			div.className = 'noresults';
			elementDisplay('result-control', 'none');
			search_active = false;
			return
		}

		elementEmpty('result-from').appendChild(
			document.createTextNode(response.from + 1)
		);
		elementEmpty('result-to').appendChild(
			document.createTextNode(response.from + response.hits.length)
		);
		elementEmpty('result-total').appendChild(
			document.createTextNode(response.total)
		);
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
document.querySelectorAll('select[name="sort"]').forEach(e => {
	e.addEventListener('change', x => doSearch());
});
document.querySelectorAll('select[name="dir"]').forEach(e => {
	e.addEventListener('change', x => doSearch());
});

// If there's a query string, use that to rebuild the search form
if(window.location.search){
	let params = new URL(window.location.toString()).searchParams;
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
	keys.forEach(k => {
		let vals = params.getAll(k);
		let els = document.querySelectorAll(`[name=${k}]`)
		if(vals.length === 0 || els.length === 0) return;

		let eli = 0;
		vals.forEach(v => {
			switch(els[eli].tagName){
				case 'SELECT':
					els[eli].querySelector(`option[value="${v}"]`).selected = true;
					break
				case 'INPUT': els[eli].value = v;
			}
			// Move the element index forward if possible,
			// to handle multiple elements with the same name
			eli = Math.min((els.length - 1),(eli + 1));
		});
	});
	doSearch();
}
search_control.getSearchBox().focus();
