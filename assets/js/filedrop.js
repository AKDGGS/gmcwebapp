const drop_zone = document.getElementById('filedrop');
const file_input = document.getElementById('filedrop-file-input');

const upload_link = document.getElementById('filedrop-upload-link');

const pb_file = document.querySelector('.filedrop-pb-file');
const pb_file_progress = document.getElementById('filedrop-pb-file-progress');
const pb_file_name = document.getElementById('filedrop-pb-file-name');

const pb_total = document.querySelector('.filedrop-pb-total');
const pb_total_progress = document.getElementById('filedrop-pb-total-progress');
const pb_total_count = document.getElementById('filedrop-pb-total-count');

let uploading = false;
let all_files = [];
let uploaded_files = [];
let count = 0;
let total_size = 0;

drop_zone.addEventListener('dragover', (e) => {
	e.preventDefault();
})

drop_zone.addEventListener('drop', (e) => {
	e.preventDefault();
	if (uploading) {
		return;
	}
	uploading = true;
	const uploadText = document.getElementById('filedrop-upload-text');
	uploadText.classList.add('uploading');
	let error_div = document.querySelector('.filedrop-error-div');
	if (error_div) {
		error_div.style.display = 'none';
		error_div.failed_files = [];
		let le = error_div.querySelector('ul');
		if (le) {
			error_div.removeChild(le);
		}
	}

	let files = e.dataTransfer.files;
	pb_total_count.textContent = count + ' / ' + files.length + " Files Transfered";
	upload_files(files, () => {
		const uploadText = document.getElementById('filedrop-upload-text');
		uploadText.classList.remove('uploading');
		uploading = false;
	});
});

upload_link.addEventListener('click', (e) => {
	e.preventDefault();
	if (uploading) {
		console.log("uploading true");
		return
	}
	let error_div = document.querySelector('.filedrop-error-div');
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
	count += 1;
	upload_files(files, () => {
		const uploadText = document.getElementById('filedrop-upload-text');
		uploadText.classList.remove('uploading');
		uploading = false;
	});
	pb_total_count.textContent = count + ' / ' + files.length + " Files Transfered";
});

function upload_files(files, uploaded_check) {
	let file_count = 0;
	let total_loaded = 0;
	count = 0;

	pb_file.style.display = 'flex';
	pb_total.style.display = 'flex';

	total_size = Array.from(files).reduce((a,b) => a + b.size, 0);

	all_files = Array.from(files).map(a => a.name);

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
				pb_file_progress.textContent = format_size(percent_completed_file / 100 *
					file.size) + ' / ' + format_size(file.size);
				pb_file_name.textContent = file.name;
				pb_total.style.width = percent_completed_total + '%';
				pb_total_progress.textContent = format_size(percent_completed_total / 100 * total_size)
					+ ' / ' + format_size(total_size);
			}
		});

		xhr.addEventListener('load', (event) => {
			if (xhr.status == 200) {
				total_loaded += file.size;
				count += 1
				pb_total_count.textContent =
					count + ' / ' + files.length + " Files Transferred";
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

function format_size(size) {
	if (size > total_size){
		size = total_size;
	}
	var i = size == 0 ? 0 : Math.floor(Math.log(size) / Math.log(1024));
    return (size / Math.pow(1024, i)).toFixed(2) * 1 + ' ' + ['B', 'kB', 'MB', 'GB', 'TB'][i];
}


function add_file_to_page(file) {
	const file_div = document.createElement('div');
	const drop_zone = document.getElementById('filedrop');
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

	const container = document.querySelector('.filedrop-container');
	const progress_bar = container.querySelector('.filedrop-pb-total');
	container.insertBefore(file_div, progress_bar.nextSibling);
}

function upload_error() {
	const uploadText = document.getElementById('filedrop-upload-text');
	uploadText.classList.remove('uploading');
	uploading = false;
	let error_div = document.querySelector('.filedrop-error-div');
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
