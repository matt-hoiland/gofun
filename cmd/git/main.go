package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const JournalRepoPath = "/Users/matthew.hoiland/prsn/journal-2022"

func CheckError(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
		os.Exit(1)
	}
}

func main() {
	// Open a local repo
	r, err := git.PlainOpen(JournalRepoPath)
	CheckError(err, fmt.Sprintf("unable to open repo at '%s'", JournalRepoPath))

	LogExample(r)
}

func LogExample(r *git.Repository) {
	// git log --grep 'fmi/albany: Clear completed tasks'

	refs, err := r.Log(&git.LogOptions{})
	CheckError(err, "trouble getting the log")
	var commits []*object.Commit
	err = refs.ForEach(func(c *object.Commit) error {
		if strings.Contains(c.Message, "fmi/albany: Clear completed tasks") {
			commits = append(commits, c)
		}
		return nil
	})
	CheckError(err, "trouble iterating over commits")

	mostRecent := true
	taskCounts := make(map[string]int)
	var tasks []string
	for _, commit := range commits {
		fmt.Printf("%v %s\n", commit.Author.When, strings.TrimSpace(commit.Message))
		file, err := commit.File("albany/_inventory.md")
		CheckError(err, "path was wrong")

		// Slurp file contents into memory
		rc, err := file.Reader()
		CheckError(err, "failed to get reader")
		defer rc.Close()

		data, err := io.ReadAll(rc)
		CheckError(err, "io error")
		ProcessInventory(mostRecent, string(data), taskCounts, &tasks)
		mostRecent = false
	}
	for _, task := range tasks {
		fmt.Printf("%2d %s\n", taskCounts[task], task)
	}
}

// ProcessInventory takes in an inventory file and produces a list of tasks
func ProcessInventory(mostRecent bool, fileContents string, taskCounts map[string]int, tasks *[]string) {
	lines := strings.Split(fileContents, "\n")
	re := regexp.MustCompile(`^(?:  )+\* \[ \] (?P<task>[\w\s]+)(?: \(\d+\))?$`)
	taskIndex := re.SubexpIndex("task")
	if taskIndex == -1 {
		panic("failed to compile regexp")
	}
	for _, line := range lines {
		if re.MatchString(line) {
			matches := re.FindStringSubmatch(line)
			task := matches[taskIndex]
			if mostRecent {
				if _, ok := taskCounts[task]; !ok {
					*tasks = append(*tasks, task)
				}
				taskCounts[task]++
			} else {
				if _, ok := taskCounts[task]; ok {
					taskCounts[task]++
				}
			}
		}
	}
}

//nolint
func ReadHeadExample(r *git.Repository) {
	// Retrieve the HEAD reference
	head, err := r.Head()
	CheckError(err, "unable to fetch HEAD reference")

	// Retrieve commit object pointed to by HEAD
	commit, err := r.CommitObject(head.Hash())
	CheckError(err, "unable to fetch HEAD commit")
	fmt.Printf("HEAD commit message: '%s'\n", strings.TrimSpace(commit.Message))

	// Get object of the albany/_inventory.md file
	file, err := commit.File("albany/_inventory.md")
	CheckError(err, "path was wrong")

	// Slurp file contents into memory
	rc, err := file.Reader()
	CheckError(err, "failed to get reader")
	defer rc.Close()

	data, err := io.ReadAll(rc)
	CheckError(err, "io error")
	fmt.Println(string(data))
}
