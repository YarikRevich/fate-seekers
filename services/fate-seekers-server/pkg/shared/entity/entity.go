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
	Started    bool       `gorm:"column:started;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	UserEntity UserEntity `gorm:"foreignKey:Issuer;references:ID"`
}

// TableName retrieves name of database table.
func (*SessionEntity) TableName() string {
	return "sessions"
}

// TableView retrieves name of database table view.
func (*SessionEntity) TableView() string {
	return "SessionEntity"
}

// BeforeCreate performs sessions cache entity eviction before sessions entity create.
func (s *SessionEntity) BeforeCreate(tx *gorm.DB) error {
	if err := tx.
		Model(&UserEntity{}).
		Where("id = ?", s.Issuer).
		First(&s.UserEntity).Error; err != nil {
		return err
	}

	cache.
		GetInstance().
		BeginSessionsTransaction()

	cache.
		GetInstance().
		EvictSessions(s.ID)

	cache.
		GetInstance().
		CommitSessionsTransaction()

	cache.
		GetInstance().
		BeginUserSessionsTransaction()

	cache.
		GetInstance().
		EvictUserSessions(s.UserEntity.Name)

	cache.
		GetInstance().
		CommitUserSessionsTransaction()

	return nil
}

// BeforeUpdate performs sessions cache entity eviction before sessions entity update.
func (s *SessionEntity) BeforeUpdate(tx *gorm.DB) error {
	if err := tx.
		Model(&UserEntity{}).
		Where("id = ?", s.Issuer).
		First(&s.UserEntity).Error; err != nil {
		return err
	}

	cache.
		GetInstance().
		BeginSessionsTransaction()

	cache.
		GetInstance().
		EvictSessions(s.ID)

	cache.
		GetInstance().
		CommitSessionsTransaction()

	cache.
		GetInstance().
		BeginUserSessionsTransaction()

	cache.
		GetInstance().
		EvictUserSessions(s.UserEntity.Name)

	cache.
		GetInstance().
		CommitUserSessionsTransaction()

	return nil
}

// LobbyEntity represents lobbies entity.
type LobbyEntity struct {
	ID            int64         `gorm:"column:id;primaryKey;auto_increment;not null"`
	UserID        int64         `gorm:"column:user_id;not null"`
	SessionID     int64         `gorm:"column:session_id;not null"`
	Skin          int64         `gorm:"column:skin;not null"`
	Health        int64         `gorm:"column:health;not null"`
	Active        bool          `gorm:"column:active;not null"`
	Host          bool          `gorm:"column:host;not null"`
	Eliminated    bool          `gorm:"column:eliminated;not null"`
	PositionX     float64       `gorm:"column:position_x;not null"`
	PositionY     float64       `gorm:"column:position_y;not null"`
	CreatedAt     time.Time     `gorm:"column:created_at;autoCreateTime"`
	UserEntity    UserEntity    `gorm:"foreignKey:UserID;references:ID"`
	SessionEntity SessionEntity `gorm:"foreignKey:SessionID;references:ID"`
}

// TableName retrieves name of database table.
func (*LobbyEntity) TableName() string {
	return "lobbies"
}

// BeforeCreate performs lobbies cache entity eviction before lobbies entity create.
func (l *LobbyEntity) BeforeCreate(tx *gorm.DB) error {
	cache.
		GetInstance().
		BeginLobbySetTransaction()

	cache.
		GetInstance().
		EvictLobbySet(l.SessionID)

	cache.
		GetInstance().
		CommitLobbySetTransaction()

	return nil
}

// BeforeUpdate performs lobbies cache entity eviction before lobbies entity update.
func (l *LobbyEntity) BeforeUpdate(tx *gorm.DB) error {
	cache.
		GetInstance().
		BeginLobbySetTransaction()

	cache.
		GetInstance().
		EvictLobbySet(l.SessionID)

	cache.
		GetInstance().
		CommitLobbySetTransaction()

	return nil
}

// MessageEntity represents messages entity.
type MessageEntity struct {
	ID         int64      `gorm:"column:id;primaryKey;auto_increment;not null"`
	Content    string     `gorm:"column:name;not null"`
	Issuer     int64      `gorm:"column:issuer;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	UserEntity UserEntity `gorm:"foreignKey:Issuer;references:ID"`
}

// TableName retrieves name of database table.
func (*MessageEntity) TableName() string {
	return "messages"
}

// BeforeCreate performs message cache entity eviction before messages entity creation.
func (m *MessageEntity) BeforeCreate(tx *gorm.DB) error {
	if err := tx.
		Model(&UserEntity{}).
		Where("id = ?", m.Issuer).
		First(&m.UserEntity).Error; err != nil {
		return err
	}

	// cache.GetInstance().EvictUserSessions()

	// cache.
	// 	GetInstance().
	// 	Evict(m.UserEntity.Name)

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

// TableView retrieves name of database table view.
func (*UserEntity) TableView() string {
	return "UserEntity"
}

// AfterFind performs user cache entity creation after user entity creation.
func (u *UserEntity) AfterFind(tx *gorm.DB) error {
	cache.
		GetInstance().
		AddUser(u.Name, u.ID)

	return nil
}
