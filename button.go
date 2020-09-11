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
	"os"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

type EventProcessor interface {
	OnEventWithCallback(evt gpiod.LineEvent, callback func())
	String() string
}

type GPIOEventSensor struct {
	chip *gpiod.Chip
	line *gpiod.Line

	lastEvent time.Time

	stop       chan struct{}
	eventQueue chan gpiod.LineEvent
	ep         []EventProcessor
}

func NewGPIOEventSensor(chip *gpiod.Chip, stop chan struct{}, ep ...EventProcessor) *GPIOEventSensor {
	sensor := &GPIOEventSensor{
		chip:       chip,
		stop:       stop,
		ep:         ep,
		eventQueue: make(chan gpiod.LineEvent, 1024),
	}
	var err error
	if sensor.line, err = requestLine(sensor); err != nil {
		log.Error("RequestLine returned error: %w", err)
		if err == syscall.Errno(22) {
			log.Info("Note that the WithPullUp option requires kernel V5.5 or later - check your kernel version.")
		}
		os.Exit(1)
	}
	return sensor
}

func requestLine(g *GPIOEventSensor) (*gpiod.Line, error) {
	offset := rpi.GPIO2
	return g.chip.RequestLine(offset,
		gpiod.AsIs,
		gpiod.WithFallingEdge(func(evt gpiod.LineEvent) {
			log.Printf("event:%3d %-7s %s", evt.Offset, evt.Timestamp)
			g.eventQueue <- evt
		}))
}

func (g *GPIOEventSensor) Close() error {
	if g.line != nil {
		log.Info("Closing line")
		_ = g.line.Close()
	}
	return nil
}

func (g *GPIOEventSensor) Start() {
	go func() {
		// fetch event from the queue
		for {
			select {
			case <-g.stop:
				log.Info("GPIOEventSensor exiting")
				return
			case evt := <-g.eventQueue:
				now := time.Now()
				// Check if an event has bot been processed in the last duration
				nextEvent := g.lastEvent.Add(5 * time.Second)
				if now.Before(nextEvent) {
					log.Infof("Skip button, next one allowed at %s (now=%s)", nextEvent, now)
					continue
				}
				log.Infof("Allowed button, next one allowed at %s (now=%s)", nextEvent, now)
				g.lastEvent = now
				g.fireProcessors(evt)
			}
		}
	}()
}

func (g *GPIOEventSensor) fireProcessors(evt gpiod.LineEvent) {
	log.Infof("Firing event %v", evt)
	for _, processor := range g.ep {
		processor := processor
		processor.OnEventWithCallback(evt, func() { log.Infof("Event %v sent to %s", evt, processor) })
	}
	log.Infof("Event %v sent to ALL the processors", evt)
}
