const pb_file = document.querySelector('#pb-file');
const pb_file_name = document.querySelector('#pb-file-name');
const pb_total_count = document.querySelector('#pb-total-count');
const drop_zone = document.querySelector('#file-drop');
const file_input = document.querySelector('#file-input');

const progress_bar_file = document.querySelector('.progress-bar-file');
const progress_bar_total = document.querySelector('.progress-bar-total');
const pb_total = document.querySelector('#pb-total');

let all_files = [];
let uploaded_files = [];

drop_zone.addEventListener('dragover', (e) => {
	e.preventDefault();
})

drop_zone.addEventListener('drop', (e) => {
	e.preventDefault();
	let error_div = document.querySelector('.error-div');
	if (error_div){
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
	upload_files(files);
});

drop_zone.addEventListener('click', () => {
	file_input.click();
});

file_input.addEventListener('change', () => {
	let files = file_input.files;
	upload_files(files);
	count = +1
	pb_total_count.textContent = count + ' / ' + files.length + " Files Transfered";
});

async function upload_files(files) {
	let total_sizes = 0;
	let total_loaded = 0;
	let count = 0;

	progress_bar_file.style.display = 'flex';
	progress_bar_total.style.display = 'flex';

	for (let i = 0; i < files.length; i++) {
		total_sizes += files[i].size;
		all_files.push(files[i].name)
	}

	for (let file of files) {
		let form_data = new FormData();
		form_data.append('file', file);
		let xhr = new XMLHttpRequest();

		xhr.upload.addEventListener('progress', (event) => {
			if (event.lengthComputable) {
				let percent_completed_file = Math.round((event.loaded / event.total) * 100);
				let percent_completed_total = Math.round(((total_loaded + event.loaded) / total_sizes) * 100);
				progress_bar_file.style.width = percent_completed_file + '%';
				pb_file.textContent = format_size(Math.round((event.loaded / event.total) * file.size), file.size) + ' / ' + format_size(file.size, file.size);
				pb_file_name.textContent = file.name;
				progress_bar_total.style.width = percent_completed_total + '%';
				pb_total.textContent = format_size(((total_loaded + event.loaded) / event.total) * file.size, file.size) + ' / ' +
					format_size(total_sizes, file.size);
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
			}
			 else {
				 errorEvent();
			}
		});

		xhr.addEventListener('error', (event) => {
			errorEvent();
		});

		xhr.open('POST', 'upload');
		await new Promise(resolve => {
			xhr.addEventListener('load', resolve);
			xhr.send(form_data);
		});
	}
}

function format_size(uploaded_amt, file_size) {
	if (file_size < 1024) {
		return uploaded_amt + ' B';
	} else if (file_size < 1048576) {
		return (uploaded_amt / 1024).toFixed(0) + ' KB';
	} else if (file_size < 1073741824) {
		return (uploaded_amt / (1048576)).toFixed(2) + ' MB';
	} else {
		return (uploaded_amt / (1073741824)).toFixed(2) + ' GB';
	}
}

function add_file_to_page(file) {
	const file_div = document.createElement('div');
	const drop_zone = document.querySelector('#file-drop');
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

function errorEvent(){
	let error_div = document.querySelector('.error-div');
	error_div.style.display = 'flex';

	let failed_files = all_files.filter(file => !uploaded_files.includes(file));
	error_div.failed_files = failed_files;
	let le = document.createElement('ul');
	for (let failedFile of failed_files) {
		let list_item = document.createElement('li');
		list_item.textContent = failedFile;
		le.appendChild(list_item);
	}
	error_div.appendChild(le);

	progress_bar_file.style.display = "none";
	progress_bar_total.style.display = "none";
}
