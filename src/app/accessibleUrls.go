package app

func (app *App) GetAccessibleUrls(includeLoopback bool) (allUrls [][]string) {
	allOrigins := app.vhostSvc.GetAccessibleURLs(includeLoopback)
	allUrls = make([][]string, len(allOrigins))

	for vhIndex, vhOrigins := range allOrigins {
		vhPrefixes := app.params[vhIndex].PrefixUrls
		if len(vhPrefixes) == 0 {
			vhPrefixes = []string{""}
		}

		allUrls[vhIndex] = make([]string, 0, len(vhOrigins)*len(vhPrefixes))

		for _, origin := range vhOrigins {
			for _, prefix := range vhPrefixes {
				allUrls[vhIndex] = append(allUrls[vhIndex], origin+prefix+"/")
			}
		}
	}

	return
}
