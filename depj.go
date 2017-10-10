package main

import (
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

// DependentJobs represents a call graph, that is,
// a collection of Kube Jobs or CronJobs with dependencies.
type DependentJobs struct {
	wg   *sync.WaitGroup `yaml:"-"`
	jobs map[string]Job  `yaml:"jobs"`
}

// New creates a new call graph.
func New() DependentJobs {
	dj := DependentJobs{}
	dj.jobs = make(map[string]Job)
	dj.wg = &sync.WaitGroup{}
	return dj
}

// FromFile reads a call graph from a YAML file.
func (dj *DependentJobs) FromFile(cgfile string) error {
	yamlFile, err := ioutil.ReadFile(cgfile)
	if err != nil {
		return err
	}
	spec := make(map[string]Job)
	err = yaml.Unmarshal(yamlFile, &spec)
	if err != nil {
		return err
	}
	// generate DJ out of spec
	for id, job := range spec {
		dj.Add(id, job.Name, countupstream(spec, id))
		dj.AddDependents(id, spec[id].Dependents...)
	}
	return nil
}

func countupstream(jobs map[string]Job, jobid string) int {
	numupstream := 0
	switch jobid {
	case "root":
		numupstream = 0
	default:
		for _, j := range jobs {
			if contains(j.Dependents, jobid) {
				numupstream++
			}
		}
	}
	return numupstream
}

func contains(list []string, element string) bool {
	for _, e := range list {
		if e == element {
			return true
		}
	}
	return false
}

// Store stores a call graph into a YAML file.
func (dj DependentJobs) Store(cgfile string) error {
	bytes, err := yaml.Marshal(dj.jobs)
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
	j.adddep(depjobs...)
	dj.jobs[id] = j
}

// Lookup retrieves a job by ID.
func (dj DependentJobs) Lookup(id string) Job {
	return dj.jobs[id]
}

// GoString return a canonical string represenation of a dependent job
func (dj DependentJobs) GoString() string {
	res := ""
	for _, j := range dj.jobs {
		res = fmt.Sprintf("%s%#v\n", res, j)
	}
	return res
}
