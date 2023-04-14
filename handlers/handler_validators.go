package handlers

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/b055/cards/models"
)

// Contains validation logic for the API endpoints
//

func validateGetDeckById(deck_id string) (string, error) {
	if deck_id == "" {
		return "", errors.New("invalid deck_id")
	}
	return deck_id, nil
}

func validateGetAllDecks(page_token string) (*time.Time, error) {
	if page_token != "" {
		token_decoded, err := base64.StdEncoding.DecodeString(page_token)
		if err != nil {
			log.Error(err)
			return nil, errors.New("invalid page token")
		}
		value, err := strconv.ParseInt(string(token_decoded), 10, 64)
		if err != nil {
			log.Error(err)
			return nil, errors.New("invalid page token")
		}
		paginator := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		if value <= paginator.UnixNano() {
			return nil, errors.New("invalid page token")
		}
		paginator = time.Unix(0, value)
		return &paginator, nil
	}
	return nil, nil
}

func validateGetCardsInDeck(deck_id string, count_param string) (string, int, error) {
	if deck_id == "" {
		return "", 0, errors.New("invalid deck_id")
	}

	count := 1
	if count_param != "" {
		if count_value, err := strconv.Atoi(count_param); err != nil {
			log.Error(err)
			return "", 0, errors.New("invalid count " + count_param)
		} else {
			if count_value < 1 {
				return "", 0, errors.New("invalid count " + count_param)
			}
			count = count_value
			if count > NUMBER_OF_CARDS {
				count = NUMBER_OF_CARDS
			}
		}
		return deck_id, count, nil
	} else {
		return "", 0, errors.New("count required")
	}
}

func validateCreateDeck(cards *[]models.Card, shuffled_param string, cards_param string) (bool, error) {
	log.Info("CreateDeck called")
	var shuffled = false
	if shuffled_param != "" {
		log.Info("shuffle parameter " + shuffled_param)
		// check if shuffled parameter is valid
		if shuffled_param == "true" || shuffled_param == "1" {
			shuffled = true
		} else if (shuffled_param != "false") && (shuffled_param != "0") {
			return shuffled, errors.New("Invalid parameter shuffled: " + shuffled_param)
		}
	}
	if cards_param != "" {
		log.Info("cards " + cards_param)
		// check if cards parameters are valid
		for _, card_param := range strings.Split(cards_param, ",") {
			card_param = strings.TrimSpace(card_param)
			if len(card_param) == 0 {
				return false, errors.New("Missing card")
			}
			if suit, value, err := models.CodeToSuitValue(card_param); err != nil {
				log.Error(err)
				return false, errors.New("Invalid card: " + card_param)
			} else {
				card_id, uuid_err := uuid.NewUUID()
				if uuid_err != nil {
					panic(uuid_err)
				}

				*cards = append(*cards, models.Card{Id: card_id.String(), Suit: suit.String(), Value: value.String(), DeckId: "place-holder"})
			}
		}
	}
	return shuffled, nil
}
