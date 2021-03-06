package db

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name string `bson:"Name"`
}

const (
	dpmAggregationTemplate = `
	[
		{
			"$match": {"player.class": {"%s": "%s"}}
		},
		{
			"$group": {
				"_id":"$player.steam_id",
				"sum_damage": {"$sum": "$stats.damage_done" },
				"sum_playtime": {"$sum": "$length"},
				"count_games": {"$sum": 1}
			}
		},
		{
			"$project": {
				"dpm": {"$round": [{"$divide": ["$sum_damage", {"$divide": ["$sum_playtime", 60]}]}, 2]},
				"games": "$count_games"
			}
		},
		{"$sort": {"dpm": -1, "games": -1}},
		{"$match": {"games": {"$gt": %d}}}
	]`
	kdrAggregationTemplate = `
	[
		{
			"$match": {"player.class": {"%s": "%s"}}
		},
		{
			"$group": {
				"_id":"$player.steam_id",
				"sum_kills": {"$sum": "$stats.kills" },
				"sum_deaths": {"$sum": "$stats.deaths"},
				"count_games": {"$sum": 1}
			}
		},
		{
			"$project": {
				"kdr": {"$round": [ {"$divide": ["$sum_kills", "$sum_deaths"]}, 1]},
				"games": "$count_games"
			}
		},
		{"$sort": {"kdr": -1, "games": -1}},
		{"$match": {"games": {"$gt": %d}}}
	]`
	healsPerMinAggregationTemplate = `
	[
		{
			"$match": {"player.class": {"$eq": "medic"}}
		},
		{
			"$group": {
				"_id": "$player.steam_id",
				"sum_heals": {"$sum": "$stats.healed"},
				"sum_playtime": {"$sum": "$length"},
				"count_games": {"$sum": 1}
			}
		},
		{
			"$project": {
				"hpm": {"$round": [{"$divide": ["$sum_heals", {"$divide": ["$sum_playtime", 60]}]}, 2]},
				"games": "$count_games"
			}
		},
		{"$sort": {"hpm": -1, "games": -1}},
		{"$match": {"games": {"$gt": %d}}}
	]`
)

type Client struct {
	database, games, names string
	ctx                    context.Context
	Conn                   *mongo.Client
}

type Player struct {
	Name   string `json:"Name"`
	Avatar string `json:"Avatar"`
}

type Result struct {
	PlayerName string   `json:"player_name"`
	Avatar     string   `json:"avatar"`
	SteamID64  string   `json:"steamid64"`
	DPM        *float64 `json:"dpm,omitempty"`
	KDR        *float64 `json:"kdr,omitempty"`
	HPM        *float64 `json:"hpm,omitempty"`
	Games      int32    `json:"games"`
}

func (r *Result) SetName(name string) {
	r.PlayerName = name
}

func NewClient(ctx context.Context, dsn, database, gamesCollection, namesCollection string) (*Client, error) {
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}
	return &Client{
		database: database,
		games:    gamesCollection,
		names:    namesCollection,
		ctx:      ctx,
		Conn:     conn,
	}, nil
}

func (c *Client) GetAverageDPM(class string, minGames int) (results []Result, err error) {
	var pipeline string
	if class == "" {
		pipeline = fmt.Sprintf(dpmAggregationTemplate, "$ne", "medic", minGames)
	} else {
		pipeline = fmt.Sprintf(dpmAggregationTemplate, "$eq", class, minGames)
	}

	var item bson.M
	opts := options.Aggregate()

	p, err := ParseMongoPipeline(pipeline)
	if err != nil {
		return nil, err
	}

	playerNames, err := c.PlayerNames()
	if err != nil {
		return nil, err
	}

	cur, err := c.Conn.
		Database(c.database).
		Collection(c.games).Aggregate(c.ctx, p, opts)
	if err != nil {
		return nil, err
	}

	for cur.Next(c.ctx) {
		r := &Result{}
		if err = cur.Decode(&item); err != nil {
			return nil, err
		}
		r.SteamID64 = item["_id"].(string)
		dpm := item["dpm"].(float64)
		r.DPM = &dpm
		r.Games = item["games"].(int32)
		r.PlayerName = playerNames[r.SteamID64].Name
		r.Avatar = playerNames[r.SteamID64].Avatar
		if err != nil {
			return nil, err
		}
		results = append(results, *r)
	}
	return results, nil
}

func (c *Client) GetAverageKDR(class string, minGames int) (results []Result, err error) {
	var pipeline string
	if class == "" {
		pipeline = fmt.Sprintf(kdrAggregationTemplate, "$ne", "medic", minGames)
	} else {
		pipeline = fmt.Sprintf(kdrAggregationTemplate, "$eq", class, minGames)
	}

	var item bson.M
	opts := options.Aggregate()

	p, err := ParseMongoPipeline(pipeline)
	if err != nil {
		return nil, err
	}

	playerNames, err := c.PlayerNames()
	if err != nil {
		return nil, err
	}

	cur, err := c.Conn.
		Database(c.database).
		Collection(c.games).Aggregate(c.ctx, p, opts)
	if err != nil {
		return nil, err
	}

	for cur.Next(c.ctx) {
		r := &Result{}
		if err = cur.Decode(&item); err != nil {
			return nil, err
		}
		r.SteamID64 = item["_id"].(string)
		kdr := item["kdr"].(float64)
		r.KDR = &kdr
		r.Games = item["games"].(int32)
		r.PlayerName = playerNames[r.SteamID64].Name
		r.Avatar = playerNames[r.SteamID64].Avatar
		results = append(results, *r)
	}
	return results, nil
}

func (c *Client) GetAverageHealsPerMin(minGames int) (results []Result, err error) {
	pipeline := fmt.Sprintf(healsPerMinAggregationTemplate, minGames)

	var item bson.M
	opts := options.Aggregate()

	p, err := ParseMongoPipeline(pipeline)
	if err != nil {
		return nil, err
	}

	cur, err := c.Conn.
		Database(c.database).
		Collection(c.games).Aggregate(c.ctx, p, opts)
	if err != nil {
		return nil, err
	}

	playerNames, err := c.PlayerNames()
	if err != nil {
		return nil, err
	}

	for cur.Next(c.ctx) {
		r := &Result{}
		if err = cur.Decode(&item); err != nil {
			return nil, err
		}
		r.SteamID64 = item["_id"].(string)
		hpm := item["hpm"].(float64)
		r.HPM = &hpm
		r.Games = item["games"].(int32)
		r.PlayerName = playerNames[r.SteamID64].Name
		r.Avatar = playerNames[r.SteamID64].Avatar
		results = append(results, *r)
	}
	return results, nil
}

func (c *Client) PlayerNames() (map[string]Player, error) {
	users := make(map[string]Player)

	cur, err := c.Conn.
		Database(c.database).
		Collection(c.names).
		Find(c.ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var item bson.M

	for cur.Next(c.ctx) {
		if err = cur.Decode(&item); err != nil {
			return nil, err
		}
		steamID := item["steam_id"].(string)
		name := item["name"].(string)
		avatar := item["avatar"].(bson.M)
		small := avatar["small"].(string)
		users[steamID] = Player{
			Name:   name,
			Avatar: small,
		}
	}
	return users, nil
}

func (c *Client) GetGamesCount() (int64, error) {
	return c.Conn.
		Database(c.database).
		Collection(c.games).
		CountDocuments(c.ctx, bson.D{})
}

func ParseMongoPipeline(str string) (pipeline mongo.Pipeline, err error) {
	str = strings.TrimSpace(str)
	if strings.Index(str, "[") != 0 {
		var doc bson.D
		if err = bson.UnmarshalExtJSON([]byte(str), false, &doc); err != nil {
			return nil, err
		}
		pipeline = append(pipeline, doc)
	} else {
		if err = bson.UnmarshalExtJSON([]byte(str), false, &pipeline); err != nil {
			return nil, err
		}
	}
	return pipeline, nil
}
