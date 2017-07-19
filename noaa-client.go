package main

func main() {
	fr := forecastRequest{
		Latitude:  39.44,
		Longitude: -84.3,
		Product:   "time-series",
		Begin:     "2017-07-19T00:00:00",
		End:       "2017-07-19T20:00:00",
		MaxT:      "maxt",
		MinT:      "mint"}

	results := callService(fr)

	formattedResult := parseResults(results)

	displayResults(formattedResult)
}
