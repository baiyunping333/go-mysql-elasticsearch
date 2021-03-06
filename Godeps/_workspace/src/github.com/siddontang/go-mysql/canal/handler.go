package canal

import (
	"errors"

	"github.com/siddontang/go/log"
)

var (
	ErrHandleInterrupted = errors.New("do handler error, interrupted")
)

type RowsEventHandler interface {
	// Handle RowsEvent, if return ErrHandleInterrupted, canal will
	// stop the sync
	Do(e *RowsEvent) error
	String() string
}

func (c *Canal) RegRowsEventHandler(h RowsEventHandler) {
	c.rsLock.Lock()
	c.rsHandlers = append(c.rsHandlers, h)
	c.rsLock.Unlock()
}

func (c *Canal) travelRowsEventHandler(e *RowsEvent) error {
	c.rsLock.Lock()
	defer c.rsLock.Unlock()

	var err error
	for _, h := range c.rsHandlers {
		if err = h.Do(e); err != nil && err != ErrHandleInterrupted {
			log.Errorf("handle %v err: %v", h, err)
		} else if err == ErrHandleInterrupted {
			log.Errorf("handle %v err, interrupted", h)
			return err
		}

	}
	return nil
}
