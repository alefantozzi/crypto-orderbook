package hyperliquid

type Config struct {
	Symbol string
}

type WsLevel struct {
	Px string `json:"px"`
	Sz string `json:"sz"`
	N  int    `json:"n"`
}

type WsBook struct {
	Coin   string      `json:"coin"`
	Levels [][]WsLevel `json:"levels"`
	Time   int64       `json:"time"`
}

type WsTrade struct {
	Coin  string   `json:"coin"`
	Side  string   `json:"side"`
	Px    string   `json:"px"`
	Sz    string   `json:"sz"`
	Hash  string   `json:"hash,omitempty"`
	Time  int64    `json:"time"`
	Tid   int64    `json:"tid,omitempty"`
	Users []string `json:"users,omitempty"`
}
