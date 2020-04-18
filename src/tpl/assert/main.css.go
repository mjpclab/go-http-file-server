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
font-family: Consolas, "Lucida Console", "San Francisco Mono", Menlo, Monaco, "Andale Mono", "DejaVu Sans Mono", monospace;
font-variant-ligatures: none;
font-kerning: none;
hyphens: none;
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
a:focus {
background: #fffaee;
}
a:hover {
background: #f5f5f5;
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
.mkdir {
margin: 1em;
padding: 1em;
background: #f7f7f7;
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
.item-list a {
display: flex;
flex-flow: row nowrap;
align-items: center;
border-bottom: 1px #f5f5f5 solid;
overflow: hidden;
zoom: 1;
}
.has-deletable .link {
padding-right: 2.2em;
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
border-bottom: 1px #f5f5f5 solid;
color: #800000;
font-weight: bold;
}
.item-list .delete:hover {
background: #fee;
}
.item-list .delete span {
margin-left: 0;
font-size: 1.5em;
}
.item-list span,
.item-list button {
margin: 0 0 0 1em;
flex-shrink: 0;
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
