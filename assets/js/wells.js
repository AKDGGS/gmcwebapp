const popup = document.getElementById('popup');
const content = document.getElementById('popup-content');
const closer = document.getElementById('popup-closure');
const prevBtn = document.getElementById('prev-btn');
const nextBtn = document.getElementById('next-btn');

const overlay = new ol.Overlay({
	element: popup,
	autoPan: { animation: { duration: 100 } }
});
let fts = [];

let point_layer = new ol.layer.Vector({
	source: new ol.source.Vector(),
	style: MAP_DEFAULTS.WellStyle
});

let label_layer = new ol.layer.Vector({
	source: new ol.source.Vector(),
	renderBuffer: 1e3,
	style: function(f) {
		MAP_DEFAULTS.LabelStyle.getText().setText(f.label);
		return MAP_DEFAULTS.LabelStyle;
	},
	declutter: true
});

//Fetch the markers
fetch('points.json')
	.then(response => {
		if (!response.ok) throw new Error(response.status + " " +
			response.statusText);
		return response.json();
	})
	.then(d => {
		if (!d) {
			console.log("No points returned");
			return;
		}
		let markers = [];
		d.forEach(({
			geog,
			name,
			well_id
		}) => {
			let f = new ol.Feature(new ol.format.WKT().readGeometry(geog));
			f.getGeometry().transform("EPSG:4326", "EPSG:3857");
			f.label = name + ' - ' + well_id;
			f.well_id = well_id;
			f.setId(well_id);
			label_layer.getSource().addFeature(f);
			markers.push(f);
		});
		point_layer.getSource().addFeatures(markers);
	})
	.catch(err => {
		alert(err);
	});

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
				case "prev-btn":
					if (currentPage > 0) {
						currentPage--;
					}
					break;
				case "next-btn":
					if (currentPage < fts.length - 1) {
						currentPage++;
					}
					break;
			}
			prevBtn.disable = true;
			prevBtn.style.color = '#868686BF';
			nextBtn.disable = true;
			nextBtn.style.color = '#868686BF';
		} else {
			currentPage = 0;
		}
		if (typeof fts[currentPage] !== "undefined") {
			let well_id = fts[currentPage].well_id;
			fetch('detail.json?id=' + well_id)
				.then(response => {
					if (!response.ok) throw new Error(response.status + " " +
						response.statusText);
					return response.json();
				})
				.then(data => {
					for (let i = 0; i < data.keywords.length; i++) {
						let arr = data.keywords[i].keywords.toString().split(",");
						let qParams = "";
						for (let j = 0; j < arr.length; j++) {
							qParams = qParams.concat("&keyword=" + arr[j]);
						}
						data.keywords[i].keywords = data.keywords[i]
							.keywords.toString().replaceAll(",", ", ");
						data.keywords[i]["keywordsURL"] = encodeURI("search#q=well_id:" +
							well_id + qParams);
					}
					data["nameURL"] = encodeURI("well/" + well_id);
					data["well_id"] = well_id;

					let t = mustache.render(document.getElementById("tmpl-popup").innerHTML, data, {}, ['[[', ']]']);
					document.getElementById("popup-content").innerHTML = t;
					document.getElementById('popup-topbar').style.visibility = 'visible';
					popup.style.visibility = 'visible';
					if (e instanceof ol.events.Event) {
						overlay.setPosition(e.coordinate);
					}
					pageNumber.innerHTML = (currentPage + 1) + " of " + fts.length;
					if (currentPage > 0) {
						prevBtn.style.visibility = 'visible';
					} else {
						prevBtn.style.visibility = 'hidden';
					}
					if (currentPage < (fts.length - 1)) {
						nextBtn.style.visibility = 'visible';
					} else {
						nextBtn.style.visibility = 'hidden';
					}
					prevBtn.style.color = "#fff";
					nextBtn.style.color = "#fff";
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
prevBtn.addEventListener("click", displayOverlayContents);
nextBtn.addEventListener("click", displayOverlayContents);
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
