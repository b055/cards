package handlers

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/b055/cards/models"
)

// Contains the handlers for the different API endpoints


const NUMBER_OF_CARDS = 52
const PAGE_SIZE = 10

func GetAllDecks(c *gin.Context) {
	log.Info("GetAllDecks called")
	paginator, validation_err := validateGetAllDecks(c.Query("page_token"))
	if validation_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": validation_err})
		return
	}

	var decks []models.Deck
	// using n + 1 pagination
	var decks_result *gorm.DB
	if paginator == nil {
		decks_result = models.DB.Order("created_at desc").Limit(PAGE_SIZE + 1).Find(&decks)
	} else {
		decks_result = models.DB.Order("created_at desc").Limit(PAGE_SIZE+1).Where("created_at < ?", paginator).Find(&decks)
	}
	if decks_result.Error != nil {
		log.Error(decks_result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to list decks"})
		return
	} else {
		token := ""
		if len(decks) == PAGE_SIZE+1 {
			token = strconv.FormatInt(decks[len(decks)-2].CreatedAt.UnixNano(), 10)
		}
		if token != "" {
			c.JSON(http.StatusOK, gin.H{
				"page_token": base64.StdEncoding.EncodeToString([]byte(token)),
				"decks":      decks[:len(decks)-1]})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"page_token": nil,
				"decks":      decks})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "getAllDecks Called"})
}

func GetDeckById(c *gin.Context) {
	log.Info("GetDeckById Called")

	deck_id, validation_err := validateGetDeckById(c.Param("deck_id"))
	if validation_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": validation_err})
		return
	}
	log.Info("GetDeckById " + deck_id + " Called")

	var deck models.Deck
	if deck_result := models.DB.First(&deck, "id = ?", deck_id); deck_result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "deck_id " + deck_id + " not found"})
		return
	} else {
		var cards []models.Card
		if cards_result := models.DB.Where("deck_id = ?", deck_id).Find(&cards); cards_result.Error != nil {
			log.Error(cards_result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get cards for deck_id " + deck_id})
			return
		} else {
			for i := 0; i < len(cards); i++ {
				cards[i].ComputeCode()
			}
			c.JSON(http.StatusOK, gin.H{"deck_id": deck_id,
				"shuffled":  deck.Shuffled,
				"remaining": deck.Remaining,
				"cards":     cards})
			return
		}
	}

}

func CreateDeck(c *gin.Context) {
	log.Info("CreateDeck Called")

	var cards []models.Card
	shuffled, validation_err := validateCreateDeck(&cards, c.PostForm("shuffled"), c.PostForm("cards"))
	if validation_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": validation_err})
		return
	}

	deck_id, uuid_err := uuid.NewUUID()
	if uuid_err != nil {
		panic(uuid_err)
	}
	card_count := NUMBER_OF_CARDS
	if len(cards) > 0 {
		card_count = len(cards)
	}
	deck := models.Deck{Id: deck_id.String(), Shuffled: shuffled, Remaining: card_count}
	if result := models.DB.Create(&deck); result.Error != nil {
		log.Errorf("Failed to create deck %v", deck)
		log.Error(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"message": result.Error})
		return
	}
	if len(cards) == 0 {
		// create whole deck of cards
		for i := 0; i <= 4; i++ {
			for j := 0; j <= 13; j++ {
				card_id, uuid_err := uuid.NewUUID()
				if uuid_err != nil {
					panic(uuid_err)
				}
				cards = append(cards, models.Card{Id: card_id.String(), Suit: models.Suit(i).String(), Value: models.Value(j).String(), DeckId: deck.Id})
			}
		}
	}
	if shuffled {
		rand.Shuffle(len(cards), func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})
	}
	// create the specided amount of cards
	for i := 0; i < len(cards); i++ {
		cards[i].DeckId = deck.Id
		if result := models.DB.Create(&cards[i]); result.Error != nil {
			log.Errorf("Failed to create card %v", cards[i])
			log.Error(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{"message": result.Error})
			return
		}
	}

	c.JSON(http.StatusOK, deck)
}

func DrawCardsInDeck(c *gin.Context) {
	log.Info("GetCardsInDeck Called")

	deck_id, count, err := validateGetCardsInDeck(c.Param("deck_id"), c.Query("count"))
	if err != nil {
		log.Error("invalid deck_id or count")
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	log.Info("GetCardsInDeck " + deck_id + " Called")

	var deck models.Deck
	if deck_result := models.DB.First(&deck, "id = ?", deck_id); deck_result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "deck_id " + deck_id + " not found"})
		return
	} else {
		var cards []models.Card
		if cards_result := models.DB.Where("deck_id = ?", deck_id).Limit(count).Find(&cards); cards_result.Error != nil {
			log.Error(cards_result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get cards for deck_id " + deck_id})
			return
		} else {
			for i := 0; i < len(cards); i++ {
				cards[i].ComputeCode()
				if delete_result := models.DB.Delete(cards[i]); delete_result.Error != nil {
					log.Errorf("Failed to delete card %v for deck_id %s", cards[i], deck_id)
					log.Error(delete_result.Error)
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get cards for deck_id " + deck_id})
					return
				}
			}
			remaining := deck.Remaining - len(cards)
			if remaining < 0 {
				remaining = 0
			}
			if update_result := models.DB.Model(deck).Update("Remaining", remaining); update_result.Error != nil {
				log.Error("Failed to update remaining cards for deck_id " + deck_id)
				log.Error(update_result.Error)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get cards for deck_id " + deck_id})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"cards": cards})
			return
		}

	}
	c.JSON(http.StatusOK, gin.H{"message": "getCardsInDeck " + deck_id + " Called"})
}
