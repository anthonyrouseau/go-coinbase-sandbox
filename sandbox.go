package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("wss://ws-feed.pro.coinbase.com", nil)
	if err != nil {
		log.Fatal("connection error: ", err)
	}
	defer conn.Close()
	initialMessage := SubscribeMessage{Type: "subscribe", Product_ids: []string{"BTC-USD"}, Feeds: []string{"ticker"}}
	if err := conn.WriteJSON(&initialMessage); err != nil {
		fmt.Println("error writing intial subscribe", err)
	}
	feedConnected := false
	for !feedConnected {
		subMessage := SubscriptionsMessage{}
		err := conn.ReadJSON(&subMessage)
		if err != nil {
			fmt.Println("error reading submessage", err)
		}
		feedConnected = true
	}
	fmt.Println("starting feed")
	go func() {
		for {
			tickerUpdate := TickerMessage{}
			err := conn.ReadJSON(&tickerUpdate)
			if err != nil {
				fmt.Println("error reading ticker message", err)
			}
			fmt.Printf("%+v\n", tickerUpdate)
		}
	}()

	timer := time.NewTimer(20 * time.Second)

	<-timer.C

}

type TickerMessage struct {
	Type       string `json:type`
	Trade_id   int    `json:trade_id`
	Sequence   int    `json:sequence`
	Time       string `json:time`
	Product_id string `json:product_id`
	Price      string `json:price`
	Side       string `json:side`
	Last_size  string `json:last_size`
	Best_bid   string `json:best_bid`
	Best_ask   string `json:best_ask`
}

type SubscribeMessage struct {
	Type        string   `json:"type"`
	Product_ids []string `json:"product_ids"`
	Feeds       []string `json:"channels"`
}

type Feed struct {
	Name        string   `json:"name"`
	Product_ids []string `json:"product_ids"`
}

type SubscriptionsMessage struct {
	Type  string `json:type`
	Feeds []Feed `json:channels`
}
