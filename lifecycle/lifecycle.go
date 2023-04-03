package lifecycle

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type LifecycleToken struct {
	Ctx         context.Context
	cancel      context.CancelFunc
	signalChan  chan os.Signal
	handlerList []ShutdownHandler
	lock        sync.Mutex
}

type ShutdownHandler func()

func InitializeLifecycle(baseContext context.Context, terminationSignals []syscall.Signal) *LifecycleToken {
	ctx, cancel := context.WithCancel(baseContext)
	signalChan := make(chan os.Signal)

	if len(terminationSignals) == 0 {
		terminationSignals = append(terminationSignals, syscall.SIGTERM)
	}

	for _, sig := range terminationSignals {
		signal.Notify(signalChan, sig)
	}

	token := LifecycleToken{
		Ctx:        ctx,
		cancel:     cancel,
		signalChan: signalChan,
	}

	go token.blockUntilTerminationSignal()

	return &token
}

func (t *LifecycleToken) RegisterShutdownHandler(handler ShutdownHandler) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.handlerList = append(t.handlerList, handler)
}

func (t *LifecycleToken) TerminateLifecycle() {
	t.signalChan <- syscall.SIGTERM
}

func (t *LifecycleToken) blockUntilTerminationSignal() {
	defer t.cancel()

	select {
	case sig := <-t.signalChan:
		log.Debugf("received signal %d", sig)
		t.lock.Lock()
		defer t.lock.Unlock()
		for _, handler := range t.handlerList {
			handler()
		}
	case <-t.Ctx.Done():
		log.Warn("non-graceful shutdown detected")
	}
}
