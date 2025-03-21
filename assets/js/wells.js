let source = new ol.source.Vector({
	url: 'points.json',
	format: new ol.format.GeoJSON({
		dataProjection: 'EPSG:4326',
		featureProjection: 'EPSG:3857'
	})
});

let point_layer = new ol.layer.Vector({
	source: source,
	style: MAP_DEFAULTS.WellStyle
});

let label_layer = new ol.layer.Vector({
	source: source,
	renderBuffer: 1e3,
	style: function(f) {
		MAP_DEFAULTS.LabelStyle.getText().setText(
			f.get('name') + (f.get('number') == undefined ? '' : ' - ') +
			f.get('number')
		);
		return MAP_DEFAULTS.LabelStyle;
	},
	declutter: true
});

const popup = new ol.Overlay({
	element: document.querySelector('#popup'),
	autoPan: { animation: { duration: 100 } }
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
				point_layer,
				label_layer
			]
		}),
	],
	overlays: [ popup ]
});

map.on('pointermove', e => {
	if (e.map.getFeaturesAtPixel(e.pixel).some(f => source.hasFeature(f))) {
		map.getViewport().style.cursor = 'pointer';
	} else {
		map.getViewport().style.cursor = 'inherit';
	}
});

let fts = new Map();

map.on('click', e => {
	fts = new Map();
	popup.getElement().querySelector('#popup-content').innerHTML = '';
	e.map.forEachFeatureAtPixel(e.pixel, f => {
		fts.set(f.get('well_id'), null);
	});
	if (fts.size) {
		if (!popup.getElement().classList.contains('show')) {
			popup.getElement().classList.add('show');
		}
		displayPopupContents(0);
		popup.setPosition(e.coordinate);
	} else {
		popup.setPosition();
	}
});

function displayPopupContents(d) {
	let fts_keys = Array.from(fts.keys());
	if (fts.get(fts_keys[0]) === null) {idx = 0;}
	let el = popup.getElement();
	try {
		if (!fts.get(fts_keys[idx + d])) {
			fetch('detail.json?id=' + fts_keys[idx + d])
			.then(response => {
				if (!response.ok) throw new Error(response.status + " " + response.statusText);
				return response.json();
			}).then(data => {
				fts.set(fts_keys[idx + d], data);
				el.querySelector('#popup-content').innerHTML = mustache.render(
					document.querySelector('#tmpl-popup').innerHTML, data, {}, ['[[', ']]']
				);
				el.querySelector('#popup-content table').classList.toggle('show');
				el.querySelector('#popup-page-number').innerHTML = ((idx + d) + 1) + ' of ' + fts.size;
				el.querySelector('#popup-prev-btn').classList.toggle('visible', idx + d > 0);
				el.querySelector('#popup-next-btn').classList.toggle('visible', idx + d < (fts.size - 1));
				idx += d;
			}).catch(error => {
				if (window.console) console.log(error);
			});
		} else {
			el.querySelector('#popup-content').innerHTML = mustache.render(
				document.querySelector('#tmpl-popup').innerHTML, fts.get(fts_keys[idx + d]), {}, ['[[', ']]']
			);
			el.querySelector('#popup-content table').classList.toggle('show');
			el.querySelector('#popup-page-number').innerHTML = ((idx + d) + 1) + ' of ' + fts.size;
			el.querySelector('#popup-prev-btn').classList.toggle('visible', idx + d > 0);
			el.querySelector('#popup-next-btn').classList.toggle('visible', idx + d < (fts.size - 1));
			idx += d;
		}
	} catch (error) {console.error(error);}
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
