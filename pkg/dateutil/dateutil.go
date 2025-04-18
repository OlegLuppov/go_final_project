package dateutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SettingsRules struct {
	Rule        string                // Правило d или y или w или m
	Days        int                   // Дни для правила по дням
	DaysOfWeek  map[time.Weekday]bool // Дни недели для правила по неделям
	DaysOfMonth []string              // Дни месяца для правила по месяцам
	Months      []string              // Месяцы для правила по месяцам
}

const (
	DateLayoutYMD = "20060102"   // Шаблон формата даты YYYMMDD
	DateLayoutDMY = "02.01.2006" // Шаблон формата даты DD.MM.YYY
)

var checkRules = map[string]string{
	"d": "d",
	"y": "y",
	"w": "w",
	"m": "m",
}

// Возращает следующую дату в зависимости от правил и ошибку
func NextDate(now string, dateStart string, repeat string) (string, error) {

	parseDateStart, err := time.Parse(DateLayoutYMD, dateStart)

	if err != nil {
		return "", fmt.Errorf("parse date start: %s", err)
	}

	parseDateNow, err := time.Parse(DateLayoutYMD, now)

	if err != nil {
		return "", fmt.Errorf("parse date now: %s", err)
	}

	settingsRules, err := ParseRepeat(repeat)

	if err != nil {
		return "", err
	}

	// По дням
	if settingsRules.Rule == "d" {
		for {

			parseDateStart = parseDateStart.AddDate(0, 0, settingsRules.Days)

			if parseDateStart.After(parseDateNow) {
				break
			}
		}
	}

	// По году
	if settingsRules.Rule == "y" {
		for {
			parseDateStart = parseDateStart.AddDate(1, 0, 0)

			if parseDateStart.After(parseDateNow) {
				break
			}
		}
	}

	//По неделям
	if settingsRules.Rule == "w" {

		for {
			_, checkWeekDay := settingsRules.DaysOfWeek[parseDateStart.Weekday()]

			if parseDateStart.After(parseDateNow) && checkWeekDay {
				break
			}

			parseDateStart = parseDateStart.AddDate(0, 0, 1)
		}
	}

	//По месяцам
	if settingsRules.Rule == "m" {
		for {
			checkDay, err := CheckCurrDayMonth(parseDateStart.Day(), parseDateStart.Month(), settingsRules.DaysOfMonth)

			if err != nil {
				return "", err
			}

			checkMonth, err := CheckCurrMonth(parseDateStart.Month(), settingsRules.Months)

			if err != nil {
				return "", err
			}

			if parseDateStart.After(parseDateNow) && checkDay && checkMonth {
				break
			}

			parseDateStart = parseDateStart.AddDate(0, 0, 1)
		}
	}

	return parseDateStart.Format(DateLayoutYMD), nil
}

// Парсит правила и возращат правила в удобном виде и ошибку если правило не верное
func ParseRepeat(repeat string) (SettingsRules, error) {
	var SettingsRules = SettingsRules{}

	if len(repeat) == 0 {
		return SettingsRules, fmt.Errorf("repeat expected a non-empty line , but got an empty one")
	}

	paramsRepeat := strings.Split(repeat, " ")

	rule, check := checkRules[paramsRepeat[0]]

	if !check {
		return SettingsRules, fmt.Errorf("repeat expected d or y or w or m , but got %s", paramsRepeat[0])
	}

	if len(paramsRepeat) == 1 && rule != "y" {
		return SettingsRules, fmt.Errorf("two repeat values expected, but got %s", repeat)
	}

	if (rule == "w" || rule == "d") && strings.Contains(paramsRepeat[1], "-") {
		return SettingsRules, fmt.Errorf("expected repeat positive value, but got negative value")
	}

	SettingsRules.Rule = rule

	// Если правило по дням
	if rule == "d" {
		days, err := strconv.Atoi(paramsRepeat[1])

		if err != nil {
			return SettingsRules, err
		}

		if days > 400 {
			return SettingsRules, fmt.Errorf("expected to last no more than 400 days, but got %d", days)
		}

		SettingsRules.Days = days
		return SettingsRules, nil
	}

	// Если правило по неделям
	if rule == "w" {
		daysOfWeek, err := GetDaysOfWeek(strings.Split(paramsRepeat[1], ","))

		if err != nil {
			return SettingsRules, err
		}
		SettingsRules.DaysOfWeek = daysOfWeek
	}

	// Если правило по месяцам
	if rule == "m" {
		SettingsRules.DaysOfMonth = strings.Split(paramsRepeat[1], ",")

		if len(paramsRepeat) == 3 {
			SettingsRules.Months = strings.Split(paramsRepeat[2], ",")
		}

		err := checkMonths(SettingsRules)

		if err != nil {
			return SettingsRules, err
		}
	}

	return SettingsRules, nil
}

// Проверка текущего дня месяца на вхождение в правило, возвращает true если правило совпало с текущим днем и ошибку
func CheckCurrDayMonth(currDay int, currMonth time.Month, days []string) (bool, error) {
	for _, day := range days {
		dayNum, err := strconv.Atoi(day)

		if err != nil {
			return false, err
		}

		if dayNum == -1 {
			nextDay := time.Date(0, currMonth+1, 0, 0, 0, 0, 0, time.UTC).Day()

			if nextDay == currDay {
				return true, nil
			}

		}

		if dayNum == -2 {
			nextDay := time.Date(0, currMonth+1, 0, 0, 0, 0, 0, time.UTC).Day() - 1

			if nextDay == currDay {
				return true, nil
			}
		}

		if dayNum == currDay {
			return true, nil
		}
	}

	return false, nil
}

// Проверка текущего месяца на вхождение в правило, возвращает true если правило совпало с текущим месяцем и ошибку
func CheckCurrMonth(currMonth time.Month, months []string) (bool, error) {
	if len(months) == 0 {
		return true, nil
	}

	for _, month := range months {

		monthNum, err := strconv.Atoi(month)

		if err != nil {
			return false, err
		}

		if monthNum == int(currMonth) {
			return true, nil
		}
	}

	return false, nil
}

// Проверяет допустимые дни месяца и месяцы и возвращает ошибку если правило не верное
func checkMonths(rules SettingsRules) error {
	for _, day := range rules.DaysOfMonth {
		num, err := strconv.Atoi(day)

		if err != nil {
			return fmt.Errorf("incorrect day of month format: %s", err)
		}

		if num > 31 || num < -2 {
			return fmt.Errorf("day of the month is not allowed, expected -2...31, but got %d", num)
		}
	}

	if len(rules.Months) > 0 {
		for _, month := range rules.Months {
			num, err := strconv.Atoi(month)

			if err != nil {
				return fmt.Errorf("incorrect month format: %s", err)
			}

			if num < 0 || num > 12 {
				return fmt.Errorf("invalid month, expected 1...12, but got %d", num)
			}
		}
	}

	return nil
}

// Возвращает дни недели из правил по дням недели и ошибку
func GetDaysOfWeek(days []string) (map[time.Weekday]bool, error) {
	daysOfWeek := make(map[time.Weekday]bool)

	for _, day := range days {
		num, err := strconv.Atoi(day)

		if err != nil {
			return nil, err
		}

		if num > 7 || num < 1 {
			return nil, fmt.Errorf("expected days of week 1...7, by got %d", num)
		}

		dayOfWeek := time.Weekday(num)

		if num == 7 {
			dayOfWeek = time.Weekday(0)
		}

		daysOfWeek[dayOfWeek] = true

	}

	return daysOfWeek, nil
}
