package controller

type Controller struct {
	Timeout  int
	Parallel int
	Pulls    int
}

type Report struct {
	FailedLayers int
	Layers       int
	Success      bool
}
