package system

import (
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

var rwLock sync.RWMutex
var callbacks []Callback
var registeredFlag uint32 = 0

type Callback interface {
	NotifyInterrupt()
}

func RegisterSignalHandler() bool {
	isRegistered := atomic.LoadUint32(&registeredFlag) == 1

	if isRegistered {
		return false
	}

	if atomic.CompareAndSwapUint32(&registeredFlag, 0, 1) {
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

		go func(chan os.Signal) {
			<-ch
			rwLock.RLock()
			defer rwLock.RUnlock()

			if len(callbacks) > 0 {
				for index := range callbacks {
					callbacks[index].NotifyInterrupt()
				}
			}
		}(ch)

		return true
	}

	return false
}

func RegisterCallback(callback Callback) {
	rwLock.Lock()
	callbacks = append(callbacks, callback)
	rwLock.Unlock()
}
