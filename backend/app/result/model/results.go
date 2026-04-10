package model

import (
	"time"

	"gorm.io/gorm"
)

// Result 筛选结果
type Result struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"createdAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	OptionID        uint           `gorm:"index;not null" json:"optionId"`   // 关联选项ID
	RoundID         int            `gorm:"not null" json:"roundId"`          // 关联轮次ID
	TopicID         int            `gorm:"index" json:"topicId"`             // 关联主题ID
	ImportanceScore int            `gorm:"default:0" json:"importanceScore"` // 重要性得分 (1-5)
	NecessityScore  int            `gorm:"default:0" json:"necessityScore"`  // 必要性得分 (0或1)
	TotalScore      int            `gorm:"default:0" json:"totalScore"`      // 总分
	Summary         string         `gorm:"type:text" json:"summary"`         // 结果总结描述
	Rank            int            `gorm:"default:0" json:"rank"`            // 排名
	IsWinner        bool           `gorm:"default:false" json:"isWinner"`    // 是否获胜选项
}

func (Result) TableName() string {
	return "results"
}
