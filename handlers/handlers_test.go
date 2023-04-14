package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/b055/cards/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func init() {
	models.ConnectDatabase()
}

func Test_CreateDeck_Shuffled(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("shuffled=true"))
	ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	CreateDeck(ctx)
	assert.EqualValues(t, http.StatusOK, w.Code)
	body, _ := io.ReadAll(w.Body)

	var result map[string]any
	json.Unmarshal(body, &result)
	fmt.Println(result)
	assert.NotNil(t, result["deck_id"])
	assert.True(t, result["shuffled"].(bool))
	assert.EqualValues(t, 52, result["remaining"])
}

func Test_CreateDeck_NotShuffled(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("shuffled=false"))
	ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	CreateDeck(ctx)
	assert.EqualValues(t, http.StatusOK, w.Code)
	body, _ := io.ReadAll(w.Body)

	var result map[string]any
	json.Unmarshal(body, &result)
	assert.NotNil(t, result["deck_id"])
	assert.False(t, result["shuffled"].(bool))
	assert.EqualValues(t, 52, result["remaining"])
}

func Test_CreateDeck_Cards(t *testing.T) {
	// create a deck and check for it
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("cards=AS,KD,AC,2C,KH"))
	ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	CreateDeck(ctx)
	assert.EqualValues(t, http.StatusOK, w.Code)
	body, _ := io.ReadAll(w.Body)

	var result map[string]any
	json.Unmarshal(body, &result)
	assert.NotNil(t, result["deck_id"])
	assert.False(t, result["shuffled"].(bool))
	assert.EqualValues(t, 5, result["remaining"])
}

func Test_DrawCards(t *testing.T) {
	var result map[string]any
	{
		// create the deck
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("shuffled=false"))
		ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		CreateDeck(ctx)
		assert.EqualValues(t, http.StatusOK, w.Code)
		body, _ := io.ReadAll(w.Body)

		json.Unmarshal(body, &result)
		assert.NotNil(t, result["deck_id"])
		assert.False(t, result["shuffled"].(bool))
		assert.EqualValues(t, 52, result["remaining"])
	}

	{
		// draw some cards fromthe created deck
		w := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%q/", result["deck_id"]), nil)
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Request.URL, _ = url.Parse("?count=20")
		ctx.Params = gin.Params{gin.Param{"deck_id", result["deck_id"].(string)}}
		DrawCardsInDeck(ctx)
		assert.EqualValues(t, http.StatusOK, w.Code)

	}

	{
		// check the number of cards remaining in the deck
		w := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%q/", result["deck_id"]), nil)
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{gin.Param{"deck_id", result["deck_id"].(string)}}
		GetDeckById(ctx)
		assert.EqualValues(t, http.StatusOK, w.Code)
		var get_deck_result map[string]any
		body, _ := io.ReadAll(w.Body)
		json.Unmarshal(body, &get_deck_result)

		assert.EqualValues(t, get_deck_result["deck_id"], result["deck_id"])
		assert.EqualValues(t, 32, get_deck_result["remaining"])
	}
}

func Test_GetDeckById(t *testing.T) {
	var result map[string]any
	{
		// create the deck
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("shuffled=false"))
		ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		CreateDeck(ctx)
		assert.EqualValues(t, http.StatusOK, w.Code)
		body, _ := io.ReadAll(w.Body)

		json.Unmarshal(body, &result)
		assert.NotNil(t, result["deck_id"])
		assert.False(t, result["shuffled"].(bool))
		assert.EqualValues(t, 52, result["remaining"])
	}

	{
		// check for the deck
		w := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%q/", result["deck_id"]), nil)
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{gin.Param{"deck_id", result["deck_id"].(string)}}
		GetDeckById(ctx)
		assert.EqualValues(t, http.StatusOK, w.Code)
		var get_deck_result map[string]any
		body, _ := io.ReadAll(w.Body)
		json.Unmarshal(body, &get_deck_result)

		assert.EqualValues(t, get_deck_result["deck_id"], result["deck_id"])

	}
}

func Test_GetAllDecks(t *testing.T) {
	{
		// creates 15 decks
		for i := 1; i <= 13; i++ {
			// create the deck
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("shuffled=false"))
			ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			CreateDeck(ctx)
			assert.EqualValues(t, http.StatusOK, w.Code)
			body, _ := io.ReadAll(w.Body)
			var result map[string]any
			json.Unmarshal(body, &result)
			assert.NotNil(t, result["deck_id"])
			assert.False(t, result["shuffled"].(bool))
			assert.EqualValues(t, 52, result["remaining"])
		}

	}

	var first_get_decks_result map[string]any
	var page_count []int
	for {
		// check for the deck
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		if val, ok := first_get_decks_result["page_token"]; ok && val != nil {
			ctx.Request.URL, _ = url.Parse("?page_token=" + val.(string))
		}
		ctx.Request.Header.Set("Content-Type", "application/json")
		GetAllDecks(ctx)
		assert.EqualValues(t, http.StatusOK, w.Code)
		body, _ := io.ReadAll(w.Body)
		json.Unmarshal(body, &first_get_decks_result)

		// assert.EqualValues(t, 10, len(first_get_decks_result["decks"].([]any)))
		assert.True(t, len(first_get_decks_result["decks"].([]any)) > 0)
		// assert.NotNil(t, first_get_decks_result["page_token"])

		page_count = append(page_count, len(first_get_decks_result["decks"].([]any)))

		if val, _ := first_get_decks_result["page_token"]; val == nil {
			break
		}
	}
	fmt.Println(page_count)
	// expect to go through at least 2 pages
	assert.True(t, len(page_count) > 1)
	// the last page shouldn't have the maximum anount of decks
	assert.True(t, page_count[len(page_count)-1] < 10)

	// {
	// 	// check for the deck
	// 	w := httptest.NewRecorder()
	// 	ctx, _ := gin.CreateTestContext(w)

	// 	ctx.Request = httptest.NewRequest(http.MethodGet, "/decks/", nil)
	// 	ctx.Request.URL, _ = url.Parse("?page_token=" + first_get_decks_result["page_token"].(string))
	// 	ctx.Request.Header.Set("Content-Type", "application/json")
	// 	GetAllDecks(ctx)
	// 	assert.EqualValues(t, http.StatusOK, w.Code)
	// 	body, _ := io.ReadAll(w.Body)
	// 	var second_get_decks_result map[string]any
	// 	json.Unmarshal(body, &second_get_decks_result)

	// 	assert.EqualValues(t, 5, len(second_get_decks_result["decks"].([]any)))
	// 	assert.Nil(t, second_get_decks_result["page_token"])

	// }
}
