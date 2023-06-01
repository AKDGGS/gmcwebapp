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
			MAP_DEFAULTS.BaseLayers,
			MAP_DEFAULTS.OverlayLayers,
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
		view: MAP_DEFAULTS.View,
		controls: ol.control.defaults.defaults({
			attribution: false
		}).extend([
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
		interactions: ol.interaction.defaults.defaults({
			mouseWheelZoom: false
		}).extend([
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

if(document.querySelector('.filedrop')){
	FileDropInit(document.querySelector('.filedrop'));
	let filedrop = new FileDrop(document.querySelector('.filedrop'));
}
