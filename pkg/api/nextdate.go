package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}

	tabRep := strings.Split(repeat, " ")
	switch tabRep[0] {
	case "d":
		if len(tabRep) != 2 {
			return "", fmt.Errorf("wrong format repeat")
		}

		addDays, err := strconv.Atoi(tabRep[1])
		if err != nil {
			return "", err
		}
		if addDays > 400 {
			return "", fmt.Errorf("wrong format repeat")
		}

		for {
			date = date.AddDate(0, 0, addDays)
			if afterNow(date, now) {
				break
			}
		}

	case "w":
		days := map[int]string{}

		if len(tabRep) <= 1 {
			return "", fmt.Errorf("wrong format repeat")
		}

		tabDays := strings.Split(tabRep[1], ",")

		for i := 0; i < len(tabDays); i++ {
			dayOfWeek, err := strconv.Atoi(tabDays[i])
			if err != nil {
				return "", err
			}

			if dayOfWeek < 1 || dayOfWeek > 7 {
				return "", fmt.Errorf("wrong format repeat")
			}

			days[dayOfWeek] = ""
		}

		for {
			date = date.AddDate(0, 0, 1)
			day := int(date.Weekday())
			if day == 0 {
				day = 7
			}
			_, exists := days[day]

			if afterNow(date, now) && exists {
				break
			}
		}

	case "y":
		if len(tabRep) > 1 {
			return "", fmt.Errorf("wrong format repeat")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "m":
		if len(tabRep) <= 1 {
			return "", fmt.Errorf("wrong format repeat")
		}

		var calend [13][32]bool
		var valMonths []string

		valDays := strings.Split(tabRep[1], ",")
		if len(tabRep) < 3 {
			valMonths = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
		} else {
			valMonths = strings.Split(tabRep[2], ",")
		}

		for i := 0; i < len(valMonths); i++ {
			month, err := strconv.Atoi(valMonths[i])
			if err != nil {
				return "", err
			}
			if month < 1 || month > 12 {
				return "", fmt.Errorf("wrong format repeat")
			}

			for j := 0; j < len(valDays); j++ {

				day, err := strconv.Atoi(valDays[j])
				if err != nil {
					return "", err
				}
				if day > 31 || day < -2 {
					return "", fmt.Errorf("wrong parameter day")
				}
				if day != -1 && day != -2 {
					calend[month][day] = true
					continue
				}

				monthTime := time.Month(month)
				firstDayOfNextMonth := time.Date(date.Year(), monthTime+1, 1, 0, 0, 0, 0, now.Location())
				lastDayOfMonth := firstDayOfNextMonth.AddDate(0, 0, day)
				calend[month][lastDayOfMonth.Day()] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)
			if afterNow(date, now) && calend[date.Month()][date.Day()] {
				break
			}
		}

	case "":
		return "", nil
	default:
		return "", fmt.Errorf("wrong format repeat")
	}

	return date.Format("20060102"), nil

}

func afterNow(date1, date2 time.Time) bool {
	return date1.Format("20060102") > date2.Format("20060102")
}

func nextDayHandler(res http.ResponseWriter, req *http.Request) {
	var now time.Time
	valNow := req.FormValue("now")

	if valNow == "" {
		now = time.Now()
	} else {
		nowFromForm, err := time.Parse(formatDate, valNow)
		if err != nil {
			http.Error(res, "incorrect now", http.StatusBadRequest)
		}
		now = nowFromForm
	}

	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(res, "error get next date", http.StatusBadRequest)
	}

	io.WriteString(res, nextDate)
}
