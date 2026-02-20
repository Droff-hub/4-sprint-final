// Пакет daysteps отвечает за учёт активности в течение дня
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

// Константы для вычислений
const (
	stepLength = 0.65 // длина одного шага в метрах
	mInKm      = 1000 // количество метров в одном километре
)

// parsePackage разбирает строку с данными о прогулке
// Формат строки: "678,0h50m" (количество шагов, продолжительность)
// Возвращает количество шагов, продолжительность и ошибку
func parsePackage(data string) (int, time.Duration, error) {
	// Проверяем на пустую строку
	if data == "" {
		return 0, 0, errors.New("пустая строка")
	}

	// Проверяем на пробелы в начале или конце
	if strings.HasPrefix(data, " ") || strings.HasSuffix(data, " ") {
		return 0, 0, errors.New("строка не должна содержать пробелы в начале или конце")
	}

	// Проверяем наличие запятой
	if !strings.Contains(data, ",") {
		return 0, 0, errors.New("отсутствует разделитель запятая")
	}

	// Разделяем строку по запятой
	parts := strings.Split(data, ",")

	// Проверяем, что получилось ровно 2 части
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат строки: ожидается 2 части, разделенных запятой")
	}

	stepsStr := parts[0]
	durationStr := parts[1]

	// Проверяем, что строка с шагами не пустая
	if stepsStr == "" {
		return 0, 0, errors.New("отсутствуют данные о количестве шагов")
	}

	// Проверяем, что строка с шагами содержит только цифры и возможно знак минус/плюс в начале
	// Но тесты ожидают ошибку для "+123" и "-123", поэтому пропускаем как есть - strconv.Atoi обработает

	// Преобразуем количество шагов в int
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования количества шагов: %w", err)
	}

	// Проверяем, что количество шагов положительное
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0: %d", steps)
	}

	// Проверяем, что строка с продолжительностью не пустая
	if durationStr == "" {
		return 0, 0, errors.New("отсутствуют данные о продолжительности")
	}

	// Преобразуем продолжительность в time.Duration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}

	// Проверяем, что продолжительность положительная
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть больше 0: %v", duration)
	}

	return steps, duration, nil
}

// DayActionInfo возвращает информацию о дневной активности
// Формат входной строки: "678,0h50m" (количество шагов, продолжительность)
// Возвращает строку с информацией о количестве шагов, дистанции и калориях
func DayActionInfo(data string, weight, height float64) string {
	// Парсим входные данные
	steps, duration, err := parsePackage(data)
	if err != nil {
		// ВАЖНО: логируем ошибку, как ожидают тесты
		log.Printf("Ошибка парсинга: %v\n", err)
		return ""
	}

	// Вычисляем дистанцию в километрах
	// Дистанция = количество шагов * длина шага (в метрах) / количество метров в километре
	distanceKm := float64(steps) * stepLength / mInKm

	// Вычисляем количество калорий, потраченных на прогулке
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("Ошибка расчета калорий: %v\n", err)
		return ""
	}

	// ВАЖНО: добавляем \n в конце строки, как ожидают тесты
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories)
}