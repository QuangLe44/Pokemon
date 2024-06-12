package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Set player name
	fmt.Print("Enter your name: ")
	name, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	name = name[:len(name)-1] // Remove newline character
	writer.WriteString(name + "\n")
	writer.Flush()

	// Set player's active Pokemon
	fmt.Println("Choose your active Pokemon:")
	fmt.Println("1. Pikachu")
	fmt.Println("2. Charmander")
	fmt.Println("3. Squirtle")
	fmt.Print("Enter the number of your choice: ")
	choiceInput, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	choiceInput = strings.TrimSpace(choiceInput)
	choice, err := strconv.Atoi(choiceInput)
	if err != nil || choice < 1 || choice > 3 {
		fmt.Println("Invalid choice. Please choose a number between 1 and 3.")
		return
	}
	fmt.Println("Sending choice:", choiceInput)

	// Send the selected Pokémon choice to the server
	writer.WriteString(choiceInput + "\n")
	writer.Flush()

	fmt.Println("Sent choice to server.")

	// Send ready signal to server
	writer.WriteString("ready\n")
	writer.Flush()

	fmt.Println("Sent ready signal to server.")

	// Wait for server messages
	for {
		message, _ := reader.ReadString('\n')
		fmt.Print(message)

		if message == "Game over! You win!\n" {
			break
		}

		// Check if it's the player's turn
		if strings.HasPrefix(message, name) {
			fmt.Println("It's your turn.")
		}
	}
}
