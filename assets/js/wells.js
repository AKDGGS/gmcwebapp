let source = new ol.source.Vector({
	url: 'points.json',
	format: new ol.format.GeoJSON({
		dataProjection: 'EPSG:4326',
		featureProjection: 'EPSG:3857'
	})
});

let map = new ol.Map({
	target: 'map',
	controls: MAP_DEFAULTS.Controls,
	interactions: MAP_DEFAULTS.Interactions,
	view: MAP_DEFAULTS.View,
	layers: [
		MAP_DEFAULTS.BaseLayers,
		MAP_DEFAULTS.OverlayLayers,
		new ol.layer.Group({
			visible: true,
			layers: [
				new ol.layer.Vector({
					source: source,
					style: MAP_DEFAULTS.WellStyle
				}),
				new ol.layer.Vector({
					source: source,
					renderBuffer: 1e3,
					style: MAP_DEFAULTS.LabelStyle,
					declutter: true
				})
			]
		}),
	]
});

let popup = new ol.Overlay({
	element: document.querySelector('#popup'),
	autoPan: { animation: { duration: 100 } }
});
map.addOverlay(popup);

map.on('pointermove', e => {
	let t = e.map.getFeaturesAtPixel(e.pixel).some(f => source.hasFeature(f));
	map.getViewport().style.cursor = (t ? 'pointer' : 'inherit');
});

let fts = [];
map.on('click', e => {
	popup.setPosition();
	delete popup.getElement().querySelector('#popup-content').dataset.well_id;
	popup.getElement().classList.add('show');
	fts = e.map.getFeaturesAtPixel(e.pixel);
	if(fts.length){
		popup.setPosition(e.coordinate);
		displayPopupContents(0);
	}
});

function displayPopupContents(d) {
	let ct = popup.getElement().querySelector('#popup-content');
	let idx = Math.max(fts.findIndex(e => e.get('well_id') == ct.dataset.well_id), 0);
	switch(d){
		case -1: if(fts[idx-1] !== undefined) idx--; break;
		case 1: if(fts[idx+1] !== undefined) idx++; break;
	}

	if(fts[idx].get('_popup') === undefined){
		fetch(`detail.json?id=${fts[idx].get('well_id')}`)
		.then(response => {
			if(!response.ok) throw new Error(`Error ${response.status}: ${response.statusText}`);
			return response.json();
		}).then(data => {
			fts[idx].set('_popup', mustache.render(
				document.querySelector('#tmpl-popup').innerHTML,
				data, {}, ['[[', ']]']
			), true);
			displayPopupContents(d);
		}).catch(err => {
			if(window.console) console.log(err);
		});
		return;
	}

	ct.innerHTML = fts[idx].get('_popup');
	ct.dataset.well_id = fts[idx].get('well_id');
	document.querySelector('#popup-page-number').innerHTML = (idx+1) + ' of ' + fts.length;
	document.querySelector('#popup-prev-btn').classList.toggle('visible', idx > 0);
	document.querySelector('#popup-next-btn').classList.toggle('visible', (idx+1) < fts.length);
}

popup.getElement().querySelector('#popup-prev-btn').addEventListener(
	'click', e => displayPopupContents(-1)
);

popup.getElement().querySelector('#popup-next-btn').addEventListener(
	'click', e => displayPopupContents(1)
);

popup.getElement().querySelector('#popup-closer').addEventListener('click', e => {
	popup.getElement().querySelector('#popup-content').innerHTML = '';
	popup.setPosition();
	return false;
});
