package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
)

var (
	url  = "https://catfact.ninja/fact"
	url2 = "https://catfact.ninja/breeds"
)

type CatFact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

type DataURL struct {
	URL string
}

type CatData struct {
	Breed   string `json:"breed"`
	Country string `json:"country"`
	Origin  string `json:"origin"`
	Coat    string `json:"coat"`
	Pattern string `json:"pattern"`
}

type ApiBreedResponse struct {
	Data []CatData
}

type BreedGroup struct {
	Country string   `json:"country"`
	Breeds  []string `json:"breeds"`
}

func main() {
	log.Println("Hello catfact docker")
	templates := template.Must(
		template.ParseFiles(
			"templates/index.html",
			"templates/cats.html",
			"templates/facts.html",
			"templates/sorted.html",
		))

	http.HandleFunc("/", handleIndex(templates, url))
	http.HandleFunc("/facts", handleReqToFacts(templates, url))
	http.HandleFunc("/cats", handleReqToCats(templates, url2))
	http.HandleFunc("/sort", handleSorted(templates, url2))

	http.ListenAndServe(":5500", nil)
}

func handleIndex(templates *template.Template, url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "<h1> Hello CatFacts</h1>")
		data := &DataURL{
			URL: url,
		}
		if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleSorted(templates *template.Template, url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data := getSortCatData(url)

		if err := templates.ExecuteTemplate(w, "sorted.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleReqToCats(templates *template.Template, url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		resp := getCatData(url)
		if err := templates.ExecuteTemplate(w, "cats.html", resp.Data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleReqToFacts(templates *template.Template, url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// ticker := time.NewTicker(2 * time.Second)
		// defer ticker.Stop()

		// for {

		fact := getCatFact(url)
		if err := templates.ExecuteTemplate(w, "facts.html", fact); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//fmt.Fprintf(w, fact)

		fmt.Println(fact)

		// 	<-ticker.C
		// }

		// go func() {
		// 	ticker := time.NewTicker(2 * time.Second)
		// 	defer ticker.Stop()

		// 	for range ticker.C {
		// 		fact := getCatFact(url)
		// 		fmt.Fprintf(w, fact.Fact)
		// 	}
		// }()

	}
}

func getCatFact(url string) CatFact {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var fact CatFact
	json.Unmarshal(body, &fact)
	return fact
}

func getCatData(url string) ApiBreedResponse {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var catData ApiBreedResponse

	json.Unmarshal(body, &catData)
	return catData
}

func getSortCatData(url string) []BreedGroup {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var catData ApiBreedResponse

	json.Unmarshal(body, &catData)

	//sorting
	breedByCountry := make(map[string][]string)

	for _, cat := range catData.Data {
		breedByCountry[cat.Country] = append(breedByCountry[cat.Country], cat.Breed)
	}

	var breedgroups []BreedGroup

	for country, breeeds := range breedByCountry {
		group := BreedGroup{
			Country: country,
			Breeds:  breeeds,
		}
		breedgroups = append(breedgroups, group)
	}

	return breedgroups
}
