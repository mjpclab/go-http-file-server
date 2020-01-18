package util

import (
	"fmt"
	"testing"
	"time"
)

func TestFormatTimeMinute(t *testing.T) {
	time := time.Now()
	fmt.Println(FormatTime(time))
}
