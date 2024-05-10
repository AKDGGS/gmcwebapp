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
	let map = new ol.Map({
		target: 'map',
		overlays: [ overlay ],
		layers: [
			MAP_DEFAULTS.BaseLayers,
			MAP_DEFAULTS.OverlayLayers,
			new ol.layer.Vector({
				style: function(feature) {
					let c = ((f) => {
						if(f.get('borehole_id')) return '99, 186, 0';
						if(f.get('outcrop_id')) return '230, 177, 1';
						if(f.get('shotline_id')) return '255, 138, 134';
						if(f.get('well_id')) return '46, 145, 230';
						return '44, 126, 167';
					})(feature);
					return new ol.style.Style({
						fill: new ol.style.Fill({ color: `rgba(${c}, 0.25)` }),
						stroke: new ol.style.Stroke({
							color: `rgba(${c}, 1)`,
							width: 5
						}),
						image: new ol.style.Circle({
							fill: new ol.style.Fill({ color: `rgba(${c}, 0.25)` }),
							stroke: new ol.style.Stroke({
								color: `rgba(${c}, 1)`,
								width: 2
							}),
							radius: 5
						})
					});
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

	map.on('pointermove', function (e) {
		map.getTargetElement().style.cursor = map.hasFeatureAtPixel(e.pixel)
			? 'pointer'
			: '';
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
		for (let i = 0; i < geojson.features.length; i++) {
			if (geojson.features[i].geometry.coordinates.length == 2){
				document.getElementById('latlon').innerHTML = geojson.features[i].geometry.coordinates[0]
				+ ", " + geojson.features[i].geometry.coordinates[1];
			}
		}

	}
}

if (document.getElementById('stash-button')){
	let stash = document.getElementById('stash-button');
	let dest = document.getElementById('stash-dest');
	stash.addEventListener('click', (e) => {
		if(stash.disabled) return false;

		if(dest.classList.contains('shown')){
			stash.innerText = 'Show Stash';
			dest.classList.remove('shown');
		} else {
			if(dest.lastChild){
				stash.innerText = 'Hide Stash';
				dest.classList.add('shown');
				return false;
			}

			stash.disabled = true;
			let hr = window.location.href;
			fetch(
				`stash.json?id=${hr.substr(hr.lastIndexOf('/')+1)}`
			).then(response => {
				if (!response.ok) throw 'response not ok';
				return response.json();
			}).then(result => {
				dest.appendChild(JSONToElement(result));
				stash.innerText = 'Hide Stash';
				dest.classList.add('shown');
				stash.disabled = false;
			}).catch(err => {
				if(window.console) console.log(err);
			});
		}
		return false;
	});
}

// Takes an arbitrary JSON object
// and renders it into a table, returning
// the table element
function JSONToElement(obj){
	let type = Object.prototype.toString.call(obj);
	switch(type){
		case '[object Boolean]':
			return document.createTextNode(obj.toString());
		case '[object String]':
			return document.createTextNode(obj);
		case '[object Number]':
			return document.createTextNode(obj.toString());
		case '[object Null]':
			return document.createTextNode('(null)');
		case '[object Object]': {
			let tbl = document.createElement('table');
			let count = 0;
			for(let i in obj){
				let tr = document.createElement('tr');

				let th = document.createElement('th');
				th.appendChild(document.createTextNode(i));
				tr.appendChild(th);

				let td = document.createElement('td');
				td.appendChild(JSONToElement(obj[i]));
				tr.appendChild(td);

				tbl.appendChild(tr);
				count++;
			}
			if(count > 0) return tbl;
			else return document.createTextNode('(Empty Object)');
		}
		case '[object Array]': {
			if(obj.length < 1) return document.createTextNode('(Empty List)');

			let tbl = document.createElement('table');
			for(let i = obj.length; i--;){
				let tr = document.createElement('tr');

				let th = document.createElement('th');
				th.appendChild(document.createTextNode(i));
				tr.appendChild(th);

				let td = document.createElement('td');
				td.appendChild(JSONToElement(obj[i]));
				tr.appendChild(td);

				tbl.appendChild(tr);
			}
			return tbl;
		}
		default:
			return document.createTextNode('Unknown - ' + type);
	}
}
