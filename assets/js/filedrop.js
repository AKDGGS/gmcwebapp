if(document.querySelector('.filedrop') && document.getElementById('file-list-container')){
	let drop_zone = new FileDrop(
		document.querySelector('.filedrop'),
		document.getElementById('file-list-container')
	);
}

function FileDrop(drop_zone, file_list_container){
	// Add header
	let header = document.createElement('div');
	header.className = 'header';
	const upload_link_span = document.createElement('span');
	upload_link_span.className = 'center-text upload-text';

	const first_text = document.createTextNode('Drag or ');
	const upload_link = document.createElement('a');
	upload_link.href = '#';
	upload_link.className = 'upload-link';
	upload_link.textContent = 'select';
	const second_text = document.createTextNode(' to upload files');
	upload_link_span.append(first_text, upload_link, second_text);
	header.appendChild(upload_link_span);
	drop_zone.appendChild(header);

	// Add manual selection
	let file_input = document.createElement('input');
	file_input.style.display = 'none';
	file_input.type = 'file';
	file_input.className = 'file-input';
	file_input.multiple = true;
	drop_zone.appendChild(file_input);

	// Add progress bars
	let upload_total = document.createElement('div');
	upload_total.className = 'progress-container';
	let upload_total_fill = document.createElement('div');
	upload_total_fill.className = 'progress-fill';
	let upload_total_text = document.createElement('span');
	upload_total_text.className = 'progress-text';
	drop_zone.appendChild(upload_total);
	upload_total.appendChild(upload_total_fill);
	upload_total.appendChild(upload_total_text);

	let upload_status = document.createElement('div');
	upload_status.className = 'progress-container';
	let upload_status_fill = document.createElement('div');
	upload_status_fill.className = 'progress-fill';
	let upload_status_text = document.createElement('span');
	upload_status_text.className = 'progress-text';
	upload_status.appendChild(upload_status_fill);
	upload_status.appendChild(upload_status_text);
	drop_zone.appendChild(upload_status);
	
	// Add error div
	let error_div = document.createElement('div');
	error_div.className = 'error-div';
	let error_span = document.createElement('span');
	error_span.className = 'center-text';
	error_span.textContent = 'These files failed to upload with the following error: ';
	error_div.appendChild(error_span);
	drop_zone.appendChild(error_div);
	error_div.style.display = 'none';

	// Move file list to end of div
	drop_zone.appendChild(file_list_container)

	let file_id = 0;
	file_name = '';
	let file_count = 0;
	let start_time = Date.now();

	// Test for "XHR2" support - AJAX style uploads
	var xhr2 = ('XMLHttpRequest' in window) && ('upload' in new XMLHttpRequest());

	// File drag and drop API?
	var filedrop = document.querySelector('.filedrop') 
	if(xhr2 && 'ondragover' in filedrop && 'ondrop' in filedrop){
		filedrop.style.display = 'block';
		filedrop.addEventListener('dragover', hoverHandler);
		filedrop.addEventListener('dragleave', hoverHandler);
		filedrop.addEventListener('drop', selectHandler);
	}

	upload_link.addEventListener('click', (e) => {
		e.preventDefault();
		let error_div = document.querySelector('.error-div');
		if (error_div){
			error_div.style.display = 'none';
			error_div.failed_files = [];
			let le = error_div.querySelector('ul');
			if (le){
				error_div.removeChild(le);
			}
		}
		file_input.click();
	});

	file_input.addEventListener('change', selectHandler);
	function statusUpdate(loaded, total, name){
		var filestatus = document.getElementById('file-status');

		// If the function is called with no arguments
		// clear the status
		if(typeof loaded === 'undefined' && typeof total === 'undefined' && name === 'undefined'){
			while(filestatus.lastChild)
			filestatus.removeChild(filestatus.lastChild);
			return;
		}

		let elapsed_time_file = (Date.now() - start_time) / 1000;
		let rate = size_formatter(Math.round(loaded / elapsed_time_file)) + '/s';

		upload_status_fill.style.width = (loaded / total) * 100 + '%';
		upload_status_text.textContent = 'Uploaded ' + size_formatter(loaded) + ' of ' + size_formatter(total) + ' ' + rate;
	}

	function hoverHandler(event){
		event.stopPropagation();
		event.preventDefault();
		filedrop.classList.toggle('drag', event.type === 'dragover');
	}

	function selectHandler(event){	
		upload_total.style.display = 'block';
		upload_status.style.display = 'block';
		file_count = 0;
		hoverHandler(event);

		var files = event.target.files || event.dataTransfer.files;
		let total_size = Array.from(files).reduce((a, b) => a + b.size, 0);
		let total_loaded = 0;
		
		for(var i = 0; i < files.length; i++){
			let xhr = new XMLHttpRequest();
			let pq = ['borehole_id', 'inventory_id', 'outcrop_id', 'prospect_id', 'shotline_id', 'well_id'].filter(v => drop_zone.dataset[v]).map(v => `${v}=${drop_zone.dataset[v]}`).join('&');
			let upload_url = '../upload' + (pq ? '?' + pq : '');
			let form_data = new FormData();
			form_data.append('file'+i, files[i]);
			form_data.append('format', 'json');
			file_name = files[i].name;
			let file = files[i];

			upload_total_fill.style.width = (file_count / files.length) * 100 + '%';
			upload_total_text.textContent = '0 / ' + files.length + " Files Transferred";
			statusUpdate();
			start_time = Date.now();
			xhr.upload.addEventListener(
				'progress',
				function(event){
					statusUpdate(event.loaded, event.total, file.name);
				},
				false
			);
			xhr.addEventListener('load', (event) => {
				if (xhr.status == 200){
					try {
						data = JSON.parse(xhr.responseText);
					} catch(e) {
						upload_error('Malformed JSON response', file.name);
						return;
					}
					file_id = data?.[0]?.file_id;
					if (file_id == undefined){
						upload_error('No file_id in response', file.name);
						return;
					}
				} else {
					upload_error(`HTTP error ${xhr.status}`, file.name);
					return;
				}

				xhr.addEventListener('error', (event) => {
					upload_error("error reading xml");
					return;
				});
				file_count += 1;
				total_loaded += file.size;
				let elapsed_time_total = (Date.now() - start_time) / 1000;

				upload_total_fill.style.width = (file_count / files.length) * 100 + '%';
				upload_total_text.textContent = file_count + ' / ' + files.length + " Files Transferred\n" 
					+ size_formatter(Math.round(total_size / elapsed_time_total)) + '/s';
				add_file_to_page(file);
			});
			xhr.open('POST', upload_url);
			xhr.send(form_data);
		}
	}

	function size_formatter(size){
		var i = size == 0 ? 0 : Math.floor(Math.log(size) / Math.log(1024));
		return (size / Math.pow(1024, i)).toFixed(2) * 1 + ' ' + ['B', 'KiB', 'MiB', 'GiB', 'TiB'][i];
	}

	function add_file_to_page(file){
		const filename_container = document.createElement('div');
		const file_link = document.createElement('a');
		file_link.href = `../file/${file_id}`;
		file_link.textContent = file.name;
		filename_container.appendChild(file_link);
		const file_size = document.createTextNode(` (${size_formatter(file.size)})`);
		filename_container.appendChild(file_size);
		file_list_container.appendChild(filename_container);
	}

	function upload_error(name, err){
		let error_div = document.querySelector('.error-div');
		error_div.style.display = 'flex';
		let le = document.createElement('ul');
		let list_item = document.createElement('li');
		list_item.textContent = err + ': ' + name;
		le.appendChild(list_item);
		error_div.appendChild(le);
	}
}