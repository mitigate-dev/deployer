package main

import "fmt"
import "flag"
import "context"
import "os"
import "os/exec"
import "time"
import "errors"
import "bytes"
import "io"
import "log"
import "regexp"
import "github.com/google/go-github/github"

func main() {
	fmt.Println("Deployer")

	username := flag.String("u",    "", "GitHub username (required)")
	password := flag.String("p",    "", "GitHub password (required)")
	org      := flag.String("org",  "", "GitHub org (required)")
	repo     := flag.String("repo", "", "GitHub repo (required)")
	env      := flag.String("env",  "", "Github deployment environment (required)")
	file     := flag.String("file", "", "File to execute when new deployment is available (required)")
	sleepInt := flag.Int(  "sleep", 30, "Time to sleep between loops (defaults to 30 seconds)")

	flag.Parse()

	if *username == "" || *password == "" || *org == "" || *repo == "" || *env == "" || *file == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println(*username, *password, *org, *repo, *env, *file)

	sleepDuration := time.Duration(*sleepInt) * time.Second
	fmt.Printf("Sleep duration: %v\n", sleepDuration)

	ctx := context.Background()
	tp := github.BasicAuthTransport{
		Username: *username,
		Password: *password,
	}

	client := github.NewClient(tp.Client())

	for {
		deployment, _, err := getDeployment(ctx, client, *org, *repo, *env)

		if err != nil {
			log.Printf("Problem in getting deployment %v\n", err)
			sleep(sleepDuration)
			continue
		}

		log.Println(*deployment.ID, *deployment.Ref, *deployment.Environment)

		gist, _, err := createDeploymentGist(ctx, client, deployment)

		if err != nil {
			log.Printf("Problem in creating gist %v\n", err)
			sleep(sleepDuration)
			continue
		}

		log.Println("Gist", *gist.HTMLURL)

		createDeploymentStatus(ctx, client, deployment, *org, *repo, "pending", *gist.HTMLURL)

		var output bytes.Buffer
		cmd := executeFile(*file, *deployment.Ref)
		cmd.Stdout = io.MultiWriter(&output, os.Stdout)
		cmd.Stderr = cmd.Stdout
		err = cmd.Run()
		if err != nil {
			log.Printf("Problem with file execution %v\n", err)
			createDeploymentStatus(ctx, client, deployment, *org, *repo, "error", *gist.HTMLURL)
			updateDeploymentGist(ctx, client, gist, output.String())
			sleep(sleepDuration)
			continue
		}

		createDeploymentStatus(ctx, client, deployment, *org, *repo, "success", *gist.HTMLURL)
		updateDeploymentGist(ctx, client, gist, output.String())

		sleep(sleepDuration)
	}
}

func sleep(duration time.Duration) {
	log.Printf("Sleeping %v\n", duration)
	time.Sleep(duration)
}

func executeFile(filePath string, ref string) *exec.Cmd {
	log.Println("Executing File", filePath)
	cmd := exec.Command(filePath, ref)
	// cmd := exec.Command("bin/deploy-stub", ref)
	return cmd
}

func getDeployment(ctx context.Context, client *github.Client, org string, repo string, env string) (*github.Deployment, *github.Response, error) {
	log.Println("Getting deployment ", env)
	opt := &github.DeploymentsListOptions{
		Environment: env,
		ListOptions: github.ListOptions{PerPage: 1},
	}
	deployments, resp, err := client.Repositories.ListDeployments(ctx, org, repo, opt)
	if err != nil {
		return nil, resp, err
	}

	if len(deployments) == 0 {
		err := errors.New("Deployment list is empty")
		return nil, resp, err
	}
	deployment := deployments[0]

	statuses, resp, err := client.Repositories.ListDeploymentStatuses(ctx, org, repo, *deployment.ID, &github.ListOptions{ PerPage: 1 })
	if err != nil {
		return deployment, resp, err
	}

	if len(statuses) > 0 {
		err := errors.New("Deployment statuses already present")
		return deployment, resp, err
	}

	return deployment, resp, err
}

func createDeploymentStatus(ctx context.Context, client *github.Client, deployment *github.Deployment, org string, repo string, state string, url string) (*github.DeploymentStatus, *github.Response, error) {
	log.Println("Deployment status create ", state)
	input := &github.DeploymentStatusRequest{
		State: github.String(state),
		LogURL: github.String(url),
	}
	status, resp, err := client.Repositories.CreateDeploymentStatus(ctx, org, repo, *deployment.ID, input)
	return status, resp, err
}

func createDeploymentGist(ctx context.Context, client *github.Client, deployment *github.Deployment) (*github.Gist, *github.Response, error) {
	title := fmt.Sprintf("ID: %d, Ref: %s, Environment: %s", *deployment.ID, *deployment.Ref, *deployment.Environment)
	log.Println("Gist create", title)
	input := &github.Gist{
		Description: github.String(title),
		Public: github.Bool(false),
		Files: map[github.GistFilename]github.GistFile{
			"output.txt": {Content: github.String("Pending")},
		},
	}
	gist, resp, err := client.Gists.Create(ctx, input)
	return gist, resp, err
}

func updateDeploymentGist(ctx context.Context, client *github.Client, gist *github.Gist, content string) (*github.Gist, *github.Response, error) {
	log.Println("Gist update")

	// Remove Terminal Escape codes - "remote:\e[1G"
	re := regexp.MustCompile(".*\x1b\\[1G")
	content = re.ReplaceAllString(content, "")

	file := gist.Files["output.txt"]
	file.Content = github.String(content)
	gist.Files["output.txt"] = file
	gist, resp, err := client.Gists.Edit(ctx, *gist.ID, gist)
	return gist, resp, err
}
