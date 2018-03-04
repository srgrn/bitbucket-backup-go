package main

import (
	"flag"
	"fmt"
	"github.com/srgrn/bitbucket-api/swagger"
	"golang.org/x/net/context"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	// "io/ioutil"
	"log"
	"os"
	"sync"
)

func main() {
	log.Println("Starting")
	var wg sync.WaitGroup
	username := flag.String("username", "", "The username to use")
	password := flag.String("password", "", "The password to use")
	flag.Parse()
	if *username == "" {
		log.Println("Using user name from environment")
		username = getEnvVarOrExit("BITBUCKET_USERNAME")
	}
	if *password == "" {
		log.Println("Using password from environment")
		password = getEnvVarOrExit("BITBUCKET_PASSWORD")
	}
	client := swagger.NewAPIClient(swagger.NewConfiguration())
	auth_ctx := context.WithValue(context.Background(), swagger.ContextBasicAuth, swagger.BasicAuth{
		UserName: *username,
		Password: *password,
	})

	// clone("bridge", "git@bitbucket.org:eranzimbler/bridge.git")
	for name, url := range get_repositories(client, auth_ctx, *username) {
		wg.Add(1)
		go func(n, u string) {
			defer wg.Done()
			log.Println(n, u)
			clone(n, u)
		}(fmt.Sprintf("%s/%s", "output", name), url)
	}
	wg.Wait()
	log.Println("Ending")
}

func clone(directory string, url string) {
	// log.Println(directory, url)
	auth, err := ssh.NewPublicKeysFromFile("git", fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME")), "")

	_, err = git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              auth,
	})
	if err != nil {
		log.Fatalln(directory, err)
	}
	// log.Println(r)
}

func get_repositories(client *swagger.APIClient, ctx context.Context, username string) map[string]string {
	repos := make(map[string]string)
	repositories, _, _ := client.RepositoriesApi.RepositoriesUsernameGet(ctx, username, nil)
	for _, repo := range repositories.Values {
		// fmt.Println(repo.Links)
		repos[repo.Name] = repo.Links.Clone[1].Href
	}
	return repos
}

func getEnvVarOrExit(varName string) *string {
	value := os.Getenv(varName)
	if value == "" {
		log.Printf("Missing environment variable %s\n", varName)
		log.Println("Set environment variable or specify on CLI")
		flag.PrintDefaults()
		os.Exit(1)
	}

	return &value
}
