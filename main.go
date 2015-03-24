package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

	for {
		_, err = fmt.Scanln(&command)
		if err != nil {
			log.Fatal(err)
		}

		// TODO implement missing steps, store events
		switch command {
		case "start", "s":
			run(config.Session["start"])
		case "break", "b":
			run(config.Session["end"])
		case "interrupt", "i":
			panic("not yet implemented")
		default:
			panic("unrecognized command")
		}
	}
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
