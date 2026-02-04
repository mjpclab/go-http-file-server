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
const Space = ' ';
const KEY_EVENT_SKIP_TAGS = ['INPUT', 'BUTTON', 'TEXTAREA'];

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

	function getFocusablePageSibling(container, isBackward, steps, startA) {
		if (!container.childElementCount) return;
		if (!startA) {
			startA = container.querySelector(':focus');
		}
		let startLI = startA && startA.closest('li');

		if (!startLI) {
			return getFirstFocusableSibling(container);
		}

		let siblingLI = startLI;
		let sib = siblingLI;
		while (true) {
			sib = isBackward ? sib.previousElementSibling : sib.nextElementSibling;
			if (!sib) break;
			const available = !sib.classList.contains(classNone) &&
				!sib.classList.contains(classHeader);
			if (available) {
				siblingLI = sib;
				--steps;
			}
			if (steps === 0) break;
		}

		const siblingA = siblingLI.querySelector('a');
		return siblingA;
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
	const PAGE_UP = 'PageUp';
	const PAGE_DOWN = 'PageDown';
	const HOME = 'Home';
	const END = 'End';

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
	const itemHeight = entryList.lastElementChild.offsetHeight;
	let itemsPerPage = 1;
	const updateItemsPerPage = () => {
		itemsPerPage = Math.max(Math.floor(visualViewport.height / itemHeight) - 2, 1);
	};
	visualViewport.addEventListener('resize', updateItemsPerPage);
	updateItemsPerPage();

	function getFocusItemByKeyPress(e) {
		if (KEY_EVENT_SKIP_TAGS.includes(e.target.tagName)) return;

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
				case PAGE_DOWN:
					return getFocusablePageSibling(entryList, false, itemsPerPage);
				case PAGE_UP:
					return getFocusablePageSibling(entryList, true, itemsPerPage);
				case END:
					return getLastFocusableSibling(entryList);
				case HOME:
					return getFirstFocusableSibling(entryList);
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
			newFocusEl.scrollIntoView({block: 'center'});
			newFocusEl.focus();
		}
	});
}

function enhanceUpload() {
	const form = document.body.querySelector('.upload form');
	if (!form) return;

	const fileInput = form.querySelector('input[type=file]');
	if (!fileInput) return;

	const submitButton = form.querySelector('button:last-of-type');
	if (submitButton) submitButton.classList.add(classNone);

	const uploadType = document.body.querySelector('.upload-type');
	if (!uploadType) return;

	const file = 'file';
	const dirFile = 'dirfile';
	const innerDirFile = 'innerdirfile';

	const itemKindFile = 'file';

	const optFile = uploadType.querySelector('.' + file);
	const optDir = uploadType.querySelector('.' + dirFile);
	const optInnerDir = uploadType.querySelector('.' + innerDirFile);
	let optActive = optFile;
	const canMkdir = Boolean(optDir);

	let itemsToFiles;
	if (location.protocol === protoHttps && typeof FileSystemHandle !== strUndef && DataTransferItem.prototype.getAsFileSystemHandle && !DataTransferItem.prototype.webkitGetAsEntry) {
		const handleKindFile = 'file';
		const handleKindDir = 'directory';
		itemsToFiles = async function (dataTransferItems) {
			async function handlesToFiles(handles, files, dirPath) {
				await Promise.all(handles.map(async handle => {
					if (handle.kind === handleKindFile) {
						const relativePath = dirPath + handle.name;
						const file = await handle.getFile();
						files.push({file, relativePath});
					} else if (handle.kind === handleKindDir) {
						const subHandles = [];
						for await(const subHandle of handle.values()) {
							subHandles.push(subHandle);
						}
						const relativePath = dirPath + handle.name + '/';
						if (subHandles.length) {
							await handlesToFiles(subHandles, files, relativePath);
						} else {
							const file = new File([''], handle.name);
							files.push({file, relativePath});
						}
					}
				}));
			}

			let hasDir = false;
			const permDescriptor = {mode: 'read'};
			const handles = await Promise.all(Array.from(dataTransferItems).filter(item => item.kind === itemKindFile).map(async item => {
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
		itemsToFiles = async function (dataTransferItems) {
			async function entriesToFiles(entries, files) {
				await Promise.all(entries.map(async (entry) => {
					if (entry.isFile) {
						const file = await new Promise((res, rej) => entry.file(res, rej));
						const relativePath = entry.fullPath.slice(1);
						files.push({file, relativePath});
					} else if (entry.isDirectory) {
						const reader = entry.createReader();
						const subTasks = [];
						for (let i = 0; ; i++) {
							const subEntries = await new Promise((res, rej) => reader.readEntries(res, rej));
							if (subEntries.length > 0) {
								subTasks.push(entriesToFiles(subEntries, files));
								continue;
							}
							if (i === 0) {
								const file = new File([''], entry.name);
								const relativePath = entry.fullPath.slice(1) + '/';
								files.push({file, relativePath});
							} else {
								await Promise.all(subTasks);
							}
							break;
						}
					}
				}));
			}

			let hasDir = false;
			const entries = Array.from(dataTransferItems, item => {
				if (item.kind !== itemKindFile) return;
				const entry = item.webkitGetAsEntry();
				if (entry.isDirectory) {
					if (!canMkdir) throw errLacksMkdir;
					hasDir = true;
				}
				return entry;
			}).filter(Boolean);
			const files = [];

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
			optFile && optFile.click();
		}

		function switchToFileMode() {
			if (optFile && optActive !== optFile) {
				optFile.focus();
				onClickOptFile(true);
			}
		}

		function switchToDirMode() {
			if (optDir && optActive !== optDir) {
				optDir.focus();
				onClickOptDir();
			}
		}

		function switchToInnerDirMode() {
			if (optInnerDir && optActive !== optInnerDir) {
				optInnerDir.focus();
				onClickOptInnerDir();
			}
		}

		return {switchToFileMode, switchToDirMode, switchToInnerDirMode};
	}

	function enableUploadProgress(switchToFileMode, switchToDirMode, switchToInnerDirMode) {
		let uploading = null;
		const classUploading = 'uploading';
		const classFailed = 'failed';
		const elUploadStatus = document.body.querySelector('.upload-status');
		if (!elUploadStatus) return;
		const elProgress = elUploadStatus.querySelector('.progress');
		const elFailedMessage = elUploadStatus.querySelector('.warn .message');

		function uploadBatch(files) {
			const maxCount = 2048;
			const fieldName = fileInput.name;

			const slices = [];
			let totalSize = 0;

			let parts = null;
			let count = Infinity;
			for (let i = 0; i <= files.length; i++) {
				if (count >= maxCount || i === files.length) {
					if (i > 0) slices.push(parts);
					if (i === files.length) break;
					parts = new FormData();
					count = 0;
				}

				const {file, relativePath} = files[i];
				totalSize += file.size;
				count += 1;
				parts.append(fieldName, file, relativePath);
			}
			if (slices.length === 0) return;

			let finishedSize = 0;
			elProgress.style.width = '';
			const onProgress = e => {
				if (e.lengthComputable) {
					const percent = 100 * (finishedSize + e.loaded) / totalSize;
					elProgress.style.width = percent + '%';
				}
			};
			const onUploadSuccess = e => {
				if (e.lengthComputable) {
					finishedSize += e.total;
				}
			};
			const {method, action} = form;
			return new Promise(function (resolve, reject) {
				const onFail = e => {
					slices.length = 0;
					reject(e);
				};
				const onDownloadSuccess = e => {
					const status = e.target.status;
					if (status < 200 || status >= 300) {
						onFail({message: e.target.statusText || status});
						return;
					}
					if (slices.length) {
						uploadSlice(slices.shift());
						return;
					}

					elProgress.style.width = '100%';
					resolve();
				};

				function uploadSlice(parts) {
					const xhr = new XMLHttpRequest();
					xhr.upload.addEventListener('progress', onProgress);
					xhr.upload.addEventListener('error', onFail);
					xhr.addEventListener('error', onFail);
					xhr.upload.addEventListener('abort', onFail);
					xhr.addEventListener('abort', onFail);
					xhr.upload.addEventListener('load', onUploadSuccess);
					xhr.addEventListener('load', onDownloadSuccess);
					xhr.open(method, action);
					xhr.setRequestHeader('accept', 'application/json');
					xhr.send(parts);
				}

				uploadSlice(slices.shift());
			});
		}

		async function tryUploadBatch(getFilesResult) {
			if (!uploading) {
				elUploadStatus.classList.remove(classFailed);
				elUploadStatus.classList.add(classUploading);
				uploading = Promise.resolve();
			}

			const filesResultTask = getFilesResult();	// must extract DataTransferItems ASAP
			const localUploading = uploading = uploading.then(async () => {
				const filesResult = await filesResultTask;
				const {files, hasDir, isInnerDir} = filesResult;

				if (hasDir) {
					isInnerDir ? switchToInnerDirMode() : switchToDirMode();
				} else {
					switchToFileMode();
				}

				await uploadBatch(files);
				if (uploading === localUploading) {
					location.reload();
				}
			}).catch(err => {
				elFailedMessage.textContent = ' - ' + err.message;
				if (err === errLacksMkdir && typeof showUploadDirFailMessage === strFunction) {
					showUploadDirFailMessage();
				} else {
					logError(err);
				}
				elUploadStatus.classList.remove(classUploading);
				elUploadStatus.classList.add(classFailed);
				throw err;
			});

			return localUploading;
		}

		async function uploadFilesProgressively(filesResult) {
			await tryUploadBatch(() => filesResult);
		}

		async function uploadItemsProgressively(dataTransferItems) {
			await tryUploadBatch(() => itemsToFiles(dataTransferItems));
		}

		return {uploadFilesProgressively, uploadItemsProgressively};
	}

	function enableFormUploadProgress(uploadFilesProgressively) {
		form.addEventListener('submit', function (e) {
			e.stopPropagation();
			e.preventDefault();
		});

		fileInput.addEventListener('change', function () {
			let hasDir = false;
			const files = Array.from(fileInput.files, file => {
				const relativePath = file.webkitRelativePath || file.name;
				if (relativePath.includes('/')) hasDir = true;
				return {file, relativePath};
			});
			uploadFilesProgressively({files, hasDir, isInnerDir: fileInput.name === innerDirFile});
		});
	}

	function enableDndUploadProgress(uploadItemsProgressively) {
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
			if (e.dataTransfer.types.includes('Files')) {
				e.stopPropagation();
				e.preventDefault();
				e.currentTarget.classList.add(classDragging);
			}
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
			uploadItemsProgressively(e.dataTransfer.items);
		}

		document.body.addEventListener('dragstart', onSelfDragStart);
		document.body.addEventListener('dragend', onDragEnd);
		const dndTarget = document.documentElement;
		dndTarget.addEventListener('dragenter', onDragEnterOver);
		dndTarget.addEventListener('dragover', onDragEnterOver);
		dndTarget.addEventListener('dragleave', onDragLeave);
		dndTarget.addEventListener('drop', onDrop);
	}

	function enablePasteUploadProgress(uploadFilesProgressively, uploadItemsProgressively) {
		const typeTextPlain = 'text/plain';

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
			const ts = getTimeStamp();
			let filename = file.name;
			let dotIndex = filename.lastIndexOf('.');
			if (dotIndex < 0) {
				dotIndex = filename.length;
			}
			filename = filename.slice(0, dotIndex) + ts + filename.slice(dotIndex);

			const files = [{file, relativePath: filename}];
			uploadFilesProgressively({files, hasDir: false});
		}

		function isTextInput(el) {
			const tagName = el.tagName;
			if (tagName === 'TEXTAREA') {
				return true;
			}
			const nonTextInputTypes = ['hidden', 'radio', 'checkbox', 'button', 'reset', 'submit', 'image'];
			if (tagName === 'INPUT' && !nonTextInputTypes.includes(el.type)) {
				return true;
			}
			return false;
		}

		document.documentElement.addEventListener('paste', function (e) {
			if (isTextInput(e.target)) return;

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
			if (dItems.length) {
				uploadItemsProgressively(dItems);
			}
		});
	}

	const {switchToFileMode, switchToDirMode, switchToInnerDirMode} = enableFileDirModeSwitch();
	const {uploadFilesProgressively, uploadItemsProgressively} = enableUploadProgress(switchToFileMode, switchToDirMode, switchToInnerDirMode);
	enableFormUploadProgress(uploadFilesProgressively);
	enableDndUploadProgress(uploadItemsProgressively);
	enablePasteUploadProgress(uploadFilesProgressively, uploadItemsProgressively);
}

function enableSelectActions() {
	const form = document.body.querySelector('form.entry-form');
	if (!form) return;

	const entryList = form.querySelector(selectorEntryList);
	if (!entryList) return;

	const pointerDownEvent = 'mousedown';
	const pointerMoveEvent = 'mousemove';
	const pointerUpEvent = 'mouseup';

	const btnDelete = form.querySelector('.action-list .delete');
	const btnToggleSelect = entryList.querySelector('.toggle-select');
	const chkSelectAll = entryList.querySelector('.select-all');

	const classSelecting = 'selecting';
	const selectorItem = 'li:not(.header)';
	const selectorVisible = `${selectorItem}${selectorNotNone}`;
	const selectorHidden = `${selectorItem}${selectorIsNone}`;
	const selectorSelectLabel = '.select';
	const selectorCheckInput = 'input[type=checkbox]';
	const selectorCheckbox = `${selectorSelectLabel} ${selectorCheckInput}`;
	const selectorUnchecked = `${selectorCheckbox}:not(:checked)`;
	const selectorChecked = `${selectorCheckbox}:checked`;
	const selectorVisibleUnchecked = `${selectorVisible} ${selectorUnchecked}`;
	const selectorHiddenChecked = `${selectorHidden} ${selectorChecked}`;

	const selectRect = document.createElement('div');
	selectRect.classList.add('select-rect');
	form.append(selectRect);

	let maxX, maxY;
	const updateMaxSize = () => {
		maxX = Math.max(document.documentElement.offsetWidth, Math.floor(visualViewport.width));
		maxY = Math.max(document.documentElement.offsetHeight, Math.floor(visualViewport.height));
	};
	visualViewport.addEventListener('resize', updateMaxSize);
	updateMaxSize();

	const classPinching = 'pinching';
	let startItem = null;
	let startX, startY;
	const getPointerPosition = e => {
		let x = e.offsetX;
		let y = e.offsetY;
		let el = e.target;
		const offsetContainer = entryList.offsetParent;
		do {
			x += el.offsetLeft;
			y += el.offsetTop;
			el = el.offsetParent;
		} while (el && el !== offsetContainer);
		return [Math.min(x, maxX), Math.min(y, maxY)];
	};
	const cleanUp = () => {
		startItem = null;
		document.documentElement.removeEventListener(pointerMoveEvent, onPointerMove);
		document.documentElement.removeEventListener(pointerUpEvent, onPointerUp);
		selectRect.classList.remove(classPinching);
	};
	const onPointerDown = e => {
		if (startItem) cleanUp();
		if (e.button !== 0) return;
		const selectLabel = e.target.closest(selectorSelectLabel);
		if (!selectLabel) return;
		startItem = selectLabel.closest(selectorItem);
		e.preventDefault();	// avoid dragging selected text

		document.documentElement.addEventListener(pointerMoveEvent, onPointerMove);
		document.documentElement.addEventListener(pointerUpEvent, onPointerUp);

		([startX, startY] = getPointerPosition(e));
		selectRect.style.left = startX + 'px';
		selectRect.style.top = startY + 'px';
		selectRect.style.width = '';
		selectRect.style.height = '';
		selectRect.classList.add(classPinching);
	};
	const onPointerMove = e => {
		const [endX, endY] = getPointerPosition(e);
		selectRect.style.left = Math.min(startX, endX) + 'px';
		selectRect.style.top = Math.min(startY, endY) + 'px';
		selectRect.style.width = Math.abs(endX - startX) + 'px';
		selectRect.style.height = Math.abs(endY - startY) + 'px';
	};
	const onPointerUp = e => {
		if (!startItem) return;	// e.g. pointer-downed & press ESC
		let fromItem = startItem;
		cleanUp();

		const [, endY] = getPointerPosition(e);
		let currentItem = e.target.closest(selectorItem);
		if (!currentItem || currentItem === fromItem) return;
		const checked = !fromItem.querySelector(selectorCheckInput).checked;
		let toItem;
		if (endY > startY) {
			toItem = currentItem;
		} else {
			toItem = fromItem;
			fromItem = currentItem;
		}
		let item = fromItem;
		while (true) {
			if (!item.classList.contains(classNone)) {
				item.querySelector(selectorCheckInput).checked = checked;
			}
			if (item === toItem) break;
			item = item.nextElementSibling;
		}
	};
	const onKeyDown = e => {
		if (e.key !== Enter && e.key !== Space) return;
		const checkbox = e.target.parentElement.querySelector(selectorCheckbox);
		if (checkbox) {
			e.preventDefault();
			checkbox.click();
		}
	};

	form.addEventListener('submit', function () {
		if (btnDelete) {
			btnDelete.disabled = true;
		}
		entryList.querySelectorAll(selectorHiddenChecked).forEach(input => input.checked = false);

		setTimeout(() => {
			form.classList.remove(classSelecting);
			entryList.querySelectorAll(selectorChecked).forEach(input => input.checked = false);
		}, 0);
	});

	if (btnToggleSelect) {
		const onToggleSelect = () => {
			form.classList.toggle(classSelecting);
			const selecting = form.classList.contains(classSelecting);
			if (btnDelete) {
				btnDelete.disabled = !selecting;
			}
			if (selecting) {
				document.documentElement.addEventListener(pointerDownEvent, onPointerDown);
				entryList.addEventListener('keydown', onKeyDown);
			} else {
				document.documentElement.removeEventListener(pointerDownEvent, onPointerDown);
				entryList.removeEventListener('keydown', onKeyDown);
				if (startItem) cleanUp();
				entryList.querySelectorAll(selectorChecked).forEach(input => input.checked = false);
			}
		};
		btnToggleSelect.addEventListener('click', onToggleSelect);
		document.body.addEventListener('keydown', function (e) {
			if (e.key !== Escape) return;
			if (KEY_EVENT_SKIP_TAGS.includes(e.target.tagName)) return;
			if (e.target === e.currentTarget) {
				onToggleSelect();
			} else if (entryList.contains(e.target)) {
				onToggleSelect();
			}
		});
	}

	if (chkSelectAll) {
		chkSelectAll.addEventListener('change', function (e) {
			const checked = e.target.checked;
			if (checked) {
				entryList.querySelectorAll(selectorVisibleUnchecked).forEach(input => input.checked = true);
			} else {
				entryList.querySelectorAll(selectorChecked).forEach(input => input.checked = false);
			}
		});
	}

	if (typeof confirmDelete === strFunction) {
		if (btnDelete) {
			btnDelete.addEventListener('click', function (e) {
				if (!confirmDelete()) e.preventDefault();
			});
		}
	}
}

enableFilter();
keepFocusOnBackwardForward();
focusChildOnNavUp();
enableKeyboardNavigate();
enhanceUpload();
enableSelectActions();
