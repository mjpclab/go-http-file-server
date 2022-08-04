TMP='/tmp'
OUTDIR='../output'
MAINNAME='ghfs'
MOD=$(go list ../src/)
source ./build.inc.version.sh
LICENSE='../LICENSE'
getLdFlags() {
	echo "-s -w -X $MOD/version.appVer=$VERSION -X $MOD/version.appArch=${ARCH:-$(go env GOARCH)}"
}
