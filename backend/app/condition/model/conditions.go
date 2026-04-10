package model

import (
	"time"

	"gorm.io/gorm"
)

// Condition 条件/筛选条件
type Condition struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"size:200;not null" json:"name"` // 条件名称
	Phase     string         `gorm:"size:20;not null" json:"phase"` // 阶段: importance(重要性), necessity(必要性)
	Score     int            `gorm:"default:0" json:"score"`        // 分数: 重要性1-5分, 必要性0或1
	SortOrder int            `gorm:"default:0" json:"sortOrder"`    // 排序
	OptionID  uint           `gorm:"index" json:"optionId"`         // 关联选项ID
	RoundID   uint           `gorm:"index" json:"roundId"`          // 关联轮次ID
	TopicID   uint           `gorm:"index" json:"topicId"`          // 关联主题ID
}

func (Condition) TableName() string {
	return "conditions"
}
