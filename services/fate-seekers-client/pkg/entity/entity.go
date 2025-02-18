package entity

import "time"

// CollectionEntity represents collections entity.
type CollectionEntity struct {
	ID        int64     `gorm:"column:id;primaryKey;auto_increment;not null"`
	Name      string    `gorm:"column:name;not null;unique"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

// TableName retrieves name of database table.
func (*CollectionEntity) TableName() string {
	return "collections"
}
