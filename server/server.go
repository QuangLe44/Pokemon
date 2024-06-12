package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Pokemon struct {
	Name       string
	Speed      int
	Attack     int
	Defense    int
	SpecialAtk int
	SpecialDef int
	Type       string
	CurrentHP  int
	TotalExp   int
	Elemental  map[string]int // Elemental damage multipliers
}

type Player struct {
	Name     string
	Pokemons []*Pokemon
	Active   *Pokemon
	Ready    bool
}

type Battle struct {
	Players     []*Player
	Turn        int
	WinnerIndex int
}

var (
	players []*Player
	mutex   sync.Mutex
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

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	player := &Player{}
	player.Name, _ = reader.ReadString('\n')
	player.Name = strings.TrimSpace(player.Name)

	fmt.Println("Received player name:", player.Name)

	// Simulate selecting pokemons
	player.Pokemons = []*Pokemon{
		&Pokemon{Name: "Pikachu", Speed: 10, Attack: 20, Defense: 15, SpecialAtk: 25, SpecialDef: 20, Type: "Electric", CurrentHP: 100, TotalExp: 0},
		&Pokemon{Name: "Charmander", Speed: 8, Attack: 22, Defense: 18, SpecialAtk: 24, SpecialDef: 16, Type: "Fire", CurrentHP: 100, TotalExp: 0},
		&Pokemon{Name: "Squirtle", Speed: 7, Attack: 18, Defense: 20, SpecialAtk: 22, SpecialDef: 18, Type: "Water", CurrentHP: 100, TotalExp: 0},
	}

	mutex.Lock()
	players = append(players, player)
	mutex.Unlock()

	fmt.Println("Player added to players list.")

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message from client:", err)
			return
		}

		// Handle the message
		// For example, you can process the message and take actions based on its content
		fmt.Println("Received message from", player.Name+":", message)

		// Send response back to the client if needed
		// For example, you can send an acknowledgment back to the client
		response := "Received message: " + message
		writer.WriteString(response)
		writer.Flush()
	}

	// 	fmt.Println("Both players are ready.")

	// 	// Game starts
	// 	battle := &Battle{Players: players, Turn: 0}

	// 	for battle.WinnerIndex == -1 {
	// 		player := battle.Players[battle.Turn%2]
	// 		enemy := battle.Players[(battle.Turn+1)%2]

	// 		sendMessage(player.Name+", it's your turn.\n", writer)
	// 		sendMessage("Your active Pokemon: "+player.Active.Name+"\n", writer)
	// 		sendMessage("Choose an action:\n1. Attack\n2. Switch Pokemon\n", writer)

	// 		action, _ := reader.ReadString('\n')
	// 		action = strings.TrimSpace(action)

	// 		switch action {
	// 		case "1":
	// 			// Simulate attack
	// 			sendMessage("Select an attack:\n1. Normal Attack\n2. Special Attack\n", writer)
	// 			attackType, _ := reader.ReadString('\n')
	// 			attackType = strings.TrimSpace(attackType)

	// 			var damage int
	// 			if attackType == "1" {
	// 				damage = player.Active.Attack - enemy.Active.Defense
	// 			} else if attackType == "2" {
	// 				damage = player.Active.SpecialAtk*enemy.Active.Elemental[player.Active.Type] - enemy.Active.SpecialDef
	// 			} else {
	// 				sendMessage("Invalid attack type.\n", writer)
	// 				continue
	// 			}

	// 			enemy.Active.CurrentHP -= damage
	// 			sendMessage("Dealt "+strconv.Itoa(damage)+" damage to "+enemy.Active.Name+".\n", writer)

	// 		case "2":
	// 			// Simulate switching Pokemon
	// 			sendMessage("Select a Pokemon to switch to:\n", writer)
	// 			for i, pokemon := range player.Pokemons {
	// 				sendMessage(strconv.Itoa(i+1)+". "+pokemon.Name+"\n", writer)
	// 			}

	// 			pokemonIndex, _ := reader.ReadString('\n')
	// 			pokemonIndex = strings.TrimSpace(pokemonIndex)
	// 			index, err := strconv.Atoi(pokemonIndex)
	// 			if err != nil || index < 1 || index > len(player.Pokemons) {
	// 				sendMessage("Invalid Pokemon index.\n", writer)
	// 				continue
	// 			}

	// 			player.Active = player.Pokemons[index-1]
	// 			sendMessage("Switched to "+player.Active.Name+".\n", writer)
	// 		}

	// 		battle.Turn++
	// 	}

	// 	winner := battle.Players[battle.WinnerIndex]
	// 	sendMessage("Game over! "+winner.Name+" wins!\n", writer)
	// }

	// func sendMessage(message string, writer *bufio.Writer) {
	// 	writer.WriteString(message)
	// 	writer.Flush()
}
