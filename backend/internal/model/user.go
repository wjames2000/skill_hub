package model

import "time"

type User struct {
	ID           int64     `xorm:"pk autoincr 'id'" json:"id"`
	Username     string    `xorm:"varchar(100) not null unique 'username'" json:"username"`
	Email        string    `xorm:"varchar(255) not null unique 'email'" json:"email"`
	PasswordHash string    `xorm:"varchar(255) not null 'password_hash'" json:"-"`
	AvatarURL    string    `xorm:"varchar(512) 'avatar_url'" json:"avatar_url"`
	Bio          string    `xorm:"text 'bio'" json:"bio"`
	GitHubID     string    `xorm:"varchar(100) unique 'github_id'" json:"github_id,omitempty"`
	GitHubToken  string    `xorm:"varchar(512) 'github_token'" json:"-"`
	Role         string    `xorm:"varchar(20) default 'user' 'role'" json:"role"`
	Status       int       `xorm:"tinyint default 1 'status'" json:"status"`
	LastLoginAt  time.Time `xorm:"datetime 'last_login_at'" json:"last_login_at"`
	CreatedAt    time.Time `xorm:"created 'created_at'" json:"created_at"`
	UpdatedAt    time.Time `xorm:"updated 'updated_at'" json:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}

type APIKey struct {
	ID         int64      `xorm:"pk autoincr 'id'" json:"id"`
	UserID     int64      `xorm:"int not null 'user_id'" json:"user_id"`
	Name       string     `xorm:"varchar(100) not null 'name'" json:"name"`
	Key        string     `xorm:"varchar(64) not null unique 'key'" json:"key"`
	LastUsedAt time.Time  `xorm:"datetime 'last_used_at'" json:"last_used_at"`
	ExpiresAt  *time.Time `xorm:"datetime 'expires_at'" json:"expires_at"`
	IsRevoked  bool       `xorm:"bool default false 'is_revoked'" json:"is_revoked"`
	CreatedAt  time.Time  `xorm:"created 'created_at'" json:"created_at"`
	UpdatedAt  time.Time  `xorm:"updated 'updated_at'" json:"updated_at"`
}

func (k *APIKey) TableName() string {
	return "api_keys"
}

type Favorite struct {
	ID        int64     `xorm:"pk autoincr 'id'" json:"id"`
	UserID    int64     `xorm:"int not null 'user_id'" json:"user_id"`
	SkillID   int64     `xorm:"int not null 'skill_id'" json:"skill_id"`
	CreatedAt time.Time `xorm:"created 'created_at'" json:"created_at"`
}

func (f *Favorite) TableName() string {
	return "favorites"
}

type Review struct {
	ID        int64     `xorm:"pk autoincr 'id'" json:"id"`
	UserID    int64     `xorm:"int not null 'user_id'" json:"user_id"`
	SkillID   int64     `xorm:"int not null 'skill_id'" json:"skill_id"`
	Score     int       `xorm:"tinyint not null 'score'" json:"score"`
	Content   string    `xorm:"text 'content'" json:"content"`
	CreatedAt time.Time `xorm:"created 'created_at'" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated 'updated_at'" json:"updated_at"`
}

func (r *Review) TableName() string {
	return "reviews"
}
