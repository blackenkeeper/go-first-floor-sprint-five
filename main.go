package main

import (
	"fmt"
	"math"
	"time"
)

// Общие константы для вычислений.
const (
	MInKm      = 1000 // количество метров в одном километре
	MinInHours = 60   // количество минут в одном часе
	LenStep    = 0.65 // длина одного шага
	CmInM      = 100  // количество сантиметров в одном метре
)

// Training общая структура для всех тренировок
type Training struct {
	TrainingType string        // тип тренировки
	Action       int           // количество повторов(шаги, гребки при плавании)
	LenStep      float64       // длина одного шага или гребка в м
	Duration     time.Duration // продолжительность тренировки
	Weight       float64       // вес пользователя в кг
}

// distance возвращает дистанцию, которую преодолел пользователь.
// Формула расчета:
// количество_повторов * длина_шага / м_в_км
func (t Training) distance() float64 {
	distance := float64(t.Action) * t.LenStep / MInKm
	return distance
}

// meanSpeed возвращает среднюю скорость бега или ходьбы.
func (t Training) meanSpeed() float64 {
	timeOfTrainingInHours := t.Duration.Hours()

	if timeOfTrainingInHours == 0 {
		return 0
	}

	meanSpeed := t.distance() / timeOfTrainingInHours
	return meanSpeed
}

// Calories возвращает количество потраченных килокалорий на тренировке.
// Пока возвращаем 0, так как этот метод будет переопределяться для каждого типа тренировки.
func (t Training) Calories() float64 {
	return 0
}

// InfoMessage содержит информацию о проведенной тренировке.
type InfoMessage struct {
	Training
	Distance float64
	Speed    float64
	Calories float64
}

// TrainingInfo возвращает труктуру InfoMessage, в которой хранится вся информация о проведенной тренировке.
func (t Training) TrainingInfo() InfoMessage {

	return InfoMessage{
		Training: t,
		Distance: t.distance(),
		Speed:    t.meanSpeed(),
		Calories: t.Calories(),
	}
}

// String возвращает строку с информацией о проведенной тренировке.
func (i InfoMessage) String() string {

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

// CaloriesCalculator интерфейс для структур: Running, Walking и Swimming.
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Константы для расчета потраченных килокалорий при беге.
const (
	CaloriesMeanSpeedMultiplier = 18   // множитель средней скорости бега
	CaloriesMeanSpeedShift      = 1.79 // коэффициент изменения средней скорости
)

// Running структура, описывающая тренировку Бег.
type Running struct {
	Training
}

// Calories возввращает количество потраченных килокалория при беге.
// Формула расчета:
// ((18 * средняя_скорость_в_км/ч * 1.79) * вес_спортсмена_в_кг / м_в_км * время_тренировки_в_часах * мин_в_часе) -
// Это переопределенный метод Calories() из Training.
func (r Running) Calories() float64 {
	runnigMeanSpeed := r.meanSpeed()
	runningTimeInMinutes := r.Duration.Hours() * MinInHours

	runningMeanSpeedModifier := CaloriesMeanSpeedMultiplier*runnigMeanSpeed + CaloriesMeanSpeedShift

	spentCaloriesWhileRunning := runningMeanSpeedModifier * r.Weight / MInKm * runningTimeInMinutes

	return spentCaloriesWhileRunning

}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
// Это переопределенный метод TrainingInfo() из Training.
func (r Running) TrainingInfo() InfoMessage {

	return InfoMessage{
		Training: r.Training,
		Distance: r.distance(),
		Speed:    r.meanSpeed(),
		Calories: r.Calories(),
	}
}

// Константы для расчета потраченных килокалорий при ходьбе.
const (
	CaloriesWeightMultiplier      = 0.035 // коэффициент для веса
	CaloriesSpeedHeightMultiplier = 0.029 // коэффициент для роста
	KmHInMsec                     = 0.278 // коэффициент для перевода км/ч в м/с
)

// Walking структура описывающая тренировку Ходьба
type Walking struct {
	Training
	Height float64
}

// Calories возвращает количество потраченных килокалорий при ходьбе.
// Формула расчета:
// ((0.035 * вес_спортсмена_в_кг + (средняя_скорость_в_метрах_в_секунду**2 / рост_в_метрах)
// * 0.029 * вес_спортсмена_в_кг) * время_тренировки_в_часах * мин_в_ч)
// Это переопределенный метод Calories() из Training.
func (w Walking) Calories() float64 {
	walkingMeanSpeedInMetresPerSecond := w.meanSpeed() * KmHInMsec
	heightInMetres := float64(w.Height / CmInM)
	trainingTimeInMinutes := w.Duration.Hours() * MinInHours

	firstWeightModifier := CaloriesWeightMultiplier * w.Weight
	secondWeightModifier := CaloriesSpeedHeightMultiplier * w.Weight

	walkingSpeedModifier := math.Pow(walkingMeanSpeedInMetresPerSecond, 2) / heightInMetres

	spentCaloriesWhileWalking := (firstWeightModifier + walkingSpeedModifier*secondWeightModifier) * trainingTimeInMinutes

	return spentCaloriesWhileWalking
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
// Это переопределенный метод TrainingInfo() из Training.
func (w Walking) TrainingInfo() InfoMessage {

	return InfoMessage{
		Training: w.Training,
		Distance: w.distance(),
		Speed:    w.meanSpeed(),
		Calories: w.Calories(),
	}
}

// Константы для расчета потраченных килокалорий при плавании.
const (
	SwimmingLenStep                  = 1.38 // длина одного гребка
	SwimmingCaloriesMeanSpeedShift   = 1.1  // коэффициент изменения средней скорости
	SwimmingCaloriesWeightMultiplier = 2    // множитель веса пользователя
)

// Swimming структура, описывающая тренировку Плавание
type Swimming struct {
	Training
	LengthPool int
	CountPool  int
}

// meanSpeed возвращает среднюю скорость при плавании.
// Формула расчета:
// длина_бассейна * количество_пересечений / м_в_км / продолжительность_тренировки_в_часах
// Это переопределенный метод Calories() из Training.
func (s Swimming) meanSpeed() float64 {
	timeOfTrainingInHours := s.Duration.Hours()

	if timeOfTrainingInHours == 0 {
		return 0
	}

	meanSpeed := float64(s.LengthPool) * float64(s.CountPool) / MInKm / timeOfTrainingInHours

	return meanSpeed
}

// Calories возвращает количество калорий, потраченных при плавании.
// Формула расчета:
// (средняя_скорость_в_км/ч + 1.1) * 2 * вес_спортсмена_в_кг * время_тренеровки_в_часах
// Это переопределенный метод Calories() из Training.
func (s Swimming) Calories() float64 {
	swimmingMeanSpeed := s.meanSpeed()
	trainingTime := s.Duration.Hours()
	spentCaloriesWhileSwimming := (swimmingMeanSpeed + SwimmingCaloriesMeanSpeedShift) *
		SwimmingCaloriesWeightMultiplier * s.Weight * trainingTime

	return spentCaloriesWhileSwimming
}

// TrainingInfo returns info about swimming training.
// Это переопределенный метод TrainingInfo() из Training.
func (s Swimming) TrainingInfo() InfoMessage {

	return InfoMessage{
		Training: s.Training,
		Distance: s.distance(),
		Speed:    s.meanSpeed(),
		Calories: s.Calories(),
	}
}

// ReadData возвращает информацию о проведенной тренировке.
func ReadData(training CaloriesCalculator) string {
	calories := training.Calories()
	info := training.TrainingInfo()

	info.Calories = calories

	return fmt.Sprint(info)
}

func main() {

	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep, // в 4-ом спринте значение lenStep для рассчёта дистанции не менялось, расчёт там был не верен:(
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}

	fmt.Println(ReadData(swimming))

	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}

	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}

	fmt.Println(ReadData(running))

}
