package reusablecode

import "time"

type ReusableCode struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string    `gorm:"column:code" json:"code"`
	IsActive  bool      `gorm:"column:is_active" json:"is_active"`
	MaxUse    *int      `gorm:"column:max_use" json:"max_use"`
	Count     *int      `gorm:"column:count" json:"count"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (u *ReusableCode) TableName() string {
	return "reusable_code"
}
