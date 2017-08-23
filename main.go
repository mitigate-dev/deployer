package main

import "fmt"
import "flag"
import "context"
import "os"
import "os/exec"
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

	repoPath := "." + *app;
	_, err := os.Stat(repoPath)
	if err != nil {
		repoURL := "https://" + *username + ":" + *password + "@github.com/" + *org + "/" + *repo + ".git"
		err = cloneRepo(repoURL, repoPath)
		if err != nil {
			fmt.Printf("Problem cloning repo %v\n", err)
			os.Exit(1)
		}
	}

	sleepDuration := time.Duration(*sleepInt) * time.Second
	fmt.Printf("Sleep duration: %v\n", sleepDuration)

	ctx := context.Background()
	tp := github.BasicAuthTransport{
		Username: *username,
		Password: *password,
	}

	client := github.NewClient(tp.Client())

	opt := &github.DeploymentsListOptions{
		Environment: *env,
		ListOptions: github.ListOptions{PerPage: 1},
	}

	for {
		deployments, _, err := client.Repositories.ListDeployments(ctx, *org, *repo, opt)

		if err != nil {
			fmt.Printf("Problem in listing deployments %v\n", err)
			sleep(sleepDuration)
			continue;
		}

		deployment := deployments[0]
		if deployment != nil {
			fmt.Println(*deployment.ID, *deployment.Ref, *deployment.Environment)

			fmt.Println("Gist it")

			title := fmt.Sprintf("ID: %d, Ref: %s, Environment: %s", *deployment.ID, *deployment.Ref, *deployment.Environment)
			input := &github.Gist{
				Description: github.String(title),
				Public: github.Bool(false),
				Files: map[github.GistFilename]github.GistFile{
					"output.txt": {Content: github.String("new file content")},
				},
			}
			gist, _, err := client.Gists.Create(ctx, input)

			if err != nil {
				fmt.Printf("Problem in creating gist %v\n", err)
				sleep(sleepDuration)
				continue;
			}

			fmt.Println(*gist.HTMLURL)

			sleep(sleepDuration)
		}
	}
}

func sleep(duration time.Duration) {
	fmt.Printf("Sleeping %v\n", duration)
	time.Sleep(duration)
}

func cloneRepo(repoURL string, repoPath string) (error) {
	fmt.Println("Cloning ", repoURL)
	cmd := exec.Command("git", "clone", repoURL, repoPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
