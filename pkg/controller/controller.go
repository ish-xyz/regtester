package controller

import (
	"time"

	"github.com/ish-xyz/regtester/pkg/config"
)

var conf config.Config

// Return a new configuration instance
func NewController(c config.Config) Controller {
	conf = c
	return Controller{}
}

//Run starts the controller which
//		will coordinate the docker pulls
func (ctrl *Controller) Run() error {

	maxParallel := make(chan int, conf.Workload.Parallel)
	for x := 0; x <= conf.Workload.Pulls; x++ {
		maxParallel <- x
		triggerDockerPull(maxParallel)
	}
	return nil
}

func pickRegistry() {
	return
}

func pickImage() {
	return
}

func triggerDockerPull(maxParallel chan int) {
	fakeDockerpull()
	<-maxParallel
}

func fakeDockerpull() {
	time.Sleep(time.Duration(10) * time.Second)
	/*
		type Report struct {
			FailedLayers int
			Layers int
			Success bool
		}
	*/
}
