class DrawBoxControl extends ol.control.Control {
	constructor(opts) {
		const options = opts || {};
		const olDrawboxControl = document.createElement('div');
		olDrawboxControl.className = 'ol-drawbox-control ol-control';

		super({ element: olDrawboxControl });

		const src = this.source = new ol.source.Vector({ wrapX: false });

		const drawboxButton = document.createElement('button');
		drawboxButton.className = 'drawbox-button';
		drawboxButton.title = "Draw an area of interest";
		olDrawboxControl.appendChild(drawboxButton);

		const olDrawboxTooltipElement = document.createElement('div');
		olDrawboxTooltipElement.className = 'ol-drawbox-tooltip';
		olDrawboxControl.appendChild(olDrawboxTooltipElement);

		const cancelButton = document.createElement('button');
		cancelButton.className = 'cancel-draw-button';
		cancelButton.innerHTML = 'Cancel';
		olDrawboxControl.appendChild(cancelButton);

		const overlay = new ol.Overlay({
			element: olDrawboxTooltipElement,
			offset: [30, 10],
			positioning: 'bottom-left'
		});

		const displayCursorTooltip = (e) => {
			overlay.setPosition(e.coordinate);
			overlay.getElement().innerHTML = 'Draw or drag a box';
		}

		drawboxButton.addEventListener('click', () => {
			olDrawboxControl.classList.add('ol-drawbox-control-active');
			src.clear();
			this.getMap().addInteraction(this.drawRectangle);
			this.getMap().getTargetElement().style.cursor = 'crosshair';
			this.getMap().on('pointermove', displayCursorTooltip);
			this.getMap().addOverlay(overlay);
		});

		drawboxButton.addEventListener('mouseenter', () => {
			olDrawboxTooltipElement.style.display = 'none';
		});

		drawboxButton.addEventListener('mouseleave', () => {
			olDrawboxTooltipElement.style.display = 'block';
		});

		cancelButton.addEventListener('mouseenter', () => {
			olDrawboxTooltipElement.style.display = 'none';
		});

		cancelButton.addEventListener('mouseleave', () => {
			olDrawboxTooltipElement.style.display = 'block';
		});

		cancelButton.addEventListener('click', () => {
			src.clear();
			this.getMap().removeInteraction(this.drawRectangle);
			this.getMap().getTargetElement().style.cursor = 'grab';
			this.getMap().un('pointermove', displayCursorTooltip);
			this.getMap().removeOverlay(overlay);
			olDrawboxControl.classList.remove('ol-drawbox-control-active');

			if (typeof options.callback === 'function') options.callback(null);
		});

		olDrawboxControl.querySelectorAll('button').forEach((button) => {
			button.addEventListener('mouseenter', () => {
				olDrawboxTooltipElement.style.display = 'none';
			});
			button.addEventListener('mouseleave', () => {
				olDrawboxTooltipElement.style.display = 'block';
			});
		});
		
		this.drawLayer = new ol.layer.Vector({
			source: src,
			style: {
				'stroke-color': '#f00',
				'stroke-width': 1,
			},
		});
		this.drawRectangle = new DrawRectangle({ source: src });

		this.drawRectangle.on('boxend', (e) => {
			this.getMap().removeInteraction(e.target);
			this.getMap().getTargetElement().style.cursor = 'grab';
			this.getMap().un('pointermove', displayCursorTooltip);
			this.getMap().removeOverlay(overlay);
			olDrawboxControl.classList.remove('ol-drawbox-control-active');

			if (typeof options.callback === 'function') {
				options.callback(e.target.getFeature());
			}
		});
	}

	getFeature() {
		let feats = this.source.getFeatures();
		if (feats.length === 0) return null;
		return feats[0];
	}

	setMap(m) {
		super.setMap(m);
		this.drawLayer.setMap(m);
		m.on('pointermove', this.drawRectangle.handleMouseMoveEvent);
	}
}

class DrawRectangle extends ol.interaction.Pointer {
	constructor(options) {
		super(options);
		this.source = options.source || new ol.source.Vector();
		this.startPoint = this.endPoint = this.feature = null;
		this.listener = this.isDragging = null;
	}

	handleDownEvent(evt) {
		const map = evt.map;
		if (evt.originalEvent.buttons === 1) {
			if (this.startPoint === null) {
				this.startPoint = map.getCoordinateFromPixel(evt.pixel);
				this.dispatchEvent('boxstart');
				return true;
			} else {
				this.endPoint = map.getCoordinateFromPixel(evt.pixel);
				const coordinates = [
					[this.startPoint, [this.endPoint[0], this.startPoint[1]], this.endPoint, [this.startPoint[0], this.endPoint[1]], this.startPoint]
				];
				const geometry = new ol.geom.Polygon(coordinates);
				this.feature = new ol.Feature(geometry);
				this.source.addFeature(this.feature);
				this.startPoint = this.endPoint = null;
				this.dispatchEvent('boxend');
				this.feature = null;
				return false;
			}
		}
	}

	handleMouseMoveEvent(evt) {
		if (this.startPoint != null) {
			const map = evt.map;
			const currentPoint = map.getCoordinateFromPixel(evt.pixel);
			if (!this.feature) {
				const coordinates = [
					[this.startPoint, [currentPoint[0], this.startPoint[1]], currentPoint, [this.startPoint[0], currentPoint[1]], this.startPoint]
				];
				const geometry = new ol.geom.Polygon(coordinates);
				this.feature = new ol.Feature(geometry);
				this.source.addFeature(this.feature);
				this.listener = (currentPoint) => {
					const coordinates = geometry.getCoordinates()[0];
					coordinates[1][0] = currentPoint[0];
					coordinates[2] = currentPoint;
					coordinates[3][1] = currentPoint[1];
					geometry.setCoordinates([coordinates]);
				};
			} else {
				this.listener(currentPoint);
			}
		}
	}

	handleDragEvent(evt) {
		const map = evt.map;
		const currentPoint = map.getCoordinateFromPixel(evt.pixel);
		if (!this.feature) {
			const coordinates = [
				[this.startPoint, [currentPoint[0], this.startPoint[1]], currentPoint, [this.startPoint[0], currentPoint[1]], this.startPoint]
			];
			const geometry = new ol.geom.Polygon(coordinates);
			this.feature = new ol.Feature(geometry);
			this.source.addFeature(this.feature);
			this.listener = (currentPoint) => {
				const coordinates = geometry.getCoordinates()[0];
				coordinates[1][0] = currentPoint[0];
				coordinates[2] = currentPoint;
				coordinates[3][1] = currentPoint[1];
				geometry.setCoordinates([coordinates]);
			};
		} else {
			this.listener(currentPoint);
		}
		this.isDragging = true;
	}

	handleUpEvent() {
		ol.Observable.unByKey(this.listener);
		if (this.isDragging) {
			this.startPoint = this.isDragging = null;
			this.dispatchEvent('boxend');
			this.feature = null;
		}
		return false;
	}

	getFeature() {
		if (this.feature) return this.feature;
		return null;
	}
}
