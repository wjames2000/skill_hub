package model

import "time"

type Skill struct {
	ID            int64     `xorm:"pk autoincr 'id'" json:"id"`
	Name          string    `xorm:"varchar(255) not null 'name'" json:"name"`
	DisplayName   string    `xorm:"varchar(255) 'display_name'" json:"display_name"`
	Description   string    `xorm:"text 'description'" json:"description"`
	ZhDescription string    `xorm:"text 'zh_description'" json:"zh_description"`
	EnDescription string    `xorm:"text 'en_description'" json:"en_description"`
	Version       string    `xorm:"varchar(50) 'version'" json:"version"`
	Author        string    `xorm:"varchar(255) 'author'" json:"author"`
	Repository    string    `xorm:"varchar(512) not null 'repository'" json:"repository"`
	RepoOwner     string    `xorm:"varchar(255) 'repo_owner'" json:"repo_owner"`
	RepoName      string    `xorm:"varchar(255) 'repo_name'" json:"repo_name"`
	DefaultBranch string    `xorm:"varchar(255) 'default_branch'" json:"default_branch"`
	SkillPath     string    `xorm:"varchar(512) 'skill_path'" json:"skill_path"`
	SkillFileSHA  string    `xorm:"varchar(64) 'skill_file_sha'" json:"skill_file_sha"`
	AvatarURL     string    `xorm:"varchar(512) 'avatar_url'" json:"avatar_url"`
	Homepage      string    `xorm:"varchar(512) 'homepage'" json:"homepage"`
	License       string    `xorm:"varchar(100) 'license'" json:"license"`
	Stars         int       `xorm:"int 'stars'" json:"stars"`
	Forks         int       `xorm:"int 'forks'" json:"forks"`
	OpenIssues    int       `xorm:"int 'open_issues'" json:"open_issues"`
	Language      string    `xorm:"varchar(100) 'language'" json:"language"`
	Topics        []string  `xorm:"json 'topics'" json:"topics"`
	Category      string    `xorm:"varchar(100) 'category'" json:"category"`
	CategoryID    int64     `xorm:"bigint 'category_id'" json:"category_id"`
	Tags          []string  `xorm:"json 'tags'" json:"tags"`
	Readme        string    `xorm:"longtext 'readme'" json:"readme"`
	Installs      int64     `xorm:"int default 0 'installs'" json:"installs"`
	Score         float64   `xorm:"decimal(10,2) default 0 'score'" json:"score"`
	IsOfficial    bool      `xorm:"bool default false 'is_official'" json:"is_official"`
	IsArchived    bool      `xorm:"bool default false 'is_archived'" json:"is_archived"`
	ScanPassed    bool      `xorm:"bool default true 'scan_passed'" json:"scan_passed"`
	ScanReport    string    `xorm:"text 'scan_report'" json:"scan_report"`
	Status        int       `xorm:"tinyint default 1 'status'" json:"status"`
	LastSyncAt    time.Time `xorm:"updated 'last_sync_at'" json:"last_sync_at"`
	CreatedAt     time.Time `xorm:"created 'created_at'" json:"created_at"`
	UpdatedAt     time.Time `xorm:"updated 'updated_at'" json:"updated_at"`

	Extra map[string]interface{} `xorm:"json 'extra'" json:"extra,omitempty"`
}

func (s *Skill) TableName() string {
	return "skills"
}

type SkillVersion struct {
	ID        int64     `xorm:"pk autoincr 'id'" json:"id"`
	SkillID   int64     `xorm:"int not null 'skill_id'" json:"skill_id"`
	Version   string    `xorm:"varchar(50) 'version'" json:"version"`
	SkillSHA  string    `xorm:"varchar(64) 'skill_sha'" json:"skill_sha"`
	ChangeLog string    `xorm:"text 'change_log'" json:"change_log"`
	CreatedAt time.Time `xorm:"created 'created_at'" json:"created_at"`
}

func (s *SkillVersion) TableName() string {
	return "skill_versions"
}

type SkillCategory struct {
	ID          int64     `xorm:"pk autoincr 'id'" json:"id"`
	Name        string    `xorm:"varchar(100) not null 'name'" json:"name"`
	ZhName      string    `xorm:"varchar(100) 'zh_name'" json:"zh_name"`
	EnName      string    `xorm:"varchar(100) 'en_name'" json:"en_name"`
	Slug        string    `xorm:"varchar(100) not null unique 'slug'" json:"slug"`
	Description string    `xorm:"text 'description'" json:"description"`
	Icon        string    `xorm:"varchar(255) 'icon'" json:"icon"`
	ParentID    int64     `xorm:"int default 0 'parent_id'" json:"parent_id"`
	SortOrder   int       `xorm:"int default 0 'sort_order'" json:"sort_order"`
	SkillCount  int       `xorm:"int default 0 'skill_count'" json:"skill_count"`
	CreatedAt   time.Time `xorm:"created 'created_at'" json:"created_at"`
	UpdatedAt   time.Time `xorm:"updated 'updated_at'" json:"updated_at"`
}

func (s *SkillCategory) TableName() string {
	return "skill_categories"
}

const (
	SkillStatusPending    = 0
	SkillStatusActive     = 1
	SkillStatusDisabled   = 2
	SkillStatusDeprecated = 3
)
