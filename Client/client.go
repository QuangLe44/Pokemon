package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	fmt.Println("Enter commands (up, down, left, right, auto):")
	for {
		fmt.Print("> ")
		var cmd string
		fmt.Scanln(&cmd)
		fmt.Fprintln(conn, cmd)
	}
}
