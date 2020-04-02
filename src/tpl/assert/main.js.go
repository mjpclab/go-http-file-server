package assert

const mainJs = `
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
`
