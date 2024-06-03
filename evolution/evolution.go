package evolution

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Up struct {
	NationalId int    `json:"nationalId"`
	Name       string `json:"name"`
	Method     string `json:"method"`
	Level      int    `json:"level"`
}

type Evolution struct {
	From []Up   `json:"from"`
	To   []Up   `json:"to"`
	ID   string `json:"_id"`
	Rev  string `json:"_rev"`
}

type InputData struct {
	Docs []Evolution `json:"docs"`
	Seq  int         `json:"seq"`
}

func Crawl() {
	url := "https://pokedex.org/assets/evolutions.txt"

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

	var allEvolutions []Evolution

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

		allEvolutions = append(allEvolutions, inputData.Docs...)
	}

	allEvolutionsJSON, err := json.MarshalIndent(allEvolutions, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal all evolutions to JSON: %s", err)
	}

	filename := "data/evolution.json"
	err = ioutil.WriteFile(filename, allEvolutionsJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write all evolutions to file: %s\nError: %s", filename, err)
	}

	fmt.Println("All evolutions have been saved to a single JSON file.")
}
