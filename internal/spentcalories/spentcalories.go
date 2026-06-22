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

var (
	ErrInvalidArgumentCount = errors.New("expected 2 arguments")
	ErrZeroSteps            = errors.New("steps must be greater than zero")
	ErrZeroDuration         = errors.New("duration must be greater than zero")
	ErrZeroWeight           = errors.New("weight must be greater than zero")
	ErrZeroHeigth           = errors.New("height must be greater than zero")
	ErrUnknownTraining      = errors.New("неизвестный тип тренировки")
)

func parseTraining(data string) (int, string, time.Duration, error) {
	dataSlice := strings.Split(data, ",")

	if len(dataSlice) != 3 {
		return 0, "", 0, ErrInvalidArgumentCount
	}

	steps, err := strconv.Atoi(dataSlice[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, ErrZeroSteps
	}

	totalTime, err := time.ParseDuration(dataSlice[2])
	if err != nil {
		return 0, "", 0, err
	}
	if totalTime <= 0 {
		return 0, "", 0, ErrZeroDuration
	}

	return steps, dataSlice[1], totalTime, nil
}

func distance(steps int, height float64) float64 {
	stepLength := stepLengthCoefficient * height
	pathLength := stepLength * float64(steps)
	return pathLength / mInKm
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
		return "", ErrUnknownTraining
	}

	distance := distance(steps, height)
	speedAvg := meanSpeed(steps, height, duration)

	result := ""

	result += fmt.Sprintf("Тип тренировки: %s\n", trainingType)
	result += fmt.Sprintf("Длительность: %.2f ч.\n", duration.Hours())
	result += fmt.Sprintf("Дистанция: %.2f км.\n", distance)
	result += fmt.Sprintf("Скорость: %.2f км/ч\n", speedAvg)
	result += fmt.Sprintf("Сожгли калорий: %.2f\n", totalCalories)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, ErrZeroSteps
	}
	if weight <= 0 {
		return 0, ErrZeroWeight
	}
	if height <= 0 {
		return 0, ErrZeroHeigth
	}
	if duration <= 0 {
		return 0, ErrZeroDuration
	}

	speedAvg := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := (weight * speedAvg * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	cal, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	return cal * walkingCaloriesCoefficient, nil
}
