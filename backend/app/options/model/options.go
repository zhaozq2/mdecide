package model

import (
	"time"

	"gorm.io/gorm"
)

type Option struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Title      string         `gorm:"size:500;not null" json:"title"`
	TopicID    uint           `gorm:"index" json:"topicId"`
	RoundID    uint           `gorm:"index" json:"roundId"`
	Importance int            `gorm:"default:0" json:"importance"`
	Necessity  int            `gorm:"default:0" json:"necessity"`
	SortOrder  int            `gorm:"default:0" json:"sortOrder"`
	IsActive   bool           `gorm:"default:true" json:"isActive"`
}

func (Option) TableName() string {
	return "options"
}
