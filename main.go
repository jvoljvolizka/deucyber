package main

import (
	"deucyber"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var Config deucyber.Config

func reterr(err error) error {

	fmt.Print(err.Error())
	return nil

}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1> yo ! wanna do some aws ?  </h1>")
}

func news(w http.ResponseWriter, r *http.Request) {

	if deucyber.DBstatus {
		var News []*deucyber.NewsItem
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		News = deucyber.GetNews()
		err := json.NewEncoder(w).Encode(News)
		if err != nil {
			reterr(err)
		}

	} else {
		fmt.Fprintf(w, "<h1> Database error -.-  </h1>")
	}

}

func events(w http.ResponseWriter, r *http.Request) {

	if deucyber.DBstatus {
		var Events []*deucyber.EventItem
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		Events = deucyber.GetEvents()
		err := json.NewEncoder(w).Encode(Events)
		if err != nil {
			reterr(err)
		}

	} else {
		fmt.Fprintf(w, "<h1> Database error -.-  </h1>")
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
	deucyber.Conf = Config
	deucyber.GetConf(Config)

	fmt.Println(Config.Botkey)
	fmt.Println(Config.DBtype)
	fmt.Println(Config.MasterID)

	go deucyber.DBcheck()
	go deucyber.Bot(Config)
	go deucyber.GitInit()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/news", news)
	router.HandleFunc("/events", events)

	err = http.ListenAndServe(":3300", handlers.LoggingHandler(os.Stdout, router))

	if err != nil {
		reterr(err)
	}

}
