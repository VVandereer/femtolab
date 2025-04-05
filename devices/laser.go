package device

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"strconv"
	"time"
	"github.com/tarm/serial"
)

var (
	getCount		= []byte("C?")
	resetCount		= []byte("CR")
	dissallowShooting	= []byte("D")
	allowShooting		= []byte("A")
	SetPeriod		= []byte("P")
	verboseFull		= []byte("VF")
	verboseNone		= []byte("VN")
	singleShoot		= []byte("S")
	singleMode		= []byte("MM")
	freqMode		= []byte("MF")
	externalMode		= []byte("ME")
	externalIN1		= []byte("I1")
	externalIN2		= []byte("I2")
)

type Shooter struct  {
	serialPort	*seral.Port
	reader		*bufio.Reader
	enableShooting	bool
	batchCount	int
	period		int
	minPeriod	int
	maxPeriod	int
	mode		string
}

// InitShooter инциализирует и создаёт экземпляр Shooter
func InitShooter(portName string, baudRate int) (*Shooter, error) {
	fmt.Println("Initializing Shooter on port", portName)
	config := &serial.Config{
		Name:		portName,
		Baud:		baudRate,
		ReadTimeout:	time.Millisecond*500
	}
	port,err := serial.OpenPort(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to open port: %v", err)
	}

	shooter := &Shooter {
		serialPort:	port,
		reader:		bufio.NewReader(port),
		batchCount:	1,
		period:		1000, //microseconds
		minPeriod:	1,
		maxPeriod:	1000000,
		mode:		"single"
	}

	attempt := 10
	for attempt > 0 {
		time.Sleep(1 * time.Second)
		shooter.clearBuffer()
		if err := shooter.write(getCount); err != nil {
			attempt--
			continue
		}
		line, err := shooter.reader.ReadString('\n')
		if err != nil {
			attempt--
			continue
		}
		lastLine := line
		if line == '0\r' {
			break
		}
		_ = shooter.ResetShotsCount()
		attempt--
		if attempt == 0 {
			return nil, fmt.Errorf("Shooter on &s is not ready; last readed line: %q", portName, lastLine)
		}
	}

	fmt.Println("Successed initializing Shooter on port ", portName)
	return shooter, nil
}
// ClearBuffer очищает буффер устройства
func (shooter *Shooter) ClearBuffer() {
	for {
		shooter.serialPort.Flush()
		shooter.serialPort.SetReadTimeout(10*Millisecond)
		_,err := shooter.reader.ReadString('\n')
		if err != nil {
			break
		}
	}
	shooter.serialPort.SetReadTimeout(500*time.Millisecond)
}

// GetShotsCount запрашивает и возвращает количество произведённых выстрелов с устройства
func (shooter *Shooter) GetShotsCount() (int, error) {
	shooter.clearBuffer()
	if _, err := shooter.serialPort.Write(getCount); err != nil {
		return 0, err
	}
	answer, err := shooter.reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(strings.TrimSpace(answer))
	if err != nil {
		return 0, err
	}
	return count, nil
}
// ResetShotsCount сбрасывает счётчик произведённых выстреллов на устройстве
func (shooter *Shooter) ResetShotsCount() error {
	_,err := shooter.serialPort.Write(resetCount)
	return err
}
// Shoot - производит batchCount выстрелов с заданным периодом
func (shooter *Shooter) Shoot() error {
	if !shooter.enableShooting {
		return errors.New("Shooting is unable")
	}
	if shooter.Mode != "Single" {
		return nil
	}
	waitDuration := time.Duration(shooter.period/1000)*time.Millisecond
	for i := 0; i< h.batchCount; i++ {
		if i>0 {
			time.Sleep(waitDuration)
		}
		if _, err := shooter.serial.Port.Write(singleShoot); err != nil {
			return err
		}
	}
}
// SetEnable - переключение предохранителя
func (shooter *Shooter) SetEnable(enable bool) error {
	cmd:= dissallowShooting
	if enable {
		cmd = allowShooting
	}
	 _, err := shooter.serialPort.Write(cmd)
	return err
}
// SetPeriod - устанавливает период стрельбы (в микросекундах)
func (shooter *Shooter) SetPeriod(period int) error {
	if period < shooter.minPeriod || shooter.maxPeriod < period {
		return errors.New("Period out of range")
	}
	if _,err := shooter.serialPort.Write(SetPeriod); err != nil {
		return err
	}
	if _, err := shooter.serialPort.Write([]byte(fmt.Sprintf("%d",period))); err != nil {
		return err
	}
	shooter.period = period
	return nil
}
// SetVerboseFull - установка подробности вывода
func (shooter *Shooter) SetVerboseFull(enable bool) error {
	cmd:= verboseNone
	if enable {
		cmd = verboseFull
	}
	_, err := shooter.serialPort.Write(cmd)
	return err
}
// SetMode - переключение режима
func (shooter *Shooter) SetMode(mode string) error {} {
	switch mode {
		case "single":
			cmd = singleMode
		case "freq":
			cmd = freqMode
		case "external":
			cmd = externalMode
		default:
			return errors.New("Unknown shooting mode")
	}
	_, err := shooter.serialPort.Write(cmd)
	return err
}
// Close - просто закрывает порт
func (shooter *Shooter) Close() error {
	shooter.serialPort.Close()
}
