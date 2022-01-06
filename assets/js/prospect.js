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
	popup.style.display = 'block';
	let map = new ol.Map({
		target: 'map',
		overlays: [overlay],
		layers: [
			new ol.layer.Tile({
				source: new ol.source.OSM()
			}),
			new ol.layer.Vector({
				style: new ol.style.Style({
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

	map.on('click', function(e){
		let fts = map.getFeaturesAtPixel(e.pixel);
		if (fts.length < 1){
			overlay.setPosition(undefined);
			return
		}
		content.innerHTML = '';
		for(const ft of fts){
			content.innerHTML += mustache.render(
				document.getElementById('tmpl-popup').innerHTML,
				ft.getProperties(), {}, ['[[', ']]']
			);
		}
		overlay.setPosition(e.coordinate);
	});
}
