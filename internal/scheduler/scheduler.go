package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/edkuzhakhmetov/go_final_project/internal/models"
)

const (
	Format = "20060102"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	tstart, err := time.Parse(Format, dstart)
	if err != nil {
		return "", fmt.Errorf("error parsing dstart: %w", err)
	}

	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}
	reps := strings.Split(repeat, " ")

	switch reps[0] {
	case "d":
		if len(reps) < 2 || reps[1] == "" {
			return "", fmt.Errorf("repeat is empty")
		}

		multiplier, err := strconv.Atoi(reps[1])
		if err != nil {
			return "", fmt.Errorf("Ошибка конвертации multiplier: %w", err)
		}
		if multiplier <= 0 {
			return "", fmt.Errorf("multiplier не является положительным числом")
		}
		if multiplier > 400 {
			return "", fmt.Errorf("multiplier слишком большой")
		}

		if tstart.After(nowDate) || tstart.Equal(nowDate) {
			res := tstart.AddDate(0, 0, multiplier)
			return res.Format(Format), nil
		}

		res := tstart
		for {
			res = res.AddDate(0, 0, multiplier)
			if res.After(nowDate) || res.Equal(nowDate) {
				break
			}
		}

		return res.Format(Format), nil

	case "y":
		if tstart.After(nowDate) || tstart.Equal(nowDate) {
			res := tstart.AddDate(1, 0, 0)
			return res.Format(Format), nil
		}

		res := tstart
		for {
			res = res.AddDate(1, 0, 0)
			if res.After(nowDate) || res.Equal(nowDate) {
				break
			}
		}
		return res.Format(Format), nil
	default:
		return "", fmt.Errorf("Некорректный формат")
	}
}

func ValidateTaskID(id string) (int, error) {
	if id == "" {
		return 0, fmt.Errorf("Не указан ID задачи")
	}
	res, err := strconv.Atoi(id)
	if err != nil || res <= 0 {
		return 0, fmt.Errorf("Некорректный ID задачи")
	}
	return res, nil
}

func ValidateTask(task models.Task) error {
	if task.Title == "" {
		return fmt.Errorf("Поле Title обязательно")
	}
	return nil
}
