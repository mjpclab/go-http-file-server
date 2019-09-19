build:
	go build main.go
tpls:
	sed -i -e '/= `/,/`/ { /`/!d }' tpl/page.go
	sed -i -e '/= `/rtpl/page.html' tpl/page.go
	sed -i -e '/= `/,/`/ { /`/! { /^\s*$$/d; s/^\s\s*// } }' tpl/page.go
rebuild: tpls build
