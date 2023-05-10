const pbFile = document.querySelector('#pb-file');
const pbFileName = document.querySelector('#pb-file-name');
const pbTotalCount = document.querySelector('#pb-total-count');
const dropZone = document.querySelector('#file-drop');
const fileInput = document.querySelector('#file-input');
const fileDropText = document.getElementById('file-drop-text');

dropZone.addEventListener('dragover', (e) => {
	e.preventDefault();
})

dropZone.addEventListener('drop', (e) => {
	e.preventDefault();
	pbFile.textContent = '';
	pbFileName.textContent = '';
	var files = e.dataTransfer.files;
	pbTotalCount.textContent = 0 + ' / ' + files.length + " Files Transfered";
	uploadFiles(files);
	fileDropText.style.display = 'none';
});

dropZone.addEventListener('click', () => {
	fileInput.click();
});

fileInput.addEventListener('change', () => {
	let files = fileInput.files;
	uploadFiles(files);
	fileDropText.style.display = 'none';
	count = +1
	pbTotalCount.textContent = count + ' / ' + files.length + " Files Transfered";
});

async function uploadFiles(files) {
	let totalSize = 0;
	let totalLoaded = 0;
	let count = 0;
	
	const progressBarFile = document.querySelector('#progress-bar-file');
	const progressBarTotal = document.querySelector('#progress-bar-total');
	const pbTotal = document.querySelector('#pb-total');

	progressBarFile.style.display = 'flex';
	progressBarTotal.style.display = 'flex';

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
				progressBarFile.style.width = percentCompletedFile + '%';
				pbFile.textContent = formatSize(Math.round((event.loaded / event.total) * file.size), file.size) + ' / ' + formatSize(file.size, file.size);
				pbFileName.textContent = file.name;
				progressBarTotal.style.width = percentCompletedTotal + '%';
				pbTotal.textContent = formatSize(((totalLoaded + event.loaded) / event.total) * file.size, file.size) + ' / ' +
					formatSize(totalSize, file.size);
			}
		});
		xhr.addEventListener('load', (event) => {
			if (xhr.status >= 200 && xhr.status < 300) {
				totalLoaded += file.size;
				count += 1
				pbTotalCount.textContent =
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

function formatSize(size, fileSize) {
	if (fileSize < 1024) {
		return size + ' B';
	} else if (fileSize < 1024 * 1024) {
		return (size / 1024).toFixed(0) + ' KB';
	} else if (fileSize < 1024 * 1024 * 1024) {
		return (size / (1024 * 1024)).toFixed(2) + ' MB';
	} else {
		return (size / (1024 * 1024 * 1024)).toFixed(2) + ' GB';
	}
}
