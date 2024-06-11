package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type server struct {
	wg         sync.WaitGroup
	listener   net.Listener
	shutdown   chan struct{}
	connection chan net.Conn
	players    map[string]*player
	pokeworld  [][]*pokemon
	pokedex    []*pokemon
}
type player struct {
	id        string
	x, y      int
	autoMode  bool
	autoTimer *time.Timer
	pokemons  []*pokemon
}

type pokemon struct {
	id    int
	level int
	ev    float64
}

func newServer(address string, size int, pokedex []*pokemon) (*server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address %s: %w", address, err)
	}
	s := &server{
		listener:   listener,
		shutdown:   make(chan struct{}),
		connection: make(chan net.Conn),
		players:    make(map[string]*player),
		pokeworld:  make([][]*pokemon, size),
		pokedex:    pokedex,
	}
	for i := range s.pokeworld {
		s.pokeworld[i] = make([]*pokemon, size)
	}
	return s, nil
}
func (s *server) acceptConnections() {
	defer s.wg.Done()
	for {
		select {
		case <-s.shutdown:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				continue
			}
			s.connection <- conn
		}
	}
}
func (s *server) handleConnections() {
	defer s.wg.Done()
	for {
		select {
		case <-s.shutdown:
			return
		case conn := <-s.connection:
			go s.handleConnection(conn)
		}
	}
}
func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()
	playerID := fmt.Sprintf("player-%d", len(s.players))
	player := &player{
		id:       playerID,
		x:        rand.Intn(len(s.pokeworld)),
		y:        rand.Intn(len(s.pokeworld)),
		pokemons: make([]*pokemon, 0, 200),
	}
	s.players[playerID] = player
	fmt.Fprintf(conn, "Welcome to Pokeworld, %s! You are at (%d, %d).\n", playerID, player.x, player.y)

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

func (s *server) Start() {
	s.wg.Add(2)
	go s.acceptConnections()
	go s.handleConnections()
	go s.spawnPokemon()
}
func (s *server) Stop() {
	close(s.shutdown)
	s.listener.Close()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return
	case <-time.After(time.Second):
		fmt.Println("Timed out waiting for connections to finish.")
		return
	}
}

func (s *server) spawnPokemon() {
	for {
		select {
		case <-s.shutdown:
			return
		default:
			for i := 0; i < 50; i++ {
				pokemon := &pokemon{
					id:    rand.Intn(len(s.pokedex)),
					level: rand.Intn(100) + 1,
					ev:    0.5 + rand.Float64(),
				}
				x, y := rand.Intn(len(s.pokeworld)), rand.Intn(len(s.pokeworld))
				s.pokeworld[x][y] = pokemon
				fmt.Println("Spawned pokemon", pokemon.id, "at (", x, ", ", y, ")")
			}
			time.Sleep(time.Minute)
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

func main() {
	pokedex := make([]*pokemon, 100)
	for i := range pokedex {
		pokedex[i] = &pokemon{id: i}
	}
	s, err := newServer(":8080", 1000, pokedex)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s.Start()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down server...")
	s.Stop()
	fmt.Println("Server stopped.")
}
