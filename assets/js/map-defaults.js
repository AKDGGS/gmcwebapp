const MAP_DEFAULTS = {
	BaseLayers: new ol.layer.Group({
		title: 'Base Maps',
		layers: [
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
	}),
	Style: function(feature) {
		let c = ((f) => {
			if(f.get('borehole_id') || f.get('borehole')){
				return '99, 186, 0';
			}
			if(f.get('outcrop_id') || f.get('outcrop')){
				return '230, 177, 1';
			}
			if(f.get('shotline_id') || f.get('shotline')){
				return '255, 138, 134';
			}
			if(f.get('well_id') || f.get('well')){
				return '46, 145, 230';
			}
			return '44, 126, 167';
		})(feature);
		return new ol.style.Style({
			fill: new ol.style.Fill({
				color: `rgba(${c}, 0.25)`
			}),
			stroke: new ol.style.Stroke({
				color: `rgba(${c}, 1)`,
				width: 2
			}),
			image: new ol.style.Circle({
				fill: new ol.style.Fill({
					color: `rgba(${c}, 0.25)`
				}),
				stroke: new ol.style.Stroke({
					color: `rgba(${c}, 1)`,
					width: 2
				}),
				radius: 5
			})
		});
	}
};
