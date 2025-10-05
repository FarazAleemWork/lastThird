package errorhandling

import (
	"regexp"
)

func IsValidCity(city string) bool {
	return len(city) <= 30 && regexp.MustCompile(`^[A-Za-z\s\-]+$`).MatchString(city)
}

func IsValidState(state string) bool {
	return len(state) <= 30 && regexp.MustCompile(`^[A-Za-z\s\-]+$`).MatchString(state)
}

func IsValidCountry(country string) bool {
	return len(country) <= 30 && regexp.MustCompile(`^[A-Za-z\s\-]+$`).MatchString(country)
}
