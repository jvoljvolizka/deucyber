package deucyber

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

/*

   title:
   description :
   link :
   date:

*/

type NewsItem struct {
	Title string `json:"Title"`
	Desc  string `json:"Desc"`
	Link  string `json:"Link"`
	//	Desc string `json:"Desc"`

}

type Config struct {
	Botkey   string `json:"Botkey"`
	MasterID string `json:"MasterID"`
	DBtype   string `json:"DBtype"`
	DBname   string `json:"DBname"`
	//	Collection string `json:"Collection"`
	Connection string `json:"Connection"`
}

var News []*NewsItem

func reterr(err error) error {

	fmt.Print(err.Error())
	return nil

}

func ParseNews(input string) (string, error) {

	parsed := strings.Split(input, " $ ")
	if len(parsed) != 3 {
		return "", errors.New("Parsing error ! ")
	}
	var item NewsItem

	item.Title = parsed[0]
	item.Desc = parsed[1]
	item.Link = parsed[2]

	InsertNews(item)

	out := "Title = " + parsed[0] + "Desc = " + parsed[1] + "Link = " + parsed[2]
	return out, nil

}

func Bot(Con Config) {

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Printf("function bot recovery error : ' %v ' \n", rec)
		}

	}()

	//key := os.Getenv("BOTTOKEN")
	key := Con.Botkey

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

				rawNew, err := ParseNews(args)
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
			News = GetNews()
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
