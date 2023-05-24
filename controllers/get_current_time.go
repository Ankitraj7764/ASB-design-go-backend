package controllers

import "time"

func GetCurrentTime() time.Time {
	currentTime := time.Now()
	return currentTime
}
