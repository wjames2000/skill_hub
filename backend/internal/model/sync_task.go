package model

import "time"

type SyncTask struct {
	ID            int64      `xorm:"pk autoincr 'id'" json:"id"`
	Type          string     `xorm:"varchar(20) not null 'type'" json:"type"`
	Strategy      string     `xorm:"varchar(30) 'strategy'" json:"strategy"`
	Status        string     `xorm:"varchar(20) not null default 'pending' 'status'" json:"status"`
	TotalRepos    int        `xorm:"int default 0 'total_repos'" json:"total_repos"`
	FoundRepos    int        `xorm:"int default 0 'found_repos'" json:"found_repos"`
	ParsedSkills  int        `xorm:"int default 0 'parsed_skills'" json:"parsed_skills"`
	NewSkills     int        `xorm:"int default 0 'new_skills'" json:"new_skills"`
	UpdatedSkills int        `xorm:"int default 0 'updated_skills'" json:"updated_skills"`
	FailedSkills  int        `xorm:"int default 0 'failed_skills'" json:"failed_skills"`
	ScannedRepos  int        `xorm:"int default 0 'scanned_repos'" json:"scanned_repos"`
	ErrorCount    int        `xorm:"int default 0 'error_count'" json:"error_count"`
	ErrorMessage  string     `xorm:"text 'error_message'" json:"error_message,omitempty"`
	StartedAt     *time.Time `xorm:"datetime 'started_at'" json:"started_at,omitempty"`
	FinishedAt    *time.Time `xorm:"datetime 'finished_at'" json:"finished_at,omitempty"`
	CreatedAt     time.Time  `xorm:"created 'created_at'" json:"created_at"`
	UpdatedAt     time.Time  `xorm:"updated 'updated_at'" json:"updated_at"`
}

func (st *SyncTask) TableName() string {
	return "sync_tasks"
}

const (
	SyncTypeFull        = "full"
	SyncTypeIncremental = "incremental"

	SyncStatusPending   = "pending"
	SyncStatusRunning   = "running"
	SyncStatusCompleted = "completed"
	SyncStatusFailed    = "failed"
	SyncStatusCancelled = "cancelled"

	StrategyTopic   = "topic"
	StrategyPath    = "path"
	StrategyAwesome = "awesome"
)
