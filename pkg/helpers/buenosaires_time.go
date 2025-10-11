package helpers

import "time"

// NowBuenosAires retorna la hora actual en la zona horaria de Buenos Aires (UTC-3)
func NowBuenosAires() time.Time {
	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		return time.Now().UTC().Add(-3 * time.Hour)
	}
	return time.Now().In(loc)
}
