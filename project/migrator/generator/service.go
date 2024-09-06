package generator

import (
	"encoding/json"
	"os"
	"regexp"

	"github.com/AAErm/otusClickHouse/project/domain"
)

// если ServiceID < 100, это "общественный транспорт"
// если  100 > ServiceID > 200, то  "аптека"
// если  200 > ServiceID > 300, то  "продуктовый магазин"
// если  300 > ServiceID > 550, то  "общепит"
// если 550  > ServiceID > 650, то "развлечения"
// если ServiceID > 650 , то "автомобильные услуги"
func GetServices() ([]domain.Service, error) {
	re := regexp.MustCompile(`^services_\d{3}\.json$`)

	files, err := os.ReadDir("./data/")
	if err != nil {
		return nil, err
	}

	services := []domain.Service{}
	for _, file := range files {
		if file.IsDir() || !re.MatchString(file.Name()) {
			continue
		}

		data, err := os.ReadFile(file.Name())
		if err != nil {
			return nil, err
		}

		var tmpservices []domain.Service
		err = json.Unmarshal(data, &tmpservices)
		if err != nil {
			return nil, err
		}
		services = append(services, tmpservices...)

	}

	return services, nil
}
