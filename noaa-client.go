package main

import (
	"flag"
)

func main() {
	var latitudeFlag float64
	var longitudeFlag float64
	var startDateTimeFlag string
	var endDateTimeFlag string

	flag.Float64Var(&latitudeFlag, "latitude", 39.44, "Latitude")
	flag.Float64Var(&longitudeFlag, "longitude", -84.3, "Longitude")
	flag.StringVar(
		&startDateTimeFlag,
		"startDateTime",
		"",
		"Starting date and time of forecast, e.g., 2017-07-29T00:00:00 (default to now)")
	flag.StringVar(
		&endDateTimeFlag,
		"endDateTime",
		"",
		"Ending date and time of forecast, e.g., 2017-07-30T00:00:00 (default to max)")
	flag.Parse()

	// "2017-07-29T00:00:00"

	fr := forecastRequest{
		Latitude:            float32(latitudeFlag),
		Longitude:           float32(longitudeFlag),
		Product:             "time-series",
		Begin:               startDateTimeFlag,
		End:                 endDateTimeFlag,
		MaxTemperature:      "maxt",
		MinTemperature:      "mint",
		ProbabilityOfPrecip: "pop12",
		SkyCover:            "sky"}

	results := fr.callService()

	formattedResult := parseResults(results)

	// formattedResult.displayResults()

	formattedResult.writeJSON()
}
