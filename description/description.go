package description

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type MonsterDescription struct {
	Description string `json:"description"`
	ID          string `json:"_id"`
	Rev         string `json:"_rev"`
}

type InputData struct {
	Docs []MonsterDescription `json:"docs"`
	Seq  int                  `json:"seq"`
}

func Crawl() {
	var allMonsterDescription []MonsterDescription
	for i := 1; i <= 3; i++ {
		url := fmt.Sprintf("https://pokedex.org/assets/descriptions-%d.txt", i)

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
			allMonsterDescription = append(allMonsterDescription, inputData.Docs...)
		}
	}

	allMonsterDescriptionJSON, err := json.MarshalIndent(allMonsterDescription, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal all MonsterDescription to JSON: %s", err)
	}

	filename := "data/MonsterDescription.json"
	err = ioutil.WriteFile(filename, allMonsterDescriptionJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write all MonsterDescription to file: %s\nError: %s", filename, err)
	}

	fmt.Println("All info have been saved to MonsterDescription.json.")
}
