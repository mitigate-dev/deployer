package main

import "fmt"
import "flag"
import "context"
import "os"
import "time"
import "github.com/google/go-github/github"

func main() {
	fmt.Println("Deployer")

	username := flag.String("u",    "",  "GitHub username (required)")
	password := flag.String("p",    "",  "GitHub password (required)")
	org      := flag.String("org",  "",  "GitHub org (required)")
	repo     := flag.String("repo", "",  "GitHub repo (required)")
	env      := flag.String("env",  "",  "Github deployment environment (required)")
	app      := flag.String("app",  "",  "Dokku application name (required)")
	sleepInt := flag.Int(   "sleep", 30, "Time to sleep between loops (defaults to 30 seconds)")

	flag.Parse()

	if *app == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println(*username, *password, *org, *repo, *env, *app)

	sleepDuration := time.Duration(*sleepInt) * time.Second
	fmt.Printf("Sleep duration: %v\n", sleepDuration)

	ctx := context.Background()
	tp := github.BasicAuthTransport{
		Username: *username,
		Password: *password,
	}

	client := github.NewClient(tp.Client())

	opt := &github.DeploymentsListOptions{}

	for {
		deployments, _, err := client.Repositories.ListDeployments(ctx, *org, *repo, opt)

		if err != nil {
			fmt.Printf("Problem in listing deployments %v\n", err)
			sleep(sleepDuration)
			continue;
		}

		for _, deployment := range deployments {
			fmt.Println(*deployment.ID, *deployment.Ref, *deployment.Environment)
		}

		sleep(sleepDuration)
	}
}


func sleep(duration time.Duration) {
	fmt.Printf("Sleeping %v\n", duration)
	time.Sleep(duration)
}
