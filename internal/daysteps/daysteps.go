package daysteps

import (
	"errors"
	"fmt"
	"log"
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

var (
	ErrInvalidArgumentCount = errors.New("expected 2 arguments")
	ErrZeroSteps            = errors.New("steps must be greater than zero")
	ErrZeroDuration         = errors.New("duration must be greater than zero")
)

func parsePackage(data string) (int, time.Duration, error) {
	dataSlice := strings.Split(data, ",")

	if len(dataSlice) != 2 {
		return 0, 0, ErrInvalidArgumentCount
	}

	steps, err := strconv.Atoi(dataSlice[0])
	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		return 0, 0, ErrZeroSteps
	}

	totalTime, err := time.ParseDuration(dataSlice[1])
	if err != nil {
		return 0, 0, err
	}
	if totalTime <= 0 {
		return 0, 0, ErrZeroDuration
	}

	return steps, totalTime, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)

	if err != nil {
		log.Println(err)
		return ""
	}

	if steps <= 0 {
		log.Println(ErrZeroSteps)
		return ""
	}

	pathLength := float64(steps) * stepLength
	pathKilometers := pathLength / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	result := ""

	result += fmt.Sprintf("Количество шагов: %v.\n", steps)
	result += fmt.Sprintf("Дистанция составила %.2f км.\n", pathKilometers)
	result += fmt.Sprintf("Вы сожгли %.2f ккал.\n", calories)

	return result
}
