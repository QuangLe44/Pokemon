package exp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func Remove() {
	file, err := os.Open("data/exp.json")
	if err != nil {
		log.Fatalf("Failed to open JSON file: %v", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	var pokemon []exp
	if err := json.Unmarshal(byteValue, &pokemon); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	seen := make(map[string]bool)
	var uniquePokemon []exp

	for _, pokemon := range pokemon {
		if !seen[pokemon.Name] {
			uniquePokemon = append(uniquePokemon, pokemon)
			seen[pokemon.Name] = true
		}
	}

	JsonFile, err := json.MarshalIndent(uniquePokemon, "", "    ")
	if err != nil {
		log.Fatalf("Failed to marshal unique pokemon to JSON: %v", err)
	}

	if err := ioutil.WriteFile("data/exp.json", JsonFile, 0644); err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	fmt.Println("Duplicates removed.")
}
