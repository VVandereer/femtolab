package devices

import (
	"fmt"
	"bufio"
	"strings"
	"time"

	"github.com/tarm/serial"
)

var (
	firmwareInfo		= []byte("?f")
	enablePower		= []byte("e")
	disablePower		= []byte("d")
	movePosition		= []byte("m")
	restorePosition		= []byte("r")
	setPosition		= []byte("p")
	getPosition		= []byte("?p")
	setMaxSpeed		= []byte("s")
	getMaxSpeed		= []byte("?s")
	setAcceleration		= []byte("a")
	getAcceleration		= []byte("?a")
	setLimitSwitchEnableBits= []byte("t")
	setLimitSwitchState	= []byte("i")
	setLimitSwitchPins	= []byte("w")
	getLimitSwitchPosition	= []byte("?l")
)


type StepperMotor struct {
	serialPort		*serial.Port
	reader		*bufio.Reader
}

// InitStepperMotor создаёт и инициализирует новый экземпляр StepperMotor
func InitStepperMotor(portName string, baudRate int) (*StepperMotor, error) {
	fmt.Println("Initializing Stepper Motor on port", portName)
	config := &serial.Config{
		Name:		portName,
		Baud:		baudRate,
		ReadTimeout:	time.Millisecond*500,
	}
	port,err := serial.OpenPort(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to open port: %v", err)
	}

	stepper := &StepperMotor{
		serialPort: port,
		reader: bufio.NewReader(port),
	}
	response, err := stepper.SendCommand("?f")
	if err != nil {
		return nil, fmt.Errorf("Asking firmware info failed: %v", err)
	}

	fmt.Println("firmware info:",strings.TrimSpace(response))
	return stepper, nil
}

// SendCommand отправляет команду шаговому мотору и возвращает полученный ответ
func (stepper *StepperMotor) SendCommand(cmd string) (string, error) {
	_, err := stepper.serialPort.Write([]byte(cmd))
	if err != nil {
		return "", err
	}
	var response strings.Builder
	for {
		line,err := stepper.reader.ReadString('\n')
		if err != nil {
			return response.String(), err
		}
		response.WriteString(line)
		if strings.Contains(line,"#") {
			break
		}
	}
	return response.String(), nil
}

// Enable включение питания моторчика
func (s *StepperMotor) Enable() error {
	_, err := s.SendCommand("e")
	return err
}

// Disable отключение питания моторчика
func (s *StepperMotor) Disable() error {
	_, err := s.SendCommand("d")
	return err
}

// SetAcceleration установка ускорения моторчика
func (s *StepperMotor) SetAcceleration(acceleration int) error {
	cmd := fmt.Sprintf("a%d", acceleration)
	_, err := s.SendCommand(cmd)
	return err
}

// AskAcceleration запрос текущего ускорения
func (s *StepperMotor) AskAcceleration() (string, error) {
	return s.SendCommand("?a")
}

// SetMaxSpeed установка максимальной скорости мотора
func (s *StepperMotor) SetMaxSpeed(speed int) error {
	cmd := fmt.Sprintf("s%d", speed)
	_, err := s.SendCommand(cmd)
	return err
}

// AskMaxSpeed запрос текущего значения максимальной скорости
func (s *StepperMotor) AskMaxSpeed() (string, error) {
	return s.SendCommand("?s")
}

// SetLimitSwitchEnableBits установка битов включения концевых выключателей; N=0,1,2,3
func (s *StepperMotor) SetLimitSwitchEnableBits(bits int) error {
	cmd := fmt.Sprintf("t%d", bits)
	_, err := s.SendCommand(cmd)
	return err
}

// SetLimitSwitchActiveState устанавливает активное состояние концевых выключателей; active=true/false
func (s *StepperMotor) SetLimitSwitchActiveState(active bool) error {
	var value int
	if active{
		value = 1
	}
	cmd := fmt.Sprintf("i%d", value)
	_, err := s.SendCommand(cmd)
	return err
}
// LimitSwitchPins устанавливает значение пинов концевых переключателей
func (s *StepperMotor) LimitSwitchPins(swap bool) error {
	var value int
	if swap {
		value = 1
	}
	cmd := fmt.Sprintf("w%d", value)
	_, err := s.SendCommand(cmd)
	return err
}

// GoTo перемещает мотор в указанную позицию
func (s *StepperMotor) GoTo(position int) error {
	cmd := fmt.Sprintf("p%d", position)
	_, err := s.SendCommand(cmd)
	return err
}

// AskPosition запрашивает текущую позицию мотора
func (s *StepperMotor) AskPosition() (int, error) {
	response, err := s.SendCommand("?p")
	if err != nil {
		return 0,err
	}
	var position int
	_, err = fmt.Sscanf(response, "%d", &position)
	if err != nil {
		return 0, err
	}
	return position, nil
}

// Move перемещает мотор на заданное количество шагов
func (s *StepperMotor) Move(steps int) error {
	cmd := fmt.Sprintf("m%d", steps)
	_, err := s.SendCommand(cmd)
	return err
}

// AskLimitSwitchCalibrationPosition запрос калибровочной позиции концевых выключателей
func (s *StepperMotor) AskLimitSwitchCalibrationPosition() (string, error) {
	return s.SendCommand("?l")
}

// RestorePosition восстановление текущей позиции мотора
func (s *StepperMotor) RestorePosition(position int) error {
	cmd := fmt.Sprintf("r%d", position)
	_, err := s.SendCommand(cmd)
	return err
}
