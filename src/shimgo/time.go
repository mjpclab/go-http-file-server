package shimgo

import "time"

func Time_AppendFormat(t time.Time, b []byte, layout string) []byte {
	formatted := t.Format(layout)
	return append(b, []byte(formatted)...)
}
