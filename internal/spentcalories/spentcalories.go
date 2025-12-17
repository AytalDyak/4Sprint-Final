package spentcalories

import (
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
	// Разделяем строку по запятой
	parts := strings.Split(data, ",")

	// Проверяем, что у нас 3 части
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных, ожидается 'шаги,активность,длительность'")
	}

	// Очищаем данные от пробелов
	stepsStr := strings.TrimSpace(parts[0])
	activity := strings.TrimSpace(parts[1])
	durationStr := strings.TrimSpace(parts[2])

	// Парсим количество шагов
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверный формат количества шагов: %v", err)
	}

	// Проверяем, что количество шагов больше 0
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	// Проверяем, что вид активности не пустой
	if activity == "" {
		return 0, "", 0, fmt.Errorf("вид активности не может быть пустым")
	}

	// Парсим длительность
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверный формат длительности: %v", err)
	}

	// Проверяем, что длительность больше 0
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("длительность должна быть больше 0")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	// Рассчитываем длину шага на основе роста
	stepLength := height * stepLengthCoefficient

	// Если рассчитанная длина шага слишком мала или отрицательная,
	// используем среднюю длину шага
	if stepLength <= 0 {
		stepLength = lenStep
	}

	// Вычисляем дистанцию в метрах и переводим в километры
	distanceMeters := float64(steps) * stepLength
	return distanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверяем, что продолжительность больше 0
	if duration <= 0 {
		return 0
	}

	// Вычисляем дистанцию
	dist := distance(steps, height)

	// Вычисляем среднюю скорость
	hours := duration.Hours()

	// Защита от деления на ноль
	if hours <= 0 {
		return 0
	}

	return dist / hours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("длительность должна быть больше 0")
	}

	// Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, fmt.Errorf("не удалось рассчитать скорость")
	}

	// Переводим продолжительность в минуты
	minutes := duration.Minutes()

	// Рассчитываем калории:
	calories := (weight * speed * minutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("длительность должна быть больше 0")
	}

	// Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, fmt.Errorf("не удалось рассчитать скорость")
	}

	// Переводим продолжительность в минуты
	minutes := duration.Minutes()

	// Рассчитываем калории
	calories := (weight * speed * minutes) / minInH
	calories = calories * walkingCaloriesCoefficient

	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Получаем данные о тренировке
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println("Ошибка парсинга данных:", err)
		return "", err
	}

	// Проверяем вес и рост
	if weight <= 0 {
		return "", fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return "", fmt.Errorf("рост должен быть больше 0")
	}

	var calories float64
	var caloriesErr error

	// Выбираем расчет калорий в зависимости от типа активности
	switch strings.ToLower(activity) {
	case "бег", "running", "run":
		calories, caloriesErr = RunningSpentCalories(steps, weight, height, duration)
	case "ходьба", "walking", "walk":
		calories, caloriesErr = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}

	// Проверяем ошибку расчета калорий
	if caloriesErr != nil {
		log.Println("Ошибка расчета калорий:", caloriesErr)
		return "", caloriesErr
	}

	// Рассчитываем дистанцию и среднюю скорость
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	// Форматируем строку результата
	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity,
		duration.Hours(),
		dist,
		speed,
		calories,
	)

	return result, nil
}
