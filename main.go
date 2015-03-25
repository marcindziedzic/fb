package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type config struct {
	Name    string
	Session map[string][]string
}

func main() {
	// TODO config location should be passed as parameter
	config := readConfig("/home/mdziedzic/.fb/example.yaml")

	fmt.Println("Hello, to kick-off your session just type `start`")
	fmt.Println("Later on you can `break` or `interrupt`")

	var command string
	var err error

	now := time.Now()
	duration := time.Minute * 90
	timer := time.NewTimer(time.Minute * 90)

	for {
		_, err = fmt.Scanln(&command)
		if err != nil {
			log.Fatal(err)
		}

		// TODO implement missing steps, store events
		switch command {
		case "start", "s":
			f1 := makeRunner(config.Session["start"])
			f2 := makeRunner(config.Session["end"])
			go runInterchangebly(f1, f2, timer)
		case "break", "b":
			timer.Reset(time.Second * 0)
		case "interrupt", "i":
			panic("not yet implemented")
		case "time", "t":
			elapsed := time.Since(now)
			fmt.Println("Time left: ", duration-elapsed)
		default:
			fmt.Println("Command not recognized, try one of `start`, `break`, `interrrupt`, `time` or their first letters")
		}
	}
}

// FIXME looping of sessions
func runInterchangebly(f1 func(), f2 func(), timer *time.Timer) {
	f1()
	<-timer.C
	f2()
	fmt.Println("session expired, for next iteration you need to start brand new session")
}

func makeRunner(sections []string) func() {
	return func() { run(sections) }
}

func run(sections []string) {
	for _, cmd := range sections {
		commands := Transform(strings.Split(cmd, "|"))

		var b bytes.Buffer
		if err := Execute(&b, commands...); err != nil {
			log.Fatal(err)
		}
		io.Copy(os.Stdout, &b)
	}
}

func readConfig(filePath string) *config {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	cfg := config{}

	err = yaml.Unmarshal([]byte(data), &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
