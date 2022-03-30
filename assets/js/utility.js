// Takes an arbitrary JSON object
// and renders it into a table, returning
// the table element
function JSONToElement(obj){
	var type = Object.prototype.toString.call(obj);

	switch(type){
		case '[object Boolean]':
			return document.createTextNode(obj.toString());

		case '[object String]':
			return document.createTextNode(obj);

		case '[object Number]':
			return document.createTextNode(obj.toString());

		case '[object Null]':
			return document.createTextNode('(null)');

		case '[object Object]':
			var tbl = document.createElement('table');
			var count = 0;
			for(var i in obj){
				var tr = document.createElement('tr');

				var th = document.createElement('th');
				th.appendChild(document.createTextNode(i));
				tr.appendChild(th);

				var td = document.createElement('td');
				td.appendChild(JSONToElement(obj[i]));
				tr.appendChild(td);

				tbl.appendChild(tr);
				count++;
			}
			if(count > 0) return tbl;
			else return document.createTextNode('(Empty Object)');

		case '[object Array]':
			if(obj.length < 1) return document.createTextNode('(Empty List)');

			var tbl = document.createElement('table');
			for(var i = obj.length; i--;){
				var tr = document.createElement('tr');

				var th = document.createElement('th');
				th.appendChild(document.createTextNode(i));
				tr.appendChild(th);

				var td = document.createElement('td');
				td.appendChild(JSONToElement(obj[i]));
				tr.appendChild(td);

				tbl.appendChild(tr);
			}
			return tbl;

		default:
			return document.createTextNode('Unknown - ' + type);
	}
}
