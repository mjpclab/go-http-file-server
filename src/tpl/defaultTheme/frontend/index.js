function logError(err) {
	console.error(err);
}

const strUndef = 'undefined';
const strFunction = 'function';
const protoHttps = 'https:';

const classNone = 'none';
const classHeader = 'header';

const selectorIsNone = '.' + classNone;
const selectorNotNone = `:not(${selectorIsNone})`;
const selectorPathList = '.path-list';
const selectorEntryList = '.entry-list';
const selectorEntry = `li:not(.${classHeader}):not(.parent)`;
const selectorEntryIsNone = selectorEntry + selectorIsNone;
const selectorEntryNotNone = selectorEntry + selectorNotNone;

const Enter = 'Enter';
const Escape = 'Escape';

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

const errLacksMkdir = new Error('mkdir permission is not enabled');

function enableFilter() {
	const filter = document.body.querySelector('.filter');
	if (!filter) return;

	const input = filter.querySelector('input');
	if (!input) return;

	const clear = filter.querySelector('button') || document.createElement('button');

	const entryList = document.querySelector(selectorEntryList);

	let timeoutId;
	const doFilter = function () {
		const filteringText = input.value.trim().toLowerCase();
		if (filteringText === filteredText) return;

		let entries;
		if (filteringText) {
			clear.style.display = 'block';

			let selector;
			if (filteringText.includes(filteredText)) {	// increment search, find in visible entries
				selector = selectorEntryNotNone;
			} else if (filteredText.includes(filteringText)) {	// decrement search, find in hidden entries
				selector = selectorEntryIsNone;
			} else {
				selector = selectorEntry;
			}
			filteredText = filteringText;

			entries = entryList.querySelectorAll(selector);
			entries.forEach(entry => {
				const name = entry.querySelector('.name');
				if (matchFilter(name.textContent)) {
					if (selector !== selectorEntryNotNone) {
						entry.classList.remove(classNone);
					}
				} else {
					if (selector !== selectorEntryIsNone) {
						entry.classList.add(classNone);
					}
				}
			});
		} else {	// filter cleared, show all entries
			clear.style.display = '';
			filteredText = '';

			entries = entryList.querySelectorAll(selectorEntryIsNone);
			entries.forEach(entry => entry.classList.remove(classNone));
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
		if (input.value) {
			clearTimeout(timeoutId);
			input.value = '';
			doFilter();
		} else {
			input.blur();
		}
	};
	filter.addEventListener('reset', function (e) {
		e.preventDefault();
	});

	input.addEventListener('keydown', function (e) {
		if (e.key === Enter) {
			onEnter();
			e.preventDefault();
		} else if (e.key === Escape) {
			onEscape();
			e.preventDefault();
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
			const inputValue = input.value;
			if (inputValue) {
				sessionStorage.setItem(location.pathname, inputValue);
			}
		});
	}
	if (input.value) {
		doFilter();
	}
}

function keepFocusOnBackwardForward() {
	const entryList = document.body.querySelector(selectorEntryList);
	entryList.addEventListener('focusin', function (e) {
		if (lastFocused !== e.target) {
			lastFocused = e.target;
		}
	});
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

	if (lastFocused) return;

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

	const entries = Array.from(document.body.querySelectorAll(selectorEntryList + '>' + selectorEntryNotNone));
	const selectorName = '.field.name';
	const selectorLink = 'a';
	for (let i = 0; i < entries.length; i++) {
		const entry = entries[i];
		const elName = entry.querySelector(selectorName);
		if (!elName) continue;

		let text = elName.textContent;
		if (text[text.length - 1] === '/') {
			text = text.slice(0, -1);
		}
		if (text !== prevChildName) continue;

		const elLink = entry.querySelector(selectorLink);
		if (!elLink) break;

		lastFocused = elLink;
		elLink.focus();
		elLink.scrollIntoView({block: 'center'});
		break;
	}
}

function enableKeyboardNavigate() {
	const pathList = document.body.querySelector(selectorPathList);
	const entryList = document.body.querySelector(selectorEntryList);
	if (!pathList && !entryList) {
		return;
	}

	function getFocusableSibling(container, isBackward, startA) {
		if (!container.childElementCount) return;
		if (!startA) {
			startA = container.querySelector(':focus');
		}
		let startLI = startA && startA.closest('li');
		if (!startLI) {
			startLI = isBackward ? container.firstElementChild : container.lastElementChild;
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
		let firstCheckedA;
		let secondCheckedA;
		let a = startA;
		do {
			if (skipRound) {
				skipRound = false;
				continue;
			}
			if (!a) {
				continue;
			}

			// firstCheckedA maybe a focused a that not belongs to the list
			// secondCheckedA must be in the list
			if (!firstCheckedA) {
				firstCheckedA = a;
			} else if (firstCheckedA === a) {
				return;
			} else if (!secondCheckedA) {
				secondCheckedA = a;
			} else if (secondCheckedA === a) {
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
	const IS_MAC_PLATFORM = PLATFORM.includes('Mac') || PLATFORM.includes('iPhone') || PLATFORM.includes('iPad') || PLATFORM.includes('iPod');

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
		};
		isToEnd = function (e) {
			return e.altKey;	// Opt key
		};
	} else {
		canArrowMove = function (e) {
			return !(e.altKey || e.shiftKey || e.metaKey);	// only allow Ctrl
		};
		isToEnd = function (e) {
			return e.ctrlKey;
		};
	}

	function getFocusItemByKeyPress(e) {
		if (SKIP_TAGS.includes(e.target.tagName)) {
			return;
		}

		if (canArrowMove(e)) {
			switch (e.key) {
				case ARROW_DOWN:
					if (isToEnd(e)) {
						return getLastFocusableSibling(entryList);
					} else {
						return getFocusableSibling(entryList, false);
					}
				case ARROW_UP:
					if (isToEnd(e)) {
						return getFirstFocusableSibling(entryList);
					} else {
						return getFocusableSibling(entryList, true);
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
			return lookup(entryList, e.key, e.shiftKey);
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
	const form = document.body.querySelector('.upload-form');
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
	const optDir = uploadType.querySelector('.' + dirFile);
	const optInnerDir = uploadType.querySelector('.' + innerDirFile);
	let optActive = optFile;
	const canMkdir = Boolean(optDir);

	let itemsToFiles;
	if (location.protocol === protoHttps && typeof FileSystemHandle !== strUndef && DataTransferItem.prototype.getAsFileSystemHandle && !DataTransferItem.prototype.webkitGetAsEntry) {
		const handleKindFile = 'file';
		const handleKindDir = 'directory';
		itemsToFiles = async function (dataTransferItems, canMkdir) {
			async function handlesToFiles(handles, files, dirPath) {
				for (let i = 0; i < handles.length; i++) {
					const handle = handles[i];
					if (handle.kind === handleKindFile) {
						const file = await handle.getFile();
						const relativePath = dirPath + file.name;
						files.push({file, relativePath});
					} else if (handle.kind === handleKindDir) {
						const subHandles = [];
						for await(const subHandle of handle.values()) {
							subHandles.push(subHandle);
						}
						if (subHandles.length) {
							await handlesToFiles(subHandles, files, dirPath + handle.name + '/');
						}
					}
				}
			}

			let hasDir = false;
			const permDescriptor = {mode: 'read'};
			const handles = await Promise.all(Array.from(dataTransferItems, async item => {
				const handle = await item.getAsFileSystemHandle();
				if (handle.kind === handleKindDir) {
					if (!canMkdir) throw errLacksMkdir;
					hasDir = true;
				}
				const permState = await handle.queryPermission(permDescriptor);
				if (permState !== 'granted') {
					await handle.requestPermission(permDescriptor);
				}
				return handle;
			}));

			const files = [];
			await handlesToFiles(handles, files, '');
			return {files, hasDir};
		};
	} else {
		itemsToFiles = async function (dataTransferItems, canMkdir) {
			async function entriesToFiles(entries, files) {
				for (let i = 0; i < entries.length; i++) {
					const entry = entries[i];
					if (entry.isFile) {
						let relativePath = entry.fullPath;
						if (relativePath[0] === '/') {
							relativePath = relativePath.slice(1);
						}
						const file = await new Promise((res, rej) => entry.file(res, rej));
						files.push({file, relativePath});
					} else if (entry.isDirectory) {
						const reader = entry.createReader();
						while (true) {
							const subEntries = await new Promise((res, rej) => reader.readEntries(res, rej));
							if (!subEntries.length) break;
							await entriesToFiles(subEntries, files);
						}
					}
				}
			}

			const entries = [];
			const files = [];

			for (let i = 0; i < dataTransferItems.length; i++) {
				const item = dataTransferItems[i];
				const entry = item.webkitGetAsEntry();
				if (entry.isFile) {
					// Safari cannot get file from entry by entry.file(), if it is a pasted image
					// so workaround is for all browsers, just get first hierarchy of files by item.getAsFile()
					const file = item.getAsFile();
					files.push({file, relativePath: file.name});
				} else if (entry.isDirectory) {
					if (!canMkdir) throw errLacksMkdir;
					entries.push(entry);
				}
			}

			const hasDir = entries.length > 0;
			await entriesToFiles(entries, files);
			return {files, hasDir};
		};
	}

	function enableFileDirModeSwitch() {
		const classActive = 'active';
		const title = form.querySelector('h4') || document.createElement('h4');

		function onClickOptAny(optTarget, clearInput) {
			if (optTarget === optActive) {
				return false;
			}

			optActive.classList.remove(classActive);
			optActive = optTarget;
			optActive.classList.add(classActive);
			title.textContent = optActive.title || optActive.textContent;

			if (clearInput) {
				fileInput.value = '';
			}

			return true;
		}

		function onClickOptFile(e) {
			if (onClickOptAny(optFile, Boolean(e))) {
				fileInput.name = file;
				fileInput.webkitdirectory = false;
			}
		}

		function onClickOptDir() {
			if (onClickOptAny(optDir, optActive === optFile)) {
				fileInput.name = dirFile;
				fileInput.webkitdirectory = true;
			}
		}

		function onClickOptInnerDir() {
			if (onClickOptAny(optInnerDir, optActive === optFile)) {
				fileInput.name = innerDirFile;
				fileInput.webkitdirectory = true;
			}
		}

		if (optFile) {
			optFile.addEventListener('click', onClickOptFile);
			fileInput.addEventListener('change', function (e) {
				// workaround fix for old browsers, select dir not work but still act like select files
				// switch back to file
				if (optActive === optFile) return;

				const files = e.target.files;
				if (!files.length) return;

				const hasDir = Array.from(files).some(file =>
					file.webkitRelativePath.includes('/')
				);
				if (!hasDir) {
					onClickOptFile();
				}
			});
		}
		if (optDir) {
			optDir.addEventListener('click', onClickOptDir);
		}
		if (optInnerDir) {
			optInnerDir.addEventListener('click', onClickOptInnerDir);
		}

		if (hasStorage) {
			const uploadTypeField = 'upload-type';
			const prevUploadType = sessionStorage.getItem(uploadTypeField);
			if (prevUploadType === dirFile) {
				optDir && optDir.click();
			} else if (prevUploadType === innerDirFile) {
				optInnerDir && optInnerDir.click();
			} else {
				optFile && optFile.click();
			}

			if (prevUploadType !== null) {
				sessionStorage.removeItem(uploadTypeField);
			}

			window.addEventListener('pagehide', function () {
				const activeUploadType = fileInput.name;
				if (activeUploadType !== file) {
					sessionStorage.setItem(uploadTypeField, activeUploadType);
				}
			});
		} else {
			optFile && optDir && optDir.click();
		}

		function switchToFileMode() {
			if (optFile && optActive !== optFile) {
				optFile.focus();
				onClickOptFile(true);
			}
		}

		function switchToDirMode() {
			if (optDir) {
				if (optActive !== optDir) {
					optDir.focus();
					onClickOptDir();
				}
			} else if (optInnerDir) {
				if (optActive !== optInnerDir) {
					optInnerDir.focus();
					onClickOptInnerDir();
				}
			}
		}

		return {switchToFileMode, switchToDirMode};
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
			const fieldName = fileInput.name;
			const parts = new FormData();
			files.forEach(file => {
				parts.append(fieldName, file.file, file.relativePath);
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
			if (!files.length) return;

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
		});

		fileInput.addEventListener('change', function () {
			const files = Array.from(fileInput.files, file => ({
				file,
				relativePath: file.webkitRelativePath || file.name
			}));
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
			if (!e.dataTransfer.files.length) return;

			itemsToFiles(e.dataTransfer.items, canMkdir).then(function (result) {
				if (result.hasDir) {
					switchToDirMode();
				} else {
					switchToFileMode();
				}
				uploadProgressively(result.files);
			}, function (err) {
				if (err === errLacksMkdir && typeof showUploadDirFailMessage === strFunction) {
					showUploadDirFailMessage();
				} else {
					logError(err);
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

		function uploadPastedFile(file) {
			switchToFileMode();

			const ts = getTimeStamp();
			let filename = file.name;
			let dotIndex = filename.lastIndexOf('.');
			if (dotIndex < 0) {
				dotIndex = filename.length;
			}
			filename = filename.slice(0, dotIndex) + ts + filename.slice(dotIndex);

			const files = [{file, relativePath: filename}];
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

			// image pasted
			if (dItems.length === 1 && dFiles.length === 1 && dFiles[0].type.startsWith('image/')) {
				uploadPastedFile(dFiles[0]);
				return;
			}

			// text pasted (with other data types in DataTransferItems
			if (dItems.length > 0 && dFiles.length === 0) {
				const textTypeIndex = data.types.findIndex(t => t === typeTextPlain);
				if (textTypeIndex < 0) return;

				const textItem = dItems[textTypeIndex];
				textItem.getAsString(function (content) {
					const file = new File([content], 'text.txt', {type: typeTextPlain});
					uploadPastedFile(file);
				});
				return;
			}

			// actual files/directories pasted
			itemsToFiles(dItems, canMkdir).then(function (result) {
				// pasted real files
				if (result.hasDir) {
					switchToDirMode();
				} else {
					switchToFileMode();
				}
				uploadProgressively(result.files);
			}, function (err) {
				if (err === errLacksMkdir && typeof showUploadDirFailMessage === strFunction) {
					showUploadDirFailMessage();
				} else {
					logError(err);
				}
			});
		});
	}

	const {switchToFileMode, switchToDirMode} = enableFileDirModeSwitch();
	const uploadProgressively = enableUploadProgress();
	enableFormUploadProgress(uploadProgressively);
	enableDndUploadProgress(uploadProgressively, switchToFileMode, switchToDirMode);
	enablePasteUploadProgress(uploadProgressively, switchToFileMode, switchToDirMode);
}

function enableNonRefreshDelete() {
	const entryList = document.body.querySelector(selectorEntryList);
	if (!entryList) return;
	if (!entryList.classList.contains('has-deletable')) return;

	entryList.addEventListener('submit', function (e) {
		if (e.defaultPrevented) return;

		const form = e.target;

		const params = Array.from(form.elements, el => {
			const {name, value} = el;
			if (name) {
				return `${name}=${encodeURIComponent(value)}`;
			}
		}).filter(Boolean).join('&');

		// follow redirect to update bf-cache
		fetch(form.action, {
			method: form.method,
			headers: {
				'Content-Type': form.enctype,
			},
			body: params
		}).then(resp => {
			const {status} = resp;
			if (status < 200 || status > 299) throw resp;
			const elItem = form.closest('li');
			elItem.remove();
		}).catch(logError);

		e.preventDefault();
	});
}

enableFilter();
keepFocusOnBackwardForward();
focusChildOnNavUp();
enableKeyboardNavigate();
enhanceUpload();
enableNonRefreshDelete();
