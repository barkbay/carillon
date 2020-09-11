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

type Singleton struct {
	isRunning *uint32
	c         chan gpiod.LineEvent

	processor EventProcessor
}

// NewSingleton takes a event processor as an argument an ensure that only one is running, if not event is discarded.
func NewSingleton(processor EventProcessor) *Singleton {
	s := &Singleton{
		processor: processor,
		isRunning: new(uint32),
		c:         make(chan gpiod.LineEvent, 1),
	}
	go func() {
		for evt := range s.c {
			if atomic.CompareAndSwapUint32(s.isRunning, 0, 1) {
				log.Infof("Start processor %s", s.processor.String())
				go s.processor.OnEventWithCallback(
					evt,
					// Restore the state to let a new goroutine run
					s.callback,
				)
				continue
			}
			log.Infof("Event discarded, goroutine for %s already running", s.processor.String())
		}
	}()
	return s
}

func (s *Singleton) callback() {
	log.Infof("Processor %s ended", s.processor.String())
	atomic.StoreUint32(s.isRunning, 0)
}
func (s *Singleton) String() string {
	return "Singleton{" + s.processor.String() + "}"
}

func (s *Singleton) OnEventWithCallback(evt gpiod.LineEvent, callback func()) {
	s.c <- evt
	callback()
}
