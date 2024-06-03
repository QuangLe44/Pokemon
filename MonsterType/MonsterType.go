package MonsterType

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type When struct {
	/*
		"multiplier":0.5,"name":"fire"
	*/
	Multiplier interface{} `json:"multiplier"`
	Name       string      `json:"name"`
}

type MonsterType struct {
	WhenDefending []When `json:"whenDefending"`
	WhenAttacking []When `json:"whenAttacking"`
	ID            string `json:"_id"`
	Rev           string `json:"_rev"`
}

// InputData represents the structure of the input text file.
type InputData struct {
	Docs []MonsterType `json:"docs"`
	Seq  int           `json:"seq"`
}

func Crawl() {
	url := "https://pokedex.org/assets/types.txt"

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

	var allMonsterType []MonsterType

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

		allMonsterType = append(allMonsterType, inputData.Docs...)
	}

	allMonsterTypeJSON, err := json.MarshalIndent(allMonsterType, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal all MonsterType to JSON: %s", err)
	}

	filename := "data/MonsterType.json"

	err = ioutil.WriteFile(filename, allMonsterTypeJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write all MonsterType to file: %s\nError: %s", filename, err)
	}

	fmt.Println("All MonsterType have been saved to a single JSON file.")
}
