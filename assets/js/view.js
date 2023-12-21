if (document.getElementById('map')){
	let fmt = new ol.format.GeoJSON({
		dataProjection: 'EPSG:4326',
		featureProjection: 'EPSG:3857'
	});
	let content = document.getElementById('popup-content');
	let popup = document.getElementById('popup');
	let overlay = new ol.Overlay({
		element: popup, autoPan: {animation:{duration:100}}
	});
	let template = document.getElementById('tmpl-popup');

	popup.style.display = 'block';
	let styles = {
	'borehole': new ol.style.Style({
		image: new ol.style.Circle({
			radius: 5,
			fill: new ol.style.Fill({
				color: 'rgba(99, 186, 0, 0.25)'
			}),
			stroke: new ol.style.Stroke({
				color: 'rgba(99, 186, 0, 255)',
				width: 3
			})
		})
	}),
	'outcrop': new ol.style.Style({
		image: new ol.style.Circle({
			radius: 5,
			fill: new ol.style.Fill({
				color: 'rgba(230, 177, 1, 0.25)'
			}),
			stroke: new ol.style.Stroke({
				color: 'rgba(230, 177, 1, 255)',
				width: 3
			})
		})
	}),
	'shotline': new ol.style.Style({
			stroke: new ol.style.Stroke({
				color: 'rgba(255, 138, 134, 255)',
				width: 3
			})
	}),
	'well': new ol.style.Style({
		image: new ol.style.Circle({
			radius: 5,
			fill: new ol.style.Fill({
				color: 'rgba(146, 203, 255, 0.25)'
			}),
			stroke: new ol.style.Stroke({
				color: 'rgba(46, 145, 230, 255)',
				width: 3
			})
		})
	}),
	'line_string': new ol.style.Style({
			stroke: new ol.style.Stroke({
				color: 'rgba(44, 126, 167, 255)',
				width: 3
			})
	}),
	'point': new ol.style.Style({
		image: new ol.style.Circle({
			radius: 5,
			fill: new ol.style.Fill({
				color: 'rgba(44, 126, 167, 0.25)'
			}),
			stroke: new ol.style.Stroke({
				color: 'rgba(44, 126, 167, 255)',
				width: 3
			})
		})
	}),
};
const style_map = {
	'borehole_id': 'borehole',
	'outcrop_id': 'outcrop',
	'shotline_id': 'shotline',
	'well_id': 'well',
	'LineString': 'line_string',
	'Point': 'point'
};

	let map = new ol.Map({
		target: 'map',
		overlays: [ overlay ],
		layers: [
			MAP_DEFAULTS.BaseLayers,
			MAP_DEFAULTS.OverlayLayers,
			new ol.layer.Vector({
				style: function(feature) {
					let style_key = Object.keys(style_map).find(key => feature.get(key) || feature.getGeometry().getType() === key);
					if (style_key) {
						feature.setStyle(styles[style_map[style_key]]);
					}
				},
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
		let title_element = document.getElementById('popup-title');
		map.on('click', function(e){
			let fts = map.getFeaturesAtPixel(e.pixel);
			if (fts.length < 1){
				overlay.setPosition(undefined);
				return
			}
			content.innerHTML = '';
			for(const ft of fts){
				console.log(ft)
				if (ft.values_.borehole_id) {
					title_element.textContent = 'Borehole(s)';
				} else if (ft.values_.outcrop_id) {
					title_element.textContent = 'Outcrop(s)';
				} else if (ft.values_.shotline_id) {
					title_element.textContent = 'Shotline(s)';
				} else if (ft.values_.well_id) {
					title_element.textContent = 'Well(s)';
				}
				content.innerHTML += mustache.render(
					template.innerHTML,
					ft.getProperties(), {}, ['[[', ']]']
				);
			}
			overlay.setPosition(e.coordinate);
		});
	}
}

if(document.getElementById('filedrop')) {
	let drop_zone = new FileDrop(document.getElementById('filedrop'), document.getElementById('file-list-container'));
}

if(document.getElementById('latlon')) {
	if(geojson != null && geojson.features.length > 0){
		document.getElementById('latlon').innerHTML = geojson.features[0].geometry.coordinates[0]
		+ ", " + geojson.features[0].geometry.coordinates[1];
	}
}
