package main

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
