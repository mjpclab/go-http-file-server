(function () {
	function enableDragUpload() {
		if (!document.querySelector || !document.addEventListener || !document.body.classList) {
			return;
		}

		var upload = document.body.querySelector('.upload');
		if (!upload) {
			return;
		}
		var fileInput = upload.querySelector('.file');

		var addClass = function (ele, className) {
			ele && ele.classList.add(className);
		};

		var removeClass = function (ele, className) {
			ele && ele.classList.remove(className);
		};

		var onDragEnterOver = function (e) {
			e.stopPropagation();
			e.preventDefault();
			addClass(e.currentTarget, 'dragging');
		};

		var onDragLeave = function (e) {
			if (e.target === e.currentTarget) {
				removeClass(e.currentTarget, 'dragging');
			}
		};

		var onDrop = function (e) {
			e.stopPropagation();
			e.preventDefault();
			removeClass(e.currentTarget, 'dragging');

			if (!e.dataTransfer.files) {
				return;
			}
			fileInput.files = e.dataTransfer.files;
		};

		upload.addEventListener('dragenter', onDragEnterOver, false);
		upload.addEventListener('dragover', onDragEnterOver, false);
		upload.addEventListener('dragleave', onDragLeave, false);
		upload.addEventListener('drop', onDrop, false);
	}

	function enableFilter() {
		if (!document.querySelector) {
			var filter = document.getElementById && document.getElementById('panel-filter');
			if (filter) {
				filter.className += ' none';
			}
			return;
		}

		// pre check
		var filter = document.body.querySelector('.filter');
		if (!filter) {
			return;
		}
		if (!filter.classList || !filter.addEventListener) {
			filter.className += ' none';
			return;
		}

		var input = filter.querySelector('input.filter-text');
		if (!input) {
			return;
		}

		// event handler
		var timeoutId;
		var lastFilterText = '';
		var doFilter = function () {
			var filterText = input.value.trim().toLowerCase();
			if (filterText === lastFilterText) {
				return;
			}

			var itemsSelector = '.item-list > li:not(.header):not(.parent)';
			var items, i;

			if (!filterText) {	// filter cleared, show all items
				itemsSelector += '.none';
				items = document.body.querySelectorAll(itemsSelector);
				for (i = items.length - 1; i >= 0; i--) {
					items[i].classList.remove('none');
				}
			} else {
				if (filterText.indexOf(lastFilterText) >= 0) {	// increment search, find in visible items
					itemsSelector += ':not(.none)';
				} else if (lastFilterText.indexOf(filterText) >= 0) {	// decrement search, find in hidden items
					itemsSelector += '.none';
				}

				items = document.body.querySelectorAll(itemsSelector);
				for (i = items.length - 1; i >= 0; i--) {
					var item = items[i];
					var name = item.querySelector('.name');
					if (name && name.textContent.toLowerCase().indexOf(filterText) < 0) {
						item.classList.add('none');
					} else {
						item.classList.remove('none');
					}
				}
			}

			lastFilterText = filterText;
		};

		var onValueMayChange = function () {
			clearTimeout(timeoutId);
			timeoutId = setTimeout(doFilter, 350);
		};
		input.addEventListener('input', onValueMayChange, false);
		input.addEventListener('change', onValueMayChange, false);
		input.addEventListener('keydown', function (e) {
			switch (e.key) {
				case 'Enter':
					clearTimeout(timeoutId);
					input.blur();
					doFilter();
					e.preventDefault();
					break;
				case 'Escape':
				case 'Esc':
					clearTimeout(timeoutId);
					input.value = '';
					doFilter();
					e.preventDefault();
					break;
			}
		}, false);

		// init
		if (sessionStorage) {
			var prevSessionFilter = sessionStorage.getItem(location.pathname);
			sessionStorage.removeItem(location.pathname);
			window.addEventListener('beforeunload', function () {
				if (input.value) {
					sessionStorage.setItem(location.pathname, input.value);
				}
			}, false);
			if (prevSessionFilter) {
				input.value = prevSessionFilter;
			}
		}
		if (input.value) {
			doFilter();
		}
	}

	function enableNonRefreshDelete() {
		if (!document.querySelector) {
			return;
		}

		var itemList = document.body.querySelector('.item-list');
		if (!itemList || !itemList.addEventListener) {
			return;
		}
		if (itemList.classList) {
			if (!itemList.classList.contains('has-deletable')) {
				return;
			}
		} else if (itemList.className.indexOf('has-deletable') < 0) {
			return;
		}

		itemList.addEventListener('click', function (e) {
			if (e.defaultPrevented || !e.target || e.target.className.indexOf('delete') < 0) {
				return;
			}

			var target = e.target;

			var xhr = new XMLHttpRequest();
			xhr.open('POST', target.href);
			xhr.onload = function () {
				var item = target;
				var parentNode = item.parentNode;
				while (item.nodeName !== 'LI') {
					if (!parentNode) {
						break;
					}
					item = parentNode;
					parentNode = item.parentNode;
				}
				if (parentNode) {
					parentNode.removeChild(item);
				}
				item = null;
				parentNode = null;
				target = null;
			};
			xhr.onerror = xhr.onabort = function () {
				target = null;
			};
			xhr.send();
			e.preventDefault();
			return false;
		}, false);
	}

	enableDragUpload();
	enableFilter();
	enableNonRefreshDelete();
})();
