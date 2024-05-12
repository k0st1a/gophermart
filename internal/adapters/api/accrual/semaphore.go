package accrual

import (
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

type Blocker interface {
	Activate(wait time.Duration)
	IsActive() bool
}

type block struct {
	active atomic.Bool
}

func NewBlock() Blocker {
	return &block{}
}

func (b *block) Activate(wait time.Duration) {
	b.active.Store(true)
	log.Printf("Active block for %s seconds", wait)
	time.Sleep(wait)
	b.active.Store(false)
	log.Printf("Deactivate block after %s seconds", wait)
}

func (b *block) IsActive() bool {
	return b.active.Load()
}
