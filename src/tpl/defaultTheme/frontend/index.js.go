package frontend

const DefaultJs = "" +
	"function logError(err) {\n" +
	"	console.error(err);\n" +
	"}\n" +
	"\n" +
	"const strUndef = 'undefined';\n" +
	"const strFunction = 'function';\n" +
	"const protoHttps = 'https:';\n" +
	"\n" +
	"const classNone = 'none';\n" +
	"const classHeader = 'header';\n" +
	"\n" +
	"const selectorIsNone = '.' + classNone;\n" +
	"const selectorNotNone = `:not(${selectorIsNone})`;\n" +
	"const selectorPathList = '.path-list';\n" +
	"const selectorEntryList = '.entry-list';\n" +
	"const selectorEntry = `li:not(.${classHeader}):not(.parent)`;\n" +
	"const selectorEntryIsNone = selectorEntry + selectorIsNone;\n" +
	"const selectorEntryNotNone = selectorEntry + selectorNotNone;\n" +
	"\n" +
	"const Enter = 'Enter';\n" +
	"const Escape = 'Escape';\n" +
	"\n" +
	"let hasStorage = false;\n" +
	"try {\n" +
	"	if (typeof sessionStorage !== strUndef) hasStorage = true;\n" +
	"} catch (err) {\n" +
	"}\n" +
	"\n" +
	"let filteredText = '';\n" +
	"\n" +
	"function matchFilter(input) {\n" +
	"	return input.toLowerCase().includes(filteredText);\n" +
	"}\n" +
	"\n" +
	"let lastFocused;\n" +
	"\n" +
	"const errLacksMkdir = new Error('mkdir permission is not enabled');\n" +
	"\n" +
	"function enableFilter() {\n" +
	"	const filter = document.body.querySelector('.filter');\n" +
	"	if (!filter) return;\n" +
	"\n" +
	"	const input = filter.querySelector('input');\n" +
	"	if (!input) return;\n" +
	"\n" +
	"	const clear = filter.querySelector('button') || document.createElement('button');\n" +
	"\n" +
	"	const entryList = document.querySelector(selectorEntryList);\n" +
	"\n" +
	"	let timeoutId;\n" +
	"	const doFilter = function () {\n" +
	"		const filteringText = input.value.trim().toLowerCase();\n" +
	"		if (filteringText === filteredText) return;\n" +
	"\n" +
	"		let entries;\n" +
	"		if (filteringText) {\n" +
	"			clear.style.display = 'block';\n" +
	"\n" +
	"			let selector;\n" +
	"			if (filteringText.includes(filteredText)) {	// increment search, find in visible entries\n" +
	"				selector = selectorEntryNotNone;\n" +
	"			} else if (filteredText.includes(filteringText)) {	// decrement search, find in hidden entries\n" +
	"				selector = selectorEntryIsNone;\n" +
	"			} else {\n" +
	"				selector = selectorEntry;\n" +
	"			}\n" +
	"			filteredText = filteringText;\n" +
	"\n" +
	"			entries = entryList.querySelectorAll(selector);\n" +
	"			entries.forEach(entry => {\n" +
	"				const name = entry.querySelector('.name');\n" +
	"				if (matchFilter(name.textContent)) {\n" +
	"					if (selector !== selectorEntryNotNone) {\n" +
	"						entry.classList.remove(classNone);\n" +
	"					}\n" +
	"				} else {\n" +
	"					if (selector !== selectorEntryIsNone) {\n" +
	"						entry.classList.add(classNone);\n" +
	"					}\n" +
	"				}\n" +
	"			});\n" +
	"		} else {	// filter cleared, show all entries\n" +
	"			clear.style.display = '';\n" +
	"			filteredText = '';\n" +
	"\n" +
	"			entries = entryList.querySelectorAll(selectorEntryIsNone);\n" +
	"			entries.forEach(entry => entry.classList.remove(classNone));\n" +
	"		}\n" +
	"	};\n" +
	"\n" +
	"	const onValueMayChange = function () {\n" +
	"		clearTimeout(timeoutId);\n" +
	"		timeoutId = setTimeout(doFilter, 350);\n" +
	"	};\n" +
	"	input.addEventListener('input', onValueMayChange);\n" +
	"	input.addEventListener('change', onValueMayChange);\n" +
	"\n" +
	"	const onEnter = function () {\n" +
	"		clearTimeout(timeoutId);\n" +
	"		input.blur();\n" +
	"		doFilter();\n" +
	"	};\n" +
	"	const onEscape = function () {\n" +
	"		if (input.value) {\n" +
	"			clearTimeout(timeoutId);\n" +
	"			input.value = '';\n" +
	"			doFilter();\n" +
	"		} else {\n" +
	"			input.blur();\n" +
	"		}\n" +
	"	};\n" +
	"\n" +
	"	input.addEventListener('keydown', function (e) {\n" +
	"		if (e.key === Enter) {\n" +
	"			onEnter();\n" +
	"			e.preventDefault();\n" +
	"		} else if (e.key === Escape) {\n" +
	"			onEscape();\n" +
	"			e.preventDefault();\n" +
	"		}\n" +
	"	});\n" +
	"	clear.addEventListener('click', function () {\n" +
	"		onEscape();\n" +
	"		input.focus();\n" +
	"	});\n" +
	"\n" +
	"	// init\n" +
	"	if (hasStorage) {\n" +
	"		const prevSessionFilter = sessionStorage.getItem(location.pathname);\n" +
	"		if (prevSessionFilter) {\n" +
	"			input.value = prevSessionFilter;\n" +
	"		}\n" +
	"		if (prevSessionFilter !== null) {\n" +
	"			sessionStorage.removeItem(location.pathname);\n" +
	"		}\n" +
	"\n" +
	"		window.addEventListener('pagehide', function () {\n" +
	"			const inputValue = input.value;\n" +
	"			if (inputValue) {\n" +
	"				sessionStorage.setItem(location.pathname, inputValue);\n" +
	"			}\n" +
	"		});\n" +
	"	}\n" +
	"	if (input.value) {\n" +
	"		doFilter();\n" +
	"	}\n" +
	"}\n" +
	"\n" +
	"function keepFocusOnBackwardForward() {\n" +
	"	const entryList = document.body.querySelector(selectorEntryList);\n" +
	"	entryList.addEventListener('focusin', function (e) {\n" +
	"		if (lastFocused !== e.target) {\n" +
	"			lastFocused = e.target;\n" +
	"		}\n" +
	"	});\n" +
	"	window.addEventListener('pageshow', function () {\n" +
	"		if (lastFocused && lastFocused !== document.activeElement) {\n" +
	"			lastFocused.focus();\n" +
	"			lastFocused.scrollIntoView({block: 'center'});\n" +
	"		}\n" +
	"	});\n" +
	"}\n" +
	"\n" +
	"function focusChildOnNavUp() {\n" +
	"	function extractCleanUrl(url) {\n" +
	"		let sepIndex = url.indexOf('?');\n" +
	"		if (sepIndex < 0) sepIndex = url.indexOf('#');\n" +
	"		if (sepIndex >= 0) {\n" +
	"			url = url.slice(0, sepIndex);\n" +
	"		}\n" +
	"		return url;\n" +
	"	}\n" +
	"\n" +
	"	if (lastFocused) return;\n" +
	"\n" +
	"	let prevUrl = document.referrer;\n" +
	"	if (!prevUrl) return;\n" +
	"	prevUrl = extractCleanUrl(prevUrl);\n" +
	"\n" +
	"	const currUrl = extractCleanUrl(location.href);\n" +
	"\n" +
	"	if (prevUrl.length <= currUrl.length) return;\n" +
	"	if (prevUrl.slice(0, currUrl.length) !== currUrl) return;\n" +
	"	const goesUp = prevUrl.slice(currUrl.length);\n" +
	"	if (currUrl[currUrl.length - 1] !== '/' && goesUp[0] !== '/') return;\n" +
	"	const matchInfo = /[^/]+/.exec(goesUp);\n" +
	"	if (!matchInfo) return;\n" +
	"	let prevChildName = matchInfo[0];\n" +
	"	if (!prevChildName) return;\n" +
	"	prevChildName = decodeURIComponent(prevChildName);\n" +
	"	if (!matchFilter(prevChildName)) return;\n" +
	"\n" +
	"	const entries = Array.from(document.body.querySelectorAll(selectorEntryList + '>' + selectorEntryNotNone));\n" +
	"	const selectorName = '.field.name';\n" +
	"	const selectorLink = 'a';\n" +
	"	for (let i = 0; i < entries.length; i++) {\n" +
	"		const entry = entries[i];\n" +
	"		const elName = entry.querySelector(selectorName);\n" +
	"		if (!elName) continue;\n" +
	"\n" +
	"		let text = elName.textContent;\n" +
	"		if (text[text.length - 1] === '/') {\n" +
	"			text = text.slice(0, -1);\n" +
	"		}\n" +
	"		if (text !== prevChildName) continue;\n" +
	"\n" +
	"		const elLink = entry.querySelector(selectorLink);\n" +
	"		if (!elLink) break;\n" +
	"\n" +
	"		lastFocused = elLink;\n" +
	"		elLink.focus();\n" +
	"		elLink.scrollIntoView({block: 'center'});\n" +
	"		break;\n" +
	"	}\n" +
	"}\n" +
	"\n" +
	"function enableKeyboardNavigate() {\n" +
	"	const pathList = document.body.querySelector(selectorPathList);\n" +
	"	const entryList = document.body.querySelector(selectorEntryList);\n" +
	"	if (!pathList && !entryList) {\n" +
	"		return;\n" +
	"	}\n" +
	"\n" +
	"	function getFocusableSibling(container, isBackward, startA) {\n" +
	"		if (!container.childElementCount) return;\n" +
	"		if (!startA) {\n" +
	"			startA = container.querySelector(':focus');\n" +
	"		}\n" +
	"		let startLI = startA && startA.closest('li');\n" +
	"		if (!startLI) {\n" +
	"			startLI = isBackward ? container.firstElementChild : container.lastElementChild;\n" +
	"		}\n" +
	"\n" +
	"		let siblingLI = startLI;\n" +
	"		do {\n" +
	"			if (isBackward) {\n" +
	"				siblingLI = siblingLI.previousElementSibling || container.lastElementChild;\n" +
	"			} else {\n" +
	"				siblingLI = siblingLI.nextElementSibling || container.firstElementChild;\n" +
	"			}\n" +
	"		} while (siblingLI !== startLI && (\n" +
	"			siblingLI.classList.contains(classNone) ||\n" +
	"			siblingLI.classList.contains(classHeader)\n" +
	"		));\n" +
	"\n" +
	"		if (siblingLI) {\n" +
	"			const siblingA = siblingLI.querySelector('a');\n" +
	"			return siblingA;\n" +
	"		}\n" +
	"	}\n" +
	"\n" +
	"	function getFirstFocusableSibling(container) {\n" +
	"		const a = container.querySelector(`li:not(.${classNone}):not(.${classHeader}) a`);\n" +
	"		return a;\n" +
	"	}\n" +
	"\n" +
	"	function getLastFocusableSibling(container) {\n" +
	"		let a = container.querySelector('li a');\n" +
	"		a = getFocusableSibling(container, true, a);\n" +
	"		return a;\n" +
	"	}\n" +
	"\n" +
	"	function getMatchedFocusableSibling(container, isBackward, startA, buf) {\n" +
	"		let skipRound = buf.length === 1;	// find next single-char prefix\n" +
	"		let firstCheckedA;\n" +
	"		let secondCheckedA;\n" +
	"		let a = startA;\n" +
	"		do {\n" +
	"			if (skipRound) {\n" +
	"				skipRound = false;\n" +
	"				continue;\n" +
	"			}\n" +
	"			if (!a) {\n" +
	"				continue;\n" +
	"			}\n" +
	"\n" +
	"			// firstCheckedA maybe a focused a that not belongs to the list\n" +
	"			// secondCheckedA must be in the list\n" +
	"			if (!firstCheckedA) {\n" +
	"				firstCheckedA = a;\n" +
	"			} else if (firstCheckedA === a) {\n" +
	"				return;\n" +
	"			} else if (!secondCheckedA) {\n" +
	"				secondCheckedA = a;\n" +
	"			} else if (secondCheckedA === a) {\n" +
	"				return;\n" +
	"			}\n" +
	"\n" +
	"			const textContent = (a.querySelector('.name') || a).textContent.toLowerCase();\n" +
	"			if (textContent.startsWith(buf)) {\n" +
	"				return a;\n" +
	"			}\n" +
	"		} while (a = getFocusableSibling(container, isBackward, a));\n" +
	"	}\n" +
	"\n" +
	"	const ARROW_UP = 'ArrowUp';\n" +
	"	const ARROW_DOWN = 'ArrowDown';\n" +
	"	const ARROW_LEFT = 'ArrowLeft';\n" +
	"	const ARROW_RIGHT = 'ArrowRight';\n" +
	"\n" +
	"	const SKIP_TAGS = ['INPUT', 'BUTTON', 'TEXTAREA'];\n" +
	"\n" +
	"	const PLATFORM = navigator.platform || navigator.userAgent;\n" +
	"	const IS_MAC_PLATFORM = PLATFORM.includes('Mac') || PLATFORM.includes('iPhone') || PLATFORM.includes('iPad') || PLATFORM.includes('iPod');\n" +
	"\n" +
	"	let lookupKey;\n" +
	"	let lookupBuffer;\n" +
	"	let lookupStartA;\n" +
	"	let lookupTimer;\n" +
	"\n" +
	"	function clearLookupContext() {\n" +
	"		lookupKey = undefined;\n" +
	"		lookupBuffer = '';\n" +
	"		lookupStartA = null;\n" +
	"	}\n" +
	"\n" +
	"	clearLookupContext();\n" +
	"\n" +
	"	function delayClearLookupContext() {\n" +
	"		clearTimeout(lookupTimer);\n" +
	"		lookupTimer = setTimeout(clearLookupContext, 850);\n" +
	"	}\n" +
	"\n" +
	"	function lookup(container, key, isBackward) {\n" +
	"		key = key.toLowerCase();\n" +
	"\n" +
	"		let currentLookupStartA;\n" +
	"		if (key === lookupKey) {\n" +
	"			// same as last key, lookup next single-char prefix\n" +
	"			currentLookupStartA = container.querySelector(':focus');\n" +
	"		} else {\n" +
	"			if (!lookupStartA) {\n" +
	"				lookupStartA = container.querySelector(':focus');\n" +
	"			}\n" +
	"			currentLookupStartA = lookupStartA;\n" +
	"			if (lookupKey === undefined) {\n" +
	"				lookupKey = key;\n" +
	"			} else {\n" +
	"				// key changed, no more single-char prefix match\n" +
	"				lookupKey = '';\n" +
	"			}\n" +
	"		}\n" +
	"		lookupBuffer += key;\n" +
	"		delayClearLookupContext();\n" +
	"		return getMatchedFocusableSibling(container, isBackward, currentLookupStartA, lookupKey || lookupBuffer);\n" +
	"	}\n" +
	"\n" +
	"	let canArrowMove;\n" +
	"	let isToEnd;\n" +
	"	if (IS_MAC_PLATFORM) {\n" +
	"		canArrowMove = function (e) {\n" +
	"			return !(e.ctrlKey || e.shiftKey || e.metaKey);	// only allow Opt\n" +
	"		};\n" +
	"		isToEnd = function (e) {\n" +
	"			return e.altKey;	// Opt key\n" +
	"		};\n" +
	"	} else {\n" +
	"		canArrowMove = function (e) {\n" +
	"			return !(e.altKey || e.shiftKey || e.metaKey);	// only allow Ctrl\n" +
	"		};\n" +
	"		isToEnd = function (e) {\n" +
	"			return e.ctrlKey;\n" +
	"		};\n" +
	"	}\n" +
	"\n" +
	"	function getFocusItemByKeyPress(e) {\n" +
	"		if (SKIP_TAGS.includes(e.target.tagName)) {\n" +
	"			return;\n" +
	"		}\n" +
	"\n" +
	"		if (canArrowMove(e)) {\n" +
	"			switch (e.key) {\n" +
	"				case ARROW_DOWN:\n" +
	"					if (isToEnd(e)) {\n" +
	"						return getLastFocusableSibling(entryList);\n" +
	"					} else {\n" +
	"						return getFocusableSibling(entryList, false);\n" +
	"					}\n" +
	"				case ARROW_UP:\n" +
	"					if (isToEnd(e)) {\n" +
	"						return getFirstFocusableSibling(entryList);\n" +
	"					} else {\n" +
	"						return getFocusableSibling(entryList, true);\n" +
	"					}\n" +
	"				case ARROW_RIGHT:\n" +
	"					if (isToEnd(e)) {\n" +
	"						return getLastFocusableSibling(pathList);\n" +
	"					} else {\n" +
	"						return getFocusableSibling(pathList, false);\n" +
	"					}\n" +
	"				case ARROW_LEFT:\n" +
	"					if (isToEnd(e)) {\n" +
	"						return getFirstFocusableSibling(pathList);\n" +
	"					} else {\n" +
	"						return getFocusableSibling(pathList, true);\n" +
	"					}\n" +
	"			}\n" +
	"		}\n" +
	"		if (!e.ctrlKey && (!e.altKey || IS_MAC_PLATFORM) && !e.metaKey && e.key.length === 1) {\n" +
	"			return lookup(entryList, e.key, e.shiftKey);\n" +
	"		}\n" +
	"	}\n" +
	"\n" +
	"	document.addEventListener('keydown', function (e) {\n" +
	"		const newFocusEl = getFocusItemByKeyPress(e);\n" +
	"		if (newFocusEl) {\n" +
	"			e.preventDefault();\n" +
	"			newFocusEl.focus();\n" +
	"		}\n" +
	"	});\n" +
	"}\n" +
	"\n" +
	"function enhanceUpload() {\n" +
	"	const form = document.body.querySelector('.upload form');\n" +
	"	if (!form) return;\n" +
	"\n" +
	"	const fileInput = form.querySelector('input[type=file]');\n" +
	"	if (!fileInput) return;\n" +
	"\n" +
	"	const submitButton = form.querySelector('button:last-of-type');\n" +
	"	if (submitButton) submitButton.classList.add(classNone);\n" +
	"\n" +
	"	const uploadType = document.body.querySelector('.upload-type');\n" +
	"	if (!uploadType) return;\n" +
	"\n" +
	"	const file = 'file';\n" +
	"	const dirFile = 'dirfile';\n" +
	"	const innerDirFile = 'innerdirfile';\n" +
	"\n" +
	"	const itemKindFile = 'file';\n" +
	"\n" +
	"	const optFile = uploadType.querySelector('.' + file);\n" +
	"	const optDir = uploadType.querySelector('.' + dirFile);\n" +
	"	const optInnerDir = uploadType.querySelector('.' + innerDirFile);\n" +
	"	let optActive = optFile;\n" +
	"	const canMkdir = Boolean(optDir);\n" +
	"\n" +
	"	let itemsToFiles;\n" +
	"	if (location.protocol === protoHttps && typeof FileSystemHandle !== strUndef && DataTransferItem.prototype.getAsFileSystemHandle && !DataTransferItem.prototype.webkitGetAsEntry) {\n" +
	"		const handleKindFile = 'file';\n" +
	"		const handleKindDir = 'directory';\n" +
	"		itemsToFiles = async function (dataTransferItems) {\n" +
	"			async function handlesToFiles(handles, files, dirPath) {\n" +
	"				await Promise.all(handles.map(async handle => {\n" +
	"					if (handle.kind === handleKindFile) {\n" +
	"						const relativePath = dirPath + handle.name;\n" +
	"						const file = await handle.getFile();\n" +
	"						files.push({file, relativePath});\n" +
	"					} else if (handle.kind === handleKindDir) {\n" +
	"						const subHandles = [];\n" +
	"						for await(const subHandle of handle.values()) {\n" +
	"							subHandles.push(subHandle);\n" +
	"						}\n" +
	"						const relativePath = dirPath + handle.name + '/';\n" +
	"						if (subHandles.length) {\n" +
	"							await handlesToFiles(subHandles, files, relativePath);\n" +
	"						} else {\n" +
	"							const file = new File([''], handle.name);\n" +
	"							files.push({file, relativePath});\n" +
	"						}\n" +
	"					}\n" +
	"				}));\n" +
	"			}\n" +
	"\n" +
	"			let hasDir = false;\n" +
	"			const permDescriptor = {mode: 'read'};\n" +
	"			const handles = await Promise.all(Array.from(dataTransferItems).filter(item => item.kind === itemKindFile).map(async item => {\n" +
	"				const handle = await item.getAsFileSystemHandle();\n" +
	"				if (handle.kind === handleKindDir) {\n" +
	"					if (!canMkdir) throw errLacksMkdir;\n" +
	"					hasDir = true;\n" +
	"				}\n" +
	"				const permState = await handle.queryPermission(permDescriptor);\n" +
	"				if (permState !== 'granted') {\n" +
	"					await handle.requestPermission(permDescriptor);\n" +
	"				}\n" +
	"				return handle;\n" +
	"			}));\n" +
	"\n" +
	"			const files = [];\n" +
	"			await handlesToFiles(handles, files, '');\n" +
	"			return {files, hasDir};\n" +
	"		};\n" +
	"	} else {\n" +
	"		itemsToFiles = async function (dataTransferItems) {\n" +
	"			async function entriesToFiles(entries, files) {\n" +
	"				await Promise.all(entries.map(async (entry) => {\n" +
	"					if (entry.isFile) {\n" +
	"						const file = await new Promise((res, rej) => entry.file(res, rej));\n" +
	"						const relativePath = entry.fullPath.slice(1);\n" +
	"						files.push({file, relativePath});\n" +
	"					} else if (entry.isDirectory) {\n" +
	"						const reader = entry.createReader();\n" +
	"						const subTasks = [];\n" +
	"						for (let i = 0; ; i++) {\n" +
	"							const subEntries = await new Promise((res, rej) => reader.readEntries(res, rej));\n" +
	"							if (subEntries.length > 0) {\n" +
	"								subTasks.push(entriesToFiles(subEntries, files));\n" +
	"								continue;\n" +
	"							}\n" +
	"							if (i === 0) {\n" +
	"								const file = new File([''], entry.name);\n" +
	"								const relativePath = entry.fullPath.slice(1) + '/';\n" +
	"								files.push({file, relativePath});\n" +
	"							} else {\n" +
	"								await Promise.all(subTasks);\n" +
	"							}\n" +
	"							break;\n" +
	"						}\n" +
	"					}\n" +
	"				}));\n" +
	"			}\n" +
	"\n" +
	"			let hasDir = false;\n" +
	"			const entries = Array.from(dataTransferItems, item => {\n" +
	"				if (item.kind !== itemKindFile) return;\n" +
	"				const entry = item.webkitGetAsEntry();\n" +
	"				if (entry.isDirectory) {\n" +
	"					if (!canMkdir) throw errLacksMkdir;\n" +
	"					hasDir = true;\n" +
	"				}\n" +
	"				return entry;\n" +
	"			}).filter(Boolean);\n" +
	"			const files = [];\n" +
	"\n" +
	"			await entriesToFiles(entries, files);\n" +
	"			return {files, hasDir};\n" +
	"		};\n" +
	"	}\n" +
	"\n" +
	"	function enableFileDirModeSwitch() {\n" +
	"		const classActive = 'active';\n" +
	"		const title = form.querySelector('h4') || document.createElement('h4');\n" +
	"\n" +
	"		function onClickOptAny(optTarget, clearInput) {\n" +
	"			if (optTarget === optActive) {\n" +
	"				return false;\n" +
	"			}\n" +
	"\n" +
	"			optActive.classList.remove(classActive);\n" +
	"			optActive = optTarget;\n" +
	"			optActive.classList.add(classActive);\n" +
	"			title.textContent = optActive.title || optActive.textContent;\n" +
	"\n" +
	"			if (clearInput) {\n" +
	"				fileInput.value = '';\n" +
	"			}\n" +
	"\n" +
	"			return true;\n" +
	"		}\n" +
	"\n" +
	"		function onClickOptFile(e) {\n" +
	"			if (onClickOptAny(optFile, Boolean(e))) {\n" +
	"				fileInput.name = file;\n" +
	"				fileInput.webkitdirectory = false;\n" +
	"			}\n" +
	"		}\n" +
	"\n" +
	"		function onClickOptDir() {\n" +
	"			if (onClickOptAny(optDir, optActive === optFile)) {\n" +
	"				fileInput.name = dirFile;\n" +
	"				fileInput.webkitdirectory = true;\n" +
	"			}\n" +
	"		}\n" +
	"\n" +
	"		function onClickOptInnerDir() {\n" +
	"			if (onClickOptAny(optInnerDir, optActive === optFile)) {\n" +
	"				fileInput.name = innerDirFile;\n" +
	"				fileInput.webkitdirectory = true;\n" +
	"			}\n" +
	"		}\n" +
	"\n" +
	"		if (optFile) {\n" +
	"			optFile.addEventListener('click', onClickOptFile);\n" +
	"			fileInput.addEventListener('change', function (e) {\n" +
	"				// workaround fix for old browsers, select dir not work but still act like select files\n" +
	"				// switch back to file\n" +
	"				if (optActive === optFile) return;\n" +
	"\n" +
	"				const files = e.target.files;\n" +
	"				if (!files.length) return;\n" +
	"\n" +
	"				const hasDir = Array.from(files).some(file =>\n" +
	"					file.webkitRelativePath.includes('/')\n" +
	"				);\n" +
	"				if (!hasDir) {\n" +
	"					onClickOptFile();\n" +
	"				}\n" +
	"			});\n" +
	"		}\n" +
	"		if (optDir) {\n" +
	"			optDir.addEventListener('click', onClickOptDir);\n" +
	"		}\n" +
	"		if (optInnerDir) {\n" +
	"			optInnerDir.addEventListener('click', onClickOptInnerDir);\n" +
	"		}\n" +
	"\n" +
	"		if (hasStorage) {\n" +
	"			const uploadTypeField = 'upload-type';\n" +
	"			const prevUploadType = sessionStorage.getItem(uploadTypeField);\n" +
	"			if (prevUploadType === dirFile) {\n" +
	"				optDir && optDir.click();\n" +
	"			} else if (prevUploadType === innerDirFile) {\n" +
	"				optInnerDir && optInnerDir.click();\n" +
	"			} else {\n" +
	"				optFile && optFile.click();\n" +
	"			}\n" +
	"\n" +
	"			if (prevUploadType !== null) {\n" +
	"				sessionStorage.removeItem(uploadTypeField);\n" +
	"			}\n" +
	"\n" +
	"			window.addEventListener('pagehide', function () {\n" +
	"				const activeUploadType = fileInput.name;\n" +
	"				if (activeUploadType !== file) {\n" +
	"					sessionStorage.setItem(uploadTypeField, activeUploadType);\n" +
	"				}\n" +
	"			});\n" +
	"		} else {\n" +
	"			optFile && optFile.click();\n" +
	"		}\n" +
	"\n" +
	"		function switchToFileMode() {\n" +
	"			if (optFile && optActive !== optFile) {\n" +
	"				optFile.focus();\n" +
	"				onClickOptFile(true);\n" +
	"			}\n" +
	"		}\n" +
	"\n" +
	"		function switchToDirMode() {\n" +
	"			if (optDir && optActive !== optDir) {\n" +
	"				optDir.focus();\n" +
	"				onClickOptDir();\n" +
	"			}\n" +
	"		}\n" +
	"\n" +
	"		function switchToInnerDirMode() {\n" +
	"			if (optInnerDir && optActive !== optInnerDir) {\n" +
	"				optInnerDir.focus();\n" +
	"				onClickOptInnerDir();\n" +
	"			}\n" +
	"		}\n" +
	"\n" +
	"		return {switchToFileMode, switchToDirMode, switchToInnerDirMode};\n" +
	"	}\n" +
	"\n" +
	"	function enableUploadProgress(switchToFileMode, switchToDirMode, switchToInnerDirMode) {\n" +
	"		let uploading = null;\n" +
	"		const classUploading = 'uploading';\n" +
	"		const classFailed = 'failed';\n" +
	"		const elUploadStatus = document.body.querySelector('.upload-status');\n" +
	"		if (!elUploadStatus) return;\n" +
	"		const elProgress = elUploadStatus.querySelector('.progress');\n" +
	"		const elFailedMessage = elUploadStatus.querySelector('.warn .message');\n" +
	"\n" +
	"		function uploadBatch(files) {\n" +
	"			const maxCount = 2048;\n" +
	"			const fieldName = fileInput.name;\n" +
	"\n" +
	"			const slices = [];\n" +
	"			let totalSize = 0;\n" +
	"\n" +
	"			let parts = null;\n" +
	"			let count = Infinity;\n" +
	"			for (let i = 0; i <= files.length; i++) {\n" +
	"				if (count >= maxCount || i === files.length) {\n" +
	"					if (i > 0) slices.push(parts);\n" +
	"					if (i === files.length) break;\n" +
	"					parts = new FormData();\n" +
	"					count = 0;\n" +
	"				}\n" +
	"\n" +
	"				const {file, relativePath} = files[i];\n" +
	"				totalSize += file.size;\n" +
	"				count += 1;\n" +
	"				parts.append(fieldName, file, relativePath);\n" +
	"			}\n" +
	"			if (slices.length === 0) return;\n" +
	"\n" +
	"			let finishedSize = 0;\n" +
	"			elProgress.style.width = '';\n" +
	"			const onProgress = e => {\n" +
	"				if (e.lengthComputable) {\n" +
	"					const percent = 100 * (finishedSize + e.loaded) / totalSize;\n" +
	"					elProgress.style.width = percent + '%';\n" +
	"				}\n" +
	"			};\n" +
	"			const onUploadSuccess = e => {\n" +
	"				if (e.lengthComputable) {\n" +
	"					finishedSize += e.total;\n" +
	"				}\n" +
	"			};\n" +
	"			const {method, action} = form;\n" +
	"			return new Promise(function (resolve, reject) {\n" +
	"				const onFail = e => {\n" +
	"					slices.length = 0;\n" +
	"					reject(e);\n" +
	"				};\n" +
	"				const onDownloadSuccess = e => {\n" +
	"					const status = e.target.status;\n" +
	"					if (status < 200 || status >= 300) {\n" +
	"						onFail({type: e.target.statusText || status});\n" +
	"						return;\n" +
	"					}\n" +
	"					if (slices.length) {\n" +
	"						uploadSlice(slices.shift());\n" +
	"						return;\n" +
	"					}\n" +
	"\n" +
	"					elProgress.style.width = '100%';\n" +
	"					resolve();\n" +
	"				};\n" +
	"\n" +
	"				function uploadSlice(parts) {\n" +
	"					const xhr = new XMLHttpRequest();\n" +
	"					xhr.upload.addEventListener('progress', onProgress);\n" +
	"					xhr.upload.addEventListener('error', onFail);\n" +
	"					xhr.addEventListener('error', onFail);\n" +
	"					xhr.upload.addEventListener('abort', onFail);\n" +
	"					xhr.addEventListener('abort', onFail);\n" +
	"					xhr.upload.addEventListener('load', onUploadSuccess);\n" +
	"					xhr.addEventListener('load', onDownloadSuccess);\n" +
	"					xhr.open(method, action);\n" +
	"					xhr.setRequestHeader('accept', 'application/json');\n" +
	"					xhr.send(parts);\n" +
	"				}\n" +
	"\n" +
	"				uploadSlice(slices.shift());\n" +
	"			});\n" +
	"		}\n" +
	"\n" +
	"		async function tryUploadBatch(getFilesResult) {\n" +
	"			if (!uploading) {\n" +
	"				elUploadStatus.classList.remove(classFailed);\n" +
	"				elUploadStatus.classList.add(classUploading);\n" +
	"				uploading = Promise.resolve();\n" +
	"			}\n" +
	"\n" +
	"			const filesResultTask = getFilesResult();	// must extract DataTransferItems ASAP\n" +
	"			const localUploading = uploading = uploading.then(async () => {\n" +
	"				const filesResult = await filesResultTask;\n" +
	"				const {files, hasDir, isInnerDir} = filesResult;\n" +
	"\n" +
	"				if (hasDir) {\n" +
	"					isInnerDir ? switchToInnerDirMode() : switchToDirMode();\n" +
	"				} else {\n" +
	"					switchToFileMode();\n" +
	"				}\n" +
	"\n" +
	"				await uploadBatch(files);\n" +
	"				if (uploading === localUploading) {\n" +
	"					location.reload();\n" +
	"				}\n" +
	"			}).catch(err => {\n" +
	"				elFailedMessage.textContent = ' - ' + err.type;\n" +
	"				if (err === errLacksMkdir && typeof showUploadDirFailMessage === strFunction) {\n" +
	"					showUploadDirFailMessage();\n" +
	"				} else {\n" +
	"					logError(err);\n" +
	"				}\n" +
	"				elUploadStatus.classList.remove(classUploading);\n" +
	"				elUploadStatus.classList.add(classFailed);\n" +
	"				throw err;\n" +
	"			});\n" +
	"\n" +
	"			return localUploading;\n" +
	"		}\n" +
	"\n" +
	"		async function uploadFilesProgressively(filesResult) {\n" +
	"			await tryUploadBatch(() => filesResult);\n" +
	"		}\n" +
	"\n" +
	"		async function uploadItemsProgressively(dataTransferItems) {\n" +
	"			await tryUploadBatch(() => itemsToFiles(dataTransferItems));\n" +
	"		}\n" +
	"\n" +
	"		return {uploadFilesProgressively, uploadItemsProgressively};\n" +
	"	}\n" +
	"\n" +
	"	function enableFormUploadProgress(uploadFilesProgressively) {\n" +
	"		form.addEventListener('submit', function (e) {\n" +
	"			e.stopPropagation();\n" +
	"			e.preventDefault();\n" +
	"		});\n" +
	"\n" +
	"		fileInput.addEventListener('change', function () {\n" +
	"			let hasDir = false;\n" +
	"			const files = Array.from(fileInput.files, file => {\n" +
	"				const relativePath = file.webkitRelativePath || file.name;\n" +
	"				if (relativePath.includes('/')) hasDir = true;\n" +
	"				return {file, relativePath};\n" +
	"			});\n" +
	"			uploadFilesProgressively({files, hasDir, isInnerDir: fileInput.name === innerDirFile});\n" +
	"		});\n" +
	"	}\n" +
	"\n" +
	"	function enableDndUploadProgress(uploadItemsProgressively) {\n" +
	"		let isSelfDragging = false;\n" +
	"		const classDragging = 'dragging';\n" +
	"\n" +
	"		function onSelfDragStart() {\n" +
	"			isSelfDragging = true;\n" +
	"		}\n" +
	"\n" +
	"		function onDragEnd() {\n" +
	"			isSelfDragging = false;\n" +
	"		}\n" +
	"\n" +
	"		function onDragEnterOver(e) {\n" +
	"			if (isSelfDragging) return;\n" +
	"			if (e.dataTransfer.items.length) {\n" +
	"				e.stopPropagation();\n" +
	"				e.preventDefault();\n" +
	"				e.currentTarget.classList.add(classDragging);\n" +
	"			}\n" +
	"		}\n" +
	"\n" +
	"		function onDragLeave(e) {\n" +
	"			if (e.target === e.currentTarget) {\n" +
	"				e.target.classList.remove(classDragging);\n" +
	"			}\n" +
	"		}\n" +
	"\n" +
	"		function onDrop(e) {\n" +
	"			e.stopPropagation();\n" +
	"			e.preventDefault();\n" +
	"			e.currentTarget.classList.remove(classDragging);\n" +
	"			fileInput.value = '';\n" +
	"			uploadItemsProgressively(e.dataTransfer.items);\n" +
	"		}\n" +
	"\n" +
	"		document.body.addEventListener('dragstart', onSelfDragStart);\n" +
	"		document.body.addEventListener('dragend', onDragEnd);\n" +
	"		const dndTarget = document.documentElement;\n" +
	"		dndTarget.addEventListener('dragenter', onDragEnterOver);\n" +
	"		dndTarget.addEventListener('dragover', onDragEnterOver);\n" +
	"		dndTarget.addEventListener('dragleave', onDragLeave);\n" +
	"		dndTarget.addEventListener('drop', onDrop);\n" +
	"	}\n" +
	"\n" +
	"	function enablePasteUploadProgress(uploadFilesProgressively, uploadItemsProgressively) {\n" +
	"		const typeTextPlain = 'text/plain';\n" +
	"\n" +
	"		function getTimeStamp() {\n" +
	"			const now = new Date();\n" +
	"			let date = String(now.getFullYear() * 10000 + (now.getMonth() + 1) * 100 + now.getDate());\n" +
	"			let time = String(now.getHours() * 10000 + now.getMinutes() * 100 + now.getSeconds());\n" +
	"			let ms = String(now.getMilliseconds());\n" +
	"			date = date.padStart(8, '0');\n" +
	"			time = time.padStart(6, '0');\n" +
	"			ms = ms.padStart(3, '0');\n" +
	"			const ts = '-' + date + '-' + time + '-' + ms;\n" +
	"			return ts;\n" +
	"		}\n" +
	"\n" +
	"		function uploadPastedFile(file) {\n" +
	"			const ts = getTimeStamp();\n" +
	"			let filename = file.name;\n" +
	"			let dotIndex = filename.lastIndexOf('.');\n" +
	"			if (dotIndex < 0) {\n" +
	"				dotIndex = filename.length;\n" +
	"			}\n" +
	"			filename = filename.slice(0, dotIndex) + ts + filename.slice(dotIndex);\n" +
	"\n" +
	"			const files = [{file, relativePath: filename}];\n" +
	"			uploadFilesProgressively({files, hasDir: false});\n" +
	"		}\n" +
	"\n" +
	"		function isTextInput(el) {\n" +
	"			const tagName = el.tagName;\n" +
	"			if (tagName === 'TEXTAREA') {\n" +
	"				return true;\n" +
	"			}\n" +
	"			const nonTextInputTypes = ['hidden', 'radio', 'checkbox', 'button', 'reset', 'submit', 'image'];\n" +
	"			if (tagName === 'INPUT' && !nonTextInputTypes.includes(el.type)) {\n" +
	"				return true;\n" +
	"			}\n" +
	"			return false;\n" +
	"		}\n" +
	"\n" +
	"		document.documentElement.addEventListener('paste', function (e) {\n" +
	"			if (isTextInput(e.target)) return;\n" +
	"\n" +
	"			const data = e.clipboardData;\n" +
	"			const dItems = data.items;\n" +
	"			const dFiles = data.files;\n" +
	"\n" +
	"			// image pasted\n" +
	"			if (dItems.length === 1 && dFiles.length === 1 && dFiles[0].type.startsWith('image/')) {\n" +
	"				uploadPastedFile(dFiles[0]);\n" +
	"				return;\n" +
	"			}\n" +
	"\n" +
	"			// text pasted (with other data types in DataTransferItems\n" +
	"			if (dItems.length > 0 && dFiles.length === 0) {\n" +
	"				const textTypeIndex = data.types.findIndex(t => t === typeTextPlain);\n" +
	"				if (textTypeIndex < 0) return;\n" +
	"\n" +
	"				const textItem = dItems[textTypeIndex];\n" +
	"				textItem.getAsString(function (content) {\n" +
	"					const file = new File([content], 'text.txt', {type: typeTextPlain});\n" +
	"					uploadPastedFile(file);\n" +
	"				});\n" +
	"				return;\n" +
	"			}\n" +
	"\n" +
	"			// actual files/directories pasted\n" +
	"			if (dItems.length) {\n" +
	"				uploadItemsProgressively(dItems);\n" +
	"			}\n" +
	"		});\n" +
	"	}\n" +
	"\n" +
	"	const {switchToFileMode, switchToDirMode, switchToInnerDirMode} = enableFileDirModeSwitch();\n" +
	"	const {uploadFilesProgressively, uploadItemsProgressively} = enableUploadProgress(switchToFileMode, switchToDirMode, switchToInnerDirMode);\n" +
	"	enableFormUploadProgress(uploadFilesProgressively);\n" +
	"	enableDndUploadProgress(uploadItemsProgressively);\n" +
	"	enablePasteUploadProgress(uploadFilesProgressively, uploadItemsProgressively);\n" +
	"}\n" +
	"\n" +
	"function enableSelectActions() {\n" +
	"	const form = document.body.querySelector('form.entry-form');\n" +
	"	if (!form) return;\n" +
	"\n" +
	"	const entryList = form.querySelector(selectorEntryList);\n" +
	"	if (!entryList) return;\n" +
	"\n" +
	"	const btnDelete = form.querySelector('.action-list .delete');\n" +
	"	const btnToggleSelect = entryList.querySelector('.toggle-select');\n" +
	"	const chkSelectAll = entryList.querySelector('.select-all');\n" +
	"\n" +
	"	const classSelecting = 'selecting';\n" +
	"	const selectorVisible = `li${selectorNotNone}`;\n" +
	"	const selectorHidden = `li${selectorIsNone}`;\n" +
	"	const selectorCheckbox = '.select input[type=checkbox]';\n" +
	"	const selectorUnchecked = `${selectorCheckbox}:not(:checked)`;\n" +
	"	const selectorChecked = `${selectorCheckbox}:checked`;\n" +
	"	const selectorVisibleUnchecked = `${selectorVisible} ${selectorUnchecked}`;\n" +
	"	const selectorHiddenChecked = `${selectorHidden} ${selectorChecked}`;\n" +
	"\n" +
	"	form.addEventListener('submit', function () {\n" +
	"		entryList.querySelectorAll(selectorHiddenChecked).forEach(input => input.checked = false);\n" +
	"\n" +
	"		setTimeout(() => {\n" +
	"			form.classList.remove(classSelecting);\n" +
	"			entryList.querySelectorAll(selectorChecked).forEach(input => input.checked = false);\n" +
	"		}, 0);\n" +
	"	});\n" +
	"\n" +
	"	if (btnToggleSelect) {\n" +
	"		btnToggleSelect.addEventListener('click', function () {\n" +
	"			form.classList.toggle(classSelecting);\n" +
	"			const selecting = form.classList.contains(classSelecting);\n" +
	"			if (btnDelete) {\n" +
	"				btnDelete.disabled = !selecting;\n" +
	"			}\n" +
	"			if (!selecting) {\n" +
	"				entryList.querySelectorAll(selectorChecked).forEach(input => input.checked = false);\n" +
	"			}\n" +
	"		});\n" +
	"	}\n" +
	"\n" +
	"	if (chkSelectAll) {\n" +
	"		chkSelectAll.addEventListener('change', function (e) {\n" +
	"			const checked = e.target.checked;\n" +
	"			if (checked) {\n" +
	"				entryList.querySelectorAll(selectorHiddenChecked).forEach(input => input.checked = false);\n" +
	"				entryList.querySelectorAll(selectorVisibleUnchecked).forEach(input => input.checked = true);\n" +
	"			} else {\n" +
	"				entryList.querySelectorAll(selectorChecked).forEach(input => input.checked = false);\n" +
	"			}\n" +
	"		});\n" +
	"	}\n" +
	"\n" +
	"	if (typeof confirmDelete === strFunction) {\n" +
	"		if (btnDelete) {\n" +
	"			btnDelete.addEventListener('click', function (e) {\n" +
	"				if (!confirmDelete()) e.preventDefault();\n" +
	"			});\n" +
	"		}\n" +
	"	}\n" +
	"}\n" +
	"\n" +
	"enableFilter();\n" +
	"keepFocusOnBackwardForward();\n" +
	"focusChildOnNavUp();\n" +
	"enableKeyboardNavigate();\n" +
	"enhanceUpload();\n" +
	"enableSelectActions();\n" +
	""
