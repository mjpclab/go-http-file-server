define embedStr
	sed -i -e '/= `/,/`/ { /`/!d }' $(1).go
	sed -i -e '/= `/r$(1)' $(1).go
endef

define embedJs
	sed -i -E -e '/= "" \+/,/^\s*""$$/ { /^\s*"(\\"|[^"])*" +/d }' $(1).go
	sed -e 's/"/\\"/g' $(1) | sed -e 's/^/\t"/' | sed -e 's/$$/\\n" +/' | sed -i -e '/= "" +/r/dev/stdin' $(1).go
endef

define embedHex
	sed -i -e '/= string(\[\]byte/,/\})/ { /[{}]/!d }' $(1).go
	od -An -tx1 -w16 -v $(1) | sed -e 's/^\s*/\t/g' | sed -E -e 's/[0-9a-fA-F]+/0x&,/g' | sed -i -e '/= string(\[\]byte/r/dev/stdin' $(1).go
endef

tpls:
	$(call embedStr, tpl/defaultTheme/frontend/index.html)
	$(call embedStr, tpl/defaultTheme/frontend/index.css)
	$(call embedJs, tpl/defaultTheme/frontend/index.js)
	$(call embedHex, tpl/defaultTheme/frontend/favicon.ico)
