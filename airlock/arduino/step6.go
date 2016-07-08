package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/firmata"
)

var button *gpio.GroveButtonDriver
var blue *gpio.GroveLedDriver
var green *gpio.GroveLedDriver
var red *gpio.GroveLedDriver
var buzzer *gpio.GroveBuzzerDriver
var touch *gpio.GroveTouchDriver
var rotary *gpio.GroveRotaryDriver

func Doorbell() {
	TurnOff()
	blue.On()
	buzzer.Tone(gpio.C4, gpio.Quarter)
	<-time.After(1 * time.Second)
	Reset()
}

func TurnOff() {
	blue.Off()
	green.Off()
}

func Reset() {
	TurnOff()
	fmt.Println("Airlock ready.")
	green.On()
}

func main() {
	gbot := gobot.NewGobot()

	a := api.NewAPI(gbot)
	a.Start()

	board := firmata.NewFirmataAdaptor("arduino", os.Args[1])

	// digital devices
	button = gpio.NewGroveButtonDriver(board, "button", "2")
	blue = gpio.NewGroveLedDriver(board, "blue", "3")
	green = gpio.NewGroveLedDriver(board, "green", "4")
	red = gpio.NewGroveLedDriver(board, "red", "5")
	buzzer = gpio.NewGroveBuzzerDriver(board, "buzzer", "7")
	touch = gpio.NewGroveTouchDriver(board, "touch", "8")

	// analog devices
	rotary = gpio.NewGroveRotaryDriver(board, "rotary", "0")

	work := func() {
		Reset()

		gobot.On(button.Event(gpio.Push), func(data interface{}) {
			TurnOff()
			fmt.Println("On!")
			blue.On()
		})

		gobot.On(button.Event(gpio.Release), func(data interface{}) {
			Reset()
		})

		gobot.On(touch.Event(gpio.Push), func(data interface{}) {
			Doorbell()
		})

		gobot.On(rotary.Event("data"), func(data interface{}) {
			fmt.Println("rotary", data)
		})
	}

	robot := gobot.NewRobot("airlock",
		[]gobot.Connection{board},
		[]gobot.Device{button, blue, green, buzzer, touch, rotary},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}