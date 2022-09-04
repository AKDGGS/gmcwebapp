const MAP_DEFAULTS = {
	BaseLayers: new ol.layer.Group({
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
	OverlayLayers: new ol.layer.Group({
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
	View: new ol.View({
		center: ol.proj.fromLonLat([ -147.77, 64.83 ]),
		zoom: 3,
		maxZoom: 19
	})
};
