package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
	Experience      int             `json:"experience,omitempty"`
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
	LearnType string   `json:"learn_type"`
	Level     int      `json:"level"`
	ID        int      `json:"id"`
	Details   MoveInfo `json:"details"`
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

type Experience struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Exp  string `json:"exp"`
}

type PokemonInfo struct {
	Pokemon      *Pokemon        `json:"pokemon"`
	Additional   *AdditionalInfo `json:"additional_info"`
	Description  *Description    `json:"description"`
	Evolution    *Evolution      `json:"evolution"`
	TypeInfo     *Mult           `json:"type_info"`
	MonsterMoves *MonsterMoves   `json:"monster_moves"`
	Experience   *Experience     `json:"experience"`
}

func LoadPokemonData(filename string) ([]Pokemon, error) {
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

func loadExperience(filename string) (map[int]Experience, error) {
	var exps []Experience
	expMap := make(map[int]Experience)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &exps)
	if err != nil {
		return nil, err
	}

	for _, exp := range exps {
		id, _ := strconv.Atoi(exp.ID)
		expMap[id] = exp
	}

	return expMap, nil
}

func Pokedex(pokemon *Pokemon, info *AdditionalInfo, desc *Description, evolution *Evolution, mult *Mult, monstermove *MonsterMoves, moveInfo map[int]MoveInfo, exp *Experience) PokemonInfo {
	if monstermove != nil {
		for i, move := range monstermove.Moves {
			if details, exists := moveInfo[move.ID]; exists {
				monstermove.Moves[i].Details = details
			}
		}
	}

	return PokemonInfo{
		Pokemon:      pokemon,
		Additional:   info,
		Description:  desc,
		Evolution:    evolution,
		TypeInfo:     mult,
		MonsterMoves: monstermove,
		Experience:   exp,
	}
}

func main() {
	pokemonFile := "data/baseInfo.json"
	pokemons, err := LoadPokemonData(pokemonFile)
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

	typeFile := "data/MonsterType.json"
	mult, err := loadType(typeFile)
	if err != nil {
		log.Fatalf("Failed to load additional information: %s", err)
	}

	monsterMovesFile := "data/monsterMoves.json"
	monsterMove, err := loadMonsterMove(monsterMovesFile)
	if err != nil {
		log.Fatalf("Failed to load additional information: %s", err)
	}

	moveInfoFile := "data/moves.json"
	moveInfo, err := loadMoveInfo(moveInfoFile)
	if err != nil {
		fmt.Println("Error loading move details:", err)
		return
	}

	experienceFile := "data/exp.json"
	experience, err := loadExperience(experienceFile)
	if err != nil {
		log.Fatalf("Failed to load experience information: %s", err)
	}

	allPokemonInfo := []PokemonInfo{}
	for _, pokemon := range pokemons {
		info := additionalInfo[pokemon.NationalID]
		desc := description[pokemon.NationalID]
		evolution := evo[pokemon.NationalID]
		multiplier := mult[pokemon.NationalID]
		monstermoves := monsterMove[pokemon.NationalID]
		exp := experience[pokemon.NationalID]

		pokemonInfo := Pokedex(&pokemon, &info, &desc, &evolution, &multiplier, &monstermoves, moveInfo, &exp)
		allPokemonInfo = append(allPokemonInfo, pokemonInfo)
	}

	jsonData, err := json.MarshalIndent(allPokemonInfo, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling to JSON: %s\n", err)
		return
	}

	err = ioutil.WriteFile("data/pokedex.json", jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON to file: %s\n", err)
		return
	}

	fmt.Println("All Pokemon information has been written to pokedex.json")
}
