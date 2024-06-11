package main

import (
	"fmt"
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
	Name string `json:"name"`
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
	go startPokeworld(s, conn)
}

func (s *server) Start() {
	s.wg.Add(2)
	go s.acceptConnections()
	go s.handleConnections()
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

func main() {
	pokemons, err := loadPokemonData("../data/baseInfo.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pokedex := make([]*pokemon, len(pokemons))
	for i, p := range pokemons {
		pokedex[i] = &pokemon{Name: p.Name}
	}

	fmt.Println("Pokédex:")
	for _, p := range pokedex {
		fmt.Println(p.Name)
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
