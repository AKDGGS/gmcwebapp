const popup = document.getElementById('popup');
const content = document.getElementById('popup-content');
const closer = document.getElementById('popup-closure');
const prevBtn = document.getElementById('prevBtn');
const nextBtn = document.getElementById('nextBtn');
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
			color: 'rgba(44, 126,167, 0.25)'
		}),
		stroke: new ol.style.Stroke({
			color: 'rgba(44, 126,167, 255)',
			width: width / 1.5
		})
	}),
	zIndex: Infinity
});

let labelStyle = new ol.style.Style({
	text: new ol.style.Text({
		offsetY: -9,
		font: '13px Calibri, sans-serif',
		fill: new ol.style.Fill({
			color: 'rgba(0,0,0,255)'
		}),
		backgroundFill: new ol.style.Fill({
			color: 'rgba(255, 255, 255, 0.1)'
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
fetch('../well_points.json')
	.then(response => {
		if (!response.ok) throw new Error(response.status + " " +
			response.statusText);
		return response.json();
	})
	.then(d => {
		let markers = [];
		d.points.forEach(({
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
		new ol.layer.Group({
			title: 'Base Maps',
			layers: [
				new ol.layer.Tile({
					title: 'Stamen Watercolor',
					type: 'base',
					source: new ol.source.Stamen({
						layer: 'watercolor'
					})
				}),
				new ol.layer.Tile({
					title: 'ESRI National Geographic',
					type: 'base',
					source: new ol.source.XYZ({
						url: '//server.arcgisonline.com/ArcGIS/rest/services/NatGeo_World_Map/MapServer/tile/{z}/{y}/{x}'
					})
				}),
				new ol.layer.Tile({
					title: 'ESRI DeLorme',
					type: 'base',
					source: new ol.source.XYZ({
						url: '//server.arcgisonline.com/ArcGIS/rest/services/Specialty/DeLorme_World_Base_Map/MapServer/tile/{z}/{y}/{x}'
					})
				}),
				new ol.layer.Tile({
					title: 'ESRI Shaded Relief',
					type: 'base',
					source: new ol.source.XYZ({
						url: '//server.arcgisonline.com/ArcGIS/rest/services/World_Shaded_Relief/MapServer/tile/{z}/{y}/{x}'
					})
				}),
				new ol.layer.Tile({
					title: 'ESRI Topographic',
					type: 'base',
					source: new ol.source.XYZ({
						url: '//server.arcgisonline.com/ArcGIS/rest/services/World_Topo_Map/MapServer/tile/{z}/{y}/{x}'
					})
				}),
				new ol.layer.Tile({
					title: 'ESRI Imagery',
					type: 'base',
					source: new ol.source.XYZ({
						url: '//server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}'
					})
				}),
				new ol.layer.Tile({
					title: 'OpenTopoMap',
					type: 'base',
					source: new ol.source.XYZ({
						url: '//server.arcgisonline.com/ArcGIS/rest/services/World_Topo_Map/MapServer/tile/{z}/{y}/{x}'
					})
				}),
				new ol.layer.Tile({
					title: 'Open Street Maps',
					type: 'base',
					source: new ol.source.OSM()
				}),
			]
		}),
		new ol.layer.Group({
			visible: true,
			layers: [
				wellPointLayer,
				wellPointLabelLayer
			]
		}),
		new ol.layer.Group({
			title: 'Overlays',
			layers: [
				new ol.layer.Image({
					title: 'PLSS (BLM)',
					visible: false,
					source: new ol.source.ImageWMS({
						url: 'https://maps.dggs.alaska.gov/arcgis/services/apps/plss/MapServer/WMSServer',
						params: {
							"LAYERS": '1,2,3',
							"TRANSPARENT": true,
							"FORMAT": 'image/png'
						},
						serverType: 'mapserver',
					})
				}),
				new ol.layer.Image({
					title: 'Quadrangles',
					visible: false,
					source: new ol.source.ImageWMS({
						url: 'https://maps.dggs.alaska.gov/arcgis/services/apps/Quad_Boundaries/MapServer/WMSServer',
						params: {
							"LAYERS": '1,2,3',
							"TRANSPARENT": true,
							"FORMAT": 'image/png'
						},
						serverType: 'mapserver',
					})
				}),
			]
		}),
	],
	overlays: [overlay],
	view: new ol.View({
		center: ol.proj.fromLonLat([ -147.77, 64.83 ]),
		zoom: 3,
		maxZoom: 19
	}),
	controls: ol.control.defaults({ attribution: false }).extend([
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
	interactions: ol.interaction.defaults({ mouseWheelZoom: false }).extend([
		new ol.interaction.MouseWheelZoom({
			condition: ol.events.condition.platformModifierKeyOnly
		})
	])
});

//Allows the overlay to be visible.
//The overlay needs to be hidden by default to prevent it being
//displayed at startup
document.getElementById('topBar').style.visibility = 'visible';
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
				case "prevBtn":
					if (!(currentPage < 0)) {
						currentPage--;
					}
					break;
				case "nextBtn":
					if (!(currentPage > fts.length - 1)) {
						currentPage++;
					}
					break;
			}
			prevBtn.disable = true;
			prevBtn.style.color = "rgba(134, 134, 134, 0.75)";
			nextBtn.disable = true;
			nextBtn.style.color = "rgba(134,134,134, 0.75)";
		} else {
			currentPage = 0;
		}
		let well_id = fts[currentPage].well_id;
		fetch('../well.json?id=' + well_id)
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
				prevBtn.style.color = "#ffffff";
				nextBtn.style.color = "#ffffff";
				running = false;
			})
			.catch(error => {
				handleError(error);
				running = false;
			});
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
	fts = map.getFeaturesAtPixel(e.pixel);
	if (fts.length < 1) {
		overlay.setPosition(undefined);
		return
	}
	displayOverlayContents(e);
	overlay.setPosition(e.coordinate);
});
