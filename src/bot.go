package deucyber

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/go-github/github"
	"gopkg.in/mgo.v2/bson"
)

func reterr(err error) error {

	fmt.Print(err.Error())
	return nil

}

func Equal(a, b []*github.PullRequest) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
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
	item.Time = time.Now().UnixNano()

	if DBstatus {
		InsertNews(item)
	} else {
		return "", errors.New("Database connection is busted -.- ")
	}

	out := "Title = " + parsed[0] + "Desc = " + parsed[1] + "Link = " + parsed[2]
	return out, nil

}

func ParseEvents(input string) (string, error) {

	parsed := strings.Split(input, " $ ")
	if len(parsed) != 6 {
		return "", errors.New("Parsing error ! ")
	}
	var item EventItem

	item.Title = parsed[0]
	item.Desc = parsed[1]
	item.Link = parsed[2]
	item.Date = parsed[3]
	item.Time = parsed[4]
	item.Location = parsed[5]

	if DBstatus {
		InsertEvents(item)
	} else {
		return "", errors.New("Database connection is busted -.- ")
	}

	out := "Title = " + parsed[0] + "Desc = " + parsed[1] + "Link = " + parsed[2] + "Date = " + parsed[3] + "Time = " + parsed[4] + "Location = " + parsed[5]
	return out, nil

}

func Bot(Con Config) {

	register := false

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

		// Create a new MessageConfig. We don't have text yet,
		// so we should leave it empty.
		msg := tg.NewMessage(update.Message.Chat.ID, "")

		/*	oldcur := CurPrs
			_ = UpdatePrs()
			if !Equal(oldcur, CurPrs) {
				msg.Text = "Yo ! we got a new Pr : " + CurPrs[len(CurPrs)-1].GetHTMLURL()
			}
		*/
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Extract the command from the Message.
		switch update.Message.Command() {
		//user commands
		case "help":
			msg.Text = "type /addnews or cry."

		case "myid":
			uid := update.Message.From.ID
			suid := fmt.Sprintf("%v", uid)
			msg.Text = suid + " this library is stupid "
		case "chatid":
			uid := update.Message.Chat.ID
			suid := fmt.Sprintf("%v", uid)
			msg.Text = suid + " this library is stupid "

		case "getnews":
			if DBstatus {
				News := GetNews()

				jsondata, err := json.Marshal(News)
				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = string(jsondata)
				}
			} else {
				msg.Text = "Database connection is busted -.-"
			}
		case "getevents":
			if DBstatus {
				Events := GetEvents()

				jsondata, err := json.Marshal(Events)
				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = string(jsondata)
				}
			} else {
				msg.Text = "Database connection is busted -.-"
			}

		//admin commands
		case "register":
			if register && update.Message.Chat.ID == Con.AdminChat {
				tmp := GetAdmin(bson.M{"id": update.Message.From.ID})
				if tmp.ID == update.Message.From.ID {
					msg.Text = "hey i know you already !"
				} else {
					var new Admin
					new.ID = update.Message.From.ID
					AddAdmin(new)
					msg.Text = "Cool Cool"
				}
			} else {
				msg.Text = "Register is closed for now"
			}
		case "killme":
			tmp := GetAdmin(bson.M{"id": update.Message.From.ID})
			if tmp.ID == update.Message.From.ID {
				DelAdmin(tmp)
				msg.Text = "Done. I will remember you..."
			} else {
				msg.Text = "Who the fuck are you ?"
			}

		case "addnews":
			args := update.Message.CommandArguments()
			if args == "" {
				msg.Text = "Yo ! give me some tasty arguments "
			} else {
				tmp := GetAdmin(bson.M{"id": update.Message.From.ID})
				if tmp.ID == update.Message.From.ID {
					rawNew, err := ParseNews(args)
					if err != nil {
						msg.Text = err.Error()
					} else {
						msg.Text = "okay! " + rawNew + " added to list"
					}
				} else {
					fmt.Println(tmp.ID)
					msg.Text = "Sorry mate you are not cool enough"
				}
			}
		case "addevents":
			args := update.Message.CommandArguments()
			if args == "" {
				msg.Text = "Yo ! give me some tasty arguments "
			} else {
				tmp := GetAdmin(bson.M{"id": update.Message.From.ID})
				if tmp.ID == update.Message.From.ID {
					rawEvent, err := ParseEvents(args)
					if err != nil {
						msg.Text = err.Error()
					} else {
						msg.Text = "okay! " + rawEvent + " added to list"

					}

				} else {
					fmt.Println(tmp.ID)
					msg.Text = "Sorry mate you are not cool enough"
				}

			}
		case "getpr":
			tmp := GetAdmin(bson.M{"id": update.Message.From.ID})
			if tmp.ID == update.Message.From.ID {
				prs, err := GetPrs()
				if err != nil {
					msg.Text = err.Error()
				} else {
					if len(prs) != 0 {
						resp := ""
						for _, pr := range prs {
							resp = resp + pr.GetHTMLURL() + fmt.Sprintf(" /merge%v", pr.GetNumber()) + "\n"
						}
						msg.Text = resp
					} else {
						msg.Text = "No new pull requests"
					}
				}

			} else {
				fmt.Println(tmp.ID)
				msg.Text = "Sorry mate you are not cool enough"
			}

			//jvol commands
		case "getconfig":
			if update.Message.From.ID == Con.MasterID {
				jsondata, _ := json.Marshal(Con)
				msg.Text = string(jsondata)

			} else {
				msg.Text = "Sorry mate you are not cool enough"
			}

		case "openregister":
			if update.Message.From.ID == Con.MasterID {
				msg.Text = "You can now add more admins"
				register = true
			} else {
				msg.Text = "Sorry mate you are not cool enough"
			}
		case "closeregister":
			if update.Message.From.ID == Con.MasterID {
				msg.Text = "Register is closed now"
				register = false
			} else {
				msg.Text = "Sorry mate you are not cool enough"
			}

		default:
			if update.Message.Command()[:5] == "merge" {
				tmp := GetAdmin(bson.M{"id": update.Message.From.ID})
				if tmp.ID == update.Message.From.ID {
					num, err := strconv.Atoi(update.Message.Command()[5:])
					if err != nil {
						msg.Text = "parsing error"
					} else {
						err = Merge(num)
						if err != nil {
							msg.Text = err.Error()
						} else {
							msg.Text = "Merge successful"
						}
					}
				} else {
					fmt.Println(tmp.ID)
					msg.Text = "Sorry mate you are not cool enough"
				}
			} else {
				msg.Text = "wat ?"
			}

		}

		_, err := bot.Send(msg)

		if err != nil {
			reterr(err)
		}
		err = nil
	}
}
