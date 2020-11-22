package asset

const mainCss = `
html, body {
margin: 0;
padding: 0;
background: #fff;
}
html {
font-family: "roboto_condensedbold", "Helvetica Neue", Helvetica, Arial, sans-serif;
}
body, input, textarea {
font-family: Consolas, "Lucida Console", "San Francisco Mono", Menlo, Monaco, "Andale Mono", "DejaVu Sans Mono", monospace;
}
body {
color: #333;
font-size: 0.625em;
font-variant-ligatures: none;
font-kerning: none;
hyphens: none;
padding-bottom: 1em;
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
padding: 0.4em 0.5em;
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
.path-list {
font-size: 1.5em;
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
.panel {
margin: 1em;
padding: 1em;
background: #f7f7f7;
}
.upload {
position: relative;
}
.upload::before {
display: none;
content: '';
position: absolute;
left: 0;
top: 0;
right: 0;
bottom: 0;
opacity: 0.7;
background: #c9c;
}
.upload.dragging::before {
display: block;
}
.upload input {
display: block;
width: 100%;
box-sizing: border-box;
}
.upload input + input {
margin-top: 0.5em;
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
top: 1.1em;
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
display: flex;
}
.filter .filter-text {
flex: 1 1 auto;
width: 100%;
box-sizing: border-box;
}
.item-list {
margin: 1em;
}
.item-list li {
position: relative;
zoom: 1;
}
.item-list li:hover {
background: #f5f5f5;
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
color: #800000;
font-weight: bold;
font-size: 1.6em;
line-height: 1em;
padding: 0.1875em 0.3125em 0.3125em;
}
.item-list .delete:hover {
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
@media (prefers-color-scheme: dark) {
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
background-color: #222;
}
a:focus {
background-color: #220;
}
a:hover:focus {
background-color: #2f2f0f;
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
.has-deletable .detail {
padding-right: 0;
}
.has-deletable .delete {
display: none;
}
}
`
