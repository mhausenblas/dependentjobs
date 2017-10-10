package main

import (
	"fmt"
	"log"
)

func main() {
	// manually create call graph:
	dj := New()
	dj.Add("root", "job 1", 0)
	dj.Add("j2", "job 2", 1)
	dj.Add("j3", "job 3", 1)
	dj.Add("j4", "job 4", 2)
	dj.Add("j5", "job 5", 2)
	dj.AddDependents("j4", "j5")
	dj.AddDependents("j2", "j4")
	dj.AddDependents("j3", "j4")
	dj.AddDependents("root", "j2", "j3", "j5")
	fmt.Printf("%#v\n", dj)

	// store the call graph:
	err := dj.Store("./examples/dump.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// read call graph from file:
	dj = New()
	err = dj.FromFile("./examples/dump.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// run the call graph and print the call sequence:
	fmt.Println("Running jobs in call graph:")
	dj.Run()
	fmt.Printf("Call sequence: %v\n", dj.CallSeq())
}
