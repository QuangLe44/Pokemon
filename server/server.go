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
	player := &Player{Conn: conn}
	mu.Lock()
	players = append(players, player)
	mu.Unlock()

	playerData, err := os.ReadFile("player.json")
	if err != nil {
		fmt.Println("Error reading player data file:", err)
		return
	}
	var allPlayers []*Player
	err = json.Unmarshal(playerData, &allPlayers)
	if err != nil {
		fmt.Println("Error unmarshaling player data:", err)
		return
	}

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
			fmt.Println("Player name is:", player.Name)
			foundPlayer := findPlayerByName(allPlayers, player.Name)
			if foundPlayer == nil {
				fmt.Println("Player not found:", player.Name)
				return
			}
			player.Pokemons = foundPlayer.Pokemons

			pokemonListMessage := "\nChoose your pokemons:\n"
			for _, pokemon := range player.Pokemons {
				pokemonListMessage += fmt.Sprintf("%s. %s\n", pokemon.ID, pokemon.Name)
			}
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
				for _, pokemon := range player.Pokemons {
					pokemonID, _ := strconv.Atoi(pokemon.ID)
					if pokemonID == choice {
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

	waitForAllPlayersReady()
	startBattle()
}

func findPlayerByName(players []*Player, name string) *Player {
	for _, p := range players {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func waitForAllPlayersReady() {
	for {
		mu.Lock()
		allReady := true
		for _, currentPlayer := range players {
			if !currentPlayer.Ready {
				allReady = false
				break
			}
		}
		mu.Unlock()
		if allReady {
			return
		}
		time.Sleep(100 * time.Millisecond) // Reduce busy-waiting CPU usage
	}
}

func startBattle() {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("Battle start")

	assignInitialTurns()

	// Start with player whose Turn is 1
	currentPlayer := findPlayerWithTurn(1)
	opponentPlayer := findPlayerWithTurn(2)

	for {
		// Send turn message to current player
		_, err := currentPlayer.Conn.Write([]byte("Your turn! Choose an action:\nSwitch: {pokemon ID}\nAttack\nForfeit\n"))
		if err != nil {
			fmt.Println("Error writing to current player:", err)
			return
		}

		// Read action from current player
		buffer := make([]byte, 1024)
		n, err := currentPlayer.Conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from current player:", err)
			return
		}
		action := strings.TrimSpace(string(buffer[:n]))
		processAction(currentPlayer, opponentPlayer, action)

		// Switch turn to the opponent
		currentPlayer, opponentPlayer = opponentPlayer, currentPlayer
	}
}

func findPlayerWithTurn(turn int) *Player {
	for _, p := range players {
		if p.Turn == turn {
			return p
		}
	}
	return nil
}

func assignInitialTurns() {
	player1 := players[0]
	player2 := players[1]

	player1Speed := getSpeedOfFirstActivePokemon(player1)
	player2Speed := getSpeedOfFirstActivePokemon(player2)

	if player1Speed >= player2Speed {
		player1.Turn = 1 // Assign turn 1 to player 1
		player2.Turn = 2 // Assign turn 2 to player 2
	} else {
		player1.Turn = 2 // Assign turn 2 to player 1
		player2.Turn = 1 // Assign turn 1 to player 2
	}
}

func getSpeedOfFirstActivePokemon(player *Player) int {
	firstActiveID := strconv.Itoa(player.Active[0])
	for _, pokemon := range player.Pokemons {
		if pokemon.ID == firstActiveID {
			fmt.Printf("%s's first pokemon speed is: %d\n", player.Name, pokemon.Speed)
			return pokemon.Speed
		}
	}
	return 0
}

func processAction(currentPlayer, opponentPlayer *Player, action string) {
	switch strings.ToLower(action) {
	case "attack":
		rand.Seed(time.Now().UnixNano())
		damage := calculateDamage(currentPlayer, opponentPlayer)
		opponentPlayer.Health[0] -= damage
		fmt.Printf("%s attacked and dealt %d damage to %s's first Pokémon. Remaining HP: %d\n", currentPlayer.Name, damage, opponentPlayer.Name, opponentPlayer.Health[0])

		// Notify players
		currentPlayer.Conn.Write([]byte(fmt.Sprintf("You attacked and dealt %d damage. Opponent's Pokémon remaining HP: %d\n", damage, opponentPlayer.Health[0])))
		opponentPlayer.Conn.Write([]byte(fmt.Sprintf("Opponent attacked and dealt %d damage to your Pokémon. Remaining HP: %d\n", damage, opponentPlayer.Health[0])))

		// Check if opponent's Pokémon fainted
		if opponentPlayer.Health[0] <= 0 {
			currentPlayer.Conn.Write([]byte("Opponent's Pokémon fainted!\n"))
			opponentPlayer.Conn.Write([]byte("Your Pokémon fainted!\n"))
			if !switchToNextPokemon(opponentPlayer) {
				// End the battle if no Pokémon left to switch to
				currentPlayer.Conn.Write([]byte("You win!\n"))
				opponentPlayer.Conn.Write([]byte("You lose!\n"))
				currentPlayer.Conn.Close()
				opponentPlayer.Conn.Close()
				return
			}
		}
	case "switch":
		// Handle Pokémon switching logic
		if len(currentPlayer.Active) < 2 {
			currentPlayer.Conn.Write([]byte("No other Pokémon to switch to.\n"))
			return
		}
		switchPokemon(currentPlayer)
		currentPlayer.Conn.Write([]byte("You switched Pokémon. Turn passed to opponent.\n"))
		opponentPlayer.Conn.Write([]byte("Opponent switched Pokémon. It's your turn.\n"))
	case "forfeit":
		// Handle player forfeiting the match
		currentPlayer.Conn.Write([]byte("You forfeited. You lose!\n"))
		opponentPlayer.Conn.Write([]byte("Opponent forfeited. You win!\n"))
		currentPlayer.Conn.Close()
		opponentPlayer.Conn.Close()
	}
}
func calculateDamage(currentPlayer, opponentPlayer *Player) int {
	currentPokemon := findPokemonByID(currentPlayer.Active[0], currentPlayer)
	opponentPokemon := findPokemonByID(opponentPlayer.Active[0], opponentPlayer)
	if currentPokemon == nil || opponentPokemon == nil {
		return 0
	}

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	attackType := r.Intn(2)

	var damage int
	if attackType == 0 {
		// Normal attack
		damage = currentPokemon.Attack - opponentPokemon.Defense
	} else {
		// Special attack
		damage = currentPokemon.SpAtk*2 - opponentPokemon.SpDef
	}

	if damage < 0 {
		damage = 0
	}

	return damage
}
func findPokemonByID(id int, player *Player) *Pokemon {
	for _, pokemon := range player.Pokemons {
		if pokemon.ID == strconv.Itoa(id) {
			return &pokemon
		}
	}
	return nil
}

func switchToNextPokemon(player *Player) bool {
	if len(player.Active) <= 1 {
		return false
	}
	player.Active = player.Active[1:]
	player.Health = player.Health[1:]
	return true
}

func switchPokemon(player *Player) {
	player.Active = append(player.Active[1:], player.Active[0])
	player.Health = append(player.Health[1:], player.Health[0])
}
