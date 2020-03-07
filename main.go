package main

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"math/rand"
	"os/exec"
	"time"
)

func runCmd(cmd *exec.Cmd) {
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func getMessage(r *rand.Rand, quotes *[196]string) string {
	randNum := r.Intn(196)
	return (*quotes)[randNum]
}

func main() {

	r := rand.New(rand.NewSource(time.Now().Unix()))

	quotes := InitQuote()

	w := watcher.New()
	w.SetMaxEvents(1)
	go func() {
		for {
			select {
			case event := <-w.Event:
				gitAdd := exec.Command("git", "add", ".")
				gitCommit := exec.Command("git", "commit", "-m", getMessage(r, &quotes))
				gitPush := exec.Command("git", "push")
				runCmd(gitAdd)
				runCmd(gitCommit)
				runCmd(gitPush)
				fmt.Println(event)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	w.Ignore(".git")
	if err := w.AddRecursive("."); err != nil {
		log.Fatalln(err)
	}

	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}
	fmt.Println()

	go func() {
		w.Wait()
		w.TriggerEvent(watcher.Create, nil)
		w.TriggerEvent(watcher.Remove, nil)
	}()

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("hello")
}
