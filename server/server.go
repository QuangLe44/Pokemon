package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Player struct {
	Name     string    `json:"name"`
	Pokemons []Pokemon `json:"pokemons"`
	Active   []int
	Health   []int
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
	players = append(players, player)
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

			for n, choiceStr := range choices {
				choice, err := strconv.Atoi(strings.TrimSpace(choiceStr))
				if err != nil {
					fmt.Println("Error converting choice to integer:", err)
					continue
				}
				player.Active = append(player.Active, choice)
				for _, pokemon := range player.Pokemons {
					pokemonID, _ := strconv.Atoi(pokemon.ID)
					if pokemonID == player.Active[n] {
						player.Health = append(player.Health, pokemon.HP)
					}
				}
			}
			fmt.Println("Player's active choices:", player.Active)
		} else if strings.TrimSpace(message) == "ready" {
			player.Ready = true
			break
		}
	}
	for {
		allReady := true
		for _, currentPlayer := range players {
			if !currentPlayer.Ready {
				mu.Lock()
				allReady = false
				mu.Unlock()
			}
		}
		if allReady {
			startBattle()
			break
		} else {
			continue
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
		processAction(currentPlayer, opponentPlayer, action)
		currentTurnIndex = 2 - currentTurnIndex + 1
	}
}
func processAction(currentPlayer, opponentPlayer *Player, action string) {
	switch strings.ToLower(action) {
	case "attack":
		rand.Seed(time.Now().UnixNano())
		damage := rand.Intn(10) + 1 // Random damage between 1 and 10
		opponentPlayer.Health[0] -= damage
		fmt.Printf("%s attacked and dealt %d damage to %s's first Pokémon. Remaining HP: %d\n", currentPlayer.Name, damage, opponentPlayer.Name, opponentPlayer.Health[0])

		// Notify players
		currentPlayer.Conn.Write([]byte(fmt.Sprintf("You attacked and dealt %d damage. Opponent's Pokémon remaining HP: %d\n", damage, opponentPlayer.Health[0])))
		opponentPlayer.Conn.Write([]byte(fmt.Sprintf("Opponent attacked and dealt %d damage to your Pokémon. Remaining HP: %d\n", damage, opponentPlayer.Health[0])))

		// Check if opponent's Pokémon fainted
		if opponentPlayer.Health[0] <= 0 {
			currentPlayer.Conn.Write([]byte("Opponent's Pokémon fainted!\n"))
			opponentPlayer.Conn.Write([]byte("Your Pokémon fainted!\n"))
			// Handle Pokémon fainting (switch to next Pokémon or end battle)
			// For simplicity, we'll end the battle here if the first Pokémon faints
			currentPlayer.Conn.Write([]byte("You win!\n"))
			opponentPlayer.Conn.Write([]byte("You lose!\n"))
			currentPlayer.Conn.Close()
			opponentPlayer.Conn.Close()
		}
	case "switch":
		// Handle Pokémon switching logic
	case "forfeit":
		// Handle player forfeiting the match
	}
}
