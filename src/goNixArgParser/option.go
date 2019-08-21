package goNixArgParser

import (
	"strings"
)

func (opt *Option) String() string {
	sb := &strings.Builder{}

	for i, flag := range opt.Flags {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(flag.Name)
	}

	if len(opt.Summary) > 0 {
		sb.WriteByte('\n')
		sb.WriteString(opt.Summary)
	}

	if len(opt.Description) > 0 {
		sb.WriteByte('\n')
		sb.WriteString(opt.Description)
	}

	db := &strings.Builder{}
	for _, d := range opt.DefaultValue {
		if len(d) > 0 {
			if db.Len() > 0 {
				db.WriteString(", ")
			}
			db.WriteString(d)
		}
	}
	if db.Len() > 0 {
		sb.WriteByte('\n')
		sb.WriteString("Default: ")
		sb.WriteString(db.String())
	}

	sb.WriteByte('\n')

	return sb.String()
}
