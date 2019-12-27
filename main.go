package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/alecthomas/kingpin"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func execCmd(command string) []byte {
	logrus.Printf("running %s\n", command)
	commandAndArgs := strings.Split(command, " ")
	cmd := exec.Command(commandAndArgs[0], commandAndArgs[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Infoln(string(out))
		logrus.Errorln(err)
	}
	return out
}

// pullIfNeeded checks if there was an update to a specified branch
func pullIfNeeded(branch string) bool {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Errorf("error getting current directory: %v", err)
	}
	r, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		logrus.Errorf("error opening repo: %v", err)
	}
	w, err := r.Worktree()
	if err != nil {
		logrus.Errorf("error getting worktree: %v", err)
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
	})
	if err != nil {
		logrus.Errorf("error checking out repo: %v", err)
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		logrus.Errorf("error pulling repo: %v", err)
	}

	return true // now always reloading
}

func main() {
	branch := kingpin.Flag("branch", "branch to sync").Required().String()
	interval := kingpin.Flag("interval", "interval for the sync loop").Default("1m").Duration()
	command := kingpin.Flag("command", "command to run after the sync").Default("make run").String()
	kingpin.Parse()

	stopChannel := make(chan int)

	for {
		select {
		case <-time.After(*interval):
			changed := pullIfNeeded(*branch)
			if changed {
				execCmd(*command)
			}

		case <-stopChannel:
			logrus.Infoln("stop requested, exiting.")
			return
		}
	}
}
