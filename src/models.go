package deucyber

var DBstatus bool = true

type NewsItem struct {
	Title string `json:"Title"`
	Desc  string `json:"Desc"`
	Link  string `json:"Link"`
	Time  int64  `json:"Time"`
}

type Config struct {
	Botkey   string `json:"Botkey"`
	MasterID string `json:"MasterID"`
	DBtype   string `json:"DBtype"`
	DBname   string `json:"DBname"`

	Connection string `json:"Connection"`
}
