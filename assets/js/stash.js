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


var stash = document.getElementById('stash-link');
if (stash !== null) {
  stash.onclick = function(evt) {
    var anchor = this;
    if (this.innerHTML === 'Show Stash') {
      var xhr = (window.ActiveXObject ? new ActiveXObject('Microsoft.XMLHTTP') : new XMLHttpRequest());
      xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && xhr.status === 200) {
          var json = JSON.parse(xhr.responseText);
          var el = JSONToElement(json.stash);
          document.getElementById('stash-dd').appendChild(el);
          anchor.innerHTML = 'Hide Stash';
        }
      };
			var inv_id = window.location.href.substr(window.location.href.lastIndexOf("/")+1);
      xhr.open('GET', '../stash.json?id=' + inv_id, true);
      xhr.send();
    } else {
      anchor.innerHTML = 'Show Stash';
      var dd = document.getElementById('stash-dd');
      if (dd !== null) {
        while (dd.lastChild) dd.removeChild(dd.lastChild);
      }
    }
    var e = evt === undefined ? window.event : evt;
    if ('preventDefault' in e) e.preventDefault();
    return false;
  };
}
