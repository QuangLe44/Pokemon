package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type ListMapObject struct {
	Name string `json:"name"`
}

type Pokemon struct {
	Descriptions    []ListMapObject `json:"descriptions"`
	Types           []ListMapObject `json:"types"`
	Abilities       []ListMapObject `json:"abilities"`
	Attack          int             `json:"attack"`
	Defense         int             `json:"defense"`
	Speed           int             `json:"speed"`
	SpAtk           int             `json:"sp_atk"`
	SpDef           int             `json:"sp_def"`
	HP              int             `json:"hp"`
	Weight          string          `json:"weight"`
	Height          string          `json:"height"`
	NationalID      int             `json:"national_id"`
	MaleFemaleRatio string          `json:"male_female_ratio"`
	CatchRate       int             `json:"catch_rate"`
	ID              string          `json:"_id"`
	Name            string          `json:"name"`
}

type AdditionalInfo struct {
	SpecialAttackEV  int    `json:"specialAttackEV"`
	HPEV             int    `json:"hpEV"`
	DefenseEV        int    `json:"defenseEV"`
	AttackEV         int    `json:"attackEV"`
	SpecialDefenseEV int    `json:"specialDefenseEV"`
	SpeedEV          int    `json:"speedEV"`
	HatchSteps       int    `json:"hatchSteps"`
	Species          string `json:"species"`
	EggGroups        string `json:"eggGroups"`
	ID               string `json:"_id"`
}

type Description struct {
	Description string `json:"description"`
}

type Evolution struct {
	From []EvolutionDetail `json:"from"`
	To   []EvolutionDetail `json:"to"`
	ID   string            `json:"_id"`
	Rev  string            `json:"_rev"`
}

type EvolutionDetail struct {
	NationalID int    `json:"nationalId"`
	Name       string `json:"name"`
	Method     string `json:"method"`
	Level      int    `json:"level"`
}

func (e *Evolution) UnmarshalJSON(data []byte) error {
	type Alias Evolution
	aux := &struct {
		From json.RawMessage `json:"from"`
		To   json.RawMessage `json:"to"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.From != nil {
		if err := json.Unmarshal(aux.From, &e.From); err != nil {
			return err
		}
	} else {
		e.From = []EvolutionDetail{}
	}

	if aux.To != nil {
		if err := json.Unmarshal(aux.To, &e.To); err != nil {
			return err
		}
	} else {
		e.To = []EvolutionDetail{}
	}

	return nil
}

func loadPokemonData(filename string) ([]Pokemon, error) {
	var pokemons []Pokemon

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &pokemons)
	if err != nil {
		return nil, err
	}

	return pokemons, nil
}

func loadAdditionalInfo(filename string) (map[int]AdditionalInfo, error) {
	var infos []AdditionalInfo
	infoMap := make(map[int]AdditionalInfo)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &infos)
	if err != nil {
		return nil, err
	}

	for _, info := range infos {
		id, _ := strconv.Atoi(info.ID)
		infoMap[id] = info
	}

	return infoMap, nil
}

func loadDescription(filename string) (map[int]Description, error) {
	var desc []Description
	descMap := make(map[int]Description)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &desc)
	if err != nil {
		return nil, err
	}

	for index, info := range desc {
		descMap[index+1] = info
	}

	return descMap, nil
}

func loadEvo(filename string) (map[int]Evolution, error) {
	var evo []Evolution
	evoMap := make(map[int]Evolution)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &evo)
	if err != nil {
		return nil, err
	}

	for _, info := range evo {
		id, _ := strconv.Atoi(info.ID)
		evoMap[id] = info
	}

	return evoMap, nil
}

func findPokemonByName(pokemons []Pokemon, name string) *Pokemon {
	for _, pokemon := range pokemons {
		if strings.ToLower(pokemon.Name) == strings.ToLower(name) {
			return &pokemon
		}
	}
	return nil
}

func displayPokemonInfo(pokemon *Pokemon, info *AdditionalInfo, desc *Description, evolution *Evolution) {
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("National ID: %d\n", pokemon.NationalID)
	fmt.Printf("Height: %s\n", pokemon.Height)
	fmt.Printf("Weight: %s\n", pokemon.Weight)
	fmt.Printf("HP: %d\n", pokemon.HP)
	fmt.Printf("Attack: %d\n", pokemon.Attack)
	fmt.Printf("Defense: %d\n", pokemon.Defense)
	fmt.Printf("Speed: %d\n", pokemon.Speed)
	fmt.Printf("Sp. Atk: %d\n", pokemon.SpAtk)
	fmt.Printf("Sp. Def: %d\n", pokemon.SpDef)
	fmt.Printf("Male/Female Ratio: %s\n", pokemon.MaleFemaleRatio)
	fmt.Printf("Catch Rate: %d\n", pokemon.CatchRate)

	fmt.Println("Descriptions:")
	for _, description := range pokemon.Descriptions {
		fmt.Printf("  - %s\n", description.Name)
	}

	fmt.Println("Types:")
	for _, typ := range pokemon.Types {
		fmt.Printf("  - %s\n", typ.Name)
	}

	fmt.Println("Abilities:")
	for _, ability := range pokemon.Abilities {
		fmt.Printf("  - %s\n", ability.Name)
	}

	if info != nil {
		fmt.Printf("\nAdditional Information:\n")
		fmt.Printf("Special Attack EV: %d\n", info.SpecialAttackEV)
		fmt.Printf("HP EV: %d\n", info.HPEV)
		fmt.Printf("Defense EV: %d\n", info.DefenseEV)
		fmt.Printf("Attack EV: %d\n", info.AttackEV)
		fmt.Printf("Special Defense EV: %d\n", info.SpecialDefenseEV)
		fmt.Printf("Speed EV: %d\n", info.SpeedEV)
		fmt.Printf("Hatch Steps: %d\n", info.HatchSteps)
		fmt.Printf("Species: %s\n", info.Species)
		fmt.Printf("Egg Groups: %s\n", info.EggGroups)
	}

	if desc != nil {
		fmt.Printf("Description: %s\n", desc.Description)
	}

	if evolution != nil {
		fmt.Printf("Evolution:\n")
		fmt.Printf("From:\n")
		if len(evolution.From) == 0 {
			fmt.Println("Pokemon does not evolve from anything")
		} else {
			for _, detail := range evolution.From {
				fmt.Println("Name:", detail.Name)
				fmt.Println("level:", detail.Level)
			}
		}
		fmt.Printf("To:\n")
		if len(evolution.To) == 0 {
			fmt.Println("Pokemon has no further evolution")
		} else {
			for _, detail := range evolution.To {
				fmt.Println("Name:", detail.Name)
				fmt.Println("level:", detail.Level)
			}
		}
	}
}

func main() {
	pokemonFile := "data/baseInfo.json"
	pokemons, err := loadPokemonData(pokemonFile)
	if err != nil {
		log.Fatalf("Failed to load Pokemon data: %s", err)
	}

	additionalInfoFile := "data/stats.json"
	additionalInfo, err := loadAdditionalInfo(additionalInfoFile)
	if err != nil {
		log.Fatalf("Failed to load additional information: %s", err)
	}

	descriptionFile := "data/MonsterDescription.json"
	description, err := loadDescription(descriptionFile)
	if err != nil {
		log.Fatalf("Failed to load additional information: %s", err)
	}

	evoFile := "data/evolution.json"
	evo, err := loadEvo(evoFile)
	if err != nil {
		log.Fatalf("Failed to load additional information: %s", err)
	}

	var name string
	fmt.Print("Enter Pokemon name: ")
	fmt.Scanln(&name)

	pokemon := findPokemonByName(pokemons, name)
	if pokemon == nil {
		fmt.Printf("Pokemon with name '%s' not found.\n", name)
		return
	}

	info := additionalInfo[pokemon.NationalID]
	desc := description[pokemon.NationalID]
	evolution := evo[pokemon.NationalID]

	displayPokemonInfo(pokemon, &info, &desc, &evolution)
}
