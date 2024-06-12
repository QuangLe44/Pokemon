package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type MonsterType struct {
	Type       string `json:"type"`
	Multiplier string `json:"multiplier"`
}

type Mult struct {
	ID           int           `json:"id"`
	MonsterTypes []MonsterType `json:"monster_types"`
}

type Move struct {
	LearnType string `json:"learn_type"`
	Level     int    `json:"level"`
	ID        int    `json:"id"`
}

type MonsterMoves struct {
	Moves []Move `json:"moves"`
	ID    string `json:"_id"`
}

type MoveInfo struct {
	TypeName    string      `json:"type_name"`
	Identifier  string      `json:"identifier"`
	Power       interface{} `json:"power"`
	PP          interface{} `json:"pp"`
	Accuracy    interface{} `json:"accuracy"`
	Description string      `json:"description"`
	Name        string      `json:"name"`
	ID          string      `json:"_id"`
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

func loadType(filename string) (map[int]Mult, error) {
	var PokeType []Mult
	typeMap := make(map[int]Mult)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &PokeType)
	if err != nil {
		return nil, err
	}

	for _, info := range PokeType {
		typeMap[info.ID] = info
	}

	return typeMap, nil
}

func loadMonsterMove(filename string) (map[int]MonsterMoves, error) {
	var moves []MonsterMoves
	moveMap := make(map[int]MonsterMoves)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &moves)
	if err != nil {
		return nil, err
	}

	for _, info := range moves {
		id, _ := strconv.Atoi(info.ID)
		moveMap[id] = info
	}

	return moveMap, nil
}

func loadMoveInfo(filename string) (map[int]MoveInfo, error) {
	var moveInfo []MoveInfo
	moveInfoMap := make(map[int]MoveInfo)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &moveInfo)
	if err != nil {
		return nil, err
	}

	for _, move := range moveInfo {
		id, _ := strconv.Atoi(move.ID)
		moveInfoMap[id] = move
	}

	return moveInfoMap, nil
}

func findPokemonByName(pokemons []Pokemon, name string) *Pokemon {
	for _, pokemon := range pokemons {
		if strings.ToLower(pokemon.Name) == strings.ToLower(name) {
			return &pokemon
		}
	}
	return nil
}

func displayPokemonInfo(pokemon *Pokemon, info *AdditionalInfo, desc *Description, evolution *Evolution, mult *Mult, monstermove *MonsterMoves, moveInfo map[int]MoveInfo) {
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

	if mult != nil {
		fmt.Printf("\nWhen attacked:\n")
		for _, mt := range mult.MonsterTypes {
			fmt.Printf("  - Type: %s, Multiplier: %s\n", mt.Type, mt.Multiplier)
		}
	}

	if monstermove != nil {
		fmt.Printf("\nMonster moves:\n")
		fmt.Println("Moves:")
		for _, move := range monstermove.Moves {
			moveInfo, exists := moveInfo[move.ID]
			if exists {
				fmt.Printf("  - Learn Type: %s, Level: %d\n", move.LearnType, move.Level)
				fmt.Printf("  - Name: %s\n", moveInfo.Name)
				fmt.Printf("  - Type: %s\n", moveInfo.TypeName)
				fmt.Printf("  - Power: %v\n", moveInfo.Power)
				fmt.Printf("  - PP: %v\n", moveInfo.PP)
				fmt.Printf("  - Accuracy: %v\n", moveInfo.Accuracy)
				fmt.Printf("  - Description: %s\n", moveInfo.Description)
				fmt.Println()
			}

		}
	}
}

// func main() {
// 	pokemonFile := "../data/baseInfo.json"
// 	pokemons, err := loadPokemonData(pokemonFile)
// 	if err != nil {
// 		log.Fatalf("Failed to load Pokemon data: %s", err)
// 	}

// 	additionalInfoFile := "../data/stats.json"
// 	additionalInfo, err := loadAdditionalInfo(additionalInfoFile)
// 	if err != nil {
// 		log.Fatalf("Failed to load additional information: %s", err)
// 	}

// 	descriptionFile := "../data/MonsterDescription.json"
// 	description, err := loadDescription(descriptionFile)
// 	if err != nil {
// 		log.Fatalf("Failed to load additional information: %s", err)
// 	}

// 	evoFile := "../data/evolution.json"
// 	evo, err := loadEvo(evoFile)
// 	if err != nil {
// 		log.Fatalf("Failed to load additional information: %s", err)
// 	}

// 	typeFile := "../data/MonsterType.json"
// 	mult, err := loadType(typeFile)
// 	if err != nil {
// 		log.Fatalf("Failed to load additional information: %s", err)
// 	}

// 	monsterMovesFile := "../data/monsterMoves.json"
// 	monsterMove, err := loadMonsterMove(monsterMovesFile)
// 	if err != nil {
// 		log.Fatalf("Failed to load additional information: %s", err)
// 	}

// 	moveInfoFIle := "../data/moves.json"
// 	moveInfo, err := loadMoveInfo(moveInfoFIle)
// 	if err != nil {
// 		fmt.Println("Error loading move details:", err)
// 		return
// 	}

// 	var name string
// 	fmt.Print("Enter Pokemon name: ")
// 	fmt.Scanln(&name)

// 	pokemon := findPokemonByName(pokemons, name)
// 	if pokemon == nil {
// 		fmt.Printf("Pokemon with name '%s' not found.\n", name)
// 		return
// 	}

// 	info := additionalInfo[pokemon.NationalID]
// 	desc := description[pokemon.NationalID]
// 	evolution := evo[pokemon.NationalID]
// 	multiplier := mult[pokemon.NationalID]
// 	monstermoves := monsterMove[pokemon.NationalID]

// 	displayPokemonInfo(pokemon, &info, &desc, &evolution, &multiplier, &monstermoves, moveInfo)
// 	// Start server
// // 	pokedex := make([]*Pokemon, 100)
// // 	for i := range pokedex {
// // 		pokedex[i] = &Pokemon{id: i}
// // 	}
// // 	s, err := newServer(":8080", 1000, pokedex)
// // 	if err != nil {
// // 		fmt.Println(err)
// // 		os.Exit(1)
// // 	}
// // 	s.Start()
// // 	sigChan := make(chan os.Signal, 1)
// // 	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
// // 	<-sigChan

// // 	fmt.Println("Shutting down server...")
// // 	s.Stop()
// // 	fmt.Println("Server stopped.")
// // }
