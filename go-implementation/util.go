package main

func appendCSVString(parentString string, value string) string {
	if value != "" {
		if parentString == "" {
			parentString = value
		} else {
			parentString += "," + value
		}
	}

	return parentString
}

func getStartStop(formattedResult dwml, layoutKey string, offset int) (string, string) {
	var returnStart string
	var returnEnd string

	var currentOffset int

	for _, timeLayout := range formattedResult.Data.TimeLayouts {
		if timeLayout.LayoutKey == layoutKey {
			currentOffset = 0
			for _, startTime := range timeLayout.StartValidTimes {
				if currentOffset == offset {
					returnStart = startTime
					break
				}
				currentOffset++
			}
			currentOffset = 0
			for _, endTime := range timeLayout.EndValidTimes {
				if currentOffset == offset {
					returnEnd = endTime
					break
				}
				currentOffset++
			}
		}
	}
	return returnStart, returnEnd
}
