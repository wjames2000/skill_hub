package model

import "time"

type RouterLog struct {
	ID               int64     `xorm:"pk autoincr 'id'" json:"id"`
	SessionID        string    `xorm:"varchar(64) index 'session_id'" json:"session_id"`
	Query            string    `xorm:"text 'query'" json:"query"`
	QueryEmbedding   []float32 `xorm:"json 'query_embedding'" json:"query_embedding,omitempty"`
	MatchedSkillID   int64     `xorm:"int 'matched_skill_id'" json:"matched_skill_id"`
	MatchedSkillName string    `xorm:"varchar(255) 'matched_skill_name'" json:"matched_skill_name"`
	MatchScore       float64   `xorm:"decimal(10,4) 'match_score'" json:"match_score"`
	MatchStrategy    string    `xorm:"varchar(50) 'match_strategy'" json:"match_strategy"`
	IsExecuted       bool      `xorm:"bool default false 'is_executed'" json:"is_executed"`
	ExecuteResult    string    `xorm:"text 'execute_result'" json:"execute_result,omitempty"`
	ExecuteDuration  int       `xorm:"int 'execute_duration'" json:"execute_duration"`
	FeedbackScore    int       `xorm:"tinyint 'feedback_score'" json:"feedback_score"`
	FeedbackComment  string    `xorm:"text 'feedback_comment'" json:"feedback_comment,omitempty"`
	UserID           int64     `xorm:"int 'user_id'" json:"user_id"`
	ClientIP         string    `xorm:"varchar(45) 'client_ip'" json:"client_ip"`
	CreatedAt        time.Time `xorm:"created 'created_at'" json:"created_at"`
}

func (r *RouterLog) TableName() string {
	return "router_logs"
}
