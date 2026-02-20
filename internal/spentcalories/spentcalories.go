// Пакет spentcalories обрабатывает информацию о тренировках и рассчитывает потраченные калории
package spentcalories

import (
	"errors"
	"fmt"
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

// parseTraining разбирает строку с данными о тренировке
// Формат строки: "3456,Ходьба,3h00m" (количество шагов, тип активности, продолжительность)
// Возвращает количество шагов, тип активности, продолжительность и ошибку
func parseTraining(data string) (int, string, time.Duration, error) {
	// Разделяем строку по запятой
	parts := strings.Split(data, ",")

	// Проверяем, что получилось ровно 3 части
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат строки: ожидается 3 части, разделенных запятой")
	}

	// Убираем лишние пробелы
	stepsStr := strings.TrimSpace(parts[0])
	activityType := strings.TrimSpace(parts[1])
	durationStr := strings.TrimSpace(parts[2])

	// Проверяем, что строка с шагами не пустая
	if stepsStr == "" {
		return 0, "", 0, errors.New("отсутствуют данные о количестве шагов")
	}

	// Проверяем, что тип активности указан
	if activityType == "" {
		return 0, "", 0, errors.New("отсутствуют данные о типе активности")
	}

	// Преобразуем количество шагов в int
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования количества шагов: %w", err)
	}

	// Проверяем, что количество шагов положительное
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше 0: %d", steps)
	}

	// Преобразуем продолжительность в time.Duration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}

	// Проверяем, что продолжительность положительная
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть больше 0: %v", duration)
	}

	return steps, activityType, duration, nil
}

// distance вычисляет дистанцию в километрах на основе количества шагов и роста пользователя
func distance(steps int, height float64) float64 {
	// Проверка входных параметров
	if steps <= 0 || height <= 0 {
		return 0
	}

	// Вычисляем длину шага: рост * коэффициент длины шага
	stepLength := height * stepLengthCoefficient

	// Вычисляем дистанцию: (количество шагов * длина шага) / количество метров в километре
	return float64(steps) * stepLength / mInKm
}

// meanSpeed вычисляет среднюю скорость в км/ч
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверяем, что продолжительность положительная
	if duration <= 0 {
		return 0
	}

	// Вычисляем дистанцию
	distanceKm := distance(steps, height)

	// Вычисляем продолжительность в часах
	durationHours := duration.Hours()

	// Проверяем, что продолжительность в часах не равна 0
	if durationHours == 0 {
		return 0
	}

	// Вычисляем и возвращаем среднюю скорость
	return distanceKm / durationHours
}

// TrainingInfo возвращает отформатированную информацию о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	// Парсим входные данные
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга тренировки: %w", err)
	}

	// Проверка входных параметров
	if weight <= 0 {
		return "", fmt.Errorf("вес должен быть больше 0: %.2f", weight)
	}
	if height <= 0 {
		return "", fmt.Errorf("рост должен быть больше 0: %.2f", height)
	}

	// Вычисляем дистанцию
	distanceKm := distance(steps, height)

	// Вычисляем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// Вычисляем калории в зависимости от типа активности
	var calories float64
	switch activityType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activityType)
	}

	if err != nil {
		return "", fmt.Errorf("ошибка расчета калорий: %w", err)
	}

	// Формируем результирующую строку
	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activityType, duration.Hours(), distanceKm, speed, calories)

	return result, nil
}

// RunningSpentCalories вычисляет количество калорий, потраченных при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0: %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0: %.2f", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0: %.2f", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0: %v", duration)
	}

	// Вычисляем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// Проверяем, что скорость не равна 0
	if speed <= 0 {
		return 0, fmt.Errorf("скорость должна быть больше 0: %.2f", speed)
	}

	// Вычисляем продолжительность в минутах
	durationMinutes := duration.Minutes()

	// Вычисляем калории: (вес * скорость * продолжительность в минутах) / минуты в часе
	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
}

// WalkingSpentCalories вычисляет количество калорий, потраченных при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0: %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0: %.2f", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0: %.2f", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0: %v", duration)
	}

	// Вычисляем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// Проверяем, что скорость не равна 0
	if speed <= 0 {
		return 0, fmt.Errorf("скорость должна быть больше 0: %.2f", speed)
	}

	// Вычисляем продолжительность в минутах
	durationMinutes := duration.Minutes()

	// Вычисляем калории для бега, затем умножаем на корректирующий коэффициент для ходьбы
	calories := (weight * speed * durationMinutes) / minInH * walkingCaloriesCoefficient

	return calories, nil
}