define embed
	sed -i -e '/= `/,/`/ { /`/!d }' $(1).go
	sed -i -e '/= `/r$(1)' $(1).go
	#sed -i -e '/= `/,/`/ { /`/! { /^\s*$$/d; s/^\s\s*// } }' $(1).go
endef

tpls:
	$(call embed, tpl/frontend/page.html)
	$(call embed, tpl/frontend/main.css)
	$(call embed, tpl/frontend/main.js)
