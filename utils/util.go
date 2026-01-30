package utils

import (
	"time"

	"github.com/rs/zerolog/log"
)

var (
	ISTLocation *time.Location
)

func init() {
	var err error
	ISTLocation, err = time.LoadLocation("Asia/Kolkata")
	if err != nil {
		ISTLocation = time.FixedZone("IST", 5.5*60*60)
		log.Info().Any("time", time.Now().In(ISTLocation)).Msg("Using fixed zone for IST")
	}
}

func GetDurationToMidnightIST() time.Duration {

	now := time.Now().In(ISTLocation)

	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, ISTLocation)
	duration := midnight.Sub(now)

	return duration
}
