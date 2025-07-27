package entity

import (
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/cache"
	"gorm.io/gorm"
)

// SessionEntity represents sessions entity.
type SessionEntity struct {
	ID         int64      `gorm:"column:id;primaryKey;auto_increment;not null"`
	Name       string     `gorm:"column:name;not null;unique"`
	Seed       int64      `gorm:"column:seed;not null"`
	Issuer     int64      `gorm:"column:issuer;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	UserEntity UserEntity `gorm:"foreignKey:Issuer;references:UserEntityID"`
}

// TableName retrieves name of database table.
func (*SessionEntity) TableName() string {
	return "sessions"
}

// BeforeCreate performs sessions cache entity eviction before sessions entity create.
func (s *SessionEntity) BeforeCreate(tx *gorm.DB) error {
	cache.
		GetInstance().
		EvictSessions(s.UserEntity.Name)

	return nil
}

// AfterDelete performs sessions cache entity eviction after sessions entity removal.
func (s *SessionEntity) AfterDelete(tx *gorm.DB) error {
	cache.
		GetInstance().
		EvictSessions(s.UserEntity.Name)

	return nil
}

// LobbyEntity represents lobbies entity.
type LobbyEntity struct {
	ID            int64         `gorm:"column:id;primaryKey;auto_increment;not null"`
	UserID        int64         `gorm:"column:user_id;not null;unique"`
	SessionID     int64         `gorm:"column:session_id;not null"`
	Skin          int64         `gorm:"column:skin;not null"`
	Health        int64         `gorm:"column:health;not null"`
	Eliminated    bool          `gorm:"column:eliminated;not null"`
	Position      float64       `gorm:"column:position;not null"`
	CreatedAt     time.Time     `gorm:"column:created_at;autoCreateTime"`
	UserEntity    UserEntity    `gorm:"foreignKey:UserID;references:UserEntityID"`
	SessionEntity SessionEntity `gorm:"foreignKey:SessionID;references:SessionEntityID"`
}

// TableName retrieves name of database table.
func (*LobbyEntity) TableName() string {
	return "lobbies"
}

// BeforeCreate performs lobbies cache entity eviction before lobbies entity create.
func (l *LobbyEntity) BeforeCreate(tx *gorm.DB) error {
	cache.
		GetInstance().
		EvictLobbySet(l.SessionID)

	return nil
}

// AfterDelete performs lobbies cache entity eviction after lobbies entity removal.
func (l *LobbyEntity) AfterDelete(tx *gorm.DB) error {
	cache.
		GetInstance().
		EvictLobbySet(l.SessionID)

	cache.
		GetInstance().
		EvictMetadata(l.UserEntity.Name)

	return nil
}

// MessageEntity represents messages entity.
type MessageEntity struct {
	ID         int64      `gorm:"column:id;primaryKey;auto_increment;not null"`
	Content    string     `gorm:"column:name;not null"`
	Issuer     int64      `gorm:"column:issuer;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	UserEntity UserEntity `gorm:"foreignKey:Issuer;references:UserEntityID"`
}

// TableName retrieves name of database table.
func (*MessageEntity) TableName() string {
	return "messages"
}

// BeforeCreate performs message cache entity eviction before messages entity creation.
func (m *MessageEntity) BeforeCreate(tx *gorm.DB) error {
	cache.
		GetInstance().
		EvictSessions(m.UserEntity.Name)

	return nil
}

// UserEntity represents users entity.
type UserEntity struct {
	ID        int64     `gorm:"column:id;primaryKey;auto_increment;not null"`
	Name      string    `gorm:"column:name;not null;unique"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

// TableName retrieves name of database table.
func (*UserEntity) TableName() string {
	return "users"
}

// AfterFind performs user cache entity creation after user entity creation.
func (u *UserEntity) AfterFind(tx *gorm.DB) error {
	cache.
		GetInstance().
		AddUser(u.Name, u.ID)

	return nil
}
