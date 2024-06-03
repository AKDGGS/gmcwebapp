class SearchControl extends ol.control.Control {
	constructor(opts) {
		const options = (opts || {});
		super({
			element: document.createElement('div'),
			target: options.target
		});

		const root = this.element;
		root.className = 'ol-search ol-control';

		this.container = document.createElement('div');
		this.container.className = 'ol-search-container';
		root.appendChild(this.container);

		if(options.moretools){
			this.menu = document.createElement('div');
			this.menu.className = 'ol-search-menu';
			this.menu.setAttribute('title', 
				options.moretitle ? options.moretitle :
				'Show/Hide Search Tools [Alt + a]'
			);
			this.menu.addEventListener('click', e => {
				this.toggleSearchTools();
			});
			this.container.appendChild(this.menu);
		}

		this.searchbox = document.createElement('input');
		this.searchbox.setAttribute('autocomplete', 'off');
		this.searchbox.setAttribute('type', 'text');
		this.searchbox.setAttribute('placeholder', 'Search for ...');
		this.searchbox.setAttribute('title', 'Text Search Box [Alt + q]');
		this.searchbox.className = 'ol-search-box';
		this.container.appendChild(this.searchbox);

		this.searchsubmit = document.createElement('input');
		this.searchsubmit.setAttribute('type', 'submit');
		this.searchsubmit.setAttribute('value', 'Search');
		this.searchsubmit.setAttribute('title', 'Submit Search');
		this.searchsubmit.className = 'ol-search-submit';
		this.container.appendChild(this.searchsubmit);

		if(options.moretools){
			const flexbreak = document.createElement('div');
			flexbreak.style.flexBasis = '100%';
			flexbreak.style.height = '0';
			this.container.appendChild(flexbreak);

			this.searchtools = document.createElement('div');
			this.searchtools.className = 'ol-search-tools';
			this.searchtools.style.minHeight = (
				options.toolsminheight ? options.toolsminheight : '100px'
			);
			this.searchtools.style.maxHeight = (
				options.toolsmaxheight ? options.toolsmaxheight : '250px'
			);

			this.container.appendChild(this.searchtools);

			window.addEventListener('keydown', e => {
				if(e.altKey && (e.key === 'a' || e.key === 'A')){
					this.toggleSearchTools();
					e.preventDefault();
					return false;
				}
			});
		}
		window.addEventListener('keydown', e => {
			if(e.altKey && (e.key === 'q' || e.key === 'Q')){
				this.getSearchBox().focus();
				e.preventDefault();
				return false;
			}
		});
	}

	getSearchBox() { return this.searchbox; }
	getSearchTools() { return this.searchtools; }
	getSearchButton() { return this.searchsubmit; }
	toggleSearchTools() {
		if(this.searchtools.style.display === 'block'){
			this.searchtools.style.display = 'none';
		} else {
			this.searchtools.style.display = 'block';
		}
	}
}
