package seeder

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"govibe/app/Models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const defaultRegionAssetsDir = "resources/assets"

func SeedRegions(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	baseDir := strings.TrimSpace(os.Getenv("REGION_ASSETS_DIR"))
	if baseDir == "" {
		baseDir = defaultRegionAssetsDir
	}

	provinces, err := loadProvinces(filepath.Join(baseDir, "provinsi", "provinsi.json"))
	if err != nil {
		return err
	}

	cities, err := loadCities(filepath.Join(baseDir, "kabupaten_kota"))
	if err != nil {
		return err
	}

	districts, err := loadDistricts(filepath.Join(baseDir, "kecamatan"))
	if err != nil {
		return err
	}

	villages, err := loadVillages(filepath.Join(baseDir, "kelurahan_desa"))
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := upsertProvinces(tx, provinces); err != nil {
			return err
		}
		if err := upsertCities(tx, cities); err != nil {
			return err
		}
		if err := upsertDistricts(tx, districts); err != nil {
			return err
		}
		if err := upsertVillages(tx, villages); err != nil {
			return err
		}
		return nil
	})
}

func loadProvinces(path string) ([]models.Province, error) {
	items, err := readRegionJSON(path)
	if err != nil {
		return nil, err
	}

	out := make([]models.Province, 0, len(items))
	for _, item := range items {
		id, err := parseRegionID(item.Code)
		if err != nil {
			return nil, fmt.Errorf("province %s: %w", item.Code, err)
		}

		out = append(out, models.Province{
			ID:   id,
			Code: item.Code,
			Name: item.Name,
		})
	}
	return out, nil
}

func loadCities(dir string) ([]models.City, error) {
	files, err := sortedGlob(filepath.Join(dir, "kab-*.json"))
	if err != nil {
		return nil, err
	}

	var out []models.City
	for _, file := range files {
		provinceCode := strings.TrimSuffix(strings.TrimPrefix(filepath.Base(file), "kab-"), ".json")
		provinceID, err := parseRegionID(provinceCode)
		if err != nil {
			return nil, fmt.Errorf("%s province code: %w", file, err)
		}

		items, err := readRegionJSON(file)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			fullCode := provinceCode + item.Code
			id, err := parseRegionID(fullCode)
			if err != nil {
				return nil, fmt.Errorf("city %s: %w", fullCode, err)
			}

			out = append(out, models.City{
				ID:         id,
				ProvinceID: provinceID,
				Code:       fullCode,
				Name:       item.Name,
			})
		}
	}
	return out, nil
}

func loadDistricts(dir string) ([]models.District, error) {
	files, err := sortedGlob(filepath.Join(dir, "kec-*.json"))
	if err != nil {
		return nil, err
	}

	var out []models.District
	for _, file := range files {
		parts := strings.Split(strings.TrimSuffix(strings.TrimPrefix(filepath.Base(file), "kec-"), ".json"), "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("%s: invalid kecamatan filename", file)
		}

		cityCode := parts[0] + parts[1]
		cityID, err := parseRegionID(cityCode)
		if err != nil {
			return nil, fmt.Errorf("%s city code: %w", file, err)
		}

		items, err := readRegionJSON(file)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			fullCode := cityCode + item.Code
			id, err := parseRegionID(fullCode)
			if err != nil {
				return nil, fmt.Errorf("district %s: %w", fullCode, err)
			}

			out = append(out, models.District{
				ID:     id,
				CityID: cityID,
				Code:   fullCode,
				Name:   item.Name,
			})
		}
	}
	return out, nil
}

func loadVillages(dir string) ([]models.Village, error) {
	files, err := sortedGlob(filepath.Join(dir, "keldesa-*.json"))
	if err != nil {
		return nil, err
	}

	var out []models.Village
	for _, file := range files {
		parts := strings.Split(strings.TrimSuffix(strings.TrimPrefix(filepath.Base(file), "keldesa-"), ".json"), "-")
		if len(parts) != 3 {
			return nil, fmt.Errorf("%s: invalid kelurahan/desa filename", file)
		}

		districtCode := parts[0] + parts[1] + parts[2]
		districtID, err := parseRegionID(districtCode)
		if err != nil {
			return nil, fmt.Errorf("%s district code: %w", file, err)
		}

		items, err := readRegionJSON(file)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			fullCode := districtCode + item.Code
			id, err := parseRegionID(fullCode)
			if err != nil {
				return nil, fmt.Errorf("village %s: %w", fullCode, err)
			}

			out = append(out, models.Village{
				ID:         id,
				DistrictID: districtID,
				Code:       fullCode,
				Name:       item.Name,
			})
		}
	}
	return out, nil
}

type regionJSONItem struct {
	Code string
	Name string
}

func readRegionJSON(path string) ([]regionJSONItem, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var data map[string]string
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	items := make([]regionJSONItem, 0, len(keys))
	for _, key := range keys {
		name := strings.TrimSpace(data[key])
		if name == "" {
			return nil, fmt.Errorf("%s: empty name for code %s", path, key)
		}

		items = append(items, regionJSONItem{
			Code: strings.TrimSpace(key),
			Name: name,
		})
	}
	return items, nil
}

func sortedGlob(pattern string) ([]string, error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func parseRegionID(code string) (uint64, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return 0, errors.New("empty code")
	}

	n, err := strconv.ParseUint(code, 10, 64)
	if err != nil || n == 0 {
		return 0, fmt.Errorf("invalid code %q", code)
	}
	return n, nil
}

func upsertProvinces(db *gorm.DB, rows []models.Province) error {
	if len(rows) == 0 {
		return nil
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"code", "name"}),
	}).CreateInBatches(rows, 1000).Error
}

func upsertCities(db *gorm.DB, rows []models.City) error {
	if len(rows) == 0 {
		return nil
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"province_id", "code", "name"}),
	}).CreateInBatches(rows, 1000).Error
}

func upsertDistricts(db *gorm.DB, rows []models.District) error {
	if len(rows) == 0 {
		return nil
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"city_id", "code", "name"}),
	}).CreateInBatches(rows, 1000).Error
}

func upsertVillages(db *gorm.DB, rows []models.Village) error {
	if len(rows) == 0 {
		return nil
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"district_id", "code", "name"}),
	}).CreateInBatches(rows, 1000).Error
}
