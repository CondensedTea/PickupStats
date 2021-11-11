package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"PickupStats/pkg/config"
	"PickupStats/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
)

const playerEndpoint = "http://api.tf2pickup.ru/players"

type PickupPlayer struct {
	SteamId string `json:"steamId" bson:"steam_id"`
	Name    string `json:"name" bson:"name"`
	Avatar  struct {
		Small  string `json:"small" bson:"small"`
		Medium string `json:"medium" bson:"medium"`
		Large  string `json:"large" bson:"large"`
	} `json:"avatar" bson:"avatar"`
	Roles          []string  `json:"roles" bson:"roles"`
	Etf2LProfileId int       `json:"etf2lProfileId" bson:"etf2l_profile_id"`
	JoinedAt       time.Time `json:"joinedAt" bson:"joined_at"`
	Id             string    `json:"id" bson:"id"`
	Links          []struct {
		Href  string `json:"href" bson:"href"`
		Title string `json:"title" bson:"title"`
	} `json:"_links" bson:"links"`
}

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	ctx := context.Background()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	client, err := db.NewClient(ctx, cfg.DSN, cfg.Database, cfg.GameCollection, cfg.NameCollection)
	if err != nil {
		log.Fatalf("Failed to init mongo client: %v", err)
	}
	players, err := GetPlayers()
	if err != nil {
		log.Fatalf("Failed to get players from API: %v", err)
	}

	log.Printf("Got %d players\n", len(players))

	for _, player := range players {
		log.Println("Processing player " + player.Name)
		bytes, err := bson.Marshal(player)
		if err != nil {
			log.Fatalf("Failed to marshall player: %v", err)
		}
		res, err := client.Conn.Database(cfg.Database).Collection("player_names_v2").InsertOne(ctx, bytes)
		if err != nil {
			log.Fatalf("Failed to insert player: %v", err)
		}
		log.Println("ID: ", res.InsertedID)
	}
	log.Println("Finished successfully")
}

func GetPlayers() ([]PickupPlayer, error) {
	players := make([]PickupPlayer, 0)
	resp, err := http.Get(playerEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned http code: %d", resp.StatusCode)
	}
	if err = json.NewDecoder(resp.Body).Decode(&players); err != nil {
		return nil, err
	}
	return players, nil
}
