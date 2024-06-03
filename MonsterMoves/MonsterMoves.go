package MonsterMoves

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type MonsterMove struct {
	Move []struct {
		LearnType string `json:"learn_type"`
		Level     int    `json:"level"`
		Id        int    `json:"id"`
	} `json:"moves"`
	ID  string `json:"_id"`
	Rev string `json:"_rev"`
}

type InputData struct {
	Docs []MonsterMove `json:"docs"`
	Seq  int           `json:"seq"`
}

func Crawl() {
	var MonsterMoves []MonsterMove

	for i := 1; i <= 3; i++ {
		url := fmt.Sprintf("https://pokedex.org/assets/monster-moves-%d.txt", i)

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
				log.Printf("Failed to unmarshal part: %s\nError: %s", part, err)
				continue
			}

			// Append each move to the MonsterMoves slice
			MonsterMoves = append(MonsterMoves, inputData.Docs...)
		}
	}

	// Marshal MonsterMoves to JSON and write to a single file
	MonsterMovesJSON, err := json.MarshalIndent(MonsterMoves, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal all moves to JSON: %s", err)
	}

	filename := "data/monsterMoves.json"
	err = ioutil.WriteFile(filename, MonsterMovesJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write all moves to file: %s\nError: %s", filename, err)
	}

	fmt.Println("All moves have been saved to a single JSON file.")
}
