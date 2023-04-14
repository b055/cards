package handlers

import (
	"encoding/base64"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/b055/cards/models"
)

// Test_validateGetDeckByIdEmpty calls handlers.validateGetDeckById with an empty string,
// checking for an error.
func Test_validateGetDeckById(t *testing.T) {
	deck_id, err := validateGetDeckById("aaldfaf")
	if deck_id == "" || err != nil {
		t.Fatalf(`validateGetDeckById("") = %q, %v, want "", error`, deck_id, err)
	}
}

// Test_validateGetDeckByIdEmpty calls handlers.validateGetDeckById with an empty string,
// checking for an error.
func Test_validateGetDeckById_Empty(t *testing.T) {
	deck_id, err := validateGetDeckById("")
	if deck_id != "" || err == nil {
		t.Fatalf(`validateGetDeckById("") = %q, %v, want "", error`, deck_id, err)
	}
}

// Test_validateGetAllDecks_InvalidToken calls handlers.validateGetAllDecks with invalid tokens,
// should return a error.
func Test_validateGetAllDecks_InvalidToken(t *testing.T) {
	for _, token := range []string{"balhdfa", base64.StdEncoding.EncodeToString([]byte("0"))} {
		paginator, err := validateGetAllDecks(token)
		if err == nil {
			t.Fatalf(`validateGetAllDecks(%q) = %q, %v, want %v, %v`, token, paginator, err, paginator, errors.New("invalid page token"))
		}
	}
}

// Test_validateGetAllDecks_InvalidToken calls handlers.validateGetAllDecks with valid tokens,
// should not return an error.
func Test_validateGetAllDecks_ValidToken(t *testing.T) {
	paginator, err := validateGetAllDecks(base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(int(time.Now().UnixNano())))))
	if err != nil {
		t.Fatalf(`validateGetAllDecks("") = %q, %v, want "", error`, paginator, err)
	}
}

// Test_validateGetAllDecks_InvalidToken calls handlers.validateGetCardsInDeck with invalid count,
// should return a error.
func Test_validateGetCardsInDeck_InvalidCount(t *testing.T) {
	deck_id := "blah"
	for _, count := range []string{"balhdfa", "-1", "0", "-.12", "1.1234", "-1.23"} {
		validated_deck_id, validated_count, err := validateGetCardsInDeck(deck_id, count)
		if err == nil {
			t.Fatalf(`validateGetCardsInDeck(%q, %q) = %q, %q, %v, want %q, %q, %v`, deck_id, count, validated_deck_id, validated_count, err, deck_id, count, errors.New("invalid count "+count))
		}
	}
}

// Test_validateGetAllDecks_InvalidToken calls handlers.validateGetCardsInDeck with valid count,
// should not return an error.
func Test_validateGetCardsInDeck_ValidCount(t *testing.T) {
	deck_id := "blah"
	for _, count := range []string{"1", "10", "20"} {
		validated_deck_id, validated_count, err := validateGetCardsInDeck(deck_id, count)
		if err != nil {
			t.Fatalf(`validateGetCardsInDeck(%q, %q) = %q, %q, %v, want %q, %q, %v`, deck_id, count, validated_deck_id, validated_count, err, deck_id, count, errors.New("invalid count "+count))
		}
	}
}

// Test_validateGetCardsInDeck_InvalidDeckId calls handlers.validateGetCardsInDeck with invalid deck_id,
// should not return an error.
func Test_validateGetCardsInDeck_InvalidDeckId(t *testing.T) {
	deck_id := ""
	count := "1"
	validated_deck_id, validated_count, err := validateGetCardsInDeck(deck_id, count)
	if err == nil {
		t.Fatalf(`validateGetCardsInDeck(%q, %q) = %q, %q, %v, want %q, %q, %v`, deck_id, count, validated_deck_id, validated_count, err, deck_id, count, errors.New("invalid deck_id"))
	}

}

// Test_validateGetAllDecks_InvalidToken calls handlers.validateGetCardsInDeck with invalid count,
// should return a error.
func Test_validateCreateDeck_InvalidShuffleParam(t *testing.T) {
	var cards []models.Card
	cards_param := "AS,KD,AC,2C,KH"
	for _, shuffle_param := range []string{"balhdfa", "-1", "-.12", "1.1234", "-1.23"} {
		shuffled, err := validateCreateDeck(&cards, shuffle_param, cards_param)
		if err == nil {
			t.Fatalf(`validateGetCardsInDeck(%v, %q, %q) = %t, %v, want %t, %v`, cards, shuffle_param, cards_param, shuffled, err, false, errors.New("Invalid parameter shuffled:  "+shuffle_param))
		}
	}
}

// Test_validateGetAllDecks_InvalidToken calls handlers.validateGetCardsInDeck with valid count,
// should not return an error.
func Test_validateCreateDeck_ValidShuffleParam(t *testing.T) {
	var cards []models.Card
	cards_param := "AS,KD,AC,2C,KH"
	for _, shuffle_param := range []string{"true", "1", "0", "false"} {
		shuffled, err := validateCreateDeck(&cards, shuffle_param, cards_param)
		if err != nil {
			t.Fatalf(`validateGetCardsInDeck(%v, %q, %q) = %t, %v, want %t, %v`, cards, shuffle_param, cards_param, shuffled, err, false, nil)
		}
	}
}

// Test_validateGetAllDecks_InvalidToken calls handlers.validateGetCardsInDeck with invalid count,
// should return a error.
func Test_validateCreateDeck_InvalidCardsParam(t *testing.T) {
	var cards []models.Card
	shuffle_param := "true"
	for _, cards_param := range []string{"AS,aa", "blah", "AS,KD,AC,,KH", "-1", "-.12", "1.1234", "-1.23"} {
		shuffled, err := validateCreateDeck(&cards, shuffle_param, cards_param)
		if err == nil {
			t.Fatalf(`validateGetCardsInDeck(%v, %q, %q) = %t, %v, want %t, %v`, cards, shuffle_param, cards_param, shuffled, err, false, errors.New("Invalid parameter shuffled:  "+shuffle_param))
		}
	}
}

// Test_validateGetAllDecks_InvalidToken calls handlers.validateGetCardsInDeck with valid count,
// should not return an error.
func Test_validateCreateDeck_ValidCardsParam(t *testing.T) {
	var cards []models.Card
	shuffle_param := "true"
	for _, cards_param := range []string{"AS,KD,AC,2C,KH", "AS,KD", "AS, KD "} {
		shuffled, err := validateCreateDeck(&cards, shuffle_param, cards_param)
		if err != nil {
			t.Fatalf(`validateGetCardsInDeck(%v, %q, %q) = %t, %v, want %t, %v`, cards, shuffle_param, cards_param, shuffled, err, false, errors.New("Invalid parameter shuffled:  "+shuffle_param))
		}
	}
}
