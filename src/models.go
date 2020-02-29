package deucyber

var DBstatus bool = true

type NewsItem struct {
	Title string `json:"Title"`
	Desc  string `json:"Desc"`
	Link  string `json:"Link"`
	Time  int64  `json:"Time"`
}

type Admin struct {
	ID int `json:"ID"`
}

type Config struct {
	Botkey     string `json:"Botkey"`
	MasterID   int    `json:"MasterID"`
	AdminChat  int64  `json:"AdminChat"`
	DBtype     string `json:"DBtype"`
	DBname     string `json:"DBname"`
	Connection string `json:"Connection"`
}

type EventItem struct {
	Title    string `json:"Title"`
	Desc     string `json:"Desc"`
	Link     string `json:"Link"`
	Date     string `json:"Date"`
	Time     string `json:"Time"`
	Location string `json:"Location"`
}
