package reflector

import (
	"sync"

	"k8s.io/client-go/tools/cache"
)

// Controller controls start and stop of multiple reflectors
type Controller struct {
	reflectors []*cache.Reflector
	mu         sync.Mutex
	running    bool
	stopCh     chan struct{}
}

// NewController creates a new Controller responsible for starting and stopping the reflectors specified in the input parameter
func NewController(reflectors ...*cache.Reflector) *Controller {
	return &Controller{
		reflectors: reflectors,
	}
}

// Start all reflectors referenced by this controller, or no-op if already started.
func (c *Controller) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.running {
		return
	}
	c.stopCh = make(chan struct{})
	for _, r := range c.reflectors {
		go r.Run(c.stopCh)
	}
}

//Stop all reflectors referenced by this controller, or no-op if not running.
func (c *Controller) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.running {
		return
	}
	close(c.stopCh)
}
