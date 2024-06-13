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
	Name     string    `json:"name"`
	Pokemons []Pokemon `json:"pokemons"`
	Active   []int
	Ready    bool
	Turn     int
	Conn     net.Conn
}

type Battle struct {
	Players     []*Player
	Turn        int
	WinnerIndex int
}

var (
	players []*Player
	mu      sync.Mutex
	wg      sync.WaitGroup
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
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	player := &Player{}
	player.Conn = conn
	mu.Lock()
	players = append(players, player)
	playerIndex := len(players) - 1
	mu.Unlock()
	defer func() {
		mu.Lock()
		defer mu.Unlock()
		// Remove the player from the players slice when the connection is closed
		players = append(players[:playerIndex], players[playerIndex+1:]...)
	}()

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
			// for _, pokemonID := range player.Pokemons {
			// 	idStr := strconv.Itoa(pokemonID)
			// 	if pokemon, found := pokemonMap[idStr]; found {
			// 		fmt.Printf("  - %s\n", pokemon.Name)
			// 	} else {
			// 		fmt.Printf("  - Unknown Pokemon ID: %d\n", pokemonID)
			// 	}
			// }
			pokemonListMessage := "\nChoose your pokemons:\n"
			index := 0
			for _, pokemon := range player.Pokemons {
				index++
				pokemonListMessage += fmt.Sprintf("%s. %s\n", pokemon.ID, pokemon.Name)
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
				go startBattle()
			}
			mu.Unlock()
		}
	}
}

func startBattle() {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("Battle start")
	var player1Speed, player2Speed int
	for n, player := range players {
		if player == nil {
			fmt.Printf("Player %d is nil\n", n)
			continue
		}
		if len(player.Active) == 0 {
			fmt.Printf("Player %d has no active Pokémon\n", n)
			continue
		}
		firstActiveID := strconv.Itoa(player.Active[0])
		for _, pokemon := range player.Pokemons {
			if pokemon.ID == firstActiveID {
				switch n {
				case 0:
					fmt.Printf("%s's first Pokémon: %s. Its speed is: %d\n", player.Name, pokemon.Name, pokemon.Speed)
					player1Speed = pokemon.Speed
				case 1:
					fmt.Printf("%s's first Pokémon: %s. Its speed is: %d\n", player.Name, pokemon.Name, pokemon.Speed)
					player2Speed = pokemon.Speed
				}
			}
		}
	}
	if player1Speed > player2Speed || player1Speed == player2Speed {
		players[0].Turn = 1
		players[1].Turn = 2
	} else {
		players[0].Turn = 2
		players[1].Turn = 1
	}
	currentTurnIndex := 1
	for {
		currentPlayer := players[currentTurnIndex-1]
		if currentPlayer == nil {
			fmt.Printf("Current player %d is nil\n", currentTurnIndex-1)
			continue
		}
		opponentPlayer := players[2-currentTurnIndex]
		if opponentPlayer == nil {
			fmt.Printf("Opponent player %d is nil\n", 2-currentTurnIndex)
			continue
		}

		// Send turn message to current player
		_, err := currentPlayer.Conn.Write([]byte("Your turn! Choose an action:\nSwitch: {pokemon ID}\nAttack\nForfeit\n"))
		if err != nil {
			fmt.Println("Error writing to current player:", err)
			return
		}

		// Read player's response
		buffer := make([]byte, 1024)
		n, err := currentPlayer.Conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from current player:", err)
			return
		}
		action := strings.TrimSpace(string(buffer[:n]))

		// Process action based on the chosen action
		processAction(currentPlayer, opponentPlayer, action)

		// Switch turns
		currentTurnIndex = 2 - currentTurnIndex + 1
	}
}

func processAction(currentPlayer, opponentPlayer *Player, action string) {
	switch action {
	case "Attack":
		// Implement attack logic here
		fmt.Printf("%s is attacking!\n", currentPlayer.Name)
		// Example: Decrease opponent's health, check for knockouts, etc.
	default:
		fmt.Printf("Processing action %s for player %s\n", action, currentPlayer.Name)
	}
}
