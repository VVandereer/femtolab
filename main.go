package main

import (
	"femtolab/devices"
	"fmt"
	"log"
)

func main() {
	DelayLine, err := devices.InitStepperMotor("COM27", 115200)
	if err != nil {
		log.Fatal(err)
	}

	// Включаем моторчик
	if err := DelayLine.Enable(); err != nil {
		log.Fatal(err)
	}

	// Играемся с параметрами
	if err := DelayLine.SetAcceleration(6400); err != nil {
		log.Fatal(err)
	}
	if err := DelayLine.SetMaxSpeed(3200); err != nil {
		log.Fatal(err)
	}
	// Перемещаем мотор
	if err := DelayLine.Move(1000); err != nil {
		log.Fatal(err)
	}

	// Запрашиваем текущую позицию моторчика
	position, err := DelayLine.AskPosition()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Текущая позиция мотора: %d\n", position)

	// Отключаем мотор
	if err := DelayLine.Disable(); err != nil {
		log.Fatal(err)
	}
}
