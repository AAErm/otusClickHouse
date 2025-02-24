package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"github.com/AAErm/otusClickHouse/project/migrator/generator"
)

const totalRecords = 1_500_000

func main() {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var wg sync.WaitGroup
	numGoroutines := 10

	recordsPerGoroutine := totalRecords / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < recordsPerGoroutine/100; j++ {
				randomUsers, err := generator.GetUsers()
				if err != nil {
					log.Println("Error fetching users:", err)
					continue
				}

				for _, randomUser := range randomUsers {
					_, err := db.Exec("INSERT INTO users (first_name, father_name, last_name, gender_code, bank, years_old) VALUES (?, ?, ?, ?, ?, ?)",
						randomUser.FirstName, randomUser.FatherName, randomUser.LastName, randomUser.GenderCode, randomUser.Bank, randomUser.YearsOld)
					if err != nil {
						log.Println("Error inserting user:", err)
					}

					query := "INSERT INTO events (price, service_id, location_id, timestamp) VALUES "
					events := generator.GenerateEventsYear(randomUser)
					values := make([]string, 0, len(events))
					for _, event := range events {
						value := fmt.Sprintf("(%d, %d, %d, '%s')",
							event.Price,
							event.ServiceID,
							event.LocationID,
							event.Timestamp.Format("2006-01-02 15:04:05"),
						)
						values = append(values, value)
					}
					query += strings.Join(values, ", ") + ";"
					_, err = db.Exec(query)
					if err != nil {
						log.Println("Error inserting events:", err)
					}
				}
			}
		}()
	}

	wg.Wait()
}
