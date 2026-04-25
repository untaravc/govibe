package models

import "time"

type Province struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Code      string    `json:"code" gorm:"type:varchar(16);not null;uniqueIndex"`
	Name      string    `json:"name" gorm:"type:varchar(191);not null;index"`
	Cities    []City    `json:"cities,omitempty" gorm:"foreignKey:ProvinceID"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Province) TableName() string {
	return "provinces"
}

type City struct {
	ID         uint64     `json:"id" gorm:"primaryKey;autoIncrement:false"`
	ProvinceID uint64     `json:"province_id" gorm:"column:province_id;not null;index"`
	Province   *Province  `json:"province,omitempty" gorm:"foreignKey:ProvinceID"`
	Code       string     `json:"code" gorm:"type:varchar(16);not null;uniqueIndex"`
	Name       string     `json:"name" gorm:"type:varchar(191);not null;index"`
	Districts  []District `json:"districts,omitempty" gorm:"foreignKey:CityID"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (City) TableName() string {
	return "cities"
}

type District struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CityID    uint64    `json:"city_id" gorm:"column:city_id;not null;index"`
	City      *City     `json:"city,omitempty" gorm:"foreignKey:CityID"`
	Code      string    `json:"code" gorm:"type:varchar(16);not null;uniqueIndex"`
	Name      string    `json:"name" gorm:"type:varchar(191);not null;index"`
	Villages  []Village `json:"villages,omitempty" gorm:"foreignKey:DistrictID"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (District) TableName() string {
	return "districts"
}

type Village struct {
	ID         uint64    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	DistrictID uint64    `json:"district_id" gorm:"column:district_id;not null;index"`
	District   *District `json:"district,omitempty" gorm:"foreignKey:DistrictID"`
	Code       string    `json:"code" gorm:"type:varchar(16);not null;uniqueIndex"`
	Name       string    `json:"name" gorm:"type:varchar(191);not null;index"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Village) TableName() string {
	return "villages"
}
