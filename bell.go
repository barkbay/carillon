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

	log "github.com/sirupsen/logrus"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

type Bell struct {
	chip *gpiod.Chip
	line *gpiod.Line
}

func NewBell(chip *gpiod.Chip) EventProcessor {
	b := &Bell{
		chip: chip,
	}
	offset := rpi.GPIO23
	var err error
	if b.line, err = chip.RequestLine(offset, gpiod.AsOutput(0)); err != nil {
		panic(err)
	}
	return b
}

func (b *Bell) String() string {
	return "Bell"
}

func (b *Bell) OnEventWithCallback(_ gpiod.LineEvent, callback func()) {
	log.Info("BEGIN BELL")
	defer func() {
		if err := b.line.SetValue(0); err != nil {
			log.Errorf("Error while setting bell line to 0: %w", err)
		}
		log.Info("END BELL")
		if callback != nil {
			callback()
		}
	}()

	for i := 0; i < 2; i++ {
		_ = b.line.SetValue(1)
		time.Sleep(500 * time.Millisecond)
		for j := 0; j < 4; j++ {
			_ = b.line.SetValue(0)
			time.Sleep(250 * time.Millisecond)
			_ = b.line.SetValue(1)
			time.Sleep(250 * time.Millisecond)
			_ = b.line.SetValue(0)
		}
	}
}
