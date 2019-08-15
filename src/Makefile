build:
	go build main.go
mktpl:
	sed -i -e '/= `/,/`/ { /`/!d }' tpl/page.go
	sed -i -e '/= `/rtpl/page.html' tpl/page.go
rebuild: mktpl build
