package entity

import "time"

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
