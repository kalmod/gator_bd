package main

import (
	"blog_agg_2/internal/config"
	"blog_agg_2/internal/database"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read() // Read config file
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DB_URL) // connect to db
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbQueries := database.New(db) // load dbQueries
	current_state := config.State{Db: dbQueries, Cfg: &cfg}

	commands := config.Commands{
		Command_map: make(map[string](func(*config.State, config.Command) error)),
	}
	commands.Command_map["login"] = config.HandlerLogin
	commands.Command_map["register"] = config.HandlerRegister
	commands.Command_map["reset"] = config.HandlerReset
	commands.Command_map["users"] = config.HandlerGetAllUsers
	commands.Command_map["agg"] = config.HandlerAgg
	commands.Command_map["addfeed"] = config.MiddlewareLoggedIn(config.HandlerAddFeed)
	commands.Command_map["feeds"] = config.HandlerGetAllFeeds
	commands.Command_map["follow"] = config.MiddlewareLoggedIn(config.HandlerCreateFollow)
	commands.Command_map["following"] = config.MiddlewareLoggedIn(config.HanlderGetFollowsForUser)
	commands.Command_map["unfollow"] = config.MiddlewareLoggedIn(config.HanlderUnfollowFeedForUser)
	commands.Command_map["browse"] = config.MiddlewareLoggedIn(config.HandlerBrowse)

	if len(os.Args) < 2 {
		fmt.Println("Too few arguments provided")
		os.Exit(1)
	}

	entered_command := config.Command{Name: os.Args[1], Args: os.Args[2:]}
	err = commands.Run(&current_state, entered_command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}
