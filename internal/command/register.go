package command

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Fraegdegjevar/Gator/internal/config"
	"github.com/Fraegdegjevar/Gator/internal/database"
	"github.com/google/uuid"
)

// 'Register' the user (name is provided) in the database and set them as the
// current user in config.
// Exits program with code 1 if user with same name already exists.
func HandlerRegister(fs config.FileSystem, s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("must supply one user to register")
	}
	// access the db query object inside state to execute our sql query
	// to create a user in DB. CreateUSer needs context.Background
	// (empty Context) + CreateUSer params (i.e see db schema - uuid,
	// creted_at, updated_at, name).
	// If fail - returns
	user, err := s.Db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Args[0],
		})
	if err != nil {
		fmt.Printf("error adding user to database: %v\n", err)
		os.Exit(1)
	}

	//Otherwise, we succeeded in adding to the db. Print details
	// and set current user
	fmt.Println("user was created in database:")
	fmt.Printf("Name: %v\n", user.Name)
	fmt.Printf("UUID: %v\n", user.ID)
	fmt.Printf("Created At: %v\n", user.CreatedAt)
	fmt.Printf("Updated At: %v\n", user.UpdatedAt)

	err = s.Config.SetUser(fs, cmd.Args[0])
	if err != nil {
		return err
	}
	return nil
}
