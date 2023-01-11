define embedStr
	sed -i -e '/= `/,/`/ { /`/!d }' $(1).go
	sed -i -e '/= `/r$(1)' $(1).go
	#sed -i -e '/= `/,/`/ { /`/! { /^\s*$$/d; s/^\s\s*// } }' $(1).go
endef

define embedHex
	sed -i -e '/= \[\]byte/,/\}/ { /[{}]/!d }' $(1).go
	od -An -tx1 -w16 -v $(1) | sed -e 's/^\s*/\t/g' | sed -E -e 's/[0-9a-fA-F]+/0x&,/g' | sed -i -e '/= \[\]byte/r/dev/stdin' $(1).go
endef

tpls:
	$(call embedStr, tpl/defaultTheme/frontend/index.html)
	$(call embedStr, tpl/defaultTheme/frontend/index.css)
	$(call embedStr, tpl/defaultTheme/frontend/index.js)
	$(call embedHex, tpl/defaultTheme/frontend/favicon.ico)
