document.querySelectorAll('.qa-results').forEach(qr => {
	qr.appendChild(document.createTextNode('Running .. '));
	fetch('qa_count.json?id=' + qr.dataset.id)
	.then(response => {
		if (!response.ok) { throw 'response not ok'; }
		return response.json();
	})
	.then(result => {
		while(qr.lastChild) qr.removeChild(qr.lastChild);

		if(result < 1) {
			let s = document.createElement('span');
			s.className = 'passed';
			s.appendChild(document.createTextNode('PASSED'));
			qr.appendChild(s);
		} else {
			let lvl = qr.dataset.level;

			let s = document.createElement('span');
			s.className = lvl;
			s.appendChild(document.createTextNode(lvl.toUpperCase()));
			qr.appendChild(s)

			qr.appendChild(document.createTextNode(' ('));
			let a = document.createElement('a');
			a.href = '#' + qr.dataset.id;
			a.addEventListener('click', function(event){
				loadReport(qr.dataset.id);
				event.preventDefault();
				return false;
			});
			a.appendChild(document.createTextNode(result));
			qr.appendChild(a);
			qr.appendChild(document.createTextNode(')'));
		}
	})
	.catch(err => {
		while(qr.lastChild) qr.removeChild(qr.lastChild);
		qr.appendChild(document.createTextNode('ERROR'));
		if(window.console){ console.log(err); }
	});
});

var isLoading = false;
function loadReport(id){
	if(isLoading){ return; }
	else { isLoading = true; }

	let out = document.getElementById('qa-output');
	while(out.lastChild) out.removeChild(out.lastChild);
	let img = document.createElement('img');
	img.src = '../img/loader.gif';
	out.appendChild(img);

	fetch('qa_run.json?id=' + id)
	.then(response => {
		if (!response.ok) { throw 'response not ok'; }
		return response.json();
	})
	.then(result => {
		let tbl = document.createElement('table');
		tbl.className = 'qa-table';
		let thead = document.createElement('thead');
		let tr = document.createElement('tr');
		result.columns.forEach(c => {
			let th = document.createElement('th');
			th.appendChild(document.createTextNode(c));
			tr.appendChild(th);
		});
		thead.appendChild(tr);
		tbl.appendChild(thead);

		let tbody = document.createElement('tbody');
		result.rows.forEach(r => {
			let tr = document.createElement('tr');
			r.forEach(k => {
				let td = document.createElement('td');
				td.appendChild(document.createTextNode(k === null ? '' : k));
				tr.appendChild(td);
			});
			tbody.appendChild(tr);
		});
		tbl.appendChild(tbody);

		while(out.lastChild) out.removeChild(out.lastChild);
		if(result.rows.length > 0){ out.appendChild(tbl); }
		else { out.appendChild(document.createTextNode("No results")); }
		isLoading = false;
	})
	.catch(err => {
		while(out.lastChild) out.removeChild(out.lastChild);
		out.appendChild(document.createTextNode('Error loading report'));
		if(window.console){ console.log(err); }
		isLoading = false;
	});
}
