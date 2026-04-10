package model

import (
	"time"

	"gorm.io/gorm"
)

type Topic struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"size:500;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	IsActive    bool           `gorm:"default:true" json:"isActive"`
	IsTemplate  bool           `gorm:"default:false" json:"isTemplate"`
}

func (Topic) TableName() string {
	return "topics"
}

type Option struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Title      string         `gorm:"size:500;not null" json:"title"`
	TopicID    uint           `gorm:"index" json:"topicId"`
	Importance int            `gorm:"default:0" json:"importance"`
	Necessity  int            `gorm:"default:0" json:"necessity"`
	SortOrder  int            `gorm:"default:0" json:"sortOrder"`
	IsActive   bool           `gorm:"default:true" json:"isActive"`
}

func (Option) TableName() string {
	return "options"
}

type Round struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"createdAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	TopicID          uint           `gorm:"index" json:"topicId"`
	RoundNumber      int            `gorm:"not null" json:"roundNumber"`
	ImportanceStatus string         `gorm:"size:20;default:'pending'" json:"importanceStatus"`
	NecessityStatus  string         `gorm:"size:20;default:'pending'" json:"necessityStatus"`
	Status           string         `gorm:"size:20;default:'pending'" json:"status"`
	Results          string         `gorm:"type:text" json:"results"`
	MatchScore       int            `gorm:"default:0" json:"matchScore"`
	IsActive         bool           `gorm:"default:true" json:"isActive"`
	OptionStrategy   string         `gorm:"size:20;default:'reuse'" json:"optionStrategy"`
	Options          string         `gorm:"type:text" json:"options"`
}

func (Round) TableName() string {
	return "rounds"
}
