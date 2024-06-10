package exp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type exp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Exp  string `json:"exp"`
}

func Crawl() {
	res, err := http.Get("https://bulbapedia.bulbagarden.net/wiki/List_of_Pok%C3%A9mon_by_effort_value_yield_(Generation_IX)")
	if err != nil {
		log.Fatal("Error fetching the URL: ", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Error loading the document: ", err)
	}

	count := 0
	const limit = 682

	baseExp := make([]exp, 0)

	doc.Find("table.sortable").Each(func(i int, table *goquery.Selection) {
		table.Find("tr").Each(func(i int, tr *goquery.Selection) {
			var ID, Name, Exp string
			tr.Find("td").Each(func(j int, td *goquery.Selection) {
				if count < limit {
					if j == 0 {
						ID = strings.TrimSpace(td.Text())
					}
					if j == 2 {
						Name = td.Find("a").AttrOr("title", "")
						Name = strings.TrimSuffix(Name, " (Pokémon)")
					}
					if j == 3 {
						Exp = strings.TrimSpace(td.Text())
					}
				}
			})
			if ID != "" && Name != "" && Exp != "" {
				baseExp = append(baseExp, exp{
					ID:   ID,
					Name: Name,
					Exp:  Exp,
				})
				count++
			}
		})
	})

	jsonData, err := json.MarshalIndent(baseExp, "", "")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	file, err := os.Create("data/exp.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Println("All exp values have been saved to a single JSON file.")
}
