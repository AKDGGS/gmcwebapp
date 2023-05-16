const pbFile = document.querySelector('#pb-file');
const pbFileName = document.querySelector('#pb-file-name');
const pbTotalCount = document.querySelector('#pb-total-count');
const dropZone = document.querySelector('#file-drop');
const fileInput = document.querySelector('#file-input');

const progressBarFile = document.querySelector('.progress-bar-file');
const progressBarTotal = document.querySelector('.progress-bar-total');
const pbTotal = document.querySelector('#pb-total');

let allFiles = [];
let uploadedFiles = [];

dropZone.addEventListener('dragover', (e) => {
	e.preventDefault();
})

dropZone.addEventListener('drop', (e) => {
	e.preventDefault();
	let errorDiv = document.querySelector('.error-div');
	if (errorDiv){
			errorDiv.style.display = 'none';
			errorDiv.failedFiles = [];
			let listElement = errorDiv.querySelector('ul');
			if (listElement) {
				errorDiv.removeChild(listElement);
			}
	}

	allFiles = [];
	uploadedFiles = [];

	var files = e.dataTransfer.files;
	pbTotalCount.textContent = 0 + ' / ' + files.length + " Files Transfered";
	uploadFiles(files);
});

dropZone.addEventListener('click', () => {
	fileInput.click();
});

fileInput.addEventListener('change', () => {
	let files = fileInput.files;
	uploadFiles(files);
	count = +1
	pbTotalCount.textContent = count + ' / ' + files.length + " Files Transfered";
});

async function uploadFiles(files) {
	let totalSize = 0;
	let totalLoaded = 0;
	let count = 0;

	progressBarFile.style.display = 'flex';
	progressBarTotal.style.display = 'flex';

	for (let i = 0; i < files.length; i++) {
		totalSize += files[i].size;
		allFiles.push(files[i].name)
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
				addFileToPage(file);
				uploadedFiles.push(file.name);
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

function addFileToPage(file) {
	const fileDiv = document.createElement('div');
	const dropZone = document.querySelector('#file-drop');
	const fileID = dropZone.dataset.fileId;

	const fileLink = document.createElement('a');
	fileLink.href = `../file/${fileID}/${file.name}`;
	if (file.name.length > 38) {
		fileLink.textContent = file.name.slice(0, 38) + '...';
	} else {
		fileLink.textContent = file.name;
	}
	fileDiv.appendChild(fileLink);

	fSize = formatSize(file.size, file.size);
	const fileSize = document.createTextNode(` (${fSize})`);
	fileDiv.appendChild(fileSize);

	const container = document.querySelector('.container');
	const progressBar = container.querySelector('.progress-bar-total');
	container.insertBefore(fileDiv, progressBar.nextSibling);
}

function errorEvent(){
	let errorDiv = document.querySelector('.error-div');
	errorDiv.style.display = 'flex';

	let failedFiles = allFiles.filter(file => !uploadedFiles.includes(file));
	errorDiv.failedFiles = failedFiles;
	let listElement = document.createElement('ul');
	for (let failedFile of failedFiles) {
		let listItemElement = document.createElement('li');
		listItemElement.textContent = failedFile;
		listElement.appendChild(listItemElement);
	}
	errorDiv.appendChild(listElement);

	progressBarFile.style.display = "none";
	progressBarTotal.style.display = "none";
}
