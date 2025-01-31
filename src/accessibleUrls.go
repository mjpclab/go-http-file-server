package src

import (
	"mjpclab.dev/ghfs/src/util"
	"strconv"
)

func printAccessibleURLs(accessibleUrls [][]string) {
	file, teardown := util.GetTTYFile()

	for vhIndex, vhUrls := range accessibleUrls {
		file.WriteString("Host " + strconv.Itoa(vhIndex) + " may be accessed by URLs:\n")
		for _, url := range vhUrls {
			file.WriteString("  " + url + "\n")
		}
	}

	teardown()
}
