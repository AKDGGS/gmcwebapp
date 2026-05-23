function FileDrop(drop_zone, file_list_container) {
	let div_header = document.createElement('div');
	div_header.className = 'header';
	drop_zone.appendChild(div_header);

	let upload_link = document.getElementById('upload-link');

	let filedrop_container = document.createElement('div');
	filedrop_container.className = 'container';
	drop_zone.appendChild(filedrop_container);

	// Add manual selection
	let file_input = document.createElement('input');
	file_input.type = 'file';
	file_input.className = 'file-input';
	file_input.multiple = true;
	filedrop_container.appendChild(file_input);

	//Add error div
	let error_div = document.createElement('div');
	error_div.className = 'error-div';
	let error_span = document.createElement('span');
	error_span.className = 'center-text';
	error_span.textContent = 'These files failed to upload with the following error: ';
	error_div.appendChild(error_span);
	filedrop_container.appendChild(error_div);

	//Add first progress bar
	let pb_file = document.createElement('div');
	pb_file.className = 'pb pb-file';

	let pb_file_progress = document.createElement('span');
	pb_file_progress.className = 'pb-text pb-file-progress';
	pb_file.appendChild(pb_file_progress);

	let pb_file_name = document.createElement('span');
	pb_file_name.className = 'pb-text pb-file-name';
	pb_file.appendChild(pb_file_name);

	filedrop_container.appendChild(pb_file);

	//Add second progress bar
	let pb_total = document.createElement('div');
	pb_total.className = 'pb pb-total';

	let pb_total_progress = document.createElement('span');
	pb_total_progress.className = 'pb-text pb-total-progress';
	pb_total.appendChild(pb_total_progress);

	let pb_total_count = document.createElement('span');
	pb_total_count.className = 'pb-text pb-total-count';
	pb_total.appendChild(pb_total_count);

	filedrop_container.appendChild(pb_total);

	let pb_uploading = false;
	let pb_all_files = [];
	let pb_uploaded_files = [];
	let pb_count = 0;
	let file_id = 0;

	drop_zone.addEventListener('dragover', (e) => {
		e.preventDefault();
	})

	drop_zone.addEventListener('drop', (e) => {
		e.preventDefault();
		if (pb_uploading) {
			return;
		}
		pb_uploading = true;
		const uploadText = document.querySelector('.upload-text');
		uploadText.classList.add('uploading');
		let error_div = document.querySelector('.error-div');
		if (error_div) {
			error_div.style.display = 'none';
			error_div.failed_files = [];
			let le = error_div.querySelector('ul');
			if (le) {
				error_div.removeChild(le);
			}
		}

		let files = e.dataTransfer.files;
		pb_total_count.textContent = pb_count + ' / ' + files.length + " Files Transfered";
		upload_files(files, () => {
			const uploadText = document.querySelector('.upload-text');
			uploadText.classList.remove('uploading');
			pb_uploading = false;
		});
	});

	upload_link.addEventListener('click', (e) => {
		e.preventDefault();
		if (pb_uploading) {
			return
		}
		let error_div = document.querySelector('.error-div');
		if (error_div) {
			error_div.style.display = 'none';
			error_div.failed_files = [];
			let le = error_div.querySelector('ul');
			if (le) {
				error_div.removeChild(le);
			}
		}
		file_input.click();
	});

	file_input.addEventListener('change', () => {
		let files = file_input.files;
		pb_count += 1;
		upload_files(files, () => {
			const uploadText = document.querySelector('.upload-text');
			uploadText.classList.remove('uploading');
			pb_uploading = false;
		});
		pb_total_count.textContent = pb_count + ' / ' + files.length + " Files Transfered";
	});

	function upload_files(files, uploaded_check) {
		let total_loaded = 0;
		pb_count = 0;
		pb_file.style.display = 'flex';
		pb_total.style.display = 'flex';

		let total_size = Array.from(files).reduce((a, b) => a + b.size, 0);

		pb_all_files = Array.from(files).map(a => a.name);
		let pb_start_time = Date.now();

		function nextFile() {
			if (pb_count >= files.length) {
				uploaded_check();
				return;
			}
			let file = files[pb_count];
			let form_data = new FormData();
			form_data.append('file', file);

			let xhr = new XMLHttpRequest();

			let pq = ['borehole_id', 'inventory_id', 'outcrop_id', 'prospect_id', 'shotline_id', 'well_id'].filter(v => drop_zone.dataset[v]).map(v => `${v}=${drop_zone.dataset[v]}`).join('&');
			let upload_url = '../upload' + (pq ? '?' + pq : '');

			xhr.upload.addEventListener('progress', (event) => {
				if (event.lengthComputable) {
					let percent_completed_file = Math.round((event.loaded / event.total) * 100);
					pb_file.style.width = percent_completed_file + '%';
					let elapsed_time_file = (Date.now() - pb_start_time) / 1000;
					pb_file_progress.textContent = pb_size_formatter(Math.round(event.loaded / elapsed_time_file)) + '/s';
					pb_file_name.textContent = file.name;

					let percent_completed_total = Math.round(((total_loaded + event.loaded) / total_size) * 100);
					pb_total.style.width = percent_completed_total + '%';
					let elapsed_time_total = (Date.now() - pb_start_time) / 1000;
					pb_total_progress.textContent = pb_size_formatter(Math.round(event.loaded / elapsed_time_total)) + '/s';
				}
			});

			xhr.addEventListener('load', (event) => {
				if (xhr.status == 200) {
					try {
						data = JSON.parse(xhr.responseText);
					} catch(e) {
						upload_error('Malformed JSON response');
						return;
					}
					file_id = data?.[0]?.file_id;
					if (file_id == undefined) {
						upload_error('No file_id in response');
						return;
					}
					total_loaded += file.size;
					let elapsed_time_total = (Date.now() - pb_start_time) / 1000;
					pb_total_progress.textContent = pb_size_formatter(Math.round(total_size / elapsed_time_total)) + '/s';
					pb_total_count.textContent =
						pb_count + ' / ' + files.length + " Files Transferred";
					add_file_to_page(file);
					pb_uploaded_files.push(file.name);
				} else {
					upload_error(`HTTP error ${xhr.status}`);
				}
			});

			xhr.addEventListener('error', (event) => {
				upload_error('Error reading XML response');
			});

			pb_count++
			xhr.open('POST', upload_url);
			xhr.addEventListener('load', nextFile);
			xhr.send(form_data);
		}
		nextFile()
	}

	function pb_size_formatter(size) {
		var i = size == 0 ? 0 : Math.floor(Math.log(size) / Math.log(1024));
		return (size / Math.pow(1024, i)).toFixed(2) * 1 + ' ' + ['B', 'kB', 'MB', 'GB', 'TB'][i];
	}

	function add_file_to_page(file) {
		const filename_container = document.createElement('div');
		const file_link = document.createElement('a');
		file_link.href = `../file/${file_id}`;
		file_link.textContent = file.name;

		filename_container.appendChild(file_link);

		const file_size = document.createTextNode(` (${pb_size_formatter(file.size)})`);
		filename_container.appendChild(file_size);

		file_list_container.appendChild(filename_container);
	}

	function upload_error(err) {
		const uploadText = document.querySelector('.upload-text');
		uploadText.classList.remove('uploading');
		pb_uploading = false;
		let error_div = document.querySelector('.error-div');
		error_div.style.display = 'flex';
		error_span.textContent += err;

		let failed_files = pb_all_files.filter(file => !pb_uploaded_files.includes(file));
		let le = document.createElement('ul');
		failed_files.forEach(failedFile => {
			let list_item = document.createElement('li');
			list_item.textContent = failedFile;
			le.appendChild(list_item);
		});
		error_div.appendChild(le);

		pb_file.style.display = "none";
		pb_total.style.display = "none";
	}
}
