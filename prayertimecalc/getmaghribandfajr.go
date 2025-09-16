package prayertimecalc

import (
	fmt "fmt"
	geocode "lastThird/geocode"
	os "os"
	time "time"

	calc "github.com/MSA-Software-LLC/adhan-go/pkg/calc"
	data "github.com/MSA-Software-LLC/adhan-go/pkg/data"
	util "github.com/MSA-Software-LLC/adhan-go/pkg/util"
)

// utilizaes the adhan go package to calculate all the prayer times of a given location
func GetPrayerTimes(city string, state string, country string) (time.Time, time.Time) {
	lat, long := geocode.ProcessGeoData(city, state, country)

	date := data.NewDateComponents(time.Date(2025, time.Month(9), 6, 0, 0, 0, 0, time.UTC))
	fmt.Println(date)

	params := calc.GetMethodParameters(calc.NORTH_AMERICA)
	params.Madhab = calc.HANAFI

	coordinates, err := util.NewCoordinates(lat, long)
	if err != nil {
		fmt.Printf("error in adding coordinates", err)
		os.Exit(1)
	}

	prayerTimes, err := calc.NewPrayerTimes(coordinates, date, params)
	if err != nil {
		fmt.Printf("error in getting prayer times", err)
		os.Exit(1)
	}

	err = prayerTimes.SetTimeZone("America/New_York")
	if err != nil {
		fmt.Printf("got error %+v", err)
		os.Exit(1)
	}

	fmt.Printf("Fajr: %+v\n", prayerTimes.Fajr)       // Fajr: 2015-07-12 04:42:00 -0400 EDT
	fmt.Printf("Sunrise: %+v\n", prayerTimes.Sunrise) // Sunrise: 2015-07-12 06:08:00 -0400 EDT
	fmt.Printf("Dhuhr: %+v\n", prayerTimes.Dhuhr)     // Dhuhr: 2015-07-12 13:21:00 -0400 EDT
	fmt.Printf("Asr: %+v\n", prayerTimes.Asr)         // Asr: 2015-07-12 18:22:00 -0400 EDT
	fmt.Printf("Maghrib: %+v\n", prayerTimes.Maghrib) // Maghrib: 2015-07-12 20:32:00 -0400 EDT
	fmt.Printf("Isha: %+v\n", prayerTimes.Isha)       // Isha: 2015-07-12 21:57:00 -0400 EDT

	return prayerTimes.Fajr, prayerTimes.Maghrib
}

// Using the prayer times from GetPrayerTimes, this method subtracts the time from isha to maghrib
func GetTahajjud(city string, state string, country string) time.Time {
	fajr, maghrib := GetPrayerTimes(city, state, country)

	nextFajr := fajr.Add(time.Hour * 24)
	duration := maghrib.Sub(nextFajr)
	fmt.Println("Duration",duration.Abs())

	twoThirdsOfNight := time.Duration(float64(duration.Abs() * 2 / 3))
	fmt.Println("Two thirds",twoThirdsOfNight)

	tahajjudStart := maghrib.Add(twoThirdsOfNight)
	fmt.Println("Tahajjud Start", tahajjudStart)

	return tahajjudStart
}
