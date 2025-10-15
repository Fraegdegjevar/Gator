package main

import (
	"fmt"

	"github.com/Fraegdegjevar/Gator/internal/config"
)

func main() {
	//fmt.Println("Initial read of config file")
	// Set our real, OSFileSystem
	fs := config.OSFileSystem{}
	conf, err := config.Read(fs)

	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("Config current_user_name: %s\n", conf.Current_user_name)

	// Set user and write to disl.
	//fmt.Printf("Setting user: %s\n", user)
	err = conf.SetUser(fs, "will")
	if err != nil {
		fmt.Println(err)
	}

	//Read from disk again and check the current user field has been set
	conf, err = config.Read(fs)
	if err != nil {
		fmt.Printf("error reading updating config: %s\n", err)
	}
	fmt.Println(conf)

}
