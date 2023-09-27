package main

import (
	"io/ioutil"
	"log"

	"gitlab.com/mefit/mefit-server/entity"
	"gitlab.com/mefit/mefit-server/utils/initializer"
	yaml "gopkg.in/yaml.v2"
)

func getMovments() (moves []entity.Movement) {
	moves = make([]entity.Movement, 0)
	yamlFile, err := ioutil.ReadFile("seed/seed.yaml")
	if err != nil {
		log.Fatalf("While reading the seed data: %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &moves)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return
}

func main() {
	//Initialize database first
	defer initializer.Initialize()()
	//We just do seed at this moment
	moves := getMovments()
	for _, m := range moves {
		if err := entity.SimpleCrud(&m).Save(); err != nil {
			panic(err)
		}
	}
}
