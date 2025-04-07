package devices
import (
	"fmt"
	"io"
	"bufio"
	"errors"
	"string"
	"time"
	"github.com/tarm/serial"
)
type Oscilloscope struct {
	serialPort      *serial.Port
	reader		*bufio.Reader
}

// InitShooter инциализирует и создаёт экземпляр Oscilloscope
func InitOscilloscope(portName string, baudRate int) (*Oscilloscope, error) {
	fmt.Println("Initializing Oscilloscope on port", portName)
	config := &serial.Config{
		Name:		portName,
		Baud:		baudRate,
		ReadTimeout:	time.Second
	}
	port,err := serial.OpenPort(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to open port: %v", err)
	}

	oscilloscope := &Oscilloscope {
		serialPort:	port,
		reader:		bufio.NewReader(port)
	}

	response, err := s.SendCommand("*IDN?")
	if err != nil {
		return nil, fmt.Errorf("Asking oscilloscope info failed: %v", err)
	}

	fmt.Println("Oscilloscope info:", strings.TrimSpace(response))
	return oscilloscope, nil
}

// SendCommand отправляет команду осциллографу и возвращает полученный ответ
func (oscilloscope *Oscilloscope) SendCommand(cmd string) (string, error) {
	_, err := oscilloscope.serialPort.Write([]byte(cmd+"\n"))
	if err != nil {
		return "", fmt.Errorf("Command send error: %w", err)
	}

	time.Sleep(300*time.Millisecond)

	response, err := oscilloscope.reader.ReadString('\n')
	if err != nil && err != io.EOF{
		return "", err
	}

	return response, err
}

// Query выполняет команду с запросом данных
func (oscilloscope *Oscilloscope) Query(cmd string) (string, error) {
	return o.SendCommand(cmd)
}

// ReadBytes читает бинарные данные
func (oscilloscope *Oscilloscope) ReadBytes(size int) ([]byte, error) {
	buffer := make([]byte, size)
	_, err := io.ReadFull(oscilloscope.reader, buffer)
	return buf, err
}

