package main

import (
    "log"
	"femtolab/devices"
	"femtolab/logger"
)

func main() {
    Logger.InitLog("", true)
    log.Println("FemtoLab was starting")

	DelayLine := devices.InitStepperMotor("COM27", 115200)
	DelayLine.Enable()
    DelayLine.SetAcceleration(6400)
	DelayLine.SetMaxSpeed(3200)
	DelayLine.Move(1000)
	position := DelayLine.AskPosition()
	fmt.Printf("Текущая позиция мотора: %d\n", position)
	DelayLine.Disable()
}
