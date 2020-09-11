/* This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 3.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"sync/atomic"

	log "github.com/sirupsen/logrus"

	"github.com/warthog618/gpiod"
)

var (
	notRunning uint32 = 0
)

type Singleton struct {
	isRunning *uint32
	c         chan gpiod.LineEvent

	processor EventProcessor
}

// NewSingleton takes a event processor as an argument an ensure that only one is running, if not event is discarded.
func NewSingleton(processor EventProcessor) *Singleton {
	s := &Singleton{
		processor: processor,
		isRunning: &notRunning,
		c:         make(chan gpiod.LineEvent, 1),
	}
	go func() {
		for evt := range s.c {
			if atomic.CompareAndSwapUint32(s.isRunning, 0, 1) {
				s.OnEvent(evt)
				// Restore the state to let a new goroutine run
				atomic.StoreUint32(s.isRunning, 0)
				continue
			}
			log.Info("Event discarded, goroutine already running")
		}
	}()
	return s
}

func (s *Singleton) OnEvent(evt gpiod.LineEvent) {
	s.processor.OnEvent(evt)
}
