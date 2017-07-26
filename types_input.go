package main

import (
	"fmt"
		"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type forecastRequest struct {
	Latitude, Longitude                                           float32
	Product, Begin, End                                           string
	MaxTemperature, MinTemperature, ProbabilityOfPrecip, SkyCover string
}

func (fr forecastRequest) callService() string {
	queryString := fmt.Sprintf("lat=%f&lon=%f&product=%s&begin=%s&end=%s&maxt=%s&mint=%s&pop12=%s&sky=%s",
		fr.Latitude, fr.Longitude, fr.Product, fr.Begin, fr.End, fr.MaxTemperature, fr.MinTemperature, fr.ProbabilityOfPrecip, fr.SkyCover)

	requestString := fmt.Sprintf("%s?%s",
		"https://graphical.weather.gov/xml/sample_products/browser_interface/ndfdXMLclient.php", queryString)

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
		MoreWeatherInformation struct {
			ApplicableLocation string `xml:"applicable-location,attr"`
			Value              string `xml:",chardata"`
		} `xml:"moreWeatherInformation"`
		TimeLayouts []struct {
			LayoutKey       string   `xml:"layout-key"`
			StartValidTimes []string `xml:"start-valid-time"`
			EndValidTimes   []string `xml:"end-valid-time"`
		} `xml:"time-layout"`
		Parameters struct {
			ApplicableLocation string `xml:"applicable-location,attr"`
			Temperatures       []struct {
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Name       string   `xml:"name"`
				Values     []string `xml:"value"`
			} `xml:"temperature"`
			ProbabilityOfPrecipitation struct {
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Name       string   `xml:"name"`
				Values     []string `xml:"value"`
			} `xml:"probability-of-precipitation"`
			CloudCoverAmount struct {
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Name       string   `xml:"name"`
				Values     []string `xml:"value"`
			} `xml:"cloud-amount"`
		} `xml:"parameters"`
	} `xml:"data"`
}

func (formattedResult dwml) displayResults() {
	var currentOffset int

	fmt.Printf("%s\n", formattedResult.Head.Product.Title)
	fmt.Printf("More info: %s (%s)\n",
		formattedResult.Data.MoreWeatherInformation.Value,
		formattedResult.Data.MoreWeatherInformation.ApplicableLocation)

	fmt.Printf("\n")

	for _, temperature := range formattedResult.Data.Parameters.Temperatures {
		fmt.Printf("%s (%s):\n", temperature.Name, temperature.TimeLayout)
		currentOffset = 0
		for _, tempValue := range temperature.Values {
			startTime, endTime := getStartStop(formattedResult, temperature.TimeLayout, currentOffset)
			fmt.Printf(" Value: %s (%s to %s)\n", tempValue, startTime, endTime)
			currentOffset++
		}
	}

	fmt.Printf("\n")

	probabilityOfPrecip := formattedResult.Data.Parameters.ProbabilityOfPrecipitation
	fmt.Printf("%s (%s):\n", probabilityOfPrecip.Name, probabilityOfPrecip.TimeLayout)
	currentOffset = 0
	for _, popValue := range probabilityOfPrecip.Values {
		startTime, endTime := getStartStop(formattedResult, probabilityOfPrecip.TimeLayout, currentOffset)
		fmt.Printf(" Value: %s (%s to %s)\n", popValue, startTime, endTime)
		currentOffset++
	}

	fmt.Printf("\n")

	cloudCover := formattedResult.Data.Parameters.CloudCoverAmount
	fmt.Printf("%s (%s):\n", cloudCover.Name, cloudCover.TimeLayout)
	currentOffset = 0
	for _, cloudValue := range cloudCover.Values {
		startTime, endTime := getStartStop(formattedResult, cloudCover.TimeLayout, currentOffset)
		_ = endTime // suppress error for the unused endTime
		fmt.Printf(" Value: %s (%s)\n", cloudValue, startTime)
		currentOffset++
	}
}

func (formattedResult dwml) writeJSON() {
	var currentOffset int

	var formattedForecast forecast

	formattedForecast.Title = formattedResult.Head.Product.Title
	formattedForecast.MoreInformation = formattedResult.Data.MoreWeatherInformation.Value
	for _, temperature := range formattedResult.Data.Parameters.Temperatures {
		if temperature.Type == "maximum" {
			currentOffset = 0
			for _, tempValue := range temperature.Values {
				startTime, endTime := getStartStop(formattedResult, temperature.TimeLayout, currentOffset)
				formattedForecast.DailyMaxTemperature =
					append(formattedForecast.DailyMaxTemperature, dailyMaxTemperature{StartDate: startTime, EndDate: endTime, Value: tempValue})
				currentOffset++
			}
		}
		if temperature.Type == "minimum" {
			currentOffset = 0
			for _, tempValue := range temperature.Values {
				startTime, endTime := getStartStop(formattedResult, temperature.TimeLayout, currentOffset)
				formattedForecast.DailyMinTemperature =
					append(formattedForecast.DailyMinTemperature, dailyMinTemperature{StartDate: startTime, EndDate: endTime, Value: tempValue})
				currentOffset++
			}
		}
	}

	probabilityOfPrecip := formattedResult.Data.Parameters.ProbabilityOfPrecipitation
	currentOffset = 0
	for _, popValue := range probabilityOfPrecip.Values {
		startTime, endTime := getStartStop(formattedResult, probabilityOfPrecip.TimeLayout, currentOffset)
		formattedForecast.ProbabilityOfPrecipitation =
			append(formattedForecast.ProbabilityOfPrecipitation, probabilityOfPrecipitation{StartDate: startTime, EndDate: endTime, Value: popValue})
		currentOffset++
	}

	cloudCover := formattedResult.Data.Parameters.CloudCoverAmount
	currentOffset = 0
	for _, cloudValue := range cloudCover.Values {
		startTime, endTime := getStartStop(formattedResult, cloudCover.TimeLayout, currentOffset)
		_ = endTime // suppress error for the unused endTime
		formattedForecast.CloudCoverAmount =
			append(formattedForecast.CloudCoverAmount, cloudCoverAmount{StartDate: startTime, Value: cloudValue})
		currentOffset++
	}

	b, err := json.MarshalIndent(formattedForecast, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	stringTemp := string(b)
	fmt.Println(stringTemp)
}
