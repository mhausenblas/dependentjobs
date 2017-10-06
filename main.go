package main

import (
	"fmt"
	"log"
)

func main() {
	dj := New()
	// manually create jobs and call graph:
	dj.Add("root", "job 1", 0)
	dj.Add("j2", "job 2", 1)
	dj.Add("j3", "job 3", 1)
	dj.Add("j4", "job 4", 2)
	dj.AddDependents("j2", "j4")
	dj.AddDependents("j3", "j4")
	dj.AddDependents("root", "j2", "j3")
	err := dj.Dump("./examples/dump.json")
	if err != nil {
		log.Fatal(err)
	}
	// reading call graph from file:
	// err := dj.FromFile("./examples/test.cg")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Printf("%+v\n", dj)
	// dj.Run()
}
