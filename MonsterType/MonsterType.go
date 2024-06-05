package MonsterType

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PokemonInfo struct {
	ID           int           `json:"id"`
	MonsterTypes []MonsterType `json:"monster_types"`
}

type MonsterType struct {
	Type       string `json:"type"`
	Multiplier string `json:"multiplier"`
}

func Crawl() {
	pokemonInfos := make([]PokemonInfo, 0)

	for i := 1; i <= 649; i++ {
		url := fmt.Sprintf("https://pokedex.org/#/pokemon/%d", i)
		doc, err := goquery.NewDocument(url)
		if err != nil {
			log.Fatalf("Error fetching URL %s: %v", url, err)
		}

		monsterTypes := make([]MonsterType, 0)

		doc.Find(".when-attacked-row").Each(func(i int, s *goquery.Selection) {
			s.Find(".monster-type").Each(func(i int, s *goquery.Selection) {
				mType := strings.TrimSpace(s.Text())
				multiplier := strings.TrimSpace(s.Next().Text())
				monsterTypes = append(monsterTypes, MonsterType{
					Type:       mType,
					Multiplier: multiplier,
				})
			})
		})

		pokemonInfo := PokemonInfo{
			ID:           i,
			MonsterTypes: monsterTypes,
		}

		pokemonInfos = append(pokemonInfos, pokemonInfo)
	}

	jsonData, err := json.MarshalIndent(pokemonInfos, "", "")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	file, err := os.Create("data/MonsterType.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Println("All info have been saved to MonsterType.json.")
}
