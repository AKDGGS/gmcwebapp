if (document.getElementById('map')){
	var fmt = new ol.format.GeoJSON({
		dataProjection: 'EPSG:4326',
		featureProjection: 'EPSG:3857'
	});
	var map = new ol.Map({
		target: 'map',
		layers: [
			new ol.layer.Tile({
				source: new ol.source.OSM()
			}),
			new ol.layer.Vector({
				style: new ol.style.Style({
					image: new ol.style.Circle({
						radius: 5,
						fill: new ol.style.Fill({
							color: 'rgba(44, 126,167, 0.25)'
						}),
						stroke: new ol.style.Stroke({
							color: 'rgba(44, 126,167, 255)',
							width: 2
						})
					})
				}),
				source: new ol.source.Vector({
					features: geojsons.map(geojson => {
						return fmt.readFeature(geojson);
					})
				})
			})
		],
		view: new ol.View({
			center: ol.proj.fromLonLat([ -147.77, 64.83 ]),
			zoom: 3,
			maxZoom: 19
		}),
		controls: ol.control.defaults({ attribution: false }).extend([
			new ol.control.ScaleLine({ units: "us" })
		]),
		interactions: ol.interaction.defaults({ mouseWheelZoom: false }).extend([
			new ol.interaction.MouseWheelZoom({
				condition: ol.events.condition.platformModifierKeyOnly
			})
		])
	});
}
