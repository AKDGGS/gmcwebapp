async function uploadFiles(files) {
	let totalSize = 0;
	let totalLoaded = 0;
	let count = 0;
	document.querySelector('#progress-bar-file').style.display = 'block';
	document.querySelector('#progress-bar-total').style.display = 'block';
	document.querySelector('#progress-bar-file').style.backgroundColor = '#cfe4fa';
	document.querySelector('#progress-bar-file-name').style.backgroundColor = '#cfe4fa';
	document.querySelector('#progress-bar-total').style.backgroundColor = '#cfe4fa';
	document.querySelector('#progress-bar-total-count').style.backgroundColor = '#cfe4fa';

	for (let i = 0; i < files.length; i++) {
		totalSize += files[i].size;
	}

	for (let file of files) {
		let formData = new FormData();
		formData.append('file', file);

		let xhr = new XMLHttpRequest();
		xhr.upload.addEventListener('progress', (event) => {
			if (event.lengthComputable) {
				let percentCompletedFile = Math.round((event.loaded / event.total) * 100);
				let percentCompletedTotal = Math.round(((totalLoaded + event.loaded) / totalSize) * 100);

				document.querySelector('#progress-bar-file').style.width = percentCompletedFile + '%';
				document.querySelector('#pb-file').textContent = formatSize(event.loaded) + ' / ' + formatSize(event.total);

				document.querySelector('#progress-bar-file-name').style.width = percentCompletedFile + '%';
				document.querySelector('#pb-file-name').textContent = file.name;

				document.querySelector('#progress-bar-total').style.width =
					percentCompletedTotal + '%';
				console.log(formatSize(totalSize));

				document.querySelector('#pb-total').textContent =
					formatSize(totalLoaded + event.loaded) + ' / ' + formatSize(totalSize);
				document.querySelector('#progress-bar-total-count').style.width =
					percentCompletedTotal + '%';
			}
		});
		xhr.addEventListener('load', (event) => {
			if (xhr.status >= 200 && xhr.status < 300) {
				totalLoaded += file.size;
				count += 1
				document.querySelector('#pb-total-count').textContent =
					count + ' / ' + files.length + " Files Transfered";
			}
		});
		xhr.open('POST', 'upload');
		await new Promise(resolve => {
			xhr.addEventListener('load', resolve);
			xhr.send(formData);
		});
	}
}

const dropZone = document.querySelector('#file-drop');
const filecount = document.querySelector('#file-count');
const fileInput = document.querySelector('#file-input');

dropZone.addEventListener('dragover', (e) => {
	e.preventDefault();
})

dropZone.addEventListener('drop', (e) => {
	e.preventDefault();
	document.querySelector('#pb-file').textContent = '';
	document.querySelector('#pb-file-name').textContent = '';
	var files = e.dataTransfer.files;
	document.querySelector('#pb-total-count').textContent = 0 + ' / ' + files.length + " Files Transfered";
	uploadFiles(files);
	document.getElementById('file-drop-text').style.display = 'none';
});

dropZone.addEventListener('click', () => {
	fileInput.click();
});

fileInput.addEventListener('change', () => {
	let files = fileInput.files;
	uploadFiles(files);
	document.getElementById('file-drop-text').style.display = 'none';
	count = +1
	document.querySelector('#pb-total-count').textContent =
		count + ' / ' + files.length + " Files Transfered";
});

function formatSize(size) {
	if (size < 1024) {
		return size + ' B';
	} else if (size < 1024 * 1024) {
		return (size / 1024).toFixed(2) + ' KB';
	} else if (size < 1024 * 1024 * 1024) {
		return (size / (1024 * 1024)).toFixed(2) + ' MB';
	} else {
		return (size / (1024 * 1024 * 1024)).toFixed(2) + ' GB';
	}
}
