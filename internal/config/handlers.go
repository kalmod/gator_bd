package config

import (
	"blog_agg_2/internal/database"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Expects a single argument, username")
	}

	userInfo, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return errors.New(fmt.Sprintf("Could not find user: %s", cmd.Args[0]))
	}
	if err := s.Cfg.SetUser(userInfo.Name); err != nil {
		return err
	}
	fmt.Printf("The user '%s' has been set!\n", s.Cfg.Current_User_Name)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Expects a single argument, username")
	}
	newUserParams := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Args[0]}
	newUser, err := s.Db.CreateUser(context.Background(), newUserParams)
	if err != nil {
		return err
	}
	fmt.Printf("User, %s, was created.\n", newUser.Name)

	s.Cfg.SetUser(newUser.Name)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.ResetUser(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func HandlerGetAllUsers(s *State, cmd Command) error {
	allUsers, err := s.Db.GetAllUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range allUsers {
		if user.Name == s.Cfg.Current_User_Name {
			fmt.Println(user.Name, "(current)")
		} else {
			fmt.Println(user.Name)
		}
	}
	return err
}

// TODO: Need to change eventually.
func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("Not enough arguments")
	} else if len(cmd.Args) > 1 {
		return errors.New("Too many arguments")
	}

	dur, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(dur)

	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func HandlerAddFeed(s *State, cmd Command, curr_user database.User) error {
	if len(cmd.Args) < 2 {
		return errors.New("Not enough arguments")
	}

	feed, err := s.Db.AddFeed(context.Background(), database.AddFeedParams{
		ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Args[0],
		Url: cmd.Args[1], UserID: curr_user.ID,
	})
	if err != nil {
		return err
	}

	new_feed_follow := database.CreateFeedFollowParams{
		ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(),
		UserID: curr_user.ID, FeedID: feed.ID,
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), new_feed_follow)
	if err != nil {
		return err
	}

	return nil
}

func HandlerGetAllFeeds(s *State, cmd Command) error {
	all_feeds, err := s.Db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range all_feeds {
		user, err := s.Db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s %s\n", feed.Name, feed.Url, user.Name)
	}

	return nil
}

func HandlerCreateFollow(s *State, cmd Command, curr_user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("Not enough arguments")
	}
	url := cmd.Args[0]
	feed, err := s.Db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	new_feed_follow := database.CreateFeedFollowParams{
		ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(),
		UserID: curr_user.ID, FeedID: feed.ID,
	}

	feed_follow, err := s.Db.CreateFeedFollow(context.Background(), new_feed_follow)
	if err != nil {
		return err
	}
	fmt.Printf("%s - followed by - %s", feed_follow.FeedName, feed_follow.UserName)

	return nil
}

func HanlderGetFollowsForUser(s *State, cmd Command, curr_user database.User) error {
	// curr_user, err := s.Db.GetUser(context.Background(), s.Cfg.Current_User_Name)
	// if err != nil {
	// 	return err
	// }

	feed_follows_list, err := s.Db.GetFeedFollowsForUser(context.Background(), curr_user.ID)
	if err != nil {
		return err
	}

	for _, feed_follow := range feed_follows_list {
		fmt.Printf("%s\n", feed_follow.FeedName)
	}

	return nil
}

func HanlderUnfollowFeedForUser(s *State, cmd Command, curr_user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("Not enough arguments")
	}
	url := cmd.Args[0]
	feed, err := s.Db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	feed_to_remove := database.RemoveFeedForUserParams{UserID: curr_user.ID, FeedID: feed.ID}
	err = s.Db.RemoveFeedForUser(context.Background(), feed_to_remove)
	if err != nil {
		return err
	}

	return nil
}

func HandlerBrowse(s *State, cmd Command, curr_user database.User) error {
	var limit int32 = 2
	if len(cmd.Args) == 1 {
		conv_int, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
		limit = int32(conv_int)
	}

	get_posts_Params := database.GetPostsForUserParams{UserID: curr_user.ID, Limit: limit}
	posts, err := s.Db.GetPostsForUser(context.Background(), get_posts_Params)
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Println(post.Title)
	}

	return nil
}

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(s *State, cmd Command) error {
	f := func(s *State, cmd Command) error {
		curr_user, err := s.Db.GetUser(context.Background(), s.Cfg.Current_User_Name)
		if err != nil {
			return err
		}
		err = handler(s, cmd, curr_user)
		if err != nil {
			return err
		}

		return nil
	}
	return f
}
