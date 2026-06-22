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
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, errors.New("Неверный формат данных!")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка преобразования количества шагов! %w", err)
	}

	activity := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка парсинга продолжительности! %w", err)
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {

	stepFloat := float64(steps)
	stepLength := stepLengthCoefficient * height
	totalDistanceMeters := stepFloat * stepLength
	distanceKm := totalDistanceMeters / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	totalDistance := distance(steps, height)
	return totalDistance / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var totalCalories float64

	switch trainingType {
	case "Ходьба":
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		} else {
			totalCalories = calories
		}
	case "Бег":
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		} else {
			totalCalories = calories
		}
	default:
		return "", fmt.Errorf("Неизвестный тип тренировки!")
	}

	distance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f",
		trainingType, duration.Hours(), distance, speed, totalCalories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("Количество шагов не может быть отрицательным!")
	}
	if weight <= 0 {
		return 0, errors.New("Вес не может быть отрицательным!")
	}
	if height <= 0 {
		return 0, errors.New("Рост не может быть отрицательным!")
	}
	if duration <= 0 {
		return 0, errors.New("Продолжительность бега не может быть отрицательной!")
	}

	speedAvg := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := (weight * speedAvg * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	return calories * walkingCaloriesCoefficient, nil
}
