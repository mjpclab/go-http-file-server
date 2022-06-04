TMP='/tmp'
OUTDIR='../output'
MAINNAME='ghfs'
MOD=$(go list ../src/)
VERSION=$(git describe --abbrev=0 --tags 2> /dev/null || git rev-parse --abbrev-ref HEAD 2> /dev/null)
LDFLAGS="-s -w -X $MOD/version.appVer=$VERSION"
LICENSE='../LICENSE'
