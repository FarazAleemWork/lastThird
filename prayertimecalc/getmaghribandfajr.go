package prayertimecalc

import (
	fmt "fmt"
	log "log"
	time "time"

	calc "github.com/MSA-Software-LLC/adhan-go/pkg/calc"
	data "github.com/MSA-Software-LLC/adhan-go/pkg/data"
	util "github.com/MSA-Software-LLC/adhan-go/pkg/util"
)

// utilizaes the adhan go package to calculate all the prayer times of a given location
func GetPrayerTimes(lat float64, long float64, timezone string) (time.Time, time.Time) {

	date := data.NewDateComponents(time.Date(2025, time.Month(9), 6, 0, 0, 0, 0, time.UTC))
	fmt.Println(date)

	params := calc.GetMethodParameters(calc.NORTH_AMERICA)
	params.Madhab = calc.HANAFI

	coordinates, err := util.NewCoordinates(lat, long)
	if err != nil {
		log.Println("error in adding coordinates", err)
	}

	prayerTimes, err := calc.NewPrayerTimes(coordinates, date, params)
	if err != nil {
		log.Println("error in getting prayer times", err)
	}

	err = prayerTimes.SetTimeZone(timezone)
	if err != nil {
		log.Printf("got error %+v", err)
	}

	/* Can use for other apps or debugging. Will also need this to account for day light savings
	fmt.Printf("Fajr: %+v\n", prayerTimes.Fajr)       // Fajr: 2015-07-12 04:42:00 -0400 EDT
	fmt.Printf("Sunrise: %+v\n", prayerTimes.Sunrise) // Sunrise: 2015-07-12 06:08:00 -0400 EDT
	fmt.Printf("Dhuhr: %+v\n", prayerTimes.Dhuhr)     // Dhuhr: 2015-07-12 13:21:00 -0400 EDT
	fmt.Printf("Asr: %+v\n", prayerTimes.Asr)         // Asr: 2015-07-12 18:22:00 -0400 EDT
	fmt.Printf("Maghrib: %+v\n", prayerTimes.Maghrib) // Maghrib: 2015-07-12 20:32:00 -0400 EDT
	fmt.Printf("Isha: %+v\n", prayerTimes.Isha)       // Isha: 2015-07-12 21:57:00 -0400 EDT*/

	return prayerTimes.Fajr, prayerTimes.Maghrib
}

// Using the prayer times from GetPrayerTimes, this method subtracts the time from isha to maghrib
func GetTahajjud(lat float64, long float64, timezone string) (string, error) {
	fajr, maghrib := GetPrayerTimes(lat, long, timezone)

	nextFajr := fajr.Add(time.Hour * 24)
	duration := maghrib.Sub(nextFajr)
	log.Println("Duration", duration.Abs())

	twoThirdsOfNight := time.Duration(float64(duration.Abs() * 2 / 3))
	log.Println("Two thirds", twoThirdsOfNight)

	tahajjudStart := maghrib.Add(twoThirdsOfNight)
	log.Println("Tahajjud Start", tahajjudStart)

	return tahajjudStart.Format("15:04:05"), nil

}
