package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	// Открываем файл для записи
	file, err := os.Create("orders.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Записываем заголовки
	_, err = file.WriteString("order_id,table_number,order_date,total_amount,status,items_count\n")
	if err != nil {
		panic(err)
	}

	// Определяем временной диапазон за последние 10 лет
	now := time.Now()
	startDate := now.AddDate(-10, 0, 0)
	today := time.Now().Truncate(24 * time.Hour)

	// Генерация 1_000_000 записей
	for i := 1; i <= 1_000_000; i++ {
		tableNumber := (i % 20) + 1             // Номера столиков от 1 до 20
		orderDate := randomDate(startDate, now) // Случайная дата в диапазоне
		totalAmount := float64(i%100) * 10      // Сумма заказа (до 990)
		itemsCount := (i % 5) + 1               // Количество позиций от 1 до 5

		// Назначаем статус в зависимости от даты
		var status string
		if orderDate.Truncate(24 * time.Hour).Equal(today) {
			// Статусы для заказов сегодняшнего дня
			if rand.Intn(2) == 0 {
				status = "Новый"
				totalAmount = 0
				itemsCount = 0
			} else {
				status = "В обработке"
			}
		} else {
			// Статусы для заказов вне сегодняшнего дня
			status = randomStatus()
		}

		sTotalAmount := ""
		if totalAmount > 0 {
			sTotalAmount = fmt.Sprintf("%.2f", totalAmount)
		}

		sItemsCount := ""
		if itemsCount > 0 {
			sItemsCount = strconv.Itoa(itemsCount)
		}
		// Форматируем строку в CSV
		line := fmt.Sprintf("%d,%d,%s,%s,%s,%s\n", i, tableNumber, orderDate.Format("2006-01-02 15:04:05"), sTotalAmount, status, sItemsCount)

		// Записываем строку в файл
		_, err = file.WriteString(line)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Файл orders.csv успешно сгенерирован.")
}

func randomDate(start, end time.Time) time.Time {
	startUnix := start.Unix()
	endUnix := end.Unix()
	randomUnix := rand.Int63n(endUnix-startUnix) + startUnix
	return time.Unix(randomUnix, 0)
}

func randomStatus() string {
	statuses := []string{"Доставлен", "Отменен", "Завершен"}
	return statuses[rand.Intn(len(statuses))]
}
