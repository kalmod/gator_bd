### Requirements
Both Postgres and Go are required to use the program.
- Postgres
- Go

### Installation
Use `go install` to install the application.

```
go install gator
```
### Configuration
To configure, you'll need to create a json file in your home directory.
`~/.gatorconfig.json`

```json
{
  "db_url": "postgres://example"
}
```
### Commands
Commands available to use with gator. 

`login [name]`: Log in as specified user. Selects a user that currently exists in the database.
`register [name]`: Registers a new user to the database
`reset`: Deletes users from the database.
`users`: Returns all users in the database.
`agg [time]`: Fetches feeds at specified interval.
- time example: `1s`,`2m`,`1h` 
`addfeed [feedname] [url]`: Adds a new feed to the database.
`feeds`: Returns all feeds in the database.
`follow [url]`: Creates a new Follow entry for the url with the currently logged in user.
`following`: Get feeds followed by the current user.
`unfollow [url]`: Remove a Follow entry for the url with the currently logged in user.
`browse [limit]`: View all the posts from the feeds the user follows. 
- Limit is the number of posts.
