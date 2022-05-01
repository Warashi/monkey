package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/Warashi/implement-interpreter-with-go/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("user.Current: %v", err)
	}
	fmt.Printf("Hello, %s!, This is the Monkey programming language!\n", user.Username)
	fmt.Println("Fell free to type in commands")
	repl.Start(os.Stdin, os.Stdout)
}
