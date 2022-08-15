package serverLog

import (
	"mjpclab.dev/ghfs/src/util"
	"time"
)

func NewBuffer(cap int) []byte {
	// prefix: 20 bytes, suffix '\n' 1 byte
	buf := make([]byte, 0, 21+cap)

	buf = util.AppendTimeSecond(buf, time.Now()) // 19 bytes
	buf = append(buf, ' ')                       // 1 byte

	return buf
}
