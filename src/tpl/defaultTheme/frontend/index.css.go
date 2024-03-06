package frontend

const DefaultCss = `
html, body {
	margin: 0;
	padding: 0;
	background: #fff;
}

html {
	font-family: "roboto_condensedbold", "Helvetica Neue", Helvetica, Arial, sans-serif;
}

body {
	color: #333;
	font-size: 0.625em;
	font-variant-ligatures: none;
	font-variant-numeric: tabular-nums;
	font-kerning: none;
	-webkit-text-size-adjust: none;
	text-size-adjust: none;
	hyphens: none;
	padding-bottom: 2em;
}

body, input, textarea, button {
	font-family: "Cascadia Mono", Consolas, "Lucida Console", "San Francisco Mono", Menlo, Monaco, "Andale Mono", "DejaVu Sans Mono", "Jetbrains Mono NL", monospace;
}

input::-ms-clear {
	display: none;
}

form {
	margin: 0;
	padding: 0;
}

ul, ol, li {
	display: block;
	margin: 0;
	padding: 0;
}

a {
	display: block;
	padding: 0.5em;
	color: #000;
	text-decoration: none;
	outline: 0;
}

a:hover {
	background: #f5f5f5;
}

a:focus {
	background: #fffae0;
}

a:hover:focus {
	background: #faf7ea;
}

input, button {
	min-width: 0;
	margin: 0;
	padding: 0.25em 0;
}

input[type=button],
input[type=submit],
input[type=reset],
button {
	cursor: pointer;
}

input:disabled[type=button],
input:disabled[type=submit],
input:disabled[type=reset],
button:disabled {
	cursor: default;
}

input[type=text] {
	padding: 0.25em;
}

em {
	font-style: normal;
	font-weight: normal;
	padding: 0 0.2em;
	border: 1px #ddd solid;
	border-radius: 3px;
}

.none, :root body .none {
	display: none;
}

.hidden {
	visibility: hidden;
}


html::before {
	display: none;
	content: '';
	position: absolute;
	position: fixed;
	z-index: 2;
	left: 0;
	top: 0;
	right: 0;
	bottom: 0;
	opacity: 0.7;
	background: #c9c;
}

html.dragging::before {
	display: block;
}


.path-list {
	font-size: 1.5em;
	line-height: 1.2;
	overflow: hidden;
	border-bottom: 1px #999 solid;
	zoom: 1;
}

.path-list li {
	position: relative;
	float: left;
	text-align: center;
	white-space: nowrap;
}

.path-list a {
	display: block;
	padding-right: 1.2em;
	min-width: 1em;
	white-space: pre-wrap;
}

.path-list a:after {
	content: '';
	position: absolute;
	top: 50%;
	right: 0.5em;
	width: 0.4em;
	height: 0.4em;
	border: 1px solid;
	border-color: #ccc #ccc transparent transparent;
	-webkit-transform: rotate(45deg) translateY(-50%);
	transform: rotate(45deg) translateY(-50%);
}

.path-list li:last-child a {
	padding-right: 0.5em;
}

.path-list li:last-child a:after {
	display: none;
}

.login {
	position: absolute;
	z-index: 1;
	right: 0;
	padding: 0.5em 1em;
}

.tab {
	display: flex;
	white-space: nowrap;
	margin: 1em 1em -1em 1em;
}

.tab label {
	flex: 0 0 auto;
	margin-right: 0.5em;
	padding: 1em;
	cursor: pointer;
}

.tab label:focus {
	outline: 0;
	text-decoration: underline;
	text-decoration-style: dotted;
}

.tab label:hover {
	background: #fbfbfb;
}

.tab label.active {
	color: #000;
	background: #f7f7f7;
}

.tab label:last-child {
	margin-right: 0;
}

.panel {
	margin: 1em;
	padding: 1em;
	background: #f7f7f7;
}

.upload-status {
	visibility: hidden;
	position: absolute;
	position: sticky;
	z-index: 1;
	left: 0;
	top: 0;
	width: 100%;
	height: 4px;
	margin-bottom: -4px;
	background: #faf5fa;
	background-color: rgba(204, 153, 204, 0.1);
	pointer-events: none;
}

.upload-status.uploading,
.upload-status.failed {
	visibility: visible;
}

.upload-status .label {
	position: absolute;
	left: 0;
	top: 0;
	width: 100%;
	color: #fff;
	text-align: center;
	opacity: 0;
	transition: transform .2s, opacity .2s;
}

.upload-status .label .content {
	position: relative;
	display: inline-block;
	vertical-align: top;
	text-align: left;
	text-align: start;
	padding: 0.5em 1em;
	box-sizing: border-box;
	overflow-wrap: break-word;
	word-break: break-word;
}

.upload-status .info .content {
	padding-left: 2.5em;
	background: #c9c;
	background-color: rgba(204, 153, 204, 0.8);
}

@keyframes wheel {
	from {
		transform: rotate(0);
	}
	to {
		transform: rotate(360deg);
	}
}

.upload-status .info .content:before,
.upload-status .info .content:after {
	content: '';
	position: absolute;
	left: 1em;
	top: 0.70em;
	width: 1em;
	height: 1em;
	box-sizing: border-box;
	border: 2px solid rgba(255, 255, 255, 0.3);
	border-radius: 50%;
	animation: wheel 1s linear infinite;
}

.upload-status .info .content:after {
	border-color: currentColor transparent transparent transparent;
}

.upload-status .warn .content {
	background: #800000;
	background-color: rgba(128, 0, 0, 0.8);
}

.upload-status.uploading .info,
.upload-status.failed .warn {
	opacity: 1;
	-webkit-transform: translateY(25%);
	transform: translateY(25%);
}

.upload-status .progress {
	position: absolute;
	left: 0;
	top: 0;
	width: 0;
	height: 100%;
	background: #c9c;
}

.upload {
	position: relative;
}

.upload input,
.upload button {
	display: block;
	width: 100%;
	box-sizing: border-box;
}

.upload button {
	position: relative;
	margin-top: 0.5em;
	overflow: hidden;
}

.upload button span {
	position: relative;
}

.archive {
	margin: 1em;
	overflow: hidden;
	zoom: 1;
}

.archive a {
	position: relative;
	float: left;
	margin: 0 0.5em;
	padding: 1em 1em 1em 3em;
	border: 2px #f5f5f5 solid;
}

.archive a:hover {
	border-color: #ddd;
}

.archive a:before {
	content: '';
	position: absolute;
	left: 1.1em;
	top: 1em;
	height: 1em;
	width: 3px;
	background: #aaa;
}

.archive a:after {
	content: '';
	position: absolute;
	left: 0.6em;
	top: 1em;
	width: 0.5em;
	height: 0.5em;
	margin-left: 1px;
	border: 3px #aaa solid;
	border-top-color: transparent;
	border-left-color: transparent;
	-webkit-transform: rotate(45deg);
	transform: rotate(45deg);
}

.mkdir form {
	display: flex;
	align-items: center;
}

.mkdir .name {
	flex: 1 1 auto;
}

.mkdir .submit {
	padding-left: 0.5em;
	padding-right: 0.5em;
}

.filter {
	display: none;
}

:root .filter {
	display: block;
}

.filter .form {
	position: relative;
	display: flex;
}

.filter input {
	flex: 1 1 auto;
	width: 97%;
	padding-right: 1.5em;
	box-sizing: border-box;
}

.filter button {
	display: none;
	position: absolute;
	right: 0;
	top: 0;
	bottom: 0;
	border: 0;
	background: none;
	padding: 0 0.5em;
}

.item-list {
	margin: 1em;
	line-height: 1.2;
}

.item-list li {
	position: relative;
	zoom: 1;
}

.item-list li:hover {
	background: #f5f5f5;
}

.item-list a {
	padding: 0.6em;
}

.item-list .detail,
.item-list .delete {
	display: flex;
	flex-flow: row nowrap;
	align-items: center;
	border-bottom: 1px #f5f5f5 solid;
	overflow: hidden;
	zoom: 1;
}

.has-deletable .detail {
	padding-right: 2.2em;
}

.item-list .field {
	margin: 0 0 0 1em;
	flex-shrink: 0;
}

.item-list .name {
	flex-grow: 1;
	flex-shrink: 1;
	flex-basis: 0;
	margin-left: 0;
	font-size: 1.5em;
	white-space: pre-wrap;
	word-break: break-all;
}

.item-list .size {
	white-space: nowrap;
	text-align: right;
	color: #666;
}

.item-list .time {
	color: #999;
	text-align: right;
	white-space: nowrap;
	overflow: hidden;
}

.item-list .delete {
	position: absolute;
	top: 0;
	right: 0;
	bottom: 0;
	display: flex;
	align-items: stretch;
}

.item-list .delete button {
	border: 0;
	color: #800000;
	background: none;
	font-weight: bold;
	font-size: 1.6em;
	line-height: 1em;
	padding: 0.1875em 0.3125em 0.3125em;
}

.item-list .delete button:hover {
	background: #fee;
}

.item-list .header:hover {
	background: none;
}

.item-list .header .detail {
	background: #fcfcfc;
}

.item-list .header .field {
	display: inline-block;
	margin: 0;
	font-size: 1.5em;
	color: #808080;
	overflow: hidden;
}

.item-list .header .time {
	width: 6.5em;
	text-align: center;
}

.error {
	margin: 1em;
	padding: 1em;
	background: #ffc;
}

@media only screen and (prefers-color-scheme: light) {
	html {
		color-scheme: light;
	}
}

@media only screen and (prefers-color-scheme: dark) {
	html {
		color-scheme: dark;
	}

	html, body {
		background: #111;
	}

	body {
		color: #ccc;
	}

	a {
		color: #ddd;
	}

	a:hover {
		background-color: #333;
	}

	a:focus {
		background-color: #330;
	}

	a:hover:focus {
		background-color: #33331a;
	}

	em {
		border-color: #555;
	}

	.path-list {
		border-bottom-color: #999;
	}

	.path-list a:after {
		border-color: #555 #555 transparent transparent;
	}

	.tab label:hover {
		background-color: #181818;
	}

	.tab label.active {
		color: #fff;
		background-color: #222;
	}

	.panel {
		background-color: #222;
	}

	.archive a {
		border-color: #222;
	}

	.archive a:hover {
		border-color: #555;
	}

	.item-list li:hover {
		background: #222;
	}

	.item-list .detail,
	.item-list .delete {
		border-bottom-color: #222;
	}

	.item-list .size {
		color: #999;
	}

	.item-list .time {
		color: #666;
	}

	.item-list .delete {
		color: #f99;
	}

	.item-list .delete:hover {
		background-color: #433;
	}

	.item-list .header .detail {
		background-color: #181818;
	}

	.error {
		background: #663;
	}
}

@media only screen and (max-width: 375px) {
	.item-list .header .time {
		width: 4.05em;
	}

	.item-list .detail .time span {
		display: none;
	}
}

@media only screen and (max-width: 350px) {
	.item-list .detail .time {
		display: none;
	}
}

@media print {
	.panel, .archive {
		display: none;
	}

	:root .panel {
		display: none;
	}

	.tab {
		display: none;
	}

	.item-list li {
		page-break-inside: avoid;
		break-inside: avoid;
	}

	.item-list li.parent {
		display: none;
	}

	.has-deletable .detail {
		padding-right: 0;
	}

	.has-deletable .delete {
		display: none;
	}
}
`
