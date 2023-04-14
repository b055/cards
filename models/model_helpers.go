package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Suit int64
type Value int64

const (
	Spades Suit = iota
	Clubs
	Hearts
	Diamonds
)

const (
	Ace Value = iota
	One
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

func (s Suit) String() string {
	switch s {
	case Spades:
		return "SPADES"
	case Clubs:
		return "CLUBS"
	case Hearts:
		return "HEARTS"
	case Diamonds:
		return "DIAMONDS"
	}
	return "unknown"
}

func (v Value) String() string {
	switch v {
	case Jack:
		return "Jack"
	case Queen:
		return "Queen"
	case King:
		return "King"
	case Ace:
		return "Ace"
	}
	values := []string{"A", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	if int(v) >= len(values) {
		return "unknown"
	}
	return values[v]
}

func ToSuit(s string) (Suit, error) {
	switch strings.ToUpper(s) {
	case "S":
		return Spades, nil
	case "C":
		return Clubs, nil
	case "H":
		return Hearts, nil
	case "D":
		return Diamonds, nil
	}
	return -1, errors.New("invalid suit " + s)
}

func ToValue(s string) (Value, error) {
	switch strings.ToUpper(s) {
	case "A":
		return Ace, nil
	case "J":
		return Jack, nil
	case "Q":
		return Queen, nil
	case "K":
		return King, nil
	}
	if value, err := strconv.Atoi(s); err != nil {
		return Value(0), err
	} else {
		return Value(value), nil
	}
}

func CodeToSuitValue(code string) (*Suit, *Value, error) {
	if len(code) != 2 && len(code) != 3 {
		return nil, nil, fmt.Errorf("invalid code: '"+code+"' length: %d", len(code))
	}
	suit, err := ToSuit(code[len(code)-1:])
	if err != nil {
		return nil, nil, err
	}
	value, err := ToValue(code[:len(code)-1])
	if err != nil {
		return nil, nil, err
	}
	return &suit, &value, nil
}

func (card *Card) ComputeCode() {
	card.Code = card.Value[:1] + card.Suit[:1]
}

func ConnectDatabase() error {
	log.Info("Connecting to database")

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		msg := fmt.Errorf("db connection error: %s", err)
		log.Error(msg)
		panic(msg)
	}

	if err := db.AutoMigrate(&Card{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&Deck{}); err != nil {
		panic(err)
	}
	DB = db
	return nil
}
