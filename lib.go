package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type forecastRequest struct {
	Latitude, Longitude                                           float32
	Product, Begin, End                                           string
	MaxTemperature, MinTemperature, ProbabilityOfPrecip, SkyCover string
}

type dwml struct {
	Head struct {
		Product struct {
			OperationalMode string `xml:"operational-mode,attr"`
			Title           string `xml:"title"`
			Category        string `xml:"category"`
		} `xml:"product"`
	} `xml:"head"`
	Data struct {
		Location struct {
			LocationKey string `xml:"location-key"`
			Point       struct {
				Latitude  string `xml:"latitude,attr"`
				Longitude string `xml:"longitude,attr"`
			} `xml:"point"`
		} `xml:"location"`
		Parameters struct {
			ApplicableLocation string `xml:"applicable-location,attr"`
			Temperatures       []struct {
				Name  string `xml:"name"`
				Value string `xml:"value"`
			} `xml:"temperature"`
		} `xml:"parameters"`
	} `xml:"data"`
}

func callService(fr forecastRequest) string {
	queryString := fmt.Sprintf("lat=%f&lon=%f&product=%s&begin=%s&end=%s&maxt=%s&mint=%s&pop12=%s&sky=%s",
		fr.Latitude, fr.Longitude, fr.Product, fr.Begin, fr.End, fr.MaxTemperature, fr.MinTemperature, fr.ProbabilityOfPrecip, fr.SkyCover)

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

func parseResults(inputXML string) dwml {
	var result dwml

	xml.Unmarshal([]byte(inputXML), &result)

	return result
}

func displayResults(formattedResult dwml) {
	fmt.Printf("%s\n", formattedResult.Head.Product.Title)

	for _, temperature := range formattedResult.Data.Parameters.Temperatures {
		fmt.Printf("%s: %s\n", temperature.Name, temperature.Value)
	}
}