package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Article struct {
	Id        string `json:"Id"`
	Setup     string `json:"Setup"`
	Punchline string `json:"Punchline"`
}

type Articles struct {
	Article []Article `json:"Article"`
}

var articles Articles
var mode string
var operationMode = &mode

// var articlesDryRun = Articles{
// 	Article: []Article{
// 		{
// 			Id:        "1",
// 			Setup:     "Reading from memory Want to hear a joke about a piece of paper?",
// 			Punchline: "Never mind...it's tearable",
// 		},
// 		{
// 			Id:        "2",
// 			Setup:     "Reading from memory Which side of the chicken has more feathers?",
// 			Punchline: "The outside.",
// 		},
// 	},
// }

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", readAllArticles)
	// myRouter.HandleFunc("/articles", readAllArticles).Queries("mode", "{[dr, fl, db]}")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	// myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

// func deleteArticle(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	for index, article := range articles.Article {
// 		if article.Id == id {
// 			articles.Article = append(articles.Article[:index], articles.Article[index+1:]...)
// 		}
// 	}
// 	file, _ := json.MarshalIndent(articles.Article, "", " ")
// 	_ = ioutil.WriteFile("jokes/jokesfile.json", file, 0644)
// }

func readAllArticles(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	for _, article := range articles.Article {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Fprintf(w, "%+v", string(reqBody))
	fmt.Println(string(reqBody))
	fmt.Println("Articles before", &articles)
	err := json.Unmarshal(reqBody, &articles)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("Articles after", &articles)
	// articles.Article = append(articles.Article, &articles)
	writeToSource()
}

func readFromSource() {
	if *operationMode == "file" {
		jokesFile := "jokes/jokesfile.json"
		jokes, err := os.Open(jokesFile)
		if err != nil {
			fmt.Printf(err.Error())
		}
		defer jokes.Close()
		byteValue, _ := ioutil.ReadAll(jokes)
		json.Unmarshal(byteValue, &articles)
	} else if *operationMode == "db" {
		fmt.Println("Connecting to the DB")
	} else {
		articles = Articles{
			Article: []Article{
				{
					Id:        "1",
					Setup:     "Reading from memory Want to hear a joke about a piece of paper?",
					Punchline: "Never mind...it's tearable",
				},
				{
					Id:        "2",
					Setup:     "Reading from memory Which side of the chicken has more feathers?",
					Punchline: "The outside.",
				},
			},
		}
	}
}

func writeToSource() {
	if *operationMode == "file" {
		file, _ := json.MarshalIndent(articles, "", " ")
		_ = ioutil.WriteFile("jokes/jokesfile.json", file, 0644)
	} else if *operationMode == "db" {
		fmt.Println("Connecting to the DB")
	}
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	// operationMode := &mode
	operationMode = flag.String("use", "dryrun", "Options : file / db")
	flag.Parse()
	fmt.Println("Operation mode", *operationMode)
	readFromSource()
	handleRequests()
}
