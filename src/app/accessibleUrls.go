package app

import (
	"../goVirtualHost"
	"../util"
	"fmt"
)

func printAccessibleURLs(vhSvc *goVirtualHost.Service) {
	vhostsUrls := vhSvc.GetAccessibleURLs(false)
	file, teardown := util.GetTTYFile()

	for vhIndex := range vhostsUrls {
		fmt.Fprintln(file, "Host", vhIndex, "may be accessed by URLs:")
		for urlIndex := range vhostsUrls[vhIndex] {
			fmt.Fprintln(file, "  ", vhostsUrls[vhIndex][urlIndex])
		}
	}

	teardown()
}
