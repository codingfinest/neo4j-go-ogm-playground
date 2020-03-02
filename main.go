package main

import (
	"fmt"

	gogm "github.com/codingfinest/neo4j-go-ogm"
)

//Define models

type Movie struct {
	gogm.Node `gogm:"label:FILM,label:PICTURE"`

	Title    string
	Released int64
	Tagline  string

	Characters []*Character `gogm:"direction:<-"`
	Directors  []*Director  `gogm:"direction:<-"`
}

type Director struct {
	gogm.Node
	Name   string
	Movies []*Movie
}

type Actor struct {
	gogm.Node
	Name       string
	Characters []*Character
}

type Character struct {
	gogm.Relationship `gogm:"reltype:ACTED_IN"`

	Roles          []string
	Name           string
	NumberOfScenes int64
	Actor          *Actor `gogm:"startNode"`
	Movie          *Movie `gogm:"endNode"`
}

func main() {

	var config = &gogm.Config{
		"uri",
		"username",
		"password",
		gogm.NONE/*log level*/,
		false/*allow cyclic ref*/}

	var ogm = gogm.New(config)

	var session, err = ogm.NewSession(true/*isWriteMode*/)
	if err != nil {
		panic(err)
	}

	theMatrix := &Movie{}
	theMatrix.Title = "The Matrix"
	theMatrix.Released = 1999

	carrieAnne := &Actor{}
	carrieAnne.Name = "Carrie-Anne Moss"

	keanu := &Actor{}
	keanu.Name = "Keanu Reeves"

	carrieAnneMatrixCharacter := &Character{Movie: theMatrix, Actor: carrieAnne, Roles: []string{"Trinity"}}
	keanuReevesMatrixCharacter := &Character{Movie: theMatrix, Actor: keanu, Roles: []string{"Neo"}}

	theMatrix.Characters = append(theMatrix.Characters, carrieAnneMatrixCharacter, keanuReevesMatrixCharacter)

	//Persist the movie. This persists the actors as well and creates the relationships with their associated properties
	if err := session.Save(&theMatrix, nil); err != nil {
		panic(err)
	}

	var loadedMatrix *Movie
	if err := session.Load(&loadedMatrix, *theMatrix.ID, nil); err != nil {
		panic(err)
	}

	for _, character := range loadedMatrix.Characters {
		fmt.Println("Actor: " + character.Actor.Name)
	}
}
