package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"time"
)

func startPokeworld(s *server, conn net.Conn) {
	playerID := fmt.Sprintf("player-%d", len(s.players))
	player := &player{
		id:       playerID,
		x:        rand.Intn(len(s.pokeworld)),
		y:        rand.Intn(len(s.pokeworld)),
		pokemons: make([]*pokemon, 0, 200),
	}
	s.players[playerID] = player
	fmt.Fprintf(conn, "Welcome to Pokeworld, %s! You are at (%d, %d).\n", playerID, player.x, player.y)
	go s.spawnPokemon()
	go func() {
		for {
			select {
			case <-s.shutdown:
				return
			default:
				if player.autoMode {
					player.move(s)
				}
				time.Sleep(time.Second)
			}
		}
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		cmd := scanner.Text()
		switch cmd {
		case "up":
			player.moveUp(s)
			fmt.Fprintf(conn, "You moved up! You are at (%d, %d).\n", player.x, player.y)
		case "down":
			player.moveDown(s)
			fmt.Fprintf(conn, "You moved down! You are at (%d, %d).\n", player.x, player.y)
		case "left":
			player.moveLeft(s)
			fmt.Fprintf(conn, "You moved left! You are at (%d, %d).\n", player.x, player.y)
		case "right":
			player.moveRight(s)
			fmt.Fprintf(conn, "You moved right! You are at (%d, %d).\n", player.x, player.y)
		case "auto":
			player.autoMode = true
			player.autoTimer = time.AfterFunc(120*time.Second, func() {
				player.autoMode = false
			})
			fmt.Fprint(conn, "Auto mode enabled.\n")
		default:
			fmt.Fprint(conn, "Invalid command.\n")
		}
	}
}
func (p *player) move(s *server) {
	dx, dy := rand.Intn(3)-1, rand.Intn(3)-1
	p.x += dx
	p.y += dy
	if p.x < 0 {
		p.x = 0
	} else if p.x >= len(s.pokeworld) {
		p.x = len(s.pokeworld) - 1
	}
	if p.y < 0 {
		p.y = 0
	} else if p.y >= len(s.pokeworld) {
		p.y = len(s.pokeworld) - 1
	}
}

func (p *player) moveUp(s *server) {
	p.y--
	if p.y < 0 {
		p.y = 0
	}
}

func (p *player) moveDown(s *server) {
	p.y++
	if p.y >= len(s.pokeworld) {
		p.y = len(s.pokeworld) - 1
	}
}

func (p *player) moveLeft(s *server) {
	p.x--
	if p.x < 0 {
		p.x = 0
	}
}

func (p *player) moveRight(s *server) {
	p.x++
	if p.x >= len(s.pokeworld) {
		p.x = len(s.pokeworld) - 1
	}
}

func (s *server) spawnPokemon() {
	for {
		select {
		case <-s.shutdown:
			return
		default:
			for i := 0; i < 50; i++ {
				pokemonName := s.pokedex[rand.Intn(len(s.pokedex))].Name
				pokemon := &pokemon{
					Name: pokemonName,
				}
				x, y := rand.Intn(len(s.pokeworld)), rand.Intn(len(s.pokeworld))
				s.pokeworld[x][y] = pokemon
				fmt.Println("Spawned pokemon", pokemon.Name, "at (", x, ", ", y, ")")
			}
			time.Sleep(time.Minute)
		}
	}
}
