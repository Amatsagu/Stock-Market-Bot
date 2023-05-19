package models

import (
	"time"
)

type DateParam struct {
	Year  int
	Month time.Month
	Day   int
}
