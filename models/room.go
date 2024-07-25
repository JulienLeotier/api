package models

type Room struct {
	ID          uint     `json:"id" gorm:"primary_key"`
	Name        string   `json:"name" gorm:"unique;not null"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
	Videos      []string `json:"videos"`
	Audios      []string `json:"audios"`
}

type VariantList struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	Variants      []User `json:"variants" gorm:"many2many:variant_lists;"`
	PriorityFirst []User `json:"priority_first" gorm:"many2many:variant_lists;"`
}

type RolePlay struct {
	ID          uint     `json:"id" gorm:"primary_key"`
	Name        string   `json:"name" gorm:"unique;not null"`
	Description string   `json:"description"`
	Photo       []string `json:"photo"`
	Audios      []string `json:"audios"`
	Videos      []string `json:"videos"`
	Room        Room     `json:"room" gorm:"foreignKey:RoomID"`
	User        User     `json:"user" gorm:"foreignKey:UserID"`
	IsVariant   bool     `json:"is_variant"`
	IsDetective bool     `json:"is_detective"`
}
