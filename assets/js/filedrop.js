const drop_zone = document.getElementById('file-drop');
const file_input = document.getElementById('file-input');

const upload_link = document.getElementById('upload-link');

const pb_file = document.querySelector('.pb-file');
const pb_file_progress = document.getElementById('pb-file-progress');
const pb_file_name = document.getElementById('pb-file-name');

const pb_total = document.querySelector('.pb-total');
const pb_total_progress = document.getElementById('pb-total-progress');
const pb_total_count = document.getElementById('pb-total-count');
let uploaded_files = [];

drop_zone.addEventListener('dragover', (e) => {
	e.preventDefault();
})

drop_zone.addEventListener('drop', (e) => {
	e.preventDefault();
	let error_div = document.querySelector('.error-div');
	if (error_div) {
		error_div.style.display = 'none';
		error_div.failed_files = [];
		let le = error_div.querySelector('ul');
		if (le) {
			error_div.removeChild(le);
		}
	}

	all_files = [];
	uploaded_files = [];

	var files = e.dataTransfer.files;
	pb_total_count.textContent = 0 + ' / ' + files.length + " Files Transfered";
	upload_files(files, () => {
		console.log('All files uploaded');
	});
});

drop_zone.addEventListener('click', () => {
	file_input.click();
});

file_input.addEventListener('change', () => {
	let files = file_input.files;
	count = +1
	pb_total_count.textContent = count + ' / ' + files.length + " Files Transfered";
});

function upload_files(files, uploaded_check) {
	let file_count = 0;
	let total_size = 0;
	let total_loaded = 0;
	let count = 0;

	pb_file.style.display = 'flex';
	pb_total.style.display = 'flex';

	for (let i = 0; i < files.length; i++) {
		total_size += files[i].size;
		all_files.push(files[i].name)
	}

	function nextFile() {
		if (file_count >= files.length) {
			uploaded_check();
			return;
		}
		let file = files[file_count];
		let form_data = new FormData();
		form_data.append('file', file);
		let borehole_id = drop_zone.getAttribute("borehole-id");
		form_data.append("borehole_id", borehole_id);
		let xhr = new XMLHttpRequest();
		['boreholeId', 'wellId'].forEach(n => {
			if (drop_zone.dataset[n]) {
				form_data.append(n, drop_zone.dataset[n]);
			}
		});
		xhr.upload.addEventListener('progress', (event) => {
			if (event.lengthComputable) {
				let percent_completed_file = Math.round((event.loaded / event.total) * 100);
				let percent_completed_total = Math.round(((total_loaded + event.loaded) / total_size) * 100);
				pb_file.style.width = percent_completed_file + '%';
				pb_file_progress.textContent = format_size(Math.round((event.loaded / event.total)
					* file.size), file.size) + ' / ' + format_size(file.size, file.size);
				pb_file_name.textContent = file.name;
				pb_total.style.width = percent_completed_total + '%';
				pb_total_progress.textContent = format_size((total_loaded + event.loaded),
					total_size) + ' / ' + format_size(total_size, total_size);
			}
		});

		xhr.addEventListener('load', (event) => {
			if (xhr.status >= 200 && xhr.status < 300) {
				total_loaded += file.size;
				count += 1
				pb_total_count.textContent =
					count + ' / ' + files.length + " Files Transfered";
				add_file_to_page(file);
				uploaded_files.push(file.name);
			} else {
				upload_error();
			}
		});

		xhr.addEventListener('error', (event) => {
			upload_error();
		});

		file_count++
		xhr.open('POST', 'upload');
		xhr.addEventListener('load', nextFile);
		xhr.send(form_data);
	}
	nextFile()
}

function format_size(uploaded_amt, unit_size) {
	if (unit_size < 1024) {
		return uploaded_amt.toFixed(2) + ' B';
	} else if (unit_size < 1048576) {
		return (uploaded_amt / 1024).toFixed(0) + ' KB';
	} else if (unit_size < 1073741824) {
		return (uploaded_amt / (1048576)).toFixed(2) + ' MB';
	} else {
		return (uploaded_amt / (1073741824)).toFixed(2) + ' GB';
	}
}

function add_file_to_page(file) {
	const file_div = document.createElement('div');
	const drop_zone = document.getElementById('file-drop');
	const file_id = drop_zone.dataset.file_id;
	const file_link = document.createElement('a');
	file_link.href = `../file/${file_id}/${file.name}`;
	if (file.name.length > 38) {
		file_link.textContent = file.name.slice(0, 38) + '...';
	} else {
		file_link.textContent = file.name;
	}
	file_div.appendChild(file_link);

	f_size = format_size(file.size, file.size);
	const file_size = document.createTextNode(` (${f_size})`);
	file_div.appendChild(file_size);

	const container = document.querySelector('.container');
	const progress_bar = container.querySelector('.pb-total');
	container.insertBefore(file_div, progress_bar.nextSibling);
}

function upload_error() {
	uploading = false;
	let error_div = document.querySelector('.error-div');
	error_div.style.display = 'flex';

	let failed_files = all_files.filter(file => !uploaded_files.includes(file));
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
