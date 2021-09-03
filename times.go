package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var longDayNames = []string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

var shortDayNames = []string{
	"Sun",
	"Mon",
	"Tue",
	"Wed",
	"Thu",
	"Fri",
	"Sat",
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

const timePattern = `.*?(?P<Day>\d{2}) (?P<Month>\D{3}) (?P<Year>\d{4}) (?P<Hours>\d{2}):(?P<Minutes>\d{2}):(?P<Seconds>\d{2})`

func parseTime(timeString string) time.Time {
	var compRegEx = regexp.MustCompile(timePattern)
	match := compRegEx.FindStringSubmatch(timeString)
	d, _ := strconv.ParseInt(match[1], 10, 32)
	y, _ := strconv.ParseInt(match[3], 10, 32)
	h, _ := strconv.ParseInt(match[4], 10, 32)
	m, _ := strconv.ParseInt(match[5], 10, 32)
	s, _ := strconv.ParseInt(match[6], 10, 32)
	return time.Date(int(y), time.Month(Find(shortMonthNames, strings.TrimSpace(match[2]))+1), int(d), int(h), int(m), int(s), 0, time.UTC)
}

func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}
