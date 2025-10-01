(function () {
	function logError(err) {
		console.error(err);
	}

	const strUndef = 'undefined';
	const protoHttps = 'https:';

	const classNone = 'none';
	const classHeader = 'header';

	const selectorIsNone = '.' + classNone;
	const selectorNotNone = `:not(${selectorIsNone})`;
	const selectorPathList = '.path-list';
	const selectorItemList = '.item-list';
	const selectorItem = `li:not(.${classHeader}):not(.parent)`;
	const selectorItemIsNone = selectorItem + selectorIsNone;
	const selectorItemNotNone = selectorItem + selectorNotNone;

	const Enter = 'Enter';
	const Escape = 'Escape';
	const Esc = 'Esc';
	const Space = ' ';

	let hasStorage = false;
	try {
		if (typeof sessionStorage !== strUndef) hasStorage = true;
	} catch (err) {
	}

	let filteredText = '';

	function matchFilter(input) {
		return input.toLowerCase().includes(filteredText);
	}

	let lastFocused;

	const errLacksMkdir = new Error("mkdir permission is not enabled")

	function enableFilter() {
		const filter = document.body.querySelector('.filter');
		if (!filter) return;

		const input = filter.querySelector('input');
		if (!input) return;

		const clear = filter.querySelector('button') || document.createElement('button');

		const itemList = document.querySelector(selectorItemList)

		let timeoutId;
		const doFilter = function () {
			const filteringText = input.value.trim().toLowerCase();
			if (filteringText === filteredText) return;

			let items
			if (filteringText) {
				clear.style.display = 'block';

				let selector
				if (filteringText.includes(filteredText)) {	// increment search, find in visible items
					selector = selectorItemNotNone;
				} else if (filteredText.includes(filteringText)) {	// decrement search, find in hidden items
					selector = selectorItemIsNone;
				} else {
					selector = selectorItem;
				}
				filteredText = filteringText;

				items = itemList.querySelectorAll(selector);
				items.forEach(item => {
					const name = item.querySelector('.name');
					if (matchFilter(name.textContent)) {
						if (selector !== selectorItemNotNone) {
							item.classList.remove(classNone);
						}
					} else {
						if (selector !== selectorItemIsNone) {
							item.classList.add(classNone);
						}
					}
				});
			} else {	// filter cleared, show all items
				clear.style.display = '';
				filteredText = '';

				items = itemList.querySelectorAll(selectorItemIsNone);
				items.forEach(item => item.classList.remove(classNone));
			}
		};

		const onValueMayChange = function () {
			clearTimeout(timeoutId);
			timeoutId = setTimeout(doFilter, 350);
		};
		input.addEventListener('input', onValueMayChange);
		input.addEventListener('change', onValueMayChange);

		const onEnter = function () {
			clearTimeout(timeoutId);
			input.blur();
			doFilter();
		};
		const onEscape = function () {
			clearTimeout(timeoutId);
			input.value = '';
			doFilter();
		};

		input.addEventListener('keydown', function (e) {
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
		});
		clear.addEventListener('click', function () {
			onEscape();
			input.focus();
		});

		// init
		if (hasStorage) {
			const prevSessionFilter = sessionStorage.getItem(location.pathname);
			if (prevSessionFilter) {
				input.value = prevSessionFilter;
			}
			if (prevSessionFilter !== null) {
				sessionStorage.removeItem(location.pathname);
			}

			window.addEventListener('pagehide', function () {
				if (input.value) {
					sessionStorage.setItem(location.pathname, input.value);
				}
			});
		}
		if (input.value) {
			doFilter();
		}
	}

	function keepFocusOnBackwardForward() {
		function onFocus(e) {
			const link = e.target.closest('a');
			if (!link || link === lastFocused) return;
			lastFocused = link;
		}

		const itemList = document.body.querySelector(selectorItemList);
		itemList.addEventListener('focusin', onFocus);
		itemList.addEventListener('click', onFocus);
		window.addEventListener('pageshow', function () {
			if (lastFocused && lastFocused !== document.activeElement) {
				lastFocused.focus();
				lastFocused.scrollIntoView({block: 'center'});
			}
		});
	}

	function focusChildOnNavUp() {
		function extractCleanUrl(url) {
			let sepIndex = url.indexOf('?');
			if (sepIndex < 0) sepIndex = url.indexOf('#');
			if (sepIndex >= 0) {
				url = url.slice(0, sepIndex);
			}
			return url;
		}

		let prevUrl = document.referrer;
		if (!prevUrl) return;
		prevUrl = extractCleanUrl(prevUrl);

		const currUrl = extractCleanUrl(location.href);

		if (prevUrl.length <= currUrl.length) return;
		if (prevUrl.slice(0, currUrl.length) !== currUrl) return;
		const goesUp = prevUrl.slice(currUrl.length);
		if (currUrl[currUrl.length - 1] !== '/' && goesUp[0] !== '/') return;
		const matchInfo = /[^/]+/.exec(goesUp);
		if (!matchInfo) return;
		let prevChildName = matchInfo[0];
		if (!prevChildName) return;
		prevChildName = decodeURIComponent(prevChildName);
		if (!matchFilter(prevChildName)) return;

		let items = document.body.querySelectorAll(selectorItemList + '>' + selectorItemNotNone);
		items = Array.prototype.slice.call(items);
		const selectorName = '.field.name';
		const selectorLink = 'a';
		for (let i = 0; i < items.length; i++) {
			const item = items[i];
			const elName = item.querySelector(selectorName);
			if (!elName) continue;

			let text = elName.textContent;
			if (text[text.length - 1] === '/') {
				text = text.slice(0, -1);
			}
			if (text !== prevChildName) continue;

			const elLink = item.querySelector(selectorLink);
			if (!elLink) break;

			lastFocused = elLink;
			elLink.focus();
			elLink.scrollIntoView({block: 'center'});
		}
	}

	function enableKeyboardNavigate() {
		const pathList = document.body.querySelector(selectorPathList);
		const itemList = document.body.querySelector(selectorItemList);
		if (!pathList && !itemList) {
			return;
		}

		function getFocusableSibling(container, isBackward, startA) {
			if (!startA) {
				startA = container.querySelector(':focus');
			}
			let startLI = startA && startA.closest('li');
			if (!startLI) {
				if (isBackward) {
					startLI = container.firstElementChild;
				} else {
					startLI = container.lastElementChild;
				}
			}
			if (!startLI) {
				return;
			}

			let siblingLI = startLI;
			do {
				if (isBackward) {
					siblingLI = siblingLI.previousElementSibling || container.lastElementChild;
				} else {
					siblingLI = siblingLI.nextElementSibling || container.firstElementChild;
				}
			} while (siblingLI !== startLI && (
				siblingLI.classList.contains(classNone) ||
				siblingLI.classList.contains(classHeader)
			));

			if (siblingLI) {
				const siblingA = siblingLI.querySelector('a');
				return siblingA;
			}
		}

		function getFirstFocusableSibling(container) {
			const a = container.querySelector(`li:not(.${classNone}):not(.${classHeader}) a`);
			return a;
		}

		function getLastFocusableSibling(container) {
			let a = container.querySelector('li a');
			a = getFocusableSibling(container, true, a);
			return a;
		}

		function getMatchedFocusableSibling(container, isBackward, startA, buf) {
			let skipRound = buf.length === 1;	// find next single-char prefix
			let firstCheckA;
			let secondCheckA;
			let a = startA;
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

				const textContent = (a.querySelector('.name') || a).textContent.toLowerCase();
				if (textContent.startsWith(buf)) {
					return a;
				}
			} while (a = getFocusableSibling(container, isBackward, a));
		}

		const ARROW_UP = 'ArrowUp';
		const ARROW_DOWN = 'ArrowDown';
		const ARROW_LEFT = 'ArrowLeft';
		const ARROW_RIGHT = 'ArrowRight';

		const SKIP_TAGS = ['INPUT', 'BUTTON', 'TEXTAREA'];

		const PLATFORM = navigator.platform || navigator.userAgent;
		const IS_MAC_PLATFORM = PLATFORM.includes('Mac') || PLATFORM.includes('iPhone') || PLATFORM.includes('iPad') || PLATFORM.includes('iPod')

		let lookupKey;
		let lookupBuffer;
		let lookupStartA;
		let lookupTimer;

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

		function lookup(container, key, isBackward) {
			key = key.toLowerCase();

			let currentLookupStartA;
			if (key === lookupKey) {
				// same as last key, lookup next single-char prefix
				currentLookupStartA = container.querySelector(':focus');
			} else {
				if (!lookupStartA) {
					lookupStartA = container.querySelector(':focus');
				}
				currentLookupStartA = lookupStartA;
				if (lookupKey === undefined) {
					lookupKey = key;
				} else {
					// key changed, no more single-char prefix match
					lookupKey = '';
				}
			}
			lookupBuffer += key;
			delayClearLookupContext();
			return getMatchedFocusableSibling(container, isBackward, currentLookupStartA, lookupKey || lookupBuffer);
		}

		let canArrowMove;
		let isToEnd;
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
			if (SKIP_TAGS.includes(e.target.tagName)) {
				return;
			}

			if (canArrowMove(e)) {
				switch (e.key) {
					case ARROW_DOWN:
						if (isToEnd(e)) {
							return getLastFocusableSibling(itemList);
						} else {
							return getFocusableSibling(itemList, false);
						}
					case ARROW_UP:
						if (isToEnd(e)) {
							return getFirstFocusableSibling(itemList);
						} else {
							return getFocusableSibling(itemList, true);
						}
					case ARROW_RIGHT:
						if (isToEnd(e)) {
							return getLastFocusableSibling(pathList);
						} else {
							return getFocusableSibling(pathList, false);
						}
					case ARROW_LEFT:
						if (isToEnd(e)) {
							return getFirstFocusableSibling(pathList);
						} else {
							return getFocusableSibling(pathList, true);
						}
				}
			}
			if (!e.ctrlKey && (!e.altKey || IS_MAC_PLATFORM) && !e.metaKey && e.key.length === 1) {
				return lookup(itemList, e.key, e.shiftKey);
			}
		}

		document.addEventListener('keydown', function (e) {
			const newFocusEl = getFocusItemByKeyPress(e);
			if (newFocusEl) {
				e.preventDefault();
				newFocusEl.focus();
			}
		});
	}

	function enhanceUpload() {
		const upload = document.body.querySelector('.upload');
		if (!upload) return;

		const form = upload.querySelector('form');
		if (!form) return;

		const fileInput = form.querySelector('input[type=file]');
		if (!fileInput) return;

		const submitButton = form.querySelector('[type=submit]');
		if (submitButton) submitButton.classList.add(classNone);

		const uploadType = document.body.querySelector('.upload-type');
		if (!uploadType) return;

		const file = 'file';
		const dirFile = 'dirfile';
		const innerDirFile = 'innerdirfile';

		const optFile = uploadType.querySelector('.' + file);
		const optDirFile = uploadType.querySelector('.' + dirFile);
		const optInnerDirFile = uploadType.querySelector('.' + innerDirFile);
		let optActive = optFile;
		const canMkdir = Boolean(optDirFile);

		function getTimeStamp() {
			const now = new Date();
			let date = String(now.getFullYear() * 10000 + (now.getMonth() + 1) * 100 + now.getDate());
			let time = String(now.getHours() * 10000 + now.getMinutes() * 100 + now.getSeconds());
			let ms = String(now.getMilliseconds());
			date = date.padStart(8, '0');
			time = time.padStart(6, '0');
			ms = ms.padStart(3, '0');
			const ts = '-' + date + '-' + time + '-' + ms;
			return ts;
		}

		let itemsToFiles;
		if (location.protocol === protoHttps && typeof FileSystemHandle !== strUndef && !DataTransferItem.prototype.webkitGetAsEntry) {
			const handleKindFile = 'file';
			const handleKindDir = 'directory';
			const permDescriptor = {mode: 'read'};
			itemsToFiles = function (dataTransferItems, canMkdir) {
				function resultsToFiles(results, files, dirPath) {
					return Promise.all(results.map(function (result) {
						const handle = result.value;
						if (handle.kind === handleKindFile) {
							return handle.queryPermission(permDescriptor).then(function (queryResult) {
								if (queryResult === 'prompt') return handle.requestPermission(permDescriptor);
							}).then(function () {
								return handle.getFile();
							}).then(function (file) {
								const relativePath = dirPath + file.name;
								files.push({file: file, relativePath: relativePath});
							}).catch(function (err) {
								logError(err);
							});
						} else if (handle.kind === handleKindDir) {
							return new Promise(function (resolve) {
								let childResults = [];
								let childIter = handle.values();

								function onLevelDone() {
									childResults = null;
									childIter = null;
									resolve();
								}

								function addChildResult() {
									childIter.next().then(function (result) {
										if (result.done) {
											if (childResults.length) {
												resultsToFiles(childResults, files, dirPath + handle.name + '/').then(onLevelDone);
											} else onLevelDone();
										} else {
											childResults.push(result);
											addChildResult();
										}
									});
								}

								addChildResult();
							});
						}
					}));
				}

				const files = [];
				let hasDir = false;

				return Promise.all(Array.from(dataTransferItems).map(function (item) {
					return item.getAsFileSystemHandle();
				})).then(function (handles) {
					handles = handles.filter(Boolean);	// undefined for pasted content
					hasDir = handles.some(function (handle) {
						return handle.kind === handleKindDir;
					});
					if (hasDir && !canMkdir) {
						throw errLacksMkdir;
					}
					const handleResults = handles.map(function (handle) {
						return {value: handle, done: false};
					});
					return resultsToFiles(handleResults, files, '').then(function () {
						return {files: files, hasDir: hasDir};
					});
				});
			}
		} else {
			itemsToFiles = function (dataTransferItems, canMkdir) {
				function entriesToFiles(entries, files) {
					return Promise.all(entries.map(function (entry) {
						return new Promise(function (resolve) {
							if (entry.isFile) {
								let relativePath = entry.fullPath;
								if (relativePath[0] === '/') {
									relativePath = relativePath.slice(1);
								}
								entry.file(function (file) {
									files.push({file: file, relativePath: relativePath});
									resolve();
								}, function (err) {
									logError(err);
									resolve()
								});
							} else if (entry.isDirectory) {
								const dirReader = entry.createReader()
								const onReadSubEntries = function (subEntries) {
									if (!subEntries.length) resolve();
									entriesToFiles(subEntries, files).then(function () {
										dirReader.readEntries(onReadSubEntries);
									});
								}
								dirReader.readEntries(onReadSubEntries);
							}
						})
					}));
				}

				const entries = [];
				const files = [];
				let hasDir = false;

				for (let i = 0; i < dataTransferItems.length; i++) {
					const item = dataTransferItems[i];
					const entry = item.webkitGetAsEntry();
					if (!entry) {	// undefined for pasted text
						continue;
					}
					if (entry.isFile) {
						// Safari cannot get file from entry by entry.file(), if it is a pasted image
						// so workaround is for all browsers, just get first hierarchy of files by item.getAsFile()
						const file = item.getAsFile();
						files.push({file: file, relativePath: file.name});
					} else if (entry.isDirectory) {
						hasDir = true;
						if (canMkdir) {
							entries.push(entry);
						} else {
							return Promise.reject(errLacksMkdir);
						}
					}
				}

				return entriesToFiles(entries, files).then(function () {
					return {files: files, hasDir: hasDir};
				});
			}
		}

		function enableFileDirModeSwitch() {
			const classActive = 'active';

			function onClickOpt(optTarget, clearInput) {
				if (optTarget === optActive) {
					return false;
				}
				optActive.classList.remove(classActive);

				optActive = optTarget;
				optActive.classList.add(classActive);

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

			if (optFile) {
				optFile.addEventListener('click', onClickOptFile);
				optFile.addEventListener('keydown', onKeydownOpt);

				fileInput.addEventListener('change', function (e) {
					// workaround fix for old browsers, select dir not work but still act like select files
					// switch back to file
					if (optActive === optFile) {
						return;
					}
					const files = e.target.files;
					if (!files.length) {
						return;
					}

					const noDir = Array.prototype.slice.call(files).every(file =>
						file.webkitRelativePath.includes('/')
					);
					if (noDir) {
						onClickOptFile();	// prevent clear input files
					}
				});
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
				const uploadTypeField = 'upload-type';
				const prevUploadType = sessionStorage.getItem(uploadTypeField);
				if (prevUploadType === dirFile) {
					optDirFile && optDirFile.click();
				} else if (prevUploadType === innerDirFile) {
					optInnerDirFile && optInnerDirFile.click();
				}

				if (prevUploadType !== null) {
					sessionStorage.removeItem(uploadTypeField);
				}

				window.addEventListener('pagehide', function () {
					const activeUploadType = fileInput.name;
					if (activeUploadType !== file) {
						sessionStorage.setItem(uploadTypeField, activeUploadType)
					}
				});
			}

			function switchToFileMode() {
				if (optFile && optActive !== optFile) {
					optFile.focus();
					onClickOptFile(true);
				}
			}

			function switchToDirMode() {
				if (optDirFile) {
					if (optActive !== optDirFile) {
						optDirFile.focus();
						onClickOptDirFile();
					}
				} else if (optInnerDirFile) {
					if (optActive !== optInnerDirFile) {
						optInnerDirFile.focus();
						onClickOptInnerDirFile();
					}
				}
			}

			return {
				switchToFileMode: switchToFileMode,
				switchToDirMode: switchToDirMode
			};
		}

		function enableUploadProgress() {	// also fix Safari upload filename has no path info
			let uploading = false;
			const batches = [];
			const classUploading = 'uploading';
			const classFailed = 'failed';
			const elUploadStatus = document.body.querySelector('.upload-status');
			if (!elUploadStatus) return;
			const elProgress = elUploadStatus.querySelector('.progress');
			const elFailedMessage = elUploadStatus.querySelector('.warn .message');

			function onProgress(e) {
				if (e.lengthComputable) {
					const percent = 100 * e.loaded / e.total;
					elProgress.style.width = percent + '%';
				}
			}

			function onFail(e) {
				elUploadStatus.classList.remove(classUploading);
				elUploadStatus.classList.add(classFailed);
				elFailedMessage.textContent = " - " + e.type;
				batches.length = 0;
			}

			function onSuccess(e) {
				const status = e.target.status
				if (status < 200 || status >= 300) {
					return onFail({type: e.target.statusText || status});
				} else if (batches.length) {
					uploadBatch(batches.shift());
				} else {
					uploading = false;
					elUploadStatus.classList.remove(classUploading);
					location.reload();
				}
			}

			function onFinally() {
				elProgress.style.width = '';
			}

			function uploadBatch(files) {
				const formName = fileInput.name;
				const parts = new FormData();
				files.forEach(function (file) {
					let relativePath
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

				const xhr = new XMLHttpRequest();
				xhr.upload.addEventListener('progress', onProgress);
				xhr.addEventListener('error', onFail);
				xhr.addEventListener('error', onFinally);
				xhr.addEventListener('abort', onFail);
				xhr.addEventListener('abort', onFinally);
				xhr.addEventListener('load', onSuccess);
				xhr.addEventListener('load', onFinally);

				xhr.open(form.method, form.action);
				xhr.send(parts);
			}

			function uploadProgressively(files) {
				if (!files.length) {
					return;
				}

				if (uploading) {
					batches.push(files);
				} else {
					uploading = true;
					elUploadStatus.classList.remove(classFailed);
					elUploadStatus.classList.add(classUploading);
					uploadBatch(files);
				}
			}

			return uploadProgressively;
		}

		function enableFormUploadProgress(uploadProgressively) {
			form.addEventListener('submit', function (e) {
				e.stopPropagation();
				e.preventDefault();

				const files = Array.prototype.slice.call(fileInput.files);
				uploadProgressively(files);
			});

			fileInput.addEventListener('change', function () {
				const files = Array.prototype.slice.call(fileInput.files);
				uploadProgressively(files);
			});
		}

		function enableDndUploadProgress(uploadProgressively, switchToFileMode, switchToDirMode) {
			let isSelfDragging = false;
			const classDragging = 'dragging';

			function onSelfDragStart() {
				isSelfDragging = true;
			}

			function onDragEnd() {
				isSelfDragging = false;
			}

			function onDragEnterOver(e) {
				if (isSelfDragging) return;
				e.stopPropagation();
				e.preventDefault();
				e.currentTarget.classList.add(classDragging);
			}

			function onDragLeave(e) {
				if (e.target === e.currentTarget) {
					e.target.classList.remove(classDragging);
				}
			}

			function onDrop(e) {
				e.stopPropagation();
				e.preventDefault();
				e.currentTarget.classList.remove(classDragging);
				fileInput.value = '';

				if (!e.dataTransfer.files.length) {
					return;
				}

				itemsToFiles(e.dataTransfer.items, canMkdir).then(function (result) {
					const files = result.files;
					if (result.hasDir) {
						switchToDirMode();
						uploadProgressively(files);
					} else {
						switchToFileMode();
						uploadProgressively(files);
					}
				}, function (err) {
					if (err === errLacksMkdir && typeof showUploadDirFailMessage !== strUndef) {
						showUploadDirFailMessage();
					}
				});
			}

			document.body.addEventListener('dragstart', onSelfDragStart);
			document.body.addEventListener('dragend', onDragEnd);
			const dndTarget = document.documentElement;
			dndTarget.addEventListener('dragenter', onDragEnterOver);
			dndTarget.addEventListener('dragover', onDragEnterOver);
			dndTarget.addEventListener('dragleave', onDragLeave);
			dndTarget.addEventListener('drop', onDrop);
		}

		function enablePasteUploadProgress(uploadProgressively, switchToFileMode, switchToDirMode) {
			const typeTextPlain = 'text/plain';
			const nonTextInputTypes = ['hidden', 'radio', 'checkbox', 'button', 'reset', 'submit', 'image'];

			function uploadPastedContent(file) {
				switchToFileMode();

				const ts = getTimeStamp();
				let filename = file.name;
				let dotIndex = filename.lastIndexOf('.');
				if (dotIndex < 0) {
					dotIndex = filename.length;
				}
				filename = filename.slice(0, dotIndex) + ts + filename.slice(dotIndex);

				const files = [{file: file, relativePath: filename}];
				uploadProgressively(files);
			}

			document.documentElement.addEventListener('paste', function (e) {
				const tagName = e.target.tagName;
				if (tagName === 'TEXTAREA') {
					return;
				} else if (tagName === 'INPUT' && !nonTextInputTypes.includes(e.target.type)) {
					return;
				}

				const data = e.clipboardData;
				const dItems = data.items;
				const dFiles = data.files;

				if (dItems.length === 1 && dFiles.length === 1 && dFiles[0].type === 'image/png') {
					// image pasted
					uploadPastedContent(dFiles[0]);
				} else if (dItems.length > 0 && dFiles.length === 0) {
					// text pasted (with other data types in DataTransferItems
					const textTypeIndex = data.types.findIndex(function (t) {
						return t === typeTextPlain;
					});
					const textItem = dItems[textTypeIndex];
					if (textItem) {
						textItem.getAsString(function (content) {
							const file = new File([content], 'text.txt', {type: typeTextPlain})
							uploadPastedContent(file);
						});
					}
				} else {
					// actual files/directories pasted
					itemsToFiles(dItems, canMkdir).then(function (result) {
						const files = result.files;

						// pasted real files
						if (result.hasDir) {
							switchToDirMode();
						} else {
							switchToFileMode();
						}
						uploadProgressively(files);
					}, function (err) {
						if (err === errLacksMkdir && typeof showUploadDirFailMessage !== strUndef) {
							showUploadDirFailMessage();
						}
					});
				}
			});
		}

		const {switchToFileMode, switchToDirMode} = enableFileDirModeSwitch();
		const uploadProgressively = enableUploadProgress();
		enableFormUploadProgress(uploadProgressively);
		enableDndUploadProgress(uploadProgressively, switchToFileMode, switchToDirMode);
		enablePasteUploadProgress(uploadProgressively, switchToFileMode, switchToDirMode);
	}

	function enableNonRefreshDelete() {
		const itemList = document.body.querySelector(selectorItemList);
		if (!itemList) return;
		if (!itemList.classList.contains('has-deletable')) return;

		itemList.addEventListener('submit', function (e) {
			if (e.defaultPrevented) return;

			const form = e.target;

			function onLoad() {
				const status = this.status;
				if (status >= 200 && status < 300) {
					const elItem = form.closest('li');
					elItem.remove();
				} else {
					logError('delete failed: ' + status + ' ' + this.statusText);
				}
			}

			let params = '';
			const els = Array.prototype.slice.call(form.elements);
			for (let i = 0; i < els.length; i++) {
				if (!els[i].name) {
					continue
				}
				if (params.length > 0) {
					params += '&'
				}
				params += els[i].name + '=' + encodeURIComponent(els[i].value)
			}
			const url = form.action;

			const xhr = new XMLHttpRequest();
			xhr.open('POST', url);	// will retrieve deleted result into bfcache
			xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
			xhr.addEventListener('load', onLoad);
			xhr.send(params);
			e.preventDefault();
			return false;
		});
	}

	enableFilter();
	keepFocusOnBackwardForward();
	focusChildOnNavUp();
	enableKeyboardNavigate();
	enhanceUpload();
	enableNonRefreshDelete();
}());
