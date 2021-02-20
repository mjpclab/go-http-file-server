define embed
	sed -i -e '/= `/,/`/ { /`/!d }' $(1).go
	sed -i -e '/= `/r$(1)' $(1).go
	sed -i -e '/= `/,/`/ { /`/! { /^\s*$$/d; s/^\s\s*// } }' $(1).go
endef

tpls:
	$(call embed, tpl/page.html)
	sed -i -e 's;/../asset/;?asset=;g' tpl/page.html.go

	$(call embed, tpl/asset/main.css)

	$(call embed, tpl/asset/main.js)
