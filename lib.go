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
		TimeLayouts []struct {
			LayoutKey      string   `xml:"layout-key"`
			StartValidTime []string `xml:"start-valid-time"`
			EndValidTime   []string `xml:"end-valid-time"`
		} `xml:"time-layout"`
		Parameters struct {
			ApplicableLocation string `xml:"applicable-location,attr"`
			Temperatures       []struct {
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Name       string   `xml:"name"`
				Value      []string `xml:"value"`
			} `xml:"temperature"`
			ProbabilityOfPrecip struct {
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Name       string   `xml:"name"`
				Value      []string `xml:"value"`
			} `xml:"probability-of-precipitation"`
			CloudCoverAmount struct {
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Name       string   `xml:"name"`
				Value      []string `xml:"value"`
			} `xml:"cloud-amount"`
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

	for _, timeLayout := range formattedResult.Data.TimeLayouts {
		fmt.Printf("\nKey: %s\n", timeLayout.LayoutKey)

		startStopTimes := make(map[int]string)

		var currentOffset int
		for _, startTime := range timeLayout.StartValidTime {
			// fmt.Printf("Start: %s\n", startTime)
			startStopTimes[currentOffset] = startTime
			currentOffset = currentOffset + 1
		}

		currentOffset = 0
		for _, endTime := range timeLayout.EndValidTime {
			// fmt.Printf("End: %s\n", endTime)
			startStopTimes[currentOffset] = fmt.Sprintf("%s - %s", startStopTimes[currentOffset], endTime)
			currentOffset = currentOffset + 1
		}

		startStopTimeList := timeLayout.StartValidTime
		currentOffset = 0
		for _, endTime := range timeLayout.EndValidTime {
			startStopTimeList[currentOffset] = fmt.Sprintf("%s - %s", startStopTimeList[currentOffset], endTime)
			currentOffset++
		}

		for _, currentTime := range startStopTimeList {
			fmt.Printf(" %s\n", currentTime)
		}
	}

	fmt.Printf("\n")

	for _, temperature := range formattedResult.Data.Parameters.Temperatures {
		fmt.Printf("%s (%s):\n", temperature.Name, temperature.TimeLayout)
		for _, tempValue := range temperature.Value {
			fmt.Printf(" Value: %s\n", tempValue)
		}
	}

	fmt.Printf("\n")

	probabilityOfPrecip := formattedResult.Data.Parameters.ProbabilityOfPrecip
	fmt.Printf("%s (%s):\n", probabilityOfPrecip.Name, probabilityOfPrecip.TimeLayout)
	for _, popValue := range probabilityOfPrecip.Value {
		fmt.Printf(" Value: %s\n", popValue)
	}

	fmt.Printf("\n")

	cloudCover := formattedResult.Data.Parameters.CloudCoverAmount
	fmt.Printf("%s (%s):\n", cloudCover.Name, cloudCover.TimeLayout)
	for _, cloudValue := range cloudCover.Value {
		fmt.Printf(" Value: %s\n", cloudValue)
	}
}
