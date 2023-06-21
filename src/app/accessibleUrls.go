package app

import (
	"mjpclab.dev/ghfs/src/goVirtualHost"
	"mjpclab.dev/ghfs/src/param"
	"mjpclab.dev/ghfs/src/util"
	"strconv"
)

func printAccessibleURLs(vhSvc *goVirtualHost.Service, params param.Params) {
	vhostsUrls := vhSvc.GetAccessibleURLs(false)
	file, teardown := util.GetTTYFile()

	for vhIndex := range vhostsUrls {
		prefix := ""
		if len(params[vhIndex].PrefixUrls) > 0 {
			prefix = params[vhIndex].PrefixUrls[0]
		}

		file.WriteString("Host " + strconv.Itoa(vhIndex) + " may be accessed by URLs:\n")
		for urlIndex := range vhostsUrls[vhIndex] {
			file.WriteString("  " + vhostsUrls[vhIndex][urlIndex] + prefix + "/\n")
		}
	}

	teardown()
}
