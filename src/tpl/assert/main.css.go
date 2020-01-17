package assert

const mainCss = `
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
font-family: Consolas, Monaco, "Andale Mono", "DejaVu Sans Mono", monospace;
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
a:focus {
background: #fffaee;
}
a:hover {
background: #f5f5f5;
}
input, button {
margin: 0;
padding: 0.25em 0;
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
.upload {
position: relative;
margin: 1em;
padding: 1em;
background: #f7f7f7;
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
.upload form {
margin: 0;
padding: 0;
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
margin: 0 1em;
overflow: hidden;
zoom: 1;
}
.archive a {
position: relative;
float: left;
margin: 1em 0.5em 0 0.5em;
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
.item-list {
margin: 1em;
}
.item-list a {
display: flex;
flex-flow: row nowrap;
align-items: center;
border-bottom: 1px #f5f5f5 solid;
zoom: 1;
}
.item-list span {
margin-left: 1em;
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
.item-list .name em {
font-style: normal;
font-weight: normal;
padding: 0 0.2em;
border: 1px #ddd solid;
border-radius: 3px;
}
.item-list .size {
white-space: nowrap;
text-align: right;
color: #666;
}
.item-list .time {
width: 10em;
color: #999;
text-align: right;
white-space: nowrap;
overflow: hidden;
float: right;
}
.error {
margin: 1em;
padding: 1em;
background: #ffc;
}
@media only screen and (max-width: 350px) {
.item-list .time {
display: none;
}
}
`
