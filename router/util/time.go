package util

import "time"

// Mongodb saved UTC time not local time
func DBTime2Local(t time.Time) time.Time {
	_, offset := t.Zone()
	return t.Add(time.Duration(offset) * time.Second)
}
