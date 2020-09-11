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
	"time"

	"github.com/warthog618/gpiod/device/rpi"

	log "github.com/sirupsen/logrus"
	"github.com/warthog618/gpiod"
)

type Light struct {
	chip *gpiod.Chip
	line *gpiod.Line
}

func NewLight(chip *gpiod.Chip) EventProcessor {
	l := &Light{
		chip: chip,
	}
	offset := rpi.GPIO26
	var err error
	if l.line, err = chip.RequestLine(offset, gpiod.AsOutput(0)); err != nil {
		panic(err)
	}
	return l
}

func (l *Light) String() string {
	return "Light"
}

func (l *Light) OnEventWithCallback(_ gpiod.LineEvent, callback func()) {
	defer func() {
		log.Info("Turn off light")
		if callback != nil {
			callback()
		}
	}()

	log.Info("Turn on light")
	for j := 0; j < 4*5; j++ {
		_ = l.line.SetValue(1)
		time.Sleep(250 * time.Millisecond)
		_ = l.line.SetValue(0)
		time.Sleep(250 * time.Millisecond)
	}
}
