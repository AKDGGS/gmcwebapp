document.querySelectorAll('.qa-results').forEach(qr => {
	fetch('qa_count.json?id=' + qr.dataset.id)
		.then(response => {
			if (!response.ok) {
				qr.innerHTML = 'Error';
				return
			}
			return response.json();
		})
		.then(result => { qr.innerHTML = result; })
		.catch(err => {
			qr.innerHTML = 'Error';
		})
});
