if(typeof geojson !== 'undefined'){
	if(document.getElementById('map')){
		let fmt = new ol.format.GeoJSON({
			dataProjection: 'EPSG:4326',
			featureProjection: 'EPSG:3857'
		});
		let content = document.getElementById('popup-content');
		let popup = document.getElementById('popup');
		let overlay = new ol.Overlay({
			element: popup,
			autoPan: { animation: { duration: 100 } }
		});
		let template = document.getElementById('tmpl-popup');
		popup.classList.add('show');

		popup.style.display = 'block';
		let map = new ol.Map({
			target: 'map',
			overlays: [ overlay ],
			layers: [
				MAP_DEFAULTS.BaseLayers,
				MAP_DEFAULTS.OverlayLayers,
				new ol.layer.Vector({
					style: MAP_DEFAULTS.DynamicStyle,
					source: new ol.source.Vector({
						features: fmt.readFeatures(geojson)
					})
				})
			],
			view: MAP_DEFAULTS.View,
			controls: MAP_DEFAULTS.Controls,
			interactions: MAP_DEFAULTS.Interactions
		});

		map.on('pointermove', function(e){
			e.map.getTargetElement().style.cursor = (
				e.map.hasFeatureAtPixel(e.pixel) ? 'pointer' : ''
			);
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
				document.querySelector("#popup-content table").classList.add('show');
				overlay.setPosition(ol.extent.getCenter(
					fts[0].getGeometry().getExtent()
				));
			});
		}
	}

	if(document.getElementById('latlon')){
		let latlon = document.getElementById('latlon');
		let findCoordinates = (x) => {
			let arr = [];
			switch(Object.prototype.toString.call(x)){
				case '[object Object]':
					for(const [k,v] of Object.entries(x)){
						if(k === 'coordinates' && Array.isArray(v) && v.length == 2){
							arr.push(v);
						} else arr = arr.concat(findCoordinates(v));
					}
				break;
				case '[object Array]':
					for(const v of x) arr = arr.concat(findCoordinates(v));
				break;
			}
			return arr;
		};
		findCoordinates(geojson).forEach((v,i) => {
			latlon.innerHTML += `${i>0?'<br>':''}${v[0]}, ${v[1]}`;
		});
	}
}

if(document.getElementById('filedrop')){
	let drop_zone = new FileDrop(
		document.getElementById('filedrop'),
		document.getElementById('file-list-container')
	);
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
