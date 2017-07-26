package main

import (
	"encoding/xml"
)

func parseResults(inputXML string) dwml {
	var result dwml

	xml.Unmarshal([]byte(inputXML), &result)

	return result
}
