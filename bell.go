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

/* Code below is inspired from https://github.com/leon-anavi/rpi-examples
 * See license here: https://github.com/leon-anavi/rpi-examples/blob/master/LICENSE
 */
func (b *Bell) playNote(note, duration uint32) {
	//This is the semiperiod of each note.
	beepDelay := 1000000 / note
	//This is how much time we need to spend on the note.
	t := (duration * 1000) / (beepDelay * 2)
	var i uint32
	for i = 0; i < t; i++ {
		//1st semiperiod
		b.write(1)
		time.Sleep(time.Duration(beepDelay) * time.Microsecond)
		//2nd semiperiod
		b.write(0)
		time.Sleep(time.Duration(beepDelay) * time.Microsecond)
	}
	//Add a little delay to separate the single notes
	b.write(0)
	time.Sleep(20 * time.Millisecond)
}

func (b *Bell) write(value int) {
	if err := b.line.SetValue(value); err != nil {
		log.Errorf("Error while setting bell line to 0: %w", err)
	}
}

func (b *Bell) OnEventWithCallback(_ gpiod.LineEvent, callback func()) {
	defer func() {
		if err := b.line.SetValue(0); err != nil {
			log.Errorf("Error while setting bell line to 0: %w", err)
		}
		if callback != nil {
			callback()
		}
	}()
	b.playNote(a, 500)
	b.playNote(a, 500)
	b.playNote(f, 350)
	b.playNote(cH, 150)

	b.playNote(a, 500)
	b.playNote(f, 350)
	b.playNote(cH, 150)
	b.playNote(a, 1000)
	b.playNote(eH, 500)

	b.playNote(eH, 500)
	b.playNote(eH, 500)
	b.playNote(fH, 350)
	b.playNote(cH, 150)
	b.playNote(gS, 500)

	b.playNote(f, 350)
	b.playNote(cH, 150)
	b.playNote(a, 1000)
	b.playNote(aH, 500)
	b.playNote(a, 350)

	b.playNote(a, 150)
	b.playNote(aH, 500)
	b.playNote(gHS, 250)
	b.playNote(gH, 250)
	b.playNote(fHS, 125)

	b.playNote(fH, 125)
	b.playNote(fHS, 250)

	time.Sleep(250 * time.Millisecond)

	b.playNote(aS, 250)
	b.playNote(dHS, 500)
	b.playNote(dH, 250)
	b.playNote(cHS, 250)
	b.playNote(cH, 125)

	b.playNote(br, 125)
	b.playNote(cH, 250)

	time.Sleep(250 * time.Millisecond)

	b.playNote(f, 125)
	b.playNote(gS, 500)
	b.playNote(f, 375)
	b.playNote(a, 125)
	b.playNote(cH, 500)

	b.playNote(a, 375)
	b.playNote(cH, 125)
	b.playNote(eH, 1000)
	b.playNote(aH, 500)
	b.playNote(a, 350)

	b.playNote(a, 150)
	b.playNote(aH, 500)
	b.playNote(gHS, 250)
	b.playNote(gH, 250)
	b.playNote(fHS, 125)

	b.playNote(fH, 125)
	b.playNote(fHS, 250)

	time.Sleep(250 * time.Millisecond)

	b.playNote(aS, 250)
	b.playNote(dHS, 500)
	b.playNote(dH, 250)
	b.playNote(cHS, 250)
	b.playNote(cH, 125)

	b.playNote(br, 125)
	b.playNote(cH, 250)

	time.Sleep(250 * time.Millisecond)

	b.playNote(f, 250)
	b.playNote(gS, 500)
	b.playNote(f, 375)
	b.playNote(cH, 125)
	b.playNote(a, 500)

	b.playNote(f, 375)
	b.playNote(c, 125)
	b.playNote(a, 1000)
}
