const drop_zone = document.getElementById('file-drop');
const file_input = document.getElementById('file-input');

const upload_link = document.getElementById('upload-link');

const progress_bar_file = document.querySelector('.progress-bar-file');
const pb_file = document.getElementById('pb-file');
const pb_file_name = document.getElementById('pb-file-name');

const progress_bar_total = document.querySelector('.progress-bar-total');
const pb_total = document.getElementById('pb-total');
const pb_total_count = document.getElementById('pb-total-count');

let total_size = 0;
let total_loaded = 0;
let count = 0;
let file_names = [];
let uploaded_files = [];
let uploading = false;

if (upload_link) {
	upload_link.addEventListener('click', (e) => {
		e.preventDefault();
		if (uploading) {
			return
		}
		file_input.click();
	});

	file_input.addEventListener('change', () => {
		let files = file_input.files;

		file_names = Array.from(files).map(file => file.name);
		total_size = Array.from(files).reduce((acc, file) => acc + file.size, 0);

		for (let i = 0; i < files.length; ++i) {
			send_file(files[i]);
		}
		count = +1
		pb_total_count.textContent = count + ' / ' + files.length + " Files Transfered";
	})
}

drop_zone.addEventListener('dragover', (e) => {
	e.preventDefault();
})

drop_zone.addEventListener('drop', (e) => {
	e.preventDefault();
	if (uploading) {
		return;
	}
	uploading = true;
	let error_div = document.querySelector('.error-div');
	if (error_div) {
		error_div.style.display = 'none';
		error_div.failed_files = [];
		let le = error_div.querySelector('ul');
		if (le) {
			error_div.removeChild(le);
		}
	}

	// reset the variables each time there is a new drop
	file_names = [];
	uploaded_files = [];
	total_size = 0;
	total_loaded = 0;
	count = 0;

	var files = e.dataTransfer.files;
	file_names = Array.from(files).map(file => file.name);
	total_size = Array.from(files).reduce((acc, file) => acc + file.size, 0);
	pb_total_count.textContent = 0 + ' / ' + files.length + " Files Transfered";
	let i = 0;

	function next() {
		if (i < files.length) {
			console.log("Total_loaded", total_loaded);
			send_file(files[i], (err) => {
				if (err) {
					upload_error();
				} else {
					i++;
					next();
				}
			});
		} else {
			uploading = false;
		}
	}
	next();
});

function send_file(file, callback) {
	progress_bar_file.style.display = 'flex';
	progress_bar_total.style.display = 'flex';

	let form_data = new FormData();
	form_data.append('file', file);
	let xhr = new XMLHttpRequest();
	let drop_zone_data = drop_zone.dataset;
	Object.keys(drop_zone_data).forEach(el => {
		form_data.append(el, drop_zone_data[el]);
	});

	xhr.upload.addEventListener('progress', (event) => {
		if (event.lengthComputable) {
			let percent_completed_file = Math.round((event.loaded / event.total) * 100);
			let percent_completed_total = Math.round(((total_loaded + event.loaded) / total_size) * 100);
			progress_bar_file.style.width = percent_completed_file + '%';
			pb_file.textContent = format_size(Math.round((event.loaded / event.total) * file.size), file.size) + ' / ' + format_size(file.size, file.size);
			pb_file_name.textContent = file.name;
			progress_bar_total.style.width = percent_completed_total + '%';

console.log("Total_loaded progress", (total_loaded + event.loaded), "Event total:", event.total, "File size", file.size);
			pb_total.textContent = format_size(((total_loaded + event.loaded) / event.total) * file.size, total_size) + ' / ' +
				format_size(total_size, total_size);
		}
	});

	xhr.addEventListener('load', (event) => {
		if (xhr.status == 200) {
			total_loaded += file.size;
			console.log("Load total_loaded", total_loaded);
			count += 1;
			pb_total_count.textContent =
				count + ' / ' + file_names.length + " Files Transfered";
			add_file_to_page(file);
			uploaded_files.push(file.name);
			callback(null);
		} else {
			callback(new Error("Upload error"));
		}
	});

	xhr.addEventListener('error', (event) => {
		callback(new Error("Upload Error"));
	});

	xhr.open('POST', 'upload');
	xhr.send(form_data);
}

function format_size(uploaded_amt, file_size) {
	if (file_size < 1024) {
		// console.log('B', uploaded_amt, file_size);
		return uploaded_amt + ' B';
	} else if (file_size < 1048576) {
		return (uploaded_amt / 1024).toFixed(0) + ' KB';
	} else if (file_size < 1073741824) {
		// console.log('MB', uploaded_amt, file_size);
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
	const progress_bar = container.querySelector('.progress-bar-total');
	container.insertBefore(file_div, progress_bar.nextSibling);
}

function upload_error() {
	uploading = false;
	let error_div = document.querySelector('.error-div');
	error_div.style.display = 'flex';

	let failed_files = file_names.filter(file => !uploaded_files.includes(file));
	let le = document.createElement('ul');
	failed_files.forEach(failedFile => {
		let list_item = document.createElement('li');
		list_item.textContent = failedFile;
		le.appendChild(list_item);
	});
	error_div.appendChild(le);

	progress_bar_file.style.display = "none";
	progress_bar_total.style.display = "none";
}
