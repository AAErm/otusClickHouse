package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/AAErm/otusClickHouse/project/domain"
	"github.com/AAErm/otusClickHouse/project/migrator/generator"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/imega/mt"
	"github.com/sirupsen/logrus"
)

var tableUsers = "oplati_users"

func main() {
	logger := newLogger()
	// MassTransport
	conf, err := mt.ParseConfig([]byte(os.Getenv("MT_CONFIG")))
	if err != nil {
		logger.Fatalf("failed to parse masstransport config, %v", err)
	}

	massT := mt.NewMT(
		mt.WithAMQP(os.Getenv("RABBITMQ_DSN")),
		mt.WithLogger(logger),
		mt.WithConfig(conf),
	)
	if err != nil {
		panic((err))
	}

	conn, err := connect()
	if err != nil {
		panic(err)
	}

	scheduler := gocron.NewScheduler(time.UTC)
	ctx := context.Background()
	scheduler.Every(5).Minutes().Do(func() {
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableUsers)

		var totalUsers int
		err = conn.QueryRow(ctx, query).Scan(&totalUsers)
		if err != nil {
			logger.Fatalf("Ошибка выполнения запроса: %v", err)
		}

		randomThirdCount := totalUsers / 2000
		query = fmt.Sprintf(`  
		SELECT id, YearsOld, GenderCode  
		FROM oplatiUsers  
		ORDER BY rand()  
		LIMIT %d  
	`, randomThirdCount)

		rows, err := conn.Query(ctx, query)
		if err != nil {
			log.Fatalf("Ошибка выполнения запроса для выборки пользователей: %v", err)
		}
		defer rows.Close()

		// Обработка выбранных пользователей
		var events []domain.Event
		for rows.Next() {
			var user domain.User
			if err := rows.Scan(&user.ID, &user.YearsOld, &user.GenderCode); err != nil {
				log.Fatalf("Ошибка чтения строки: %v", err)
			}
			events = append(events, generator.GenerateEvent(user))
		}

		body, err := json.Marshal(events)
		if err != nil {
			logger.Fatalf("failed to marshal events with error %v", err)
		}
		massT.Cast("clickhouse", mt.Request{
			Header: mt.Header{
				"method": "addEvents",
			},
			Body: body,
		})
	})

	scheduler.Every(10).Minutes().Do(func() {
		newUsers, err := generator.GetUsers()
		if err != nil {
			logger.Fatalf("failed to get users %v", err)
		}
		body, err := json.Marshal(newUsers)
		if err != nil {
			logger.Fatalf("failed to marshal events with error %v", err)
		}

		massT.Cast("clickhouse", mt.Request{
			Header: mt.Header{
				"method": "addUsers",
			},
			Body: body,
		})
	})
}

func newLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
	})

	return logger
}

func connect() (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{"clickhouse:9440"},
			Auth: clickhouse.Auth{
				Database: os.Getenv("CLICKHOUSE_DB"),
				Username: os.Getenv("CLICKHOUSE_USER"),
				Password: os.Getenv("CLICKHOUSE_PASSWORD"),
			},
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "an-example-go-client", Version: "0.1"},
				},
			},

			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
			TLS: &tls.Config{
				InsecureSkipVerify: true,
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}
