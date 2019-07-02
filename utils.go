package main

import "time"

const TimeStringLayout = "2006-01-02T15:04:05.000000"

/*
#
# Parse string like: "2019-07-01T18:11:07.000000" to time.Time
# with current system Location
#
*/
func ParseTimeInCurrentLocation(str string) (time.Time, error) {
	parsedTime, err := time.ParseInLocation(TimeStringLayout, str, time.Now().Location())
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}
