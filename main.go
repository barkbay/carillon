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
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/warthog618/gpiod"

	log "github.com/sirupsen/logrus"
)

func init() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "2006-01-02T15:04:05.999999999Z07:00"
	Formatter.FullTimestamp = true
	Formatter.ForceColors = true
	log.SetFormatter(Formatter)
	log.SetLevel(log.DebugLevel)
}

func main() {
	sigs := make(chan os.Signal, 1)
	stop := make(chan struct{}, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		close(stop)
	}()

	chip, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		panic(err)
	}
	defer func() {
		log.Info("Closing chip")
		chip.Close()
	}()

	sensor := NewGPIOEventSensor(
		chip,
		stop,
		NewSingleton(NewBell(chip)),
	)
	defer sensor.Close()

	// Listen for events and publish events to processors
	sensor.Start()

	<-stop
	fmt.Println("exiting...")
}
