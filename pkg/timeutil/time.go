package timeutil

import "time"

const dateLayoutStr = "20060102"

func toDate(s string) (time.Time, error) {
	return time.Parse(dateLayoutStr, s)
}

func TodaySpecialDate(s string) (bool, error) {
	t, err := toDate(s)
	if err != nil {
		return false, err
	}
	today := time.Now()
	return today.Month() == t.Month() && today.Day() == t.Day(), nil
}
