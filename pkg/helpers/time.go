package helpers

import "time"

func GetTime() time.Time {
	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		return time.Now().UTC().Add(-3 * time.Hour)
	}
	return time.Now().In(loc)
}
