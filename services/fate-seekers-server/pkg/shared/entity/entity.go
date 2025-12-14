package entity

import (
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
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
		EvictSessionsByName(s.Name)

	cache.
		GetInstance().
		EvictUserSessions(s.UserEntity.Name)

	cache.
		GetInstance().
		EvictLobbySet(s.ID)

	return nil
}

// GenerationsEntity represents generations entity.
type GenerationsEntity struct {
	ID            int64         `gorm:"column:id;primaryKey;auto_increment;not null"`
	SessionID     int64         `gorm:"column:session_id;not null"`
	Instance      string        `gorm:"column:instance;not null"`
	Name          string        `gorm:"column:name;not null"`
	Type          string        `gorm:"column:type;not null"`
	Active        bool          `gorm:"column:active;not null"`
	PositionX     float64       `gorm:"column:position_x;not null"`
	PositionY     float64       `gorm:"column:position_y;not null"`
	CreatedAt     time.Time     `gorm:"column:created_at;autoCreateTime"`
	SessionEntity SessionEntity `gorm:"foreignKey:SessionID;references:ID"`
}

// TableName retrieves name of database table.
func (*GenerationsEntity) TableName() string {
	return "generations"
}

// TableView retrieves name of database table view.
func (*GenerationsEntity) TableView() string {
	return "GenerationsEntity"
}

// BeforeCreate performs generations cache entity eviction before generations entity create.
func (g *GenerationsEntity) BeforeCreate(tx *gorm.DB) error {
	if err := tx.
		Model(&SessionEntity{}).
		Where("id = ?", g.SessionID).
		First(&g.SessionEntity).Error; err != nil {
		return err
	}

	switch g.Name {
	case dto.ChestGenerationType:
		cache.
			GetInstance().
			EvictGeneratedChests(g.SessionEntity.Name)
	case dto.HealthPackGenerationType:
		cache.
			GetInstance().
			EvictGeneratedHealthPacks(g.SessionEntity.Name)
	}

	return nil
}

// AssociationsEntity represents associations entity.
type AssociationsEntity struct {
	ID                int64             `gorm:"column:id;primaryKey;auto_increment;not null"`
	SessionID         int64             `gorm:"column:session_id;not null"`
	GenerationID      int64             `gorm:"column:generation_id;not null"`
	Instance          string            `gorm:"column:instance;not null"`
	Name              string            `gorm:"column:name;not null"`
	Active            bool              `gorm:"column:active;not null"`
	CreatedAt         time.Time         `gorm:"column:created_at;autoCreateTime"`
	SessionEntity     SessionEntity     `gorm:"foreignKey:SessionID;references:ID"`
	GenerationsEntity GenerationsEntity `gorm:"foreignKey:GenerationID;references:ID"`
}

// TableName retrieves name of database table.
func (*AssociationsEntity) TableName() string {
	return "associations"
}

// TableView retrieves name of database table view.
func (*AssociationsEntity) TableView() string {
	return "AssociationsEntity"
}

// BeforeCreate performs associations cache entity eviction before generations entity create.
func (a *AssociationsEntity) BeforeCreate(tx *gorm.DB) error {
	if err := tx.
		Model(&SessionEntity{}).
		Where("id = ?", a.SessionID).
		First(&a.SessionEntity).Error; err != nil {
		return err
	}

	cache.
		GetInstance().
		EvictGeneratedChests(a.SessionEntity.Name)

	return nil
}

// LobbyEntity represents lobbies entity.
type LobbyEntity struct {
	ID              int64             `gorm:"column:id;primaryKey;auto_increment;not null"`
	UserID          int64             `gorm:"column:user_id;not null"`
	SessionID       int64             `gorm:"column:session_id;not null"`
	InventoryID     int64             `gorm:"column:inventory_id;not null"`
	Skin            int64             `gorm:"column:skin;not null"`
	Health          int64             `gorm:"column:health;not null;default:100"`
	Active          bool              `gorm:"column:active;not null"`
	Host            bool              `gorm:"column:host;not null"`
	Eliminated      bool              `gorm:"column:eliminated;not null"`
	PositionX       float64           `gorm:"column:position_x;not null"`
	PositionY       float64           `gorm:"column:position_y;not null"`
	PositionStatic  bool              `gorm:"column:position_static;not null"`
	Ammo            int64             `gorm:"column:position_static;not null"`
	CreatedAt       time.Time         `gorm:"column:created_at;autoCreateTime"`
	UserEntity      UserEntity        `gorm:"foreignKey:UserID;references:ID"`
	SessionEntity   SessionEntity     `gorm:"foreignKey:SessionID;references:ID"`
	InventoryEntity []InventoryEntity `gorm:"foreignKey:InventoryID;references:UserID"`
}

// TableName retrieves name of database table.
func (*LobbyEntity) TableName() string {
	return "lobbies"
}

// TableView retrieves name of database table view.
func (*LobbyEntity) TableView() string {
	return "LobbyEntity"
}

// BeforeCreate performs lobbies cache entity eviction before lobbies entity create.
func (l *LobbyEntity) BeforeCreate(tx *gorm.DB) error {
	cache.
		GetInstance().
		EvictLobbySet(l.SessionID)

	return nil
}

// InventoryEntity represents inventory entity.
type InventoryEntity struct {
	ID            int64         `gorm:"column:id;primaryKey;auto_increment;not null"`
	UserID        int64         `gorm:"column:user_id;not null"`
	SessionID     int64         `gorm:"column:session_id;not null"`
	Name          string        `gorm:"column:name;not null"`
	CreatedAt     time.Time     `gorm:"column:created_at;autoCreateTime"`
	UserEntity    UserEntity    `gorm:"foreignKey:UserID;references:ID"`
	SessionEntity SessionEntity `gorm:"foreignKey:SessionID;references:ID"`
}

// TableName retrieves name of database table.
func (*InventoryEntity) TableName() string {
	return "inventory"
}

// TableView retrieves name of database table view.
func (*InventoryEntity) TableView() string {
	return "InventoryEntity"
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
