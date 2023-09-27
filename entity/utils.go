package entity

import (
	"time"

	"github.com/dc0d/persical"
)

const threshold = 1

func priceAndExpDateProcessor(price float64) (int64, int64) {
	a, b, c := persical.GregorianToPersian(
		time.Now().Year(),
		int(time.Now().Month()),
		time.Now().Day())

	var remaining = monthDays(b) - c
	if remaining < threshold {
		remaining += monthDays(b + 1)
		a, b, c = persical.PersianToGregorian(a, b+2, 0)
		expDay := time.Date(a, time.Month(b), c, 0, 0, 0, 0, time.Now().Location())
		return int64(price), expDay.Unix()
	}

	total := (float64(remaining) / float64(monthDays(b))) * float64(price)
	expDay := time.Now().Add(time.Hour * 24 * time.Duration(remaining))

	a, b, c = persical.PersianToGregorian(a, b+1, 0)
	expDay = time.Date(a, time.Month(b), c, 0, 0, 0, 0, time.Now().Location())

	return int64(total), expDay.Unix()
}

func monthDays(index int) (monthDays int) {
	monthDays = 31
	if index > 6 {
		monthDays = 30
	}
	return
}

