package main

import (
	"context"
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
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var tableUsers = "users"

func main() {
	logger := newLogger()
	// MassTransport
	connAMQP, err := amqp.Dial(os.Getenv("RABBITMQ_DSN"))
	if err != nil {
		log.Fatalf("Ошибка подключения к RabbitMQ: %s", err)
	}
	defer connAMQP.Close()

	// Создайте канал
	ch, err := connAMQP.Channel()
	if err != nil {
		log.Fatalf("Ошибка создания канала: %s", err)
	}
	defer ch.Close()

	// Объявите очередь
	queueName := "clickhouse"
	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Ошибка объявления очереди: %s", err)
	}

	conn, err := connect()
	if err != nil {
		panic(err)
	}

	scheduler := gocron.NewScheduler(time.UTC)
	ctx := context.Background()
	scheduler.Every(5).Minutes().Do(func() {
		logger.Println("start generate events")
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableUsers)

		var totalUsers uint64
		err = conn.QueryRow(ctx, query).Scan(&totalUsers)
		if err != nil {
			logger.Fatalf("Ошибка выполнения запроса: %v", err)
		}

		randomThirdCount := totalUsers / 2000
		query = fmt.Sprintf(`
		SELECT id, years_old, gender_code
		FROM users
		ORDER BY rand()
		LIMIT %d
	`, randomThirdCount)

		rows, err := conn.Query(ctx, query)
		if err != nil {
			log.Fatalf("Ошибка выполнения запроса для выборки пользователей: %v", err)
		}
		defer rows.Close()

		var events []domain.Event
		for rows.Next() {
			var (
				user  domain.User
				ID    uint64
				years uint8
			)
			if err := rows.Scan(&ID, &years, &user.GenderCode); err != nil {
				log.Fatalf("Ошибка чтения строки: %v", err)
			}
			user.ID = int(ID)
			user.YearsOld = int(years)
			events = append(events, generator.GenerateEvent(user))
		}

		body, err := json.Marshal(events)
		if err != nil {
			logger.Fatalf("failed to marshal events with error %v", err)
		}
		err = ch.Publish(
			"clickhouse_bridge",
			queueName+".events",
			false,
			false,
			amqp.Publishing{
				Headers: amqp.Table{
					"method": "AddEvents",
				},
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			logger.Printf("failed to send events to rabbit %s", err)
		}
	})

	scheduler.Every(10).Minutes().Do(func() {
		logger.Println("start generate users")

		newUsers, err := generator.GetUsers()
		if err != nil {
			logger.Fatalf("failed to get users %v", err)
		}
		body, err := json.Marshal(newUsers)
		if err != nil {
			logger.Fatalf("failed to marshal events with error %v", err)
		}

		err = ch.Publish(
			"clickhouse_bridge",
			queueName+".users",
			false,
			false,
			amqp.Publishing{
				Headers: amqp.Table{
					"method": "AddUsers",
				},
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			logger.Printf("failed to send users to rabbit %s", err)
		}
	})
	logger.Print("publisher started")
	scheduler.StartBlocking()
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
			Addr: []string{"clickhouse:9000"},
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
