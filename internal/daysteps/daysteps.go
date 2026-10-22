package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage принимает строку формата "шаги,продолжительность".
// Возвращает количество шагов и продолжительность прогулки.
// При неверном формате данных, ошибке преобразования, неположительном
// количестве шагов или продолжительности возвращает 0, 0 и ошибку.
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("expected two comma-separated values")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("parse steps: %w", err)
	}

	if steps <= 0 {
		return 0, 0, errors.New("steps must be greater than zero")
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("parse duration: %w", err)
	}

	if duration <= 0 {
		return 0, 0, errors.New("walk duration must be greater than zero")
	}

	return steps, duration, nil
}

// DayActionInfo принимает строку формата "шаги,продолжительность",
// вес пользователя в килограммах и рост в метрах.
// Возвращает строку с количеством шагов, дистанцией в километрах
// и потраченными килокалориями.
// При некорректных данных выводит ошибку на экран и возвращает пустую строку.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	info := fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		steps,
		distanceKm,
		calories,
	)

	return info
}
