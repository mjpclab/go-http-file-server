build:
	go build main.go
tpls:
	sed -i -e '/= `/,/`/ { /`/!d }' tpl/page.html.go
	sed -i -e '/= `/rtpl/page.html' tpl/page.html.go
	sed -i -e '/= `/,/`/ { /`/! { /^\s*$$/d; s/^\s\s*// } }' tpl/page.html.go
	sed -i -e 's;/../assert/;?assert=;g' tpl/page.html.go

	sed -i -e '/= `/,/`/ { /`/!d }' tpl/assert/main.css.go
	sed -i -e '/= `/rtpl/assert/main.css' tpl/assert/main.css.go
	sed -i -e '/= `/,/`/ { /`/! { /^\s*$$/d; s/^\s\s*// } }' tpl/assert/main.css.go

	sed -i -e '/= `/,/`/ { /`/!d }' tpl/assert/main.js.go
	sed -i -e '/= `/rtpl/assert/main.js' tpl/assert/main.js.go
	sed -i -e '/= `/,/`/ { /`/! { /^\s*$$/d; s/^\s\s*// } }' tpl/assert/main.js.go
rebuild: tpls build
