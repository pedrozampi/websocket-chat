package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {

	file, err := os.OpenFile("messages.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	logger := zerolog.New(file).With().Timestamp().Logger()

	conn, err := upgrader.Upgrade(w, r, nil)

	name := fmt.Sprintf("user%d", int(math.Floor(rand.Float64()*10000)))

	if err != nil {
		logger.Error().Err(err)

		return
	}
	defer conn.Close()

	for {

		_, message, err := conn.ReadMessage()

		if err != nil {
			logger.Error().Err(err)
			break
		}

		logger.Debug().Str("user", name).Msg(string(message))

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			logger.Error().Err(err)
			break
		}

	}
}

func main() {
	http.HandleFunc("/websocket", websocketHandler)
	log.Fatal().Err(http.ListenAndServe(":8080", nil))
}
