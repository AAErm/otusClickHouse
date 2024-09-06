package usergenerator

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/AAErm/otusClickHouse/project/domain"
)

const (
	apiURL = "https://api.randomdatatools.ru/?count=100"
)

var banks = []string{"Беларусбанк", "Белагропромбанк", "БПС-Сбербанк", "Приорбанк", "Банк ВТБ (Беларусь)"}
var bankDistribution = []int{40, 20, 15, 15, 10}

func GetUsers() ([]domain.User, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var randomUsers []domain.User

	if err := json.NewDecoder(resp.Body).Decode(&randomUsers); err != nil {
		return nil, err
	}

	for _, v := range randomUsers {
		bank := assignBank()
		v.Bank = bank
	}

	return randomUsers, nil
}

func assignBank() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	n := rand.Intn(100)
	cumulative := 0

	for i, percent := range bankDistribution {
		cumulative += percent
		if n < cumulative {
			return banks[i]
		}
	}

	return banks[len(banks)-1]
}
