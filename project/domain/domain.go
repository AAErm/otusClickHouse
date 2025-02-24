package domain

import "time"

// от 10 до 300 ивентов на использование
type Event struct {
	ID         int
	Price      int
	ServiceID  int
	LocationID int
	Timestamp  time.Time
}

// примерно 800сервисов по 5 категориями,
type Service struct {
	ID       int
	Name     string
	Category string
}

// Минск: ~20% (около 2 миллионов жителей) 46.95
// Гомель: ~4.5% (около 500 тысяч жителей) 10.56
// Могилёв: ~3.6% (около 380 тысяч жителей) 8.45
// Витебск: ~3.5% (около 370 тысяч жителей) 8.22
// Гродно: ~3.4% (около 360 тысяч жителей) 7.98
// Брест: ~3.3% (около 350 тысяч жителей) 7.75
// Барановичи: ~1.2% (около 140 тысяч жителей) 2.82
// Борисов: ~1.1% (около 135 тысяч жителей) 2.58
// Пинск: ~1% (около 130 тысяч жителей) 2.35
// Орша: ~1% (около 120 тысяч жителей) 2.35
type Location struct {
	ID   int
	City string
}

// около 1.5млн юзеров
type User struct {
	ID         int
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	FatherName string `json:"father_name"`
	YearsOld   int    `json:"years_old"`
	GenderCode string `json:"gender_code"`
	Bank       string `json:"bank"`
}
