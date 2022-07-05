if (document.getElementById('map')){
	let fmt = new ol.format.GeoJSON({
		dataProjection: 'EPSG:4326',
		featureProjection: 'EPSG:3857'
	});
	let content = document.getElementById('popup-content');
	let popup = document.getElementById('popup');
	let overlay = new ol.Overlay({
		element: popup, autoPan: { animation: { duration: 100 }}
	});
	let template = document.getElementById('tmpl-popup');

	popup.style.display = 'block';
	let map = new ol.Map({
		target: 'map',
		overlays: [ overlay ],
		layers: [
			new ol.layer.Group({
				title: 'Base Maps',
				layers: [
					new ol.layer.Tile({
						title: 'Stamen Watercolor',
						type: 'base',
						source: new ol.source.Stamen({ layer: 'watercolor' })
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
					})
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
			new ol.layer.Vector({
				style: new ol.style.Style({
					fill: new ol.style.Fill({
						color: 'rgba(44, 126, 167, 0.25)'
					}),
					stroke: new ol.style.Stroke({
						color: 'rgba(44, 126, 167, 255)',
						width: 5
					}),
					image: new ol.style.Circle({
						radius: 5,
						fill: new ol.style.Fill({
							color: 'rgba(44, 126, 167, 0.25)'
						}),
						stroke: new ol.style.Stroke({
							color: 'rgba(44, 126, 167, 255)',
							width: 2
						})
					})
				}),
				source: new ol.source.Vector({
					features: fmt.readFeatures(geojson)
				})
			})
		],
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

	let closer = document.getElementById('popup-closer');
	closer.addEventListener("click", function(){
		overlay.setPosition(undefined);
		return false;
	});

	// Only enable popups if a template is provided
	if(template != null){
		map.on('click', function(e){
			let fts = map.getFeaturesAtPixel(e.pixel);
			if (fts.length < 1){
				overlay.setPosition(undefined);
				return
			}
			content.innerHTML = '';
			for(const ft of fts){
				content.innerHTML += mustache.render(
					template.innerHTML,
					ft.getProperties(), {}, ['[[', ']]']
				);
			}
			overlay.setPosition(e.coordinate);
		});
	}
}
