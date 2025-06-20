package timeutil

import "time"

const dateLayoutStr = "20060102"

func ToTime(s string) (time.Time, error) {
	return time.Parse(dateLayoutStr, s)
}
