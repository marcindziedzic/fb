package main

import (
	"bytes"
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
	config := readConfig("example.yaml")
	sections := config.Session["start"]

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
