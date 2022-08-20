package app

import (
	"mjpclab.dev/ghfs/src/goVirtualHost"
	"mjpclab.dev/ghfs/src/util"
	"strconv"
)

func printAccessibleURLs(vhSvc *goVirtualHost.Service) {
	vhostsUrls := vhSvc.GetAccessibleURLs(false)
	file, teardown := util.GetTTYFile()

	for vhIndex := range vhostsUrls {
		file.WriteString("Host " + strconv.Itoa(vhIndex) + " may be accessed by URLs:\n")
		for urlIndex := range vhostsUrls[vhIndex] {
			file.WriteString("  " + vhostsUrls[vhIndex][urlIndex] + "\n")
		}
	}

	teardown()
}
