package models

import (
	"time"
)

type Card struct {
	Id        string    `gorm:"primaryKey" json:"-"`
	Suit      string    `json:"suit"`
	Value     string    `json:"value"`
	DeckId    string    `gorm:"foreignKey" json:"-"`
	Code      string    `gorm:"-:all" json:"code"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Deck struct {
	Id        string    `gorm:"primaryKey" json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
	CreatedAt time.Time `json:"-" gorm:"index"`
	UpdatedAt time.Time `json:"-"`
}
