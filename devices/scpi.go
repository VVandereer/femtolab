package devices

import (
	"fmt"
)

// Базовая конфигурация
func (oscilloscope *Oscilloscope) Reset() error {
	_, err := oscilloscope.SendCommand("*RST")
	return err
}

func (oscilloscope *Oscilloscope) GetIDN() (string, error) {
	return oscilloscope.Query("*IDN?")
}

// Управление каналами
func (oscilloscope *Oscilloscope) EnableChannel(ch int) error {
	_, err := oscilloscope.SendCommand(fmt.Sprintf("SELECT:CH%d ON", ch))
	return err
}

func (oscilloscope *Oscilloscope) SetVerticalScale(ch int, scale float64) error {
	_, err := oscilloscope.SendCommand(fmt.Sprintf("CH%d:VOLTS/DIV %f", ch, scale))
	return err
}

// Измерения
func (oscilloscope *Oscilloscope) MeasureVpp(ch int) (string, error) {
	return oscilloscope.Query(fmt.Sprintf("MEASURE:VPP? CH%d", ch))
}

func (oscilloscope *Oscilloscope) MeasureFrequency(ch int) (string, error) {
	return oscilloscope.Query(fmt.Sprintf("MEASURE:FREQUENCY? CH%d", ch))
}

// Управление захватом
func (oscilloscope *Oscilloscope) SingleCapture() error {
	_, err := oscilloscope.SendCommand("ACQUIRE:STOPAFTER SEQUENCE")
	if err != nil {
		return err
	}
	_, err = oscilloscope.SendCommand("ACQUIRE:STATE RUN")
	return err
}

// Настройка waveform
func (oscilloscope *Oscilloscope) ConfigureWaveform(ch int) error {
	commands := []string{
		fmt.Sprintf("DATA:SOURCE CH%d", ch),
		"DATA:WIDTH 2",
		"DATA:ENC RIBINARY",
	}

	for _, cmd := range commands {
		if _, err := oscilloscope.SendCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}

func (oscilloscope *Oscilloscope) GetWaveform() (string, error) {
	return oscilloscope.Query("CURVE?")
}
