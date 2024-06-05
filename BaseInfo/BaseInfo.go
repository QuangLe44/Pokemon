package BaseInfo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type list struct {
	Name        string `json:"name"`
	ResourceURI string `json:"resource_uri"`
}

type baseInfo struct {
	Description     []list `json:"descriptions"`
	Type            []list `json:"types"`
	Abilities       []list `json:"abilities"`
	Attack          int    `json:"attack"`
	Defense         int    `json:"defense"`
	Speed           int    `json:"speed"`
	SpAtk           int    `json:"sp_atk"`
	SpDef           int    `json:"sp_def"`
	HP              int    `json:"hp"`
	Weight          string `json:"weight"`
	Height          string `json:"height"`
	NationalID      int    `json:"national_id"`
	MaleFemaleRatio string `json:"male_female_ratio"`
	CatchRate       int    `json:"catch_rate"`
	ID              string `json:"_id"`
	Name            string `json:"name"`
}

type InputData struct {
	Docs []baseInfo `json:"docs"`
	Seq  int        `json:"seq"`
}

func Crawl() {
	var allbaseInfo []baseInfo

	for i := 1; i <= 3; i++ {
		url := fmt.Sprintf("https://pokedex.org/assets/skim-monsters-%d.txt", i)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalf("Failed to create request: %s", err)
		}

		req.Header.Set("Referer", "https://pokedex.org/js/worker.js")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to send request: %s", err)
		}
		defer resp.Body.Close()

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %s", err)
		}

		parts := strings.Split(string(content), "\n")

		for _, part := range parts {
			if strings.TrimSpace(part) == "" {
				continue
			}

			var inputData InputData
			err := json.Unmarshal([]byte(part), &inputData)
			if err != nil {
				log.Printf("Failed to unmarshal: %s\nError: %s", part, err)
				continue
			}

			for _, monster := range inputData.Docs {
				name, err := strconv.Atoi(monster.ID)
				if err != nil {
					log.Printf("Failed to convert monster ID to int value: %s\nError: %s", monster.ID, err)
					continue
				}

				monster.ID = strconv.Itoa(name)
				allbaseInfo = append(allbaseInfo, monster)
			}
		}
	}

	InfoFile, err := json.MarshalIndent(allbaseInfo, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal all baseInfo to JSON: %s", err)
	}

	filename := "data/baseInfo.json"
	err = ioutil.WriteFile(filename, InfoFile, 0644)
	if err != nil {
		log.Fatalf("Failed to write all baseInfo to file: %s\nError: %s", filename, err)
	}

	fmt.Println("All Info have been saved to baseInfo.json.")
}
