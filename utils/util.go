package utils

import "time"

var (
	ISTLocation *time.Location
)

func init() {
	var err error
	ISTLocation, err = time.LoadLocation("Asia/Kolkata")
	if err != nil {
		panic(err)
	}
}

func GetDurationToMidnightIST() time.Duration {

	now := time.Now().In(ISTLocation)

	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, ISTLocation)
	duration := midnight.Sub(now)

	return duration
}
