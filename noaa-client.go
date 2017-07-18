package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type forecastRequest struct {
	Latitude, Longitude             float32
	Product, Begin, End, MaxT, MinT string
}

func callService(fr forecastRequest) string {
	queryString := fmt.Sprintf("lat=%f&lon=%f&product=%s&begin=%s&end=%s&maxt=%s&mint=%s",
		fr.Latitude, fr.Longitude, fr.Product, fr.Begin, fr.End, fr.MaxT, fr.MinT)

	requestString := fmt.Sprintf("%s?%s",
		"http://graphical.weather.gov/xml/sample_products/browser_interface/ndfdXMLclient.php",
		queryString)

	res, err := http.Get(requestString)
	if err != nil {
		log.Fatal(err)
	}

	results, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s", results)
}

func main() {
	fr := forecastRequest{
		Latitude:  39.44,
		Longitude: -84.3,
		Product:   "time-series",
		Begin:     "2017-07-18T00:00:00",
		End:       "2017-07-18T20:00:00",
		MaxT:      "maxt",
		MinT:      "mint"}

	results := callService(fr)

	fmt.Printf("%s", results)
}
