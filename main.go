package main

import (
	"flag"
	"github.com/src-d/go-git/transport/http"
	"github.com/srgrn/bitbucket-api/swagger"
	"golang.org/x/net/context"
	"gopkg.in/src-d/go-git.v4"
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
	// client := swagger.NewAPIClient(swagger.NewConfiguration())
	// auth_ctx := context.WithValue(context.Background(), swagger.ContextBasicAuth, swagger.BasicAuth{
	// 	UserName: *username,
	// 	Password: *password,
	// })

	clone(username, password, "bridge", "https://bitbucket.org/eranzimbler/bridge.git")
	// for name, url := range get_repositories(client, auth_ctx, *username) {
	// 	log.Println(name, url)
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		clone(username, password, name, url)
	// 	}()
	// }
	wg.Wait()
	log.Println("Ending")
}

func clone(username *string, password *string, directory string, url string) {
	// log.Println(*username, *password, directory, url)
	r, err := git.PlainClone(directory, true, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              &http.BasicAuth{Username: *username, Password: *password},
	})
	if err != nil {
		log.Fatalln(directory, err)
	}
	log.Println(r)
}

func get_repositories(client *swagger.APIClient, ctx context.Context, username string) map[string]string {
	repos := make(map[string]string)
	repositories, _, _ := client.RepositoriesApi.RepositoriesUsernameGet(ctx, username, nil)
	for _, repo := range repositories.Values {
		repos[repo.Name] = repo.Links.Clone[0].Href
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
