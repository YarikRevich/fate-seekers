package entity

import "time"

// Collection represents collection entity.
type Collection struct {
	ID        int64     `gorm:"column:id;primaryKey;auto_increment;not null"`
	Name      string    `gorm:"column:name;not null;unique"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

// TableName retrieves name of database table.
func (*Collection) TableName() string {
	return "collection"
}
