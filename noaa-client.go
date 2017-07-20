package main

func main() {
	fr := forecastRequest{
		Latitude:            39.44,
		Longitude:           -84.3,
		Product:             "time-series",
		Begin:               "2017-07-20T00:00:00",
		End:                 "2017-07-20T20:00:00",
		MaxTemperature:      "maxt",
		MinTemperature:      "mint",
		ProbabilityOfPrecip: "pop12",
		SkyCover:            "sky"}

	results := callService(fr)

	formattedResult := parseResults(results)

	displayResults(formattedResult)
}
