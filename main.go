package main

import (
	"deucyber"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

/*

   title:
   description :
   link :
   date:

*/

type newsItem struct {
	Title string `json:"Title"`
	Desc  string `json:"Desc"`
	Link  string `json:"Link"`
	//	Desc string `json:"Desc"`

}

type config struct {
	Botkey   string `json:"Botkey"`
	MasterID string `json:"MasterID"`
	DBtype   string `json:"DBtype"`
	DBname   string `json:"DBname"`
	//	Collection string `json:"Collection"`
	Connection string `json:"Connection"`
}

var News []*newsItem
var Config config

func parseNews(input string) (string, error) {

	parsed := strings.Split(input, " $ ")
	if len(parsed) != 3 {
		return "", errors.New("Parsing error ! ")
	}
	var item newsItem

	item.Title = parsed[0]
	item.Desc = parsed[1]
	item.Link = parsed[2]

	News = append(News, &item)
	out := "Title = " + parsed[0] + "Desc = " + parsed[1] + "Link = " + parsed[2]
	return out, nil

}

func reterr(err error) error {

	fmt.Print(err.Error())
	return nil

}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1> yo ! wanna do some aws ?  </h1>")
}

var newsList []string

func news(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(News)

	if err != nil {
		reterr(err)
	}
}

func bot() {

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Printf("function bot recovery error : ' %v ' \n", rec)
		}

	}()

	//key := os.Getenv("BOTTOKEN")
	key := Config.Botkey

	if key == "" {
		panic("Token variable is empty -.- check config file")
	}
	bot, err := tg.NewBotAPI(key)

	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we should leave it empty.
		msg := tg.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "type /addnews or cry."
		case "addnews":
			args := update.Message.CommandArguments()
			if args == "" {
				msg.Text = "Yo ! give me some tasty arguments "
			} else {
				newsList = append(newsList, args)
				rawNew, err := parseNews(args)
				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = "okay! " + rawNew + " added to list"

				}
			}
		case "myid":
			uid := update.Message.From.ID
			suid := fmt.Sprintf("%v", uid)
			msg.Text = suid + " this library is stupid "
		case "fuckmeup":
			jsondata, err := json.Marshal(News)
			if err != nil {
				msg.Text = err.Error()
			} else {
				msg.Text = string(jsondata)
			}
		default:
			msg.Text = "wat ?"
		}

		_, err := bot.Send(msg)

		if err != nil {
			reterr(err)
		}
		err = nil
	}
}

func main() {

	file := os.Args[1]

	conFile, err := os.Open(file)

	if err != nil {
		fmt.Println("file open error")
		log.Panic(err)
	}

	conData, err := ioutil.ReadAll(conFile)

	if err != nil {
		fmt.Println("file read error")
		log.Panic(err)
	}

	err = json.Unmarshal(conData, &Config)

	if err != nil {
		fmt.Println("json parsing error")
		log.Panic(err)
	}

	///clean this shit please
	var dbconf deucyber.Config
	dbconf.DBname = Config.DBname
	dbconf.DBtype = Config.DBtype
	dbconf.Connection = Config.Connection

	deucyber.GetConf(dbconf)

	fmt.Println(Config.Botkey)
	fmt.Println(Config.DBtype)
	fmt.Println(Config.MasterID)

	go bot()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/news", news)

	err = http.ListenAndServe(":3300", handlers.LoggingHandler(os.Stdout, router))

	if err != nil {
		reterr(err)
	}

}
