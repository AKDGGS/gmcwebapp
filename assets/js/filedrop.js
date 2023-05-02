async function uploadFiles(files) {
    let formData = new FormData();
    for (let i = 0; i < files.length; i++) {
        formData.append('files', files[i]);
    }
    try {
        let response = await fetch('upload', {
            method: 'POST',
            body: formData
        });
    } catch (error) {
        console.error(error);
    }
	}

const dropZone = document.querySelector('#file-drop');

dropZone.addEventListener('dragover', (e)=>{
	e.preventDefault();
})

dropZone.addEventListener('drop', (e)=>{
	e.preventDefault();
	var files = e.dataTransfer.files;
	uploadFiles(files)
})
