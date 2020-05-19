package deucyber

var DBstatus bool = true
var Conf Config

type NewsItem struct {
	Title       string `json:"Title"`
	Desc        string `json:"Desc"`
	Link        string `json:"Link"`
	PublishTime int64  `json:"PublishTime"`
}

type Admin struct {
	ID int `json:"ID"`
}

type Config struct {
	Botkey         string `json:"Botkey"`
	MasterID       int    `json:"MasterID"`
	AdminChat      int64  `json:"AdminChat"`
	DBtype         string `json:"DBtype"`
	DBname         string `json:"DBname"`
	Connection     string `json:"Connection"`
	GithubUsername string `json:"GithubUsername"`
	GithubRepo     string `json:"GithubRepo"`
	GithubApiKey   string `json:"GithubApiKey"`
	SplitString    string `json:"SplitString"`
}

type EventItem struct {
	Title       string `json:"Title"`
	Desc        string `json:"Desc"`
	Link        string `json:"Link"`
	Date        string `json:"Date"`
	Time        string `json:"Time"`
	Location    string `json:"Location"`
	PublishTime int64  `json:"PublishTime"`
}
