(function () {
	var strUndef = 'undefined';

	var classNone = 'none';
	var classHeader = 'header';
	var leavingEvent = typeof window.onpagehide !== strUndef ? 'pagehide' : 'beforeunload';

	var Enter = 'Enter';
	var Escape = 'Escape';
	var Esc = 'Esc';
	var Space = ' ';

	var hasClass, addClass, removeClass;
	if (document.body.classList) {
		hasClass = function (el, className) {
			return el && el.classList.contains(className);
		}
		addClass = function (el, className) {
			el && el.classList.add(className);
		}
		removeClass = function (el, className) {
			el && el.classList.remove(className);
		}
	} else {
		hasClass = function (el, className) {
			if (!el) return;
			var reClassName = new RegExp('\\b' + className + '\\b');
			return reClassName.test(el.className);
		}
		addClass = function (el, className) {
			if (!el) return;
			var originalClassName = el.className;
			var reClassName = new RegExp('\\b' + className + '\\b');
			if (!reClassName.test(originalClassName)) {
				el.className = originalClassName + ' ' + className;
			}
		}
		removeClass = function (el, className) {
			if (!el) return;
			var originalClassName = el.className;
			var reClassName = new RegExp('^\\s*' + className + '\\s+|\\s+' + className + '\\b', 'g');
			var newClassName = originalClassName.replace(reClassName, '');
			if (originalClassName !== newClassName) {
				el.className = newClassName;
			}
		}
	}

	var hasStorage = false;
	try {
		if (typeof sessionStorage !== strUndef) hasStorage = true;
	} catch (err) {
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
		if (!filter.addEventListener) {
			filter.className += ' none';
			return;
		}

		var input = filter.querySelector('input');
		if (!input) {
			return;
		}
		var clear = filter.querySelector('button');

		var selectorNone = '.' + classNone;
		var selectorNotNone = ':not(' + selectorNone + ')';
		var selectorItem = '.item-list > li:not(.' + classHeader + '):not(.parent)';
		var selectorItemNone = selectorItem + selectorNone;
		var selectorItemNotNone = selectorItem + selectorNotNone;

		// event handler
		var timeoutId;
		var lastFilterText = '';
		var doFilter = function () {
			var filterText = input.value.trim().toLowerCase();
			if (filterText === lastFilterText) {
				return;
			}

			var selector, items, i;

			if (!filterText) {	// filter cleared, show all items
				if (clear) {
					clear.style.display = '';
				}
				selector = selectorItemNone;
				items = document.body.querySelectorAll(selector);
				for (i = items.length - 1; i >= 0; i--) {
					removeClass(items[i], classNone);
				}
			} else {
				if (clear) {
					clear.style.display = 'block';
				}
				if (filterText.indexOf(lastFilterText) >= 0) {	// increment search, find in visible items
					selector = selectorItemNotNone;
				} else if (lastFilterText.indexOf(filterText) >= 0) {	// decrement search, find in hidden items
					selector = selectorItemNone;
				} else {
					selector = selectorItem;
				}

				items = document.body.querySelectorAll(selector);
				for (i = items.length - 1; i >= 0; i--) {
					var item = items[i];
					var name = item.querySelector('.name');
					if (name && name.textContent.toLowerCase().indexOf(filterText) < 0) {
						addClass(item, classNone);
					} else {
						removeClass(item, classNone);
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

		var onEnter = function () {
			clearTimeout(timeoutId);
			input.blur();
			doFilter();
		};
		var onEscape = function () {
			clearTimeout(timeoutId);
			input.value = '';
			doFilter();
		};

		var ENTER_CODE = 13;
		var ESCAPE_CODE = 27;

		input.addEventListener('keydown', function (e) {
			if (e.key) {
				switch (e.key) {
					case Enter:
						onEnter();
						e.preventDefault();
						break;
					case Escape:
					case Esc:
						onEscape();
						e.preventDefault();
						break;
				}
			} else if (e.keyCode) {
				switch (e.keyCode) {
					case ENTER_CODE:
						onEnter();
						e.preventDefault();
						break;
					case ESCAPE_CODE:
						onEscape();
						e.preventDefault();
						break;
				}
			}
		}, false);

		clear && clear.addEventListener('click', function () {
			clearTimeout(timeoutId);
			input.value = '';
			input.focus();
			doFilter();
		});

		// init
		if (hasStorage) {
			var prevSessionFilter = sessionStorage.getItem(location.pathname);
			sessionStorage.removeItem(location.pathname);

			window.addEventListener(leavingEvent, function () {
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

	function enableKeyboardNavigate() {
		if (
			!document.querySelector ||
			!document.addEventListener ||
			!document.body.parentElement
		) {
			return;
		}

		var pathList = document.body.querySelector('.path-list');
		var itemList = document.body.querySelector('.item-list');
		if (!pathList && !itemList) {
			return;
		}

		function getFocusableSibling(container, isPrev, startA) {
			if (!container) {
				return
			}
			if (!startA) {
				startA = container.querySelector(':focus');
			}
			var startLI = startA;
			while (startLI && startLI.tagName !== 'LI') {
				startLI = startLI.parentElement;
			}
			if (!startLI) {
				if (isPrev) {
					startLI = container.firstElementChild;
				} else {
					startLI = container.lastElementChild;
				}
			}
			if (!startLI) {
				return;
			}

			var siblingLI = startLI;
			do {
				if (isPrev) {
					siblingLI = siblingLI.previousElementSibling;
					if (!siblingLI) {
						siblingLI = container.lastElementChild;
					}
				} else {
					siblingLI = siblingLI.nextElementSibling;
					if (!siblingLI) {
						siblingLI = container.firstElementChild;
					}
				}
			} while (siblingLI !== startLI && (
				hasClass(siblingLI, classNone) ||
				hasClass(siblingLI, classHeader)
			));

			if (siblingLI) {
				var siblingA = siblingLI.querySelector('a');
				return siblingA;
			}
		}

		var selectorFirstAvailLi = 'li:not(.' + classNone + '):not(.' + classHeader + ')';

		function getFirstFocusableSibling(container) {
			var li = container.querySelector(selectorFirstAvailLi);
			var a = li && li.querySelector('a');
			return a;
		}

		function getLastFocusableSibling(container) {
			var a = container.querySelector('li a');
			a = getFocusableSibling(container, true, a);
			return a;
		}

		function getMatchedFocusableSibling(container, isPrev, startA, buf) {
			var skipRound = buf.length === 1;	// find next prefix
			var matchKeyA;
			var firstCheckA;
			var secondCheckA;
			var a = startA;
			do {
				if (skipRound) {
					skipRound = false;
					continue;
				}
				if (!a) {
					continue;
				}

				// firstCheckA maybe a focused a that not belongs to the list
				// secondCheckA must be in the list
				if (!firstCheckA) {
					firstCheckA = a;
				} else if (firstCheckA === a) {
					return;
				} else if (!secondCheckA) {
					secondCheckA = a;
				} else if (secondCheckA === a) {
					return;
				}

				var textContent = (a.querySelector('.name') || a).textContent.toLowerCase();
				if (buf.length <= textContent.length && textContent.substring(0, buf.length) === buf) {
					return a;
				}
			} while (a = getFocusableSibling(container, isPrev, a));
			return matchKeyA;
		}

		var UP = 'Up';
		var DOWN = 'Down';
		var LEFT = 'Left';
		var RIGHT = 'Right';

		var ARROW_UP = 'ArrowUp';
		var ARROW_DOWN = 'ArrowDown';
		var ARROW_LEFT = 'ArrowLeft';
		var ARROW_RIGHT = 'ArrowRight';

		var ARROW_UP_CODE = 38;
		var ARROW_DOWN_CODE = 40;
		var ARROW_LEFT_CODE = 37;
		var ARROW_RIGHT_CODE = 39;

		var SKIP_TAGS = ['INPUT', 'BUTTON', 'TEXTAREA'];

		var PLATFORM = navigator.platform;
		var IS_MAC_PLATFORM = PLATFORM.indexOf('Mac') >= 0 || PLATFORM.indexOf('iPhone') >= 0 || PLATFORM.indexOf('iPad') >= 0 || PLATFORM.indexOf('iPod') >= 0

		var lookupKey;
		var lookupBuffer;
		var lookupStartA;
		var lookupTimer;

		function clearLookupContext() {
			lookupKey = undefined;
			lookupBuffer = '';
			lookupStartA = null;
		}

		clearLookupContext();

		function delayClearLookupContext() {
			clearTimeout(lookupTimer);
			lookupTimer = setTimeout(clearLookupContext, 850);
		}

		function lookup(key) {
			key = key.toLowerCase();

			var currentLookupStartA;
			if (key === lookupKey) {
				// same as last key, lookup next for the same key as prefix
				currentLookupStartA = itemList.querySelector(':focus');
			} else {
				if (!lookupStartA) {
					lookupStartA = itemList.querySelector(':focus');
				}
				currentLookupStartA = lookupStartA;
				if (lookupKey === undefined) {
					lookupKey = key;
				} else {
					// key changed, no more prefix match
					lookupKey = '';
				}
			}
			lookupBuffer += key;
			delayClearLookupContext();
			return getMatchedFocusableSibling(itemList, false, currentLookupStartA, lookupKey || lookupBuffer);
		}

		var canArrowMove;
		var isToEnd;
		if (IS_MAC_PLATFORM) {
			canArrowMove = function (e) {
				return !(e.ctrlKey || e.shiftKey || e.metaKey);	// only allow Opt
			}
			isToEnd = function (e) {
				return e.altKey;	// Opt key
			}
		} else {
			canArrowMove = function (e) {
				return !(e.altKey || e.shiftKey || e.metaKey);	// only allow Ctrl
			}
			isToEnd = function (e) {
				return e.ctrlKey;
			}
		}

		function getFocusItemByKeyPress(e) {
			if (SKIP_TAGS.indexOf(e.target.tagName) >= 0) {
				return;
			}

			if (e.key) {
				if (canArrowMove(e)) {
					switch (e.key) {
						case LEFT:
						case ARROW_LEFT:
							if (isToEnd(e)) {
								return getFirstFocusableSibling(pathList);
							} else {
								return getFocusableSibling(pathList, true);
							}
						case RIGHT:
						case ARROW_RIGHT:
							if (isToEnd(e)) {
								return getLastFocusableSibling(pathList);
							} else {
								return getFocusableSibling(pathList, false);
							}
						case UP:
						case ARROW_UP:
							if (isToEnd(e)) {
								return getFirstFocusableSibling(itemList);
							} else {
								return getFocusableSibling(itemList, true);
							}
						case DOWN:
						case ARROW_DOWN:
							if (isToEnd(e)) {
								return getLastFocusableSibling(itemList);
							} else {
								return getFocusableSibling(itemList, false);
							}
					}
				}
				if (!e.ctrlKey && (!e.altKey || IS_MAC_PLATFORM) && !e.metaKey && e.key.length === 1) {
					return lookup(e.key);
				}
			} else if (e.keyCode) {
				if (canArrowMove(e)) {
					switch (e.keyCode) {
						case ARROW_LEFT_CODE:
							if (isToEnd(e)) {
								return getFirstFocusableSibling(pathList);
							} else {
								return getFocusableSibling(pathList, true);
							}
						case ARROW_RIGHT_CODE:
							if (isToEnd(e)) {
								return getLastFocusableSibling(pathList);
							} else {
								return getFocusableSibling(pathList, false);
							}
						case ARROW_UP_CODE:
							if (isToEnd(e)) {
								return getFirstFocusableSibling(itemList);
							} else {
								return getFocusableSibling(itemList, true);
							}
						case ARROW_DOWN_CODE:
							if (isToEnd(e)) {
								return getLastFocusableSibling(itemList);
							} else {
								return getFocusableSibling(itemList, false);
							}
					}
				}
				if (!e.ctrlKey && (!e.altKey || IS_MAC_PLATFORM) && !e.metaKey && e.keyCode >= 32 && e.keyCode <= 126) {
					return lookup(String.fromCharCode(e.keyCode));
				}
			}
		}

		document.addEventListener('keydown', function (e) {
			var newFocusEl = getFocusItemByKeyPress(e);
			if (newFocusEl) {
				e.preventDefault();
				newFocusEl.focus();
			}
		});
	}

	function enhanceUpload() {
		if (!document.querySelector || !document.addEventListener) {
			return;
		}

		var upload = document.body.querySelector('.upload');
		if (!upload) {
			return;
		}
		var form = upload.querySelector('form');
		if (!form) {
			return;
		}
		var fileInput = form.querySelector('.file');
		if (!fileInput) {
			return;
		}

		var uploadType = document.body.querySelector('.upload-type');
		if (!uploadType) {
			return;
		}

		var file = 'file';
		var dirFile = 'dirfile';
		var innerDirFile = 'innerdirfile';

		var optFile = uploadType.querySelector('.' + file);
		var optDirFile = uploadType.querySelector('.' + dirFile);
		var optInnerDirFile = uploadType.querySelector('.' + innerDirFile);
		var optActive = optFile;

		var padStart = String.prototype.padStart ? function (sourceString, targetLength, padTemplate) {
			return sourceString.padStart(targetLength, padTemplate);
		} : function (sourceString, targetLength, padTemplate) {
			var sourceLength = sourceString.length;
			if (sourceLength >= targetLength) {
				return sourceString;
			}
			var padLength = targetLength - sourceLength
			var repeatCount = Math.ceil(padLength / padTemplate.length);
			var padString;
			if (String.prototype.repeat) {
				padString = padTemplate.repeat(repeatCount);
			} else {
				padString = '';
				for (var i = 0; i < repeatCount; i++) {
					padString += padTemplate;
				}
			}
			if (padString.length > padLength) {
				padString = padString.substring(0, padLength);
			}

			return padString + sourceString;
		}

		function getTimeStamp() {
			var now = new Date();
			var date = String(now.getFullYear() * 10000 + (now.getMonth() + 1) * 100 + now.getDate());
			var time = String(now.getHours() * 10000 + now.getMinutes() * 100 + now.getSeconds());
			var ms = String(now.getMilliseconds());
			date = padStart(date, 8, '0');
			time = padStart(time, 6, '0');
			var ms = padStart(ms, 3, '0');
			var ts = '-' + date + '-' + time + '-' + ms;
			return ts;
		}

		function entriesToFiles(entries, onDone) {
			var files = [];
			var len = entries.length;
			var cb = 0;
			if (!len) {
				onDone(files);
			}

			function increaseCb() {
				cb++;
				if (cb === len) {
					onDone(files);
				}
			}

			entries.forEach(function (entry) {
				if (entry.isFile) {
					var relativePath = entry.fullPath;
					if (relativePath[0] === '/') {
						relativePath = relativePath.substring(1);
					}
					entry.file(function (file) {
						files.push({file: file, relativePath: relativePath});
						increaseCb();
					}, function (err) {
						increaseCb();
						typeof console !== strUndef && console.error(err);
					});
				} else if (entry.isDirectory) {
					var reader = entry.createReader();
					reader.readEntries(function (subEntries) {
						if (subEntries.length) {
							entriesToFiles(subEntries, function (subFiles) {
								Array.prototype.push.apply(files, subFiles);
								increaseCb();
							});
						} else {
							increaseCb();
						}
					});
				}
			});
		}

		function itemsToFiles(dataTransferItems, onDone) {
			var files = [];
			var len = dataTransferItems.length;
			if (!len) {
				onDone(files);
			}

			var entries = [];
			for (var i = 0; i < len; i++) {
				var item = dataTransferItems[i];
				var entry = item.webkitGetAsEntry();
				if (!entry) {
					continue;
				}
				if (entry.isFile) {
					// Safari cannot get file from entry by entry.file(), if it is a pasted image
					// so workaround is for all browsers, just get first hierarchy of files by item.getAsFile()
					var file = item.getAsFile();
					files.push({file: file, relativePath: file.name});
				} else if (entry.isDirectory) {
					entries.push(entry);
				}
			}

			entriesToFiles(entries, function (entryFiles) {
				files.push.apply(files, entryFiles);
				onDone(files);
			});
		}

		function itemsHasDir(dataTransferItems) {
			if (!dataTransferItems) {
				return false;
			}
			var hasDir = false;
			var items = Array.prototype.slice.call(dataTransferItems);
			if (items.length && items[0].webkitGetAsEntry) {
				for (var i = 0, len = items.length; i < len; i++) {
					var entry = items[i].webkitGetAsEntry();
					if (entry && entry.isDirectory) {
						hasDir = true;
						break;
					}
				}
			}
			return hasDir;
		}

		function switchToFileMode() {
			if (optFile && optActive !== optFile) {
				optFile.focus();
				optFile.click();
			}
		}

		function switchToDirMode() {
			if (optDirFile) {
				if (optActive !== optDirFile) {
					optDirFile.focus();
					optDirFile.click();
				}
			} else if (optInnerDirFile) {
				if (optActive !== optInnerDirFile) {
					optInnerDirFile.focus();
					optInnerDirFile.click();
				}
			}
		}

		function switchToAnyDirMode() {
			if (optActive === optFile) {
				if (optDirFile) {
					optDirFile.focus();
					optDirFile.click();
				} else if (optInnerDirFile) {
					optInnerDirFile.focus();
					optInnerDirFile.click();
				}
			}
		}

		function enableAddDirFile() {
			var classHidden = 'hidden';
			var classActive = 'active';

			function onClickOpt(optTarget, clearInput) {
				if (optTarget === optActive) {
					return;
				}
				removeClass(optActive, classActive);

				optActive = optTarget;
				addClass(optActive, classActive);

				if (clearInput) {
					fileInput.value = '';
				}
				return true;
			}

			function onClickOptFile(e) {
				if (onClickOpt(optFile, Boolean(e))) {
					fileInput.name = file;
					fileInput.webkitdirectory = false;
				}
			}

			function onClickOptDirFile() {
				if (onClickOpt(optDirFile, optActive === optFile)) {
					fileInput.name = dirFile;
					fileInput.webkitdirectory = true;
				}
			}

			function onClickOptInnerDirFile() {
				if (onClickOpt(optInnerDirFile, optActive === optFile)) {
					fileInput.name = innerDirFile;
					fileInput.webkitdirectory = true;
				}
			}

			function onKeydownOpt(e) {
				switch (e.key) {
					case Enter:
					case Space:
						if (e.ctrlKey || e.altKey || e.metaKey || e.shiftKey) {
							break;
						}
						e.preventDefault();
						e.stopPropagation();
						if (e.target === optActive) {
							break;
						}
						e.target.click();
						break;
				}
			}

			if (typeof fileInput.webkitdirectory === strUndef) {
				addClass(uploadType, classNone);
				return;
			}
			optDirFile && removeClass(optDirFile, classHidden);
			optInnerDirFile && removeClass(optInnerDirFile, classHidden);

			if (optFile) {
				optFile.addEventListener('click', onClickOptFile);
				optFile.addEventListener('keydown', onKeydownOpt);
			}
			if (optDirFile) {
				optDirFile.addEventListener('click', onClickOptDirFile);
				optDirFile.addEventListener('keydown', onKeydownOpt);
			}
			if (optInnerDirFile) {
				optInnerDirFile.addEventListener('click', onClickOptInnerDirFile);
				optInnerDirFile.addEventListener('keydown', onKeydownOpt);
			}

			if (hasStorage) {
				var uploadTypeField = 'upload-type';
				var prevUploadType = sessionStorage.getItem(uploadTypeField);
				sessionStorage.removeItem(uploadTypeField);

				window.addEventListener(leavingEvent, function () {
					var activeUploadType = fileInput.name;
					if (activeUploadType !== file) {
						sessionStorage.setItem(uploadTypeField, activeUploadType)
					}
				}, false);

				if (prevUploadType === dirFile) {
					optDirFile && optDirFile.click();
				} else if (prevUploadType === innerDirFile) {
					optInnerDirFile && optInnerDirFile.click();
				}
			}

			optFile && fileInput.addEventListener('change', function (e) {
				// workaround fix for mobile device, select dir not work but still act like select files
				// switch back to file
				if (optActive === optFile) {
					return;
				}
				var files = e.target.files;
				if (!files || !files.length) {
					return;
				}

				var nodir = Array.prototype.slice.call(files).every(function (file) {
					return !file.webkitRelativePath;
				});
				if (nodir) {
					onClickOptFile();	// prevent clear input files
				}
			});
		}

		function enableUploadProgress() {	// also fix Safari upload filename has no path info
			if (!FormData) {
				return;
			}

			var uploading = false;
			var batches = [];
			var classUploading = 'uploading';
			var classFailed = 'failed';
			var elUploadStatus = document.body.querySelector('.upload-status');
			var elProgress = elUploadStatus && elUploadStatus.querySelector('.progress');
			var elFailedMessage = elUploadStatus && elUploadStatus.querySelector('.failed .message');

			function onComplete() {
				if (elProgress) {
					elProgress.style.width = '';
				}
			}

			function onSuccess() {
				if (batches.length) {
					return uploadBatch(batches.shift());	// use "return" for tail call optimize
				} else {
					uploading = false;
					removeClass(elUploadStatus, classUploading);
				}
			}

			function onFail(e) {
				removeClass(elUploadStatus, classUploading);
				addClass(elUploadStatus, classFailed);
				if (elFailedMessage) {
					elFailedMessage.textContent = " - " + e.type;
				}
			}

			function onLoad() {
				!uploading && location.reload();
			}

			function onProgress(e) {
				if (e.lengthComputable) {
					var percent = 100 * e.loaded / e.total;
					elProgress.style.width = percent + '%';
				}
			}

			function uploadProgressively(files) {
				if (!files || !files.length) {
					return;
				}

				if (uploading) {
					batches.push(files);
				} else {
					uploading = true;
					addClass(elUploadStatus, classUploading);
					uploadBatch(files);
				}
			}

			function uploadBatch(files) {
				var formName = fileInput.name;
				var parts = new FormData();
				files.forEach(function (file) {
					var relativePath
					if (file.file) {
						// unwrap object {file, relativePath}
						relativePath = file.relativePath;
						file = file.file;
					} else if (file.webkitRelativePath) {
						relativePath = file.webkitRelativePath
					}
					if (!relativePath) {
						relativePath = file.name;
					}

					parts.append(formName, file, relativePath);
				});

				var xhr = new XMLHttpRequest();
				xhr.upload.addEventListener('error', onComplete);
				xhr.upload.addEventListener('error', onFail);
				xhr.upload.addEventListener('abort', onComplete);
				xhr.upload.addEventListener('abort', onFail);
				xhr.upload.addEventListener('load', onComplete);
				xhr.upload.addEventListener('load', onSuccess);
				xhr.upload.addEventListener('load', onLoad);
				if (elProgress) {
					xhr.upload.addEventListener('progress', onProgress);
				}

				xhr.open(form.method, form.action);
				xhr.send(parts);
			}

			return uploadProgressively;
		}

		function enableFormUploadProgress(uploadProgressively) {
			form.addEventListener('submit', function (e) {
				e.stopPropagation();
				e.preventDefault();

				var files = Array.prototype.slice.call(fileInput.files);
				uploadProgressively(files);
			});

			fileInput.addEventListener('change', function () {
				var files = Array.prototype.slice.call(fileInput.files);
				uploadProgressively(files);
			});
		}

		function enableAddDragDrop(uploadProgressively) {
			var isSelfDragging = false;
			var classDragging = 'dragging';

			function onSelfDragStart() {
				isSelfDragging = true;
			}

			function onDragEnd() {
				isSelfDragging = false;
			}

			function onDragEnterOver(e) {
				if (!isSelfDragging) {
					e.stopPropagation();
					e.preventDefault();
					addClass(e.currentTarget, classDragging);
				}
			}

			function onDragLeave(e) {
				if (e.target === e.currentTarget) {
					removeClass(e.currentTarget, classDragging);
				}
			}

			function onDrop(e) {
				e.stopPropagation();
				e.preventDefault();
				removeClass(e.currentTarget, classDragging);
				fileInput.value = '';

				if (!e.dataTransfer || !e.dataTransfer.files || !e.dataTransfer.files.length) {
					return;
				}

				var items = e.dataTransfer.items;
				if (itemsHasDir(items)) {
					if (!uploadProgressively) {
						// must use progressive upload by JS if has directory
						return;
					}
					var itemsCount = items.length;	// save items count earlier, items will be lost after calling FileSystemFileEntry.file()
					itemsToFiles(items, function (files) {
						itemsCount > 1 ? switchToDirMode() : switchToAnyDirMode();
						uploadProgressively(files);
					});
				} else {
					switchToFileMode();

					if (uploadProgressively) {
						var files = Array.prototype.slice.call(e.dataTransfer.files);
						uploadProgressively(files);
					} else {
						fileInput.files = e.dataTransfer.files;
						form.submit();
					}
				}
			}

			document.body.addEventListener('dragstart', onSelfDragStart);
			document.body.addEventListener('dragend', onDragEnd);
			var dragDropEl = document.documentElement;
			dragDropEl.addEventListener('dragenter', onDragEnterOver);
			dragDropEl.addEventListener('dragover', onDragEnterOver);
			dragDropEl.addEventListener('dragleave', onDragLeave);
			dragDropEl.addEventListener('drop', onDrop);
		}

		function enableAddPaste(uploadProgressively) {
			if (!uploadProgressively) {
				document.documentElement.addEventListener('paste', function (e) {
					var data = e.clipboardData;
					if (data && data.files && data.files.length) {
						switchToFileMode();
						fileInput.files = data.files;
						form.submit();
					}
				});
				return;
			}

			var typeTextPlain = 'text/plain';
			var createTextFile;
			var textFilename = 'text.txt';
			if (Blob && Blob.prototype.msClose) {	// legacy Edge
				createTextFile = function (content) {
					var file = new Blob([content], {type: typeTextPlain});
					file.name = textFilename;
					return file;
				};
			} else if (File) {
				createTextFile = function (content) {
					return new File([content], textFilename, {type: typeTextPlain});
				}
			}

			var nonTextInputTypes = ['hidden', 'radio', 'checkbox', 'button', 'reset', 'submit', 'image'];

			function uploadPastedFiles(files) {
				switchToFileMode();
				var ts = getTimeStamp();
				files = files.map(function (f, i) {
					var filename = f.name;
					var dotIndex = filename.lastIndexOf('.');
					if (dotIndex < 0) {
						dotIndex = filename.length;
					}
					filename = filename.substring(0, dotIndex) + ts + '-' + i + filename.substring(dotIndex);
					return {
						file: f,
						relativePath: filename
					}
				});
				uploadProgressively(files);
			}

			function generatePastedFiles(data) {
				var files;
				var items;
				if (data.files && data.files.length) {
					// pasted content is image
					files = Array.prototype.slice.call(data.files);
				} else if (data.items && data.items.length) {
					// pasted content is text
					items = Array.prototype.slice.call(data.items);
					files = items.map(function (item) {
						return item.getAsFile();
					}).filter(Boolean);
				} else {
					files = [];
				}

				if (files.length) {
					uploadPastedFiles(files);
					return;
				}

				if (!createTextFile) {
					return;
				}
				if (!items) {
					return;
				}
				var plainTextFiles = 0;
				for (var i = 0, itemsCount = items.length; i < itemsCount; i++) {
					if (data.types[i] !== typeTextPlain) {
						continue
					}
					plainTextFiles++;
					items[i].getAsString(function (content) {
						var file = createTextFile(content);
						files.push(file);
						if (files.length === plainTextFiles) {
							uploadPastedFiles(files);
						}
					});
				}
			}

			document.documentElement.addEventListener('paste', function (e) {
				var tagName = e.target.tagName;
				if (tagName === 'INPUT') {
					if (nonTextInputTypes.indexOf(e.target.type) < 0) {
						return;
					}
				}
				if (tagName === 'TEXTAREA') {
					return;
				}
				var data = e.clipboardData;
				if (!data) {
					return;
				}

				var items = data.items;
				var itemsCount = items.length;	// save items count earlier, items will be lost after calling FileSystemFileEntry.file()
				if (!items || !itemsCount) {
					generatePastedFiles(data);
					return;
				}
				var hasDir = itemsHasDir(items);
				itemsToFiles(items, function (files) {
					if (!files.length) {
						// for pasted text
						generatePastedFiles(data);
						return;
					}
					if (files.length === 1 && files[0].file.type === 'image/png') {
						// suppose for pasted image
						files = files.map(function (fileInfo) {
							return fileInfo && fileInfo.file;
						});
						generatePastedFiles({files: files});
						return;
					}
					// pasted real files, not image binary
					if (hasDir) {
						itemsCount > 1 ? switchToDirMode() : switchToAnyDirMode();
					} else {
						switchToFileMode();
					}
					uploadProgressively(files);
				});
			});
		}

		enableAddDirFile();
		var uploadProgressively = enableUploadProgress();
		enableFormUploadProgress(uploadProgressively);
		enableAddDragDrop(uploadProgressively);
		enableAddPaste(uploadProgressively);
	}

	function enableNonRefreshDelete() {
		if (!document.querySelector) {
			return;
		}

		var itemList = document.body.querySelector('.item-list');
		if (!itemList || !itemList.addEventListener) {
			return;
		}
		if (!hasClass(itemList, 'has-deletable')) return;

		itemList.addEventListener('submit', function (e) {
			if (e.defaultPrevented) {
				return;
			}

			function onLoad() {
				var elItem = e.target;
				while (elItem && elItem.nodeName !== 'LI') {
					elItem = elItem.parentNode;
				}
				if (!elItem) {
					return;
				}
				var elItemParent = elItem.parentNode;
				elItemParent && elItemParent.removeChild(elItem);
			}

			var params = '';
			var els = [];
			Array.prototype.push.apply(els, e.target.elements)
			for (var i = 0, len = els.length; i < len; i++) {
				if (!els[i].name) {
					continue
				}
				if (params.length > 0) {
					params += '&'
				}
				params += els[i].name + '=' + encodeURIComponent(els[i].value)
			}
			var url = e.target.action + '?' + params

			var xhr = new XMLHttpRequest();
			xhr.open('POST', url);	// will retrieve deleted result into bfcache
			xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
			xhr.addEventListener('load', onLoad);
			xhr.send();
			e.preventDefault();
			return false;
		}, false);
	}

	enableFilter();
	enableKeyboardNavigate();
	enhanceUpload();
	enableNonRefreshDelete();
})();
