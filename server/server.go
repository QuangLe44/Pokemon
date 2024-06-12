package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Player struct {
	Name     string `json:"name"`
	Pokemons []int  `json:"pokemons"`
	Active   []int
	Ready    bool
}

type Battle struct {
	Players     []*Player
	Turn        int
	WinnerIndex int
}

var (
	players     []*Player
	pokemonList []Pokemon
	mu          sync.Mutex
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	pokemonList, err := loadPokemonData("../data/baseInfo.json")
	if err != nil {
		fmt.Println("Error reading pokemon data:", err)
		return
	}

	pokemonMap := make(map[string]Pokemon)
	for _, pokemon := range pokemonList {
		pokemonMap[pokemon.ID] = pokemon
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		fmt.Println("Accepted new connection.")
		go handleClient(conn, pokemonMap)
	}
}

func handleClient(conn net.Conn, pokemonMap map[string]Pokemon) {
	defer conn.Close()
	player := &Player{}
	mu.Lock()
	players = append(players, player)
	mu.Unlock()

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		message := string(buffer[:n])
		if strings.HasPrefix(message, "Player name:") {
			player.Name = strings.TrimSpace(message[len("Player name: "):])
			fmt.Println("Player name is: " + player.Name)
			fmt.Println("Received player name:", player.Name)
			playerData, err := os.ReadFile("player.json")
			if err != nil {
				fmt.Println("Error reading player data file: ", err)
				os.Exit(1)
			}
			var allPlayers []*Player
			err = json.Unmarshal(playerData, &allPlayers)
			if err != nil {
				fmt.Println("Error unmarshaling player data:", err)
				return
			}
			var foundPlayer *Player
			for _, p := range allPlayers {
				if p.Name == player.Name {
					foundPlayer = p
					break
				}
			}
			if foundPlayer == nil {
				fmt.Println("Player not found in player data file:", player.Name)
				return
			}
			player.Pokemons = foundPlayer.Pokemons
			for _, pokemonID := range player.Pokemons {
				idStr := strconv.Itoa(pokemonID)
				if pokemon, found := pokemonMap[idStr]; found {
					fmt.Printf("  - %s\n", pokemon.Name)
				} else {
					fmt.Printf("  - Unknown Pokemon ID: %d\n", pokemonID)
				}
			}
			pokemonListMessage := "\nChoose your pokemons:\n"
			index := 0
			for _, pokemonID := range player.Pokemons {
				idStr := strconv.Itoa(pokemonID)
				if pokemon, found := pokemonMap[idStr]; found {
					index++
					pokemonListMessage += fmt.Sprintf("%d. %s\n", index, pokemon.Name)
				} else {
					pokemonListMessage += fmt.Sprintf("  - Unknown Pokemon ID: %d\n", pokemonID)
				}
			}
			fmt.Print(pokemonListMessage)
			conn.Write([]byte(pokemonListMessage))
			fmt.Println("Sent pokemon list to client")
		} else if strings.HasPrefix(message, "Player choice:") {
			choicesStr := strings.TrimSpace(message[len("Player choice:"):])
			choices := strings.Split(choicesStr, ",")

			for _, choiceStr := range choices {
				choice, err := strconv.Atoi(strings.TrimSpace(choiceStr))
				if err != nil {
					fmt.Println("Error converting choice to integer:", err)
					continue
				}
				player.Active = append(player.Active, choice)
			}
			fmt.Println("Player's active choices:", player.Active)
		} else if strings.TrimSpace(message) == "ready" {
			mu.Lock()
			player.Ready = true
			allReady := true
			for _, p := range players {
				if !p.Ready {
					allReady = false
					break
				}
			}
			if allReady {
				fmt.Println("All players are ready. Starting the game...")
				go startBattle(pokemonMap)
			}
			mu.Unlock()
		}
	}
}

func startBattle(pokemonMap map[string]Pokemon) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("Battle start")
	for _, player := range players {
		if len(player.Active) > 0 {
			firstPokemonID := player.Active[0]
			fmt.Printf("%s's first Pokémon: %s\n", player.Name, pokemonMap[strconv.Itoa(firstPokemonID)].Name)
		} else {
			fmt.Printf("No active Pokémon for player %s\n", player.Name)
		}
	}
}
