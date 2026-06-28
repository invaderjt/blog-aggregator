package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/invaderjt/blog-aggregator/internal/config"
	"github.com/invaderjt/blog-aggregator/internal/database"
)

func initialize() *state {
	s := &state{}
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Could not read config file")
	}
	s.Cfg = &cfg
	return s
}

func updateState(s *state, db *sql.DB) {
	dbQueries := database.New(db)
	s.Db = dbQueries

	var cmds commands
	cmds.Commands = make(map[string]func(*state, command) error)
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)

	input := os.Args
	if len(input) < 2 {
		log.Fatalln("Not enough arguments")
	}
	var cmd command
	cmd.Name = input[1]
	cmd.Args = input[2:]
	err := cmds.run(s, cmd)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s.Cfg)

}

type state struct {
	Db  *database.Queries
	Cfg *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	Commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	return c.Commands[cmd.Name](s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.Commands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		log.Fatalln("Login requires username argument")
	}

	name := cmd.Args[0]
	user, exists := s.Db.GetUser(context.Background(), name)
	if exists != nil {
		log.Fatalf("User with name %s does not exist\n", name)
	}

	err := s.Cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("%s set as username.", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		log.Fatalln("Register requires username argument")
	}

	ctx := context.Background()
	uuid := uuid.New()
	created_at := time.Now()
	updated_at := created_at
	name := cmd.Args[0]

	if _, exists := s.Db.GetUser(ctx, name); exists == nil {
		log.Fatalf("User with name %s already exists\n", name)
	}

	params := database.CreateUserParams{
		ID:        uuid,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		Name:      name,
	}

	user, err := s.Db.CreateUser(context.Background(), params)
	if err != nil {
		log.Fatalln("Error creating user")
	}

	err = s.Cfg.SetUser(user.Name)
	if err != nil {
		log.Fatalln("Error setting username")
	}
	fmt.Printf("User %s registered and logged in\n", name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.Db.ResetDB(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}
