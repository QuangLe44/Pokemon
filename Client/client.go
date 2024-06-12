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
	reader := bufio.NewReader(os.Stdin)
	defer conn.Close()

	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name) // Remove newline character
	conn.Write([]byte("Player name: " + name))

	go func() {
		for {
			buffer := make([]byte, 1024)
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading from server:", err.Error())
				os.Exit(1)
			}
			message := string(buffer)
			fmt.Println(message)
		}
	}()
	fmt.Print("Enter 'ready' when you are ready: ")
	var clientChoice []string
	for {
		message, _ := reader.ReadString('\n')
		switch strings.TrimSpace(message) {
		case "ready":
			pokemonList := strings.Join(clientChoice, ",")
			conn.Write([]byte("Player choice: " + pokemonList))
			fmt.Println("You are ready")
			conn.Write([]byte("ready"))

		default:
			if isInteger(message) {
				fmt.Println(len(clientChoice))
				if len(clientChoice) < 3 {
					clientChoice = append(clientChoice, strings.TrimSpace(message))
				} else {
					fmt.Println("You have chosen all pokemons. Please enter 'ready'")
				}
				fmt.Println(message)
			} else {
				fmt.Println(message)
				message = strings.TrimSpace(message)
				conn.Write([]byte(message))
			}
		}
	}
}
func isInteger(s string) bool {
	_, err := strconv.Atoi(strings.TrimSpace(s))
	return err == nil
}
