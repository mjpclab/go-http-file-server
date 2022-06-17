package util

var hexSeq = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

func ByteToHex(b byte) (high, low byte) {
	high = hexSeq[b>>4]
	low = hexSeq[b&15]
	return
}
