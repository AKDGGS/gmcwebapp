let search_control = new SearchControl({ moretools: true });
let drawbox_control = new DrawBoxControl({ callback: e => doSearch(e) });
let result_source = new ol.source.Vector();
let fmt = new ol.format.GeoJSON({
	dataProjection: 'EPSG:4326',
	featureProjection: 'EPSG:3857'
});
let map;

// Compare to URLSearchParams for equality, optionally ignoring 'from'
URLSearchParams.prototype.equals = function(u2, igf){
	const u1 = this;
	const keys = [... new Set(
		Array.from(u1.keys()).concat(Array.from(u2.keys()))
	)];
	return keys.every(k => {
		if(igf === true && k === 'from') return true;
		const a1 = u1.getAll(k).sort();
		const a2 = u2.getAll(k).sort();
		if(a1.length !== a2.length) return false;
		return a1.every((e,i) => (e === a2[i]));
	});
};

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

	let feat = drawbox_control.getFeature();
	if(feat !== null){
		new_sp.append('geojson', fmt.writeGeometry(
			feat.getGeometry(), { decimals: 5 }
		));
	}

	let old_sp = new URLSearchParams(window.location.search);
	let nfrom = Number(old_sp.get('from'));

	if(!new_sp.equals(old_sp, true)) nfrom = 0;
	else if(dir === 1) nfrom += Math.max(Number(new_sp.get('size')), 25);
	else if(dir === -1) {
		nfrom = Math.max(nfrom - Math.max(Number(new_sp.get('size')), 25), 0);
	}
	if(nfrom > 0) new_sp.append('from', nfrom);

	fetch(`search.json?${new_sp.toString()}`).then(response => {
		if(!response.ok){ throw 'response not ok'; }
		return response.json();
	}).then(response => {
		if(!new_sp.equals(old_sp)){
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

		document.getElementById('result-csv').href = `search.csv?${new_sp.toString()}`;
		document.querySelector('#result-from').textContent = (response.from + 1);
		document.getElementById('result-geojson').href = `search.geojson?${new_sp.toString()}`;
		document.getElementById('result-pdf').href = `search.pdf?${new_sp.toString()}`;
		document.querySelector('#result-to').textContext = (
			response.from + response.hits.length
		);
		document.querySelector('#result-total').textContent = response.total;
		document.querySelector('#result-prev').disabled = (response.from === 0);
		document.querySelector('#result-next').disabled = (
			(response.from + response.hits.length) >= response.total
		);

		response.hits.forEach(hit => {
			if ('geometries' in hit) {
				let feat = new ol.Feature({
					geometry: new ol.geom.GeometryCollection(
						hit.geometries.map(g => fmt.readFeature(g).getGeometry())
					)
				});
				feat.setProperties(hit, true);
				result_source.addFeature(feat);
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

	// Handle the contents of the search box
	let q = pr.get('q');
	search_control.getSearchBox().value = (q == null ? '' : q);
	pr.delete('q');

	// Handle draw box
	drawbox_control.source.clear();
	let geojson = pr.get('geojson');
	if(geojson !== null){
		drawbox_control.source.addFeature(
			new ol.Feature({geometry: fmt.readGeometry(geojson)})
		);
	}

	// Fetch distinct keys
	let keys = Array.from(pr.keys()).filter((v,i,a) => {
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
						else vals.splice(i, 1);
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
	let popup = new ol.Overlay({
		element: document.querySelector('#popup'),
		autoPan: { animation: { duration: 100 } }
	});
	map.addOverlay(popup);

	map.on('pointermove', e => {
		const features = e.map.getFeaturesAtPixel(e.pixel);
		if (features.some(feature => result_source.hasFeature(feature))) {
			map.getViewport().style.cursor = 'pointer';
		} else {
			map.getViewport().style.cursor = 'inherit';
		}
	});

	map.on('click', e => {
		let fts = [];
		e.map.forEachFeatureAtPixel(e.pixel, (feature) => {
			if (result_source.hasFeature(feature)) {
				fts.push(feature);
			}
		});
		if (!fts.length) return popup.setPosition();
		popup.setPosition(e.coordinate);
		popup.getElement().querySelector('#popup-content').innerHTML = mustache.render(
			document.querySelector('#tmpl-popup').innerHTML,
			fts.map(f => f.getProperties()), {}, ['[[', ']]']
		);
		displayPopupContents(0);
	});

	// The popup's display is initially none to prevent the popup from being briefly
	// displayed when the map first loads. Adding show changes the
	// display to block allowing the popup to be visible when a feature is clicked.
	popup.getElement().classList.add('show');

	function displayPopupContents(d) {
		let el = popup.getElement();
		let tables = Array.from(el.querySelectorAll('#popup-content table'));
		let cur = el.querySelector('#popup-content .show');
		if (!cur) cur = tables[0];
		let idx = tables.indexOf(cur);
		tables.scrollTop = 0;
		if (idx >= 0 && idx < tables.length) {
			tables[idx].classList.toggle('show');
		}
		if (d != 0)	tables[idx + d].classList.toggle('show');
		el.querySelector('#popup-page-number').innerHTML = ((idx + d) + 1) + ' of ' + tables.length;
		el.querySelector('#popup-prev-btn').classList.toggle('visible', idx + d > 0);
		el.querySelector('#popup-next-btn').classList.toggle('visible', idx + d < (tables.length - 1));
	}

	popup.getElement().querySelector('#popup-prev-btn').addEventListener(
		'click', e => displayPopupContents(-1)
	);

	popup.getElement().querySelector('#popup-next-btn').addEventListener(
		'click', e => displayPopupContents(1)
	);

	popup.getElement().querySelector('#popup-closer').addEventListener('click', e => {
		popup.setPosition();
		return false;
	});

	search_control.getSearchBox().addEventListener('focus', e => {
		popup.setPosition();
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
	document.querySelectorAll(
		'select[name^="sort"], select[name^="dir"]'
	).forEach(e => e.addEventListener('change', x => doSearch()));

	// If there's a query string, use that to rebuild the search form
	if(window.location.href.includes('?')) updateFromURL();
	// Refresh search if user moves backwards or forwards in history
	window.addEventListener('popstate', updateFromURL);

	search_control.getSearchBox().focus();
});
