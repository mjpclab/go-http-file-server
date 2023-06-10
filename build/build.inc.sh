TMP='/tmp'
OUTDIR='../output'
MAINNAME='ghfs'
MOD=$(go list ../src/)
source ./build.inc.version.sh
LICENSE='../LICENSE'
LICENSE_GO='../src/shimgo/LICENSE_GO'
getLdFlags() {
	echo "-s -w"
}
