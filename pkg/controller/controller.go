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

	maxConcPulls := make(chan int, conf.Workload.MaxConcPulls)
	for x := 0; x <= conf.Workload.Pulls; x++ {
		maxConcPulls <- x
		go triggerDockerPull(maxConcPulls)
	}
	return nil
}

/*
func selectRegistry() {
	return
}

func selectImage() {
	return
}
*/

func triggerDockerPull(maxConcPulls chan int) {
	fakeDockerPull()
	<-maxConcPulls
}

func fakeDockerPull() {
	time.Sleep(time.Duration(10) * time.Second)
	/*
		type Report struct {
			FailedLayers int
			Layers int
			Success bool
		}
	*/
}
