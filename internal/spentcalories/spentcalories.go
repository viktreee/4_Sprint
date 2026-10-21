package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	//	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining принимает строку формата "шаги,вид тренировки,продолжительность".
// Возвращает количество шагов, вид тренировки и продолжительность.
// При неверном формате данных, ошибке преобразования, неположительном
// количестве шагов или продолжительности возвращает 0, пустую строку, 0 и ошибку.
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("данные должны содержать три значения через запятую")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше нуля")
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность тренировки должна быть больше нуля")
	}

	return steps, parts[1], duration, nil
}

// distance принимает количество шагов и рост пользователя в метрах
// а возвращает пройденную дистанцию в километрах
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLength

	return distanceMeters / mInKm
}

// meanSpeed принимает количество шагов steps, рост пользователя height и продолжительность активности duration
// и возвращает среднюю скорость.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distanceKm := distance(steps, height)

	return distanceKm / duration.Hours()
}

// TrainingInfo принимает строку с данными тренировки, вес пользователя
// в килограммах и рост в метрах. Возвращает строку с видом тренировки,
// продолжительностью, дистанцией, средней скоростью и потраченными калориями.
// При некорректных данных или неизвестном виде тренировки
// возвращает пустую строку и ошибку.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64

	switch trainingType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		log.Println(err)
		return "", err
	}

	distanceKm := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	info := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		trainingType,
		duration.Hours(),
		distanceKm,
		speed,
		calories,
	)

	return info, nil
}

// RunningSpentCalories рассчитывает и возвращает количество потраченных при беге килокалорий.
// Принимает количество шагов(steps), вес(weight) в килограммах, рост(height) в метрах и продолжительность бега(duration).
// Если любой из параметров меньше или равен нулю, возвращает 0 и ошибку.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше нуля")
	}

	if weight <= 0 {
		return 0, errors.New("вес должен быть больше нуля")
	}

	if height <= 0 {
		return 0, errors.New("рост должен быть больше нуля")
	}

	if duration <= 0 {
		return 0, errors.New("продолжительность бега должна быть больше нуля")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	calories := weight * speed * durationMinutes / minInH

	return calories, nil
}

// WalkingSpentCalories()
// принимает:
// steps int — количество шагов, weight float64 — вес(кг.), height float64 рост(м.) пользователя, duration time.Duration — продолжительность ходьбы;
// высчитывает и возвращает:
// количество калорий, потраченных при ходьбе(float64);
// если любой из параметров меньше или равен нулю, возвращает 0 и ошибку.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше нуля")
	}

	if weight <= 0 {
		return 0, errors.New("вес должен быть больше нуля")
	}

	if height <= 0 {
		return 0, errors.New("рост должен быть больше нуля")
	}

	if duration <= 0 {
		return 0, errors.New("продолжительность ходьбы должна быть больше нуля")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	calories := weight * speed * durationMinutes / minInH

	return calories * walkingCaloriesCoefficient, nil
}
