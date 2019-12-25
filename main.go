package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/alecthomas/kingpin"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// pullIfNeeded checks if there was an update to a specified branch
func pullIfNeeded(branch string) {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Errorf("error getting current directory: %s, err")
	}
	r, err := git.PlainOpen(dir)
	if err != nil {
		logrus.Errorf("error opening repo: %s, err")
	}
	w, err := r.Worktree()
	if err != nil {
		logrus.Errorf("error getting worktree: %s, err")
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
	})
	if err != nil {
		logrus.Errorf("error checking out repo: %s, err")
	}
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		logrus.Errorf("error pulling repo: %s, err")
	}
}

func main() {
	branch := kingpin.Flag("branch", "branch to sync").Required().String()
	interval := kingpin.Flag("interval", "interval for the sync loop").Default("1m").Duration()
	kingpin.Parse()

	stopChannel := make(chan int)

	for {
		select {
		case <-time.After(*interval):
			pullIfNeeded(*branch)
		case <-stopChannel:
			logrus.Infoln("stop requested, exiting.")
			return
		}
	}
}
