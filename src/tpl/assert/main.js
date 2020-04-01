(function enableDragUpload() {
	if (!document.querySelector) {
		return;
	}

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
})();

(function enableDelete() {
	if (!canDelete || !document.querySelector) {
		return;
	}

	var onClick = function (e) {
		var target = e ? e.target :
			event ? event.srcElement :
				null;

		if (target && confirm('Delete?')) {
			// extract and assembly delete url
			// display name is not reliable since non-displayable char is escaped
			var pathName = target.parentNode.pathname;
			if (pathName.charAt(pathName.length - 1) === '/') {
				pathName = pathName.substr(0, pathName.length - 1);
			}

			var index = pathName.lastIndexOf('/') + 1;
			var deleteUrl = pathName.substr(0, index) + '?delete&name=' + pathName.substr(index);
			location.href = deleteUrl;
		}

		return false;
	};

	var itemList = document.querySelector('.item-list');
	if (!itemList) {
		return;
	}

	var buttons = itemList.querySelectorAll('.delete');
	var buttonCount = buttons.length;
	for (var i = 0; i < buttonCount; i++) {
		buttons[i].onclick = onClick;
	}

	itemList.className += ' can-delete';
})();
