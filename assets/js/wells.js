const popup = document.getElementById('popup');
const content = document.getElementById('popup-content');
const closer = document.getElementById('popup-closure');
const prevBtn = document.getElementById('prev-btn');
const nextBtn = document.getElementById('next-btn');
let template = document.getElementById('tmpl-popup');

const overlay = new ol.Overlay({
		element: popup, autoPan: { animation: { duration: 100 }}
});
let fts = [];
let width = 2.5;

function handleError(err) {
	alert(err);
}

let markerStyle = new ol.style.Style({
	image: new ol.style.Circle({
		radius: width * 2,
		fill: new ol.style.Fill({
			color: '#2C7EA740'
		}),
		stroke: new ol.style.Stroke({
			color: '#2C7EA7FF',
			width: width / 1.5
		})
	}),
	zIndex: Infinity
});

let labelStyle = new ol.style.Style({
	text: new ol.style.Text({
		offsetY: -9,
		font: '12px Open Sans',
		fill: new ol.style.Fill({
			color: '#000000FF'
		}),
		backgroundFill: new ol.style.Fill({
			color: '#FFFFFF1A'
		})
	})
});

let wellPointLayer = new ol.layer.Vector({
	source: new ol.source.Vector(),
	style: markerStyle
});

let wellPointLabelLayer = new ol.layer.Vector({
	source: new ol.source.Vector(),
	renderBuffer: 1e3,
	style: function(feature) {
		labelStyle.getText().setText(feature.label);
		return labelStyle;
	},
	declutter: true
});

//Fetch the makers
fetch('points.json')
	.then(response => {
		if (!response.ok) throw new Error(response.status + " " +
			response.statusText);
		return response.json();
	})
	.then(d => {
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
			wellPointLabelLayer.getSource().addFeature(f);
			markers.push(f);
		});
		wellPointLayer.getSource().addFeatures(markers);
	})
	.catch(error => {
		handleError(error);
	});
let map = new ol.Map({
	target: 'map',
	layers: [
		MAP_DEFAULTS.BaseLayers,
		MAP_DEFAULTS.OverlayLayers,
		new ol.layer.Group({
			visible: true,
			layers: [
				wellPointLayer,
				wellPointLabelLayer
			]
		}),
	],
	overlays: [ overlay ],
	view: MAP_DEFAULTS.View,
	controls: ol.control.defaults.defaults({ attribution: false }).extend([
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
	])
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
				.catch(error => {
					handleError(error);
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
