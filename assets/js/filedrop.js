async function uploadFiles(files) {
    let formData = new FormData();
    for (let i = 0; i < files.length; i++) {
        formData.append('files', files[i]);
    }
    try {
        let xhr = new XMLHttpRequest();
        xhr.upload.addEventListener('progress', (event) => {
            if (event.lengthComputable) {
                let percentCompleted = Math.round((event.loaded / event.total) * 100);
                document.querySelector('#progress-bar').style.width = percentCompleted + '%';
            }
        });
        xhr.addEventListener('load', (event) => {
					if (xhr.status >= 200 && xhr.status < 300) {
							document.querySelector('#progress-bar').textContent = 'Uploads succeeded';
					} else {
							document.querySelector('#progress-bar').textContent = 'An error occurred while uploading the files';
							document.querySelector('#progress-bar').style.backgroundColor = 'Tomato';
					}
        });
        xhr.open('POST', 'upload');
        xhr.send(formData);
    } catch (error) {
        console.log(error)
    }
}

const dropZone = document.querySelector('#file-drop');

dropZone.addEventListener('dragover', (e) => {
    e.preventDefault();
})

dropZone.addEventListener('drop', (e) => {
    e.preventDefault();
    var files = e.dataTransfer.files;
    uploadFiles(files)
});
