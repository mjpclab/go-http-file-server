TMP='/tmp'
OUTDIR='../output'
MAINNAME='ghfs'
MOD=$(go list ../src/)
source ./build.inc.version.sh
LDFLAGS="-s -w"
LICENSE='../LICENSE'
