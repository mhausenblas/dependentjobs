package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

// DependentJobs represents a call graph, that is,
// a collection of Kube Jobs or CronJobs with dependencies.
type DependentJobs struct {
	wg   *sync.WaitGroup
	jobs map[string]Job
}

// New creates a new call graph.
func New() DependentJobs {
	dj := DependentJobs{}
	dj.jobs = make(map[string]Job)
	dj.wg = &sync.WaitGroup{}
	return dj
}

// FromFile reads a call graph from a YAML manifest file.
func (dj *DependentJobs) FromFile(cgfile string) error {
	yamlFile, err := ioutil.ReadFile(cgfile)
	if err != nil {
		return err
	}
	var spec []struct {
		ID         string   `yaml:"id"`
		Name       string   `yaml:"name"`
		Dependents []string `yaml:"dependents"`
	}
	err = yaml.Unmarshal(yamlFile, &spec)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", spec)
	// generate DJ out of spec
	// Add for all, then deps with root last
	// for _, e := range spec {
	// 	if e.ID == "root"
	// }

	return nil
}

// Dump stores a raw call graph into a file as JSON.
func (dj DependentJobs) Dump(cgfile string) error {
	bytes, err := json.Marshal(dj.jobs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cgfile, bytes, 0644)
}

// Run takes a call graph of jobs
// and runs it in order of its dependencies.
func (dj DependentJobs) Run() {
	dj.wg.Add(len(dj.jobs))
	go dj.jobs["root"].launch(dj, dj.wg)
	dj.wg.Wait()
}

// Add adds a job to the call graph.
func (dj *DependentJobs) Add(id, name string, numupstream int) {
	j := newjob(id, name, numupstream)
	dj.jobs[id] = j
}

// AddDependents adds one or more dependent jobs to a job.
func (dj DependentJobs) AddDependents(id string, depjobs ...string) {
	j := dj.jobs[id]
	depj := []Job{}
	for _, d := range depjobs {
		depj = append(depj, dj.jobs[d])
	}
	j.adddep(depj...)
	dj.jobs[id] = j
}

// Lookup retrieves a job by ID.
func (dj DependentJobs) Lookup(id string) Job {
	return dj.jobs[id]
}
