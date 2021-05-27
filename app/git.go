package app

import (
	"fmt"
	"os"
	"regexp"

	"github.com/xanzy/go-gitlab"
)

var (
	gitClient *gitlab.Client

	gitIsSetup bool
)

func SwitchBranch(branch string) (string, error) {
	return run(Config.GitPath, "checkout", branch)
}

func CreateMergeBranch() error {
	if err := gitSetup(); err != nil {
		return err
	}
	if _, err := runQuiet(Config.GitPath, "checkout", "-b", Config.MRBranch); err != nil {
		return err
	}
	if _, err := runQuiet(Config.GitPath, "add", "."); err != nil {
		return err
	}
	if _, err := runQuiet(Config.GitPath, "commit", "-m", "$ composer update"); err != nil {
		return err
	}
	if _, err := runQuiet(Config.GitPath, "push", "origin", Config.MRBranch); err != nil {
		return err
	}

	return nil
}

// DeleteOriginBranch will delete a branch from origin
func deleteOriginBranch(branch string) error {
	if err := gitSetup(); err != nil {
		return err
	}

	fmt.Printf("Deleting older branch/MR: %s\n", branch)
	if _, err := runQuiet(Config.GitPath, "push", "origin", ":"+branch); err != nil {
		return err
	}

	return nil
}

func gitSetup() error {
	if gitIsSetup {
		return nil
	}
	if _, err := run(Config.GitPath, "config", "user.name", Config.GitUser); err != nil {
		return err
	}

	if _, err := run(Config.GitPath, "config", "user.email", Config.GitEmail); err != nil {
		return err
	}

	if getAPIToken() != "" &&
		os.Getenv("CI_REPOSITORY_URL") != "" {
		var re = regexp.MustCompile(`^https:\/\/gitlab-ci-token:(.*)@(.*)`)
		var str = os.Getenv("CI_REPOSITORY_URL")

		match := re.FindStringSubmatch(str)
		originURL := fmt.Sprintf("https://gitlab-ci-token:%s@%s",
			getAPIToken(),
			match[2],
		)

		if _, err := run(Config.GitPath, "remote", "set-url", "origin", originURL); err != nil {
			fmt.Println("Error setting remote")
			return err
		}
	}

	gitIsSetup = true

	return nil
}
