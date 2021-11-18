package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/xanzy/go-gitlab"
)

func getGitlab() *gitlab.Client {
	url := os.Getenv("GITLAB_URL")
	token := os.Getenv("GITLAB_TOKEN")
	gl, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return gl
}

func getNewFileContent(localFilePath string) string {
	contents, err := os.ReadFile(localFilePath)
	if err != nil {
		log.Fatal("Could not open file")
	}
	return string(contents)
}

type args struct {
	file          string
	repo          string
	repoFile      string
	commitMessage string
	force         bool
	branchName    string
}

func parseArgs() *args {
	fileFlag := flag.String("file", "", "Specify the location of the local file")
	repoFlag := flag.String("repo", "", "The location of the GitLab repo. E.g. 'gitlab-org/gitlab'")
	repoFileFlag := flag.String("repo-file", "", "The repository file to replace")
	branchNameFlag := flag.String("branch", "", "The branch name to use. Optional. Defaults to refactor-change-file-<filename>")
	commitMessageFlag := flag.String("m", "", "The commit message")
	forceFlag := flag.Bool("force", false, "Force branch creation")
	flag.Parse()
	arguments := args{}
	arguments.file = *fileFlag
	arguments.repo = *repoFlag
	arguments.repoFile = *repoFileFlag
	arguments.commitMessage = *commitMessageFlag
	arguments.force = *forceFlag
	if *branchNameFlag == "" {
		arguments.branchName = fmt.Sprintf("refactor-change-file-%s", filepath.Base(*fileFlag))
	} else {
		arguments.branchName = *branchNameFlag
	}
	if arguments.file == "" || arguments.repo == "" || arguments.repoFile == "" || arguments.commitMessage == "" {
		log.Fatal("file, repo, repo-file and commitMessage (-m) are required.")
	}

	return &arguments
}

func getBranchOptions(branchName string, defaultBranch string) *gitlab.CreateBranchOptions {
	return &gitlab.CreateBranchOptions{Branch: &branchName, Ref: &defaultBranch}
}

func getDefaultBranch(gl *gitlab.Client, repo string) string {
	project, _, err := gl.Projects.GetProject(repo, nil, nil)
	if err != nil {
		panic(err)
	}
	return project.DefaultBranch
}

func createBranch(gl *gitlab.Client, branchName string, defaultBranch string, arguments *args) {
	if arguments.force {
		gl.Branches.DeleteBranch(arguments.repo, branchName)
	}
	_, _, err := gl.Branches.CreateBranch(arguments.repo, getBranchOptions(branchName, defaultBranch))
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not create branch '%s'. Does it already exist? Use -force to ignore.", branchName))
		os.Exit(1)
	}
}

func main() {
	arguments := parseArgs()
	gl := getGitlab()
	defaultBranch := getDefaultBranch(gl, arguments.repo)
	createBranch(gl, arguments.branchName, defaultBranch, arguments)
	options := &gitlab.UpdateFileOptions{
		Branch:        gitlab.String(arguments.branchName),
		Content:       gitlab.String(getNewFileContent(arguments.file)),
		CommitMessage: gitlab.String(arguments.commitMessage),
	}
	_, _, err := gl.RepositoryFiles.UpdateFile(arguments.repo, arguments.repoFile, options)
	if err != nil {
		panic(err)
	}
	gl.MergeRequests.CreateMergeRequest(arguments.repo, &gitlab.CreateMergeRequestOptions{Title: &arguments.commitMessage, SourceBranch: &arguments.branchName, TargetBranch: &defaultBranch, RemoveSourceBranch: gitlab.Bool(true)})
}
