package main

func main() {
	dj := New()
	dj.Add("job 1", 0)
	dj.Add("job 2", 1)
	dj.Add("job 3", 1)
	dj.Add("job 4", 2)
	dj.AddDependents(1, 3)
	dj.AddDependents(2, 3)
	dj.AddDependents(0, 1, 2)
	dj.Run()
}
