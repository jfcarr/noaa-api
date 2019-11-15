package main

type forecast struct {
	Title                      string
	MoreInformation            string
	DailyMaxTemperature        []dailyMaxTemperature
	DailyMinTemperature        []dailyMinTemperature
	ProbabilityOfPrecipitation []probabilityOfPrecipitation
	CloudCoverAmount           []cloudCoverAmount
}

type dailyMaxTemperature struct {
	StartDate string
	EndDate   string
	Value     string
}

type dailyMinTemperature struct {
	StartDate string
	EndDate   string
	Value     string
}

type probabilityOfPrecipitation struct {
	StartDate string
	EndDate   string
	Value     string
}

type cloudCoverAmount struct {
	StartDate string
	Value     string
}
