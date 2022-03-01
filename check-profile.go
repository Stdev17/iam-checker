package main

import (
	"time"
)

// CheckProfileExpired filters out the valid profiles from the originally fetched data.
func CheckProfileExpired(hour time.Duration, given []IAMProfile) []IAMProfile {
	var filtered []IAMProfile
	for _, val := range given {
		if val.CreatedDate.Add(hour).Before(time.Now()) {
			filtered = append(filtered, val)
		}
	}

	return filtered
}
