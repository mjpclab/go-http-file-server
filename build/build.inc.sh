export CGO_ENABLED=0
OUTDIR='../output'
MAINNAME='ghfs'
MOD=$(go list ../src/)
source ./build.inc.version.sh
getLdFlags() {
	echo "-s -w"
}
