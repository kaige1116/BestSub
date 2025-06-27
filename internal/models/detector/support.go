package detector

type DetectorType string

const (
	DetectorTypeAlive        DetectorType = "alive"
	DetectorTypeSpeed        DetectorType = "speed"
	DetectorTypeCountry      DetectorType = "country_api"
	DetectorTypeCountryRegex DetectorType = "country_regex"
)
