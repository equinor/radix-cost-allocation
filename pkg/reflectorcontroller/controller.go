package reflectorcontroller

import (
	"sync"

	"k8s.io/client-go/tools/cache"
)

// ReflectorController controls start and stop of multiple reflectors
type ReflectorController struct {
	reflectors []*cache.Reflector
	mu         sync.Mutex
	running    bool
	stopCh     chan struct{}
}

// New creates a new ReflectorController responsible for starting and stopping the reflectors specified in the input parameter
func New(reflectors ...*cache.Reflector) *ReflectorController {
	return &ReflectorController{
		reflectors: reflectors,
	}
}

// Start all reflectors referenced by this controller, or no-op if already started.
func (c *ReflectorController) Start() {
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
func (c *ReflectorController) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.running {
		return
	}
	close(c.stopCh)
}
