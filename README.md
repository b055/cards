# CARDS

Creates an In-Memory database that forms the versioned backend for a few APIs.

Fetch from the repository using
`go get github.com/b055/cards`

Build using
`go build`


Run the unit tests from the root directory with

`go test -v ./...`

Run the server with the command

`./cards`


```
INFO[0000] Connecting to database                       
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /api/v1/decks             --> github.com/b055/cards/handlers.GetAllDecks (3 handlers)
[GIN-debug] GET    /api/v1/decks/:deck_id    --> github.com/b055/cards/handlers.GetDeckById (3 handlers)
[GIN-debug] POST   /api/v1/decks             --> github.com/b055/cards/handlers.CreateDeck (3 handlers)
[GIN-debug] GET    /api/v1/decks/:deck_id/*draw --> github.com/b055/cards/handlers.DrawCardsInDeck (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Environment variable PORT is undefined. Using port :8080 by default
[GIN-debug] Listening and serving HTTP on :8080
```


The server can be reached at `http://localhost:8080`
## APIs
### Create a new Deck

POST   /api/v1/decks

It would return a given deck by its UUID. If the deck was not passed over or is invalid it should return an error. This method will "open the deck", meaning that it will list all cards by the order it was created.

#### Params
shuffled
: true/false or 1/0 boolean that determines if the created deck should be shuffled or not.

cards
: comma-separated codes to create a custom deck

Example request:
`
curl --location --request POST 'http://localhost:8080/api/v1/decks' \
--form 'shuffled="1"' \
--form 'cards="AS,KH,8C"'
`

Example response:
```
{
    "deck_id": "a251071b-662f-44b6-ba11-e24863039c59",
    "shuffled": false,
    "remaining": 30
}
```

### Open a Deck
GET    /api/v1/decks/:deck_id

Returns a given deck by its UUID. If the deck was not passed over or is invalid it should return an error.
This method lists all cards by the order it was created.

Example request:
`
curl --location --request GET 'http://localhost:8080/api/v1/decks/a251071b-662f-44b6-ba11-e24863039c59'
`

Example response:
```
{
    "deck_id": "a251071b-662f-44b6-ba11-e24863039c59",
    "shuffled": false,
    "remaining": 3,
    "cards": [
        {
            "value": "ACE",
            "suit": "SPADES",
            "code": "AS"
        },
				{
            "value": "KING",
            "suit": "HEARTS",
            "code": "KH"
        },
        {
            "value": "8",
            "suit": "CLUBS",
            "code": "8C"
        }
    ]
}
```

### Draw from a Deck
GET    /api/v1/decks/:deck_id/draw/

Example request:
`curl --location --request GET 'http://localhost:8080/api/v1/decks/74c6e0a8-dac6-11ed-b2bf-865a7a4b8830/draw?count=2'`

Example response:

```
{
    "cards": [
        {
            "suit": "HEARTS",
            "value": "King",
            "code": "KH"
        },
        {
            "suit": "SPADES",
            "value": "Ace",
            "code": "AS"
        }
    ],
    "deck_id": "74c6e0a8-dac6-11ed-b2bf-865a7a4b8830",
    "remaining": 3,
    "shuffled": true
}
```


### List all Decks
GET    /api/v1/decks

This was beyond the scope of the assignment, however I found it very useful during testing. So I decided to keep it in.
Returns a paginated list of all the decks that have been created.

#### Params
page_token
: The token required to obtain the next page. For example, `curl --location --request GET 'http://localhost:8080/api/v1/decks/?page_token=MTY4MTQ3MzcyMg=='`

Example request:
`curl --location --request GET 'http://localhost:8080/api/v1/decks/'`

Example response:

```
{
    "decks": [
        {
            "deck_id": "f89a47b0-da96-11ed-846a-865a7a4b8830",
            "shuffled": false,
            "remaining": 52
        },
        {
            "deck_id": "eff516ee-da96-11ed-8440-865a7a4b8830",
            "shuffled": false,
            "remaining": 52
        },
        {
            "deck_id": "ef831b52-da96-11ed-8414-865a7a4b8830",
            "shuffled": false,
            "remaining": 52
        },
        {
            "deck_id": "ef19d818-da96-11ed-83e6-865a7a4b8830",
            "shuffled": false,
            "remaining": 52
        },
        {
            "deck_id": "eeb71c6e-da96-11ed-83ba-865a7a4b8830",
            "shuffled": false,
            "remaining": 52
        },
        {
            "deck_id": "ee5935e0-da96-11ed-8391-865a7a4b8830",
            "shuffled": false,
            "remaining": 52
        },
        {
            "deck_id": "edf9ae36-da96-11ed-8363-865a7a4b8830",
            "shuffled": false,
            "remaining": 52
        },
        {
            "deck_id": "ed905da0-da96-11ed-8338-865a7a4b8830",
            "shuffled": false,
            "remaining": 52
        },
        {
            "deck_id": "ed287ff0-da96-11ed-830c-865a7a4b8830",
            "shuffled": false,
            "remaining": 52
        },
        {
            "deck_id": "ec7b91fa-da96-11ed-82dc-865a7a4b8830",
            "shuffled": false,
            "remaining": 52
        }
    ],
    "page_token": "MTY4MTQ1NzcyNg=="
}
```

