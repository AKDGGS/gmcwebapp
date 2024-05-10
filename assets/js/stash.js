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

let stash = document.getElementById('stash-button');
if (stash !== null) {
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
				`../stash.json?id=${hr.substr(hr.lastIndexOf('/')+1)}`
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
