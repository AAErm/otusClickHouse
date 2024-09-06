package generator

import (
	"math/rand"
	"time"

	"github.com/AAErm/otusClickHouse/project/domain"
)

var cities = []string{
	"Минск",
	"Гомель",
	"Могилёв",
	"Витебск",
	"Гродно",
	"Брест",
	"Барановичи",
	"Борисов",
	"Пинск",
	"Орша",
}

func generateEvents(user domain.User) []domain.Event {
	numEvents := getNumberOfEvents(user.YearsOld)
	events := make([]domain.Event, numEvents)

	for i := 0; i < numEvents; i++ {
		serviceID, price := generateServiceIDAndPrice(user.YearsOld, user.GenderCode)
		events[i] = domain.Event{
			ID:         i + 1,
			Price:      price, // Assuming random price generation, can be modified
			ServiceID:  serviceID,
			LocationID: generateLocationID(),
			Timestamp:  time.Now(),
		}
	}

	return events
}

func getNumberOfEvents(yearsOld int) int {
	switch {
	case yearsOld <= 16:
		return rand.Intn(51) + 50
	case yearsOld >= 17 && yearsOld <= 23:
		return rand.Intn(126) + 75
	case yearsOld >= 24 && yearsOld <= 35:
		return rand.Intn(201) + 100
	case yearsOld >= 36 && yearsOld <= 55:
		return rand.Intn(76) + 75
	default:
		return rand.Intn(91) + 10
	}
}

// если YearsOld < 16:
// - 50%  "общественный транспорт"
// - 25% "общепит"
// - 15% "развлечения"
// - 10% ServiceID случайное число от 1 до 650
// если 16 < YearsOld < 23 и GenderCode = "man":
// - 30%  "общественный транспорт"
// - 25% "общепит"
// - 10% "развлечения"
// - 10% "продуктовый магазин"
// - 25% ServiceID случайное число от 1 до 800
// если 16 < YearsOld < 23 и GenderCode = "woman":
// - 30%  "общественный транспорт"
// - 25% "общепит"
// - 20% "развлечения"
// - 15% "продуктовый магазин"
// - 10% ServiceID случайное число от 1 до 800
// если 23 < YearsOld < 35 и GenderCode = "man":
// - 10%  "общественный транспорт"
// - 15% "общепит"
// - 20% "развлечения"
// - 30% "продуктовый магазин"
// - 15% "автомобильные услуги"
// - 10% ServiceID случайное число от 1 до 800
// если 23 < YearsOld < 35 и GenderCode = "woman":
// - 20%  "общественный транспорт"
// - 15% "общепит"
// - 20% "развлечения"
// - 30% "продуктовый магазин"
// - 5% "автомобильные услуги"
// - 10% ServiceID случайное число от 1 до 800
// если 35 < YearsOld < 55 и GenderCode = "man":
// - 10%  "общественный транспорт"
// - 5% "общепит"
// - 20% "развлечения"
// - 30% "продуктовый магазин"
// - 15% "автомобильные услуги"
// - 10% аптека
// - 10% ServiceID случайное число от 1 до 800
// если 35 < YearsOld < 55 и GenderCode = "woman":
// - 10%  "общественный транспорт"
// - 5% "общепит"
// - 20% "развлечения"
// - 10% аптека
// - 30% "продуктовый магазин"
// - 5% "автомобильные услуги"
// - 10% ServiceID случайное число от 1 до 800
// если YearsOld > 55
// - 20%  "общественный транспорт"
// - 2% "общепит"
// - 10% "развлечения"
// - 30% аптека
// - 30% "продуктовый магазин"
// - 8% "автомобильные услуги"

func generateServiceIDAndPrice(yearsOld int, genderCode string) (int, int) {
	var serviceCategoryProbabilities []int
	switch {
	case yearsOld < 16:
		serviceCategoryProbabilities = []int{50, 25, 15, 10}
	case yearsOld < 23 && genderCode == "man":
		serviceCategoryProbabilities = []int{30, 25, 10, 10, 25}
	case yearsOld < 23 && genderCode == "woman":
		serviceCategoryProbabilities = []int{30, 25, 20, 15, 10}
	case yearsOld < 35 && genderCode == "man":
		serviceCategoryProbabilities = []int{10, 15, 20, 30, 15, 10}
	case yearsOld < 35 && genderCode == "woman":
		serviceCategoryProbabilities = []int{20, 15, 20, 30, 5, 10}
	case yearsOld < 55 && genderCode == "man":
		serviceCategoryProbabilities = []int{10, 5, 20, 10, 30, 15, 10}
	case yearsOld < 55 && genderCode == "woman":
		serviceCategoryProbabilities = []int{10, 5, 20, 10, 30, 5, 10}
	default: // YearsOld > 55
		serviceCategoryProbabilities = []int{20, 2, 10, 30, 30, 8}
	}

	serviceCategory := chooseCategory(serviceCategoryProbabilities)

	var serviceID, price int
	switch serviceCategory {
	case 0: // "общественный транспорт"
		if rand.Intn(100) < 90 {
			serviceID = rand.Intn(99) + 1
			price = 1
			if serviceID == 1 {
				price = rand.Intn(25) + 1
			}
		}
	case 1: // "общепит"
		serviceID = rand.Intn(250) + 301
		price = rand.Intn(36) + 15
	case 2: // "развлечения"
		serviceID = rand.Intn(101) + 550
		price = rand.Intn(46) + 5
	case 3: // "продуктовый магазин"
		serviceID = rand.Intn(101) + 200
		price = rand.Intn(396) + 5
	case 4: // "аптека"
		serviceID = rand.Intn(101) + 100
		price = rand.Intn(96) + 5
	case 5: // "автомобильные услуги"
		serviceID = rand.Intn(350) + 651
		price = rand.Intn(3998) + 3
	case 6: // "случайное число от 1 до 800"
		serviceID = rand.Intn(800) + 1
		price = generatePriceBasedOnServiceID(serviceID)
	}

	return serviceID, price
}

func chooseCategory(probabilities []int) int {
	total := 0
	for _, prob := range probabilities {
		total += prob
	}

	randValue := rand.Intn(total)
	current := 0
	for i, prob := range probabilities {
		current += prob
		if randValue < current {
			return i
		}
	}
	return -1
}

// если ServiceID = 1 -- Price случайный от 1 до 25. это электричка
// если категория услуги "общественный транспорт" и ServiceID != 1, Price =1
// если категория "аптека", Price от 5 до 100
// если "продуктовый магазин", Price от 5 до 400
// если "общепит", Price от 15 до 50
// если "развлечения", Price от 5 до 50
// если "автомобильные услуги", от 3 до 4000
func generatePriceBasedOnServiceID(serviceID int) int {
	switch {
	case serviceID == 1:
		return rand.Intn(25) + 1
	case serviceID < 100:
		return 1
	case serviceID < 200:
		return rand.Intn(96) + 5
	case serviceID < 300:
		return rand.Intn(396) + 5
	case serviceID < 550:
		return rand.Intn(36) + 15
	case serviceID < 650:
		return rand.Intn(46) + 5
	default: // ServiceID > 650
		return rand.Intn(3998) + 3
	}
}

// Минск: ~20% (около 2 миллионов жителей) 46.95
// Гомель: ~4.5% (около 500 тысяч жителей) 10.56
// Могилёв: ~3.6% (около 380 тысяч жителей) 8.45
// Витебск: ~3.5% (около 370 тысяч жителей) 8.22
// Гродно: ~3.4% (около 360 тысяч жителей) 7.98
// Брест: 	~3.3% (около 350 тысяч жителей) 7.75
// Барановичи: ~1.2% (около 140 тысяч жителей) 2.82
// Борисов: ~1.1% (около 135 тысяч жителей) 2.58
// Пинск: ~1% (около 130 тысяч жителей) 2.35
// Орша: ~1% 	(около 120 тысяч жителей) 2.35
func generateLocationID() int {
	randomValue := rand.Float64() * 100

	switch {
	case randomValue < 46.95:
		return 1
	case randomValue < 46.95+10.56:
		return 2
	case randomValue < 46.95+10.56+8.45:
		return 3
	case randomValue < 46.95+10.56+8.45+8.22:
		return 4
	case randomValue < 46.95+10.56+8.45+8.22+7.98:
		return 5
	case randomValue < 46.95+10.56+8.45+8.22+7.98+7.75:
		return 6
	case randomValue < 46.95+10.56+8.45+8.22+7.98+7.75+2.82:
		return 7
	case randomValue < 46.95+10.56+8.45+8.22+7.98+7.75+2.82+2.58:
		return 8
	case randomValue < 46.95+10.56+8.45+8.22+7.98+7.75+2.82+2.58+2.35:
		return 9
	default:
		return 10
	}
}
