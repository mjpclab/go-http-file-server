(function () {
	var classNone = 'none';
	var classHeader = 'header';

	function enableKeyboardNavigate() {
		if (
			!document.querySelector ||
			!document.addEventListener ||
			!document.body.classList ||
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
				siblingLI.classList.contains(classNone) ||
				siblingLI.classList.contains(classHeader)
			));

			if (siblingLI) {
				var siblingA = siblingLI.querySelector('a');
				return siblingA;
			}
		}

		function getMatchedFocusableSibling(container, isPrev, startA, buf, key) {
			var skipRound = buf === key;
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
				} else if (firstCheckA && !secondCheckA) {
					secondCheckA = a
				} else if (secondCheckA === a) {
					return;
				}

				var textContent = (a.querySelector('.name') || a).textContent.toLowerCase();
				if (buf.length <= textContent.length && textContent.substring(0, buf.length) === buf) {
					return a;
				}
				if (!matchKeyA && textContent[0] === key) {
					matchKeyA = a;
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

		var lookupKey = '';
		var lookupBuffer = '';
		var lookupStartA = null;
		var lookupTimer;

		function delayClearLookupContext() {
			clearTimeout(lookupTimer);
			lookupTimer = setTimeout(function () {
				lookupBuffer = '';
				lookupStartA = null;
			}, 850);
		}

		function lookup(key) {
			key = key.toLowerCase();

			if (key === lookupKey && key === lookupBuffer) {
				// same as last key, lookup next for the same key as prefix
				lookupStartA = itemList.querySelector(':focus');
				lookupBuffer = lookupKey;
			} else {
				if (!lookupStartA) {
					lookupStartA = itemList.querySelector(':focus');
				}
				lookupKey = key;
				lookupBuffer += key;
			}
			delayClearLookupContext();
			return getMatchedFocusableSibling(itemList, false, lookupStartA, lookupBuffer, key);
		}

		document.addEventListener('keydown', function (e) {
			if (
				e.ctrlKey ||
				e.altKey ||
				SKIP_TAGS.indexOf(e.target.tagName) >= 0
			) {
				return;
			}

			var newFocusEl;

			if (e.key) {
				switch (e.key) {
					case LEFT:
					case ARROW_LEFT:
						if (!e.shiftKey && !e.metaKey) {
							newFocusEl = getFocusableSibling(pathList, true);
						}
						break;
					case RIGHT:
					case ARROW_RIGHT:
						if (!e.shiftKey && !e.metaKey) {
							newFocusEl = getFocusableSibling(pathList, false);
						}
						break;
					case UP:
					case ARROW_UP:
						if (!e.shiftKey && !e.metaKey) {
							newFocusEl = getFocusableSibling(itemList, true);
						}
						break;
					case DOWN:
					case ARROW_DOWN:
						if (!e.shiftKey && !e.metaKey) {
							newFocusEl = getFocusableSibling(itemList, false);
						}
						break;
					default:
						if (e.key.length === 1) {
							newFocusEl = lookup(e.key);
						}
						break;
				}
			} else if (e.keyCode) {
				switch (e.keyCode) {
					case ARROW_LEFT_CODE:
						if (!e.shiftKey && !e.metaKey) {
							newFocusEl = getFocusableSibling(pathList, true);
						}
						break;
					case ARROW_RIGHT_CODE:
						if (!e.shiftKey && !e.metaKey) {
							newFocusEl = getFocusableSibling(pathList, false);
						}
						break;
					case ARROW_UP_CODE:
						if (!e.shiftKey && !e.metaKey) {
							newFocusEl = getFocusableSibling(itemList, true);
						}
						break;
					case ARROW_DOWN_CODE:
						if (!e.shiftKey && !e.metaKey) {
							newFocusEl = getFocusableSibling(itemList, false);
						}
						break;
					default:
						if (e.keyCode >= 32 && e.keyCode <= 126) {
							newFocusEl = lookup(String.fromCharCode(e.keyCode));
						}
				}
			}
			if (newFocusEl) {
				e.preventDefault();
				newFocusEl.focus();
			}
		});
	}

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
				selector = selectorItemNone;
				items = document.body.querySelectorAll(selector);
				for (i = items.length - 1; i >= 0; i--) {
					items[i].classList.remove(classNone);
				}
			} else {
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
						item.classList.add(classNone);
					} else {
						item.classList.remove(classNone);
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

	enableKeyboardNavigate();
	enableDragUpload();
	enableFilter();
	enableNonRefreshDelete();
})();
