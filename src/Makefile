define embed
	sed -i -e '/= `/,/`/ { /`/!d }' $(1).go
	sed -i -e '/= `/r$(1)' $(1).go
	#sed -i -e '/= `/,/`/ { /`/! { /^\s*$$/d; s/^\s\s*// } }' $(1).go
endef

tpls:
	$(call embed, tpl/frontend/index.html)
	$(call embed, tpl/frontend/index.css)
	$(call embed, tpl/frontend/index.js)
