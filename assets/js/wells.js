const popup = document.getElementById('popup');
const content = document.getElementById('popup-content');
const closer = document.getElementById('popup-closer');
const prev_btn = document.getElementById('popup-prev-btn');
const next_btn = document.getElementById('popup-next-btn');

const overlay = new ol.Overlay({
	element: popup,
	autoPan: { animation: { duration: 100 } }
});
let fts = [];

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

popup.style.display = 'block';
let map = new ol.Map({
	target: 'map',
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
	overlays: [ overlay ],
	view: MAP_DEFAULTS.View,
	controls: MAP_DEFAULTS.Controls,
	interactions: MAP_DEFAULTS.Interactions
});

map.on('pointermove', function(e){
	e.map.getTargetElement().style.cursor = (
		e.map.hasFeatureAtPixel(e.pixel) ? 'pointer' : ''
	);
});

//Allows the overlay to be visible.
//The overlay needs to be hidden by default to prevent it being
//displayed at startup
document.getElementById('popup-topbar').style.visibility = 'visible';
popup.style.visibility = 'visible';

//Pagination Code
//Fetch the well data
let currentPage;
let running = false;

function displayOverlayContents(e) {
	if (!running) {
		content.scrollTop = 0;
		running = true;
		if (e instanceof MouseEvent) {
			switch (e.target.id) {
				case "popup-prev-btn":
					if (currentPage > 0) {
						currentPage--;
					}
					break;
				case "popup-next-btn":
					if (currentPage < fts.length - 1) {
						currentPage++;
					}
					break;
			}
			prev_btn.disable = true;
			prev_btn.style.color = '#868686BF';
			next_btn.disable = true;
			next_btn.style.color = '#868686BF';
		} else {
			currentPage = 0;
		}
		if (typeof fts[currentPage] !== "undefined") {
			let well_id = fts[currentPage].get('well_id');
			fetch('detail.json?id=' + well_id)
				.then(response => {
					if (!response.ok) throw new Error(response.status + " " +
						response.statusText);
					return response.json();
				})
				.then(data => {
					data["well_id"] = well_id;
					let t = mustache.render(document.getElementById("tmpl-popup").innerHTML, data, {}, ['[[', ']]']);
					document.getElementById("popup-content").innerHTML = t;
					document.getElementById('popup-topbar').style.visibility = 'visible';
					popup.style.visibility = 'visible';
					if (e instanceof ol.events.Event) {
						overlay.setPosition(e.coordinate);
					}
					document.getElementById("popup-page-number").innerHTML = (currentPage + 1) + " of " + fts.length;
					if (currentPage > 0) {
						prev_btn.style.visibility = 'visible';
					} else {
						prev_btn.style.visibility = 'hidden';
					}
					if (currentPage < (fts.length - 1)) {
						next_btn.style.visibility = 'visible';
					} else {
						next_btn.style.visibility = 'hidden';
					}
					prev_btn.style.color = "#fff";
					next_btn.style.color = "#fff";
					running = false;
				})
				.catch(err => {
					alert(err);
					running = false;
				});
		}
	}
}
//Popup
prev_btn.addEventListener("click", displayOverlayContents);
next_btn.addEventListener("click", displayOverlayContents);
closer.addEventListener("click", function() {
	overlay.setPosition(undefined);
	closer.blur();
	return false;
});

map.on('click', function(e) {
	popup.style.visibility = "hidden";
	document.getElementById('popup-topbar').style.visibility = 'hidden';
	fts = map.getFeaturesAtPixel(e.pixel);
	if (fts.length < 1) {
		overlay.setPosition(undefined);
		return
	}
	overlay.setPosition(e.coordinate);
	displayOverlayContents(e);
});
