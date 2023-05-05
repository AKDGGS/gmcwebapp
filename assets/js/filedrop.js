async function uploadFiles(files) {
	let totalSize = 0;
	let totalLoaded = 0;
	let count = 0;
	document.querySelector('#progress-bar-file').style.display = 'block';
	document.querySelector('#progress-bar-total').style.display = 'block';

	for (let i = 0; i < files.length; i++) {
		totalSize += files[i].size
	}

	for (let i = 0; i < files.length; i++) {
		let file = files[i]
		let formData = new FormData();
		//must be 'files' because web/upload.go looks for files...might want to change this???
		formData.append('files', file);

		let xhr = new XMLHttpRequest();
		xhr.upload.addEventListener('progress', (event) => {
			if (event.lengthComputable) {
				let percentCompletedFile = Math.round((event.loaded / event.total) * 100);
				let percentCompletedTotal = Math.round(((totalLoaded + event.loaded) / totalSize) * 100);
				document.querySelector('#progress-bar-file').style.width = percentCompletedFile + '%';
				document.querySelector('#pb-file').textContent = formatSize(event.loaded) + ' / ' +  formatSize(event.total);
				document.querySelector('#progress-bar-total').style.width = percentCompletedTotal + '%';
				document.querySelector('#pb-total').textContent = formatSize(totalLoaded + event.loaded) + ' / ' +  formatSize(totalSize);
			}
		});
		xhr.addEventListener('load', (event) => {
			if (xhr.status >= 200 && xhr.status < 300) {
				totalLoaded += file.size;
				count += 1
			// 	if (files.length == 1){
			// 	document.querySelector('#progress-bar-total').textContent = count + ' file uploaded';
			// } else {
			// 	document.querySelector('#progress-bar-total').textContent = count + ' files uploaded';
			// }
			}
		});
		xhr.open('POST', 'upload');
		xhr.send(formData);
		document.querySelector('#progress-bar-file').style.backgroundColor = '#194a6b';
		document.querySelector('#progress-bar-total').style.backgroundColor = '#194a6b';
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
	var files = e.dataTransfer.files;
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
	document.querySelector('#progress-bar-total').textContent = count + ' files uploaded';
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
