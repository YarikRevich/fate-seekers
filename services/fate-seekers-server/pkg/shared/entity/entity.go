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
	Issuer     int64      `gorm:"column:issuer;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	UserEntity UserEntity `gorm:"foreignKey:Issuer;references:UserEntityID"`
}

// TableName retrieves name of database table.
func (*SessionEntity) TableName() string {
	return "sessions"
}

// AfterCreate performs sessions cache entity eviction after sessions entity create.
func (s *SessionEntity) AfterCreate(tx *gorm.DB) error {
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

// AfterCreate performs message cache entity eviction after messages entity creation.
func (m *MessageEntity) AfterCreate(tx *gorm.DB) error {
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
