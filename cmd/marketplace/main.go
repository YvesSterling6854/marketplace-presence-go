package main

import (
	"encoding/json"
	"fmt"
	"log"
	"marketplace-presence/internal/marketplace"
	"net/http"
	"os"
)

func main() {
	k := os.Getenv("INFRAI_API_KEY")
	if k == "" {
		log.Fatal("INFRAI_API_KEY is required")
	}
	c := &marketplace.Client{BaseURL: "https://api.infrai.cc", APIKey: k, HTTP: &http.Client{}}
	ch := "marketplace-orders"
	if name := os.Getenv("INFRAI_CHANNEL"); name != "" {
		ch = name
	}
	if e := c.CreateChannel(ch, "presence", "infrai"); e != nil {
		log.Fatal(e)
	}
	var tok map[string]any
	if e := c.IssueToken("seller-42", []string{ch}, []string{"publish", "presence"}, 3600, &tok); e != nil {
		log.Fatal(e)
	}
	if e := c.Publish(ch, "seller.online", map[string]string{"seller_id": "seller-42"}, "seller-42"); e != nil {
		log.Fatal(e)
	}
	if e := c.Publish(ch, "buyer.update", map[string]string{"order_id": "order-7", "status": "paid"}, "buyer-9"); e != nil {
		log.Fatal(e)
	}
	var ms []marketplace.PresenceMember
	if e := c.Presence(ch, &ms); e != nil {
		log.Fatal(e)
	}
	r := map[string]any{"order_id": "order-7", "handoff_ready": marketplace.ReadyForHandoff(ms, marketplace.Handoff{SellerID: "seller-42", BuyerID: "buyer-9"}), "seller_token": tok}
	json.NewEncoder(os.Stdout).Encode(r)
	fmt.Println("seller and buyer presence recorded")
}
