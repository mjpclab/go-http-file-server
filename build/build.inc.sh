TMP='/tmp'
OUTDIR='../output'
MAINNAME='ghfs'
MOD=$(go list ../src/)
source ./build.inc.version.sh
LICENSE='../LICENSE'
getLdFlags() {
	echo "-s -w"
}
