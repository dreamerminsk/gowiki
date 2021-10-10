package utils

import (
	"regexp"
	"strconv"
	"time"
)

var longDayNames = []string{
	"Воскресенье",
	"Понедельник",
	"Вторник",
	"Среда",
	"Четверг",
	"Пятница",
	"Суббота",
}

var shortDayNames = []string{
	"Вс",
	"Пн",
	"Вт",
	"Ср",
	"Чт",
	"Пт",
	"Сб",
}

var shortMonthNames = []string{
	"Янв",
	"Фев",
	"Мар",
	"Апр",
	"Май",
	"Июн",
	"Июл",
	"Авг",
	"Сен",
	"Окт",
	"Ноя",
	"Дек",
}

var longMonthNames = []string{
	"Январь",
	"Февраль",
	"Март",
	"Апрель",
	"Май",
	"Июнь",
	"Июль",
	"Август",
	"Сентябрь",
	"Октябрь",
	"Ноябрь",
	"Декабрь",
}

const timePattern = `.*?(?P<Day>\d{2}) (?P<Month>\D{3}) (?P<Year>\d{4}) (?P<Hours>\d{2}):(?P<Minutes>\d{2}):(?P<Seconds>\d{2}).*?`

var timeRegEx = regexp.MustCompile(timePattern)

func ParseTime(timeString string) time.Time {
	match := timeRegEx.FindStringSubmatch(timeString)

	d, _ := strconv.Atoi(match[1])
	y, _ := strconv.Atoi(match[3])
	h, _ := strconv.Atoi(match[4])
	m, _ := strconv.Atoi(match[5])
	s, _ := strconv.Atoi(match[6])

	return time.Date(
		y, time.Month(Find(shortMonthNames, match[2])+1), d,
		h, m, s, 0, time.UTC)
}

func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}
