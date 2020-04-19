(function () {
	function enableDragUpload() {
		var upload = document.querySelector('.upload');
		if (!upload || !upload.addEventListener) {
			return;
		}
		var fileInput = upload.querySelector('.files');

		var addClass = function (ele, className) {
			ele && ele.classList && ele.classList.add(className);
		};

		var removeClass = function (ele, className) {
			ele && ele.classList && ele.classList.remove(className);
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

		upload.addEventListener('dragenter', onDragEnterOver);
		upload.addEventListener('dragover', onDragEnterOver);
		upload.addEventListener('dragleave', onDragLeave);
		upload.addEventListener('drop', onDrop);
	}

	function enableNonRefreshDelete() {
		var itemList = document.querySelector('.item-list.has-deletable');
		if (!itemList || !itemList.addEventListener) {
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
		});
	}

	if (!document.querySelector) {
		return;
	}
	enableDragUpload();
	enableNonRefreshDelete();
})();
