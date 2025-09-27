export CGO_ENABLED=0
OUTDIR='../output'
MAINNAME='ghfs'
MOD=$(go list ../src/)
source ./build.inc.version.sh
getLdFlags() {
	echo "-s -w -X $MOD/version.appVer=$VERSION -X $MOD/version.appArch=${ARCH:-$(go env GOARCH)}"
}
TAR=${TAR:-tar}
