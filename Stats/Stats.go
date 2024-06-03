package Stats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type MonsterStats struct {
	SpecialAttackEV  int         `json:"specialAttackEV"`
	HpEV             int         `json:"hpEV"`
	HatchSteps       int         `json:"hatchSteps"`
	DefenseEV        int         `json:"defenseEV"`
	AttackEV         int         `json:"attackEV"`
	SpecialDefenseEV int         `json:"specialDefenseEV"`
	SpeedEV          int         `json:"speedEV"`
	GenderRatio      interface{} `json:"genderRatio"`
	Species          string      `json:"species"`
	JapaneseName     string      `json:"japaneseName"`
	HepburnName      string      `json:"Guranburu"`
	EggGroups        string      `json:"eggGroups"`

	ID  string `json:"_id"`
	Rev string `json:"_rev"`
}

// InputData represents the structure of the input text file.
type InputData struct {
	Docs []MonsterStats `json:"docs"`
	Seq  int            `json:"seq"`
}

func Crawl() {
	var allStats []MonsterStats // Initialize a slice to hold all supplementals

	for i := 1; i <= 3; i++ {
		url := fmt.Sprintf("https://pokedex.org/assets/monsters-supplemental-%d.txt", i)

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

			// Append each supplemental to the allStats slice
			allStats = append(allStats, inputData.Docs...)
		}
	}

	// Marshal allStats to JSON and write to a single file
	allStatsJSON, err := json.MarshalIndent(allStats, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal all supplementals to JSON: %s", err)
	}

	filename := "data/stats.json"
	err = ioutil.WriteFile(filename, allStatsJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write all supplementals to file: %s\nError: %s", filename, err)
	}

	fmt.Println("All supplementals have been saved to a single JSON file.")
}
