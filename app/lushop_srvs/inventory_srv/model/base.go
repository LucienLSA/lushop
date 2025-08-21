package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type GoodsDetail struct {
	Goods int32
	Num   int32
}

type GormDetailList []GoodsDetail

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (g GormDetailList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

type BaseModel struct {
	ID        int32          `gorm:"primarykey;type:int" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time" json:"-"`
	IsDeleted bool           `json:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
