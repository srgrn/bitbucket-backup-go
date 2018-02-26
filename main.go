package main

import (
	"fmt"
	"github.com/ktrysmt/go-bitbucket"
	"gopkg.in/src-d/go-git.v4"
	"os"
)

func main() {
	fmt.Println("Welcome to the new code")
	user := os.Getenv("BITBUCKET_USERNAME")
	pass := os.Getenv("BITBUCKET_PASSWORD")
	c := bitbucket.NewBasicAuth(user, pass)

}
