// Пакет devices реализует управление лазерным устройством (HWShooter) через последовательный порт.
package devices

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/tarm/serial"
	
	"strconv"
	"strings"
	"sync"
	"time"
)

// ModeEnum определяет режимы стрельбы.
type ModeEnum int

const (
	// ModeFrequency – режим с фиксированной частотой.
	ModeFrequency ModeEnum = iota
	// ModeExternalSync – режим синхронизации по внешнему сигналу.
	ModeExternalSync
	// ModeSingleShot – одиночный выстрел.
	ModeSingleShot
)

// ExternalModeEnum определяет режимы внешней синхронизации.
type ExternalModeEnum int

const (
	// ExternalIN1 – соответствует внешнему входу 1.
	ExternalIN1 ExternalModeEnum = 1
	// ExternalIN2 – соответствует внешнему входу 2.
	ExternalIN2 ExternalModeEnum = 2
)

// Определение команд – каждая команда представлена массивом байтов.
var (
	singleShoot      = []byte("S")
	singleMode       = []byte("MM")
	externalMode     = []byte("ME")
	freqMode         = []byte("MF")
	setPeriod        = []byte("P")
	externalIN1      = []byte("I1")
	externalIN2      = []byte("I2")
	getCount         = []byte("C?")
	resetCount       = []byte("CR")
	disallowShooting = []byte("D")
	allowShooting    = []byte("A")
	verboseFull      = []byte("VF")
	verboseNone      = []byte("VN")
)

// HWShooter представляет лазерное устройство для стрельбы.
type HWShooter struct {
	// serialPort – соединение с устройством.
	serialPort *serial.Port
	// reader – для чтения строк из последовательного порта.
	reader *bufio.Reader
	// mu — мьютекс для синхронизации доступа к порту.
	mu sync.Mutex

	// shotsCount – количество выстрелов, полученное с устройства.
	shotsCount int
	// enableShooting – флаг, указывающий, разрешена ли стрельба.
	enableShooting bool
	// batchCount – количество выстрелов в одной серии.
	batchCount int

	// period – время ожидания между выстрелами в микросекундах.
	period int
	// maxPeriod и minPeriod – максимальное и минимальное значение периода.
	maxPeriod int
	minPeriod int

	// lastMode и lastExternalMode – текущие режимы стрельбы.
	lastMode         ModeEnum
	lastExternalMode ExternalModeEnum

	// stopChan используется для остановки фоновой горутины обновления счетчика выстрелов.
	stopChan chan struct{}
	// wg – группа ожидания для фонового обновления.
	wg sync.WaitGroup
}

// NewHWShooter создаёт и инициализирует новый экземпляр HWShooter.
// Открывает последовательный порт, очищает буферы и проверяет готовность устройства.
func NewHWShooter(portName string, baudRate int) (*HWShooter, error) {
	// Открываем последовательный порт с заданными параметрами.
	c := &serial.Config{Name: portName, Baud: baudRate, ReadTimeout: time.Millisecond * 500}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть порт: %v", err)
	}

	shooter := &HWShooter{
		serialPort: s,
		reader:     bufio.NewReader(s),
		batchCount: 1,
		period:     1000,    // период по умолчанию в микросекундах
		maxPeriod:  1000000, // максимальное значение периода
		minPeriod:  1,       // минимальное значение периода
		lastMode:   ModeSingleShot,
		// По умолчанию внешний режим – IN1.
		lastExternalMode: ExternalIN1,
		stopChan:         make(chan struct{}),
	}

	// Очищаем буферы порта от старых данных.
	shooter.discardBuffers()

	// Проверяем готовность устройства (до 10 попыток).
	attempt := 10
	var lastLine string
	for attempt > 0 {
		time.Sleep(1 * time.Second)
		shooter.discardBuffers()
		if err := shooter.write(getCount); err != nil {
			attempt--
			continue
		}
		// Читаем строку из устройства.
		line, err := shooter.reader.ReadString('\n')
		if err != nil {
			attempt--
			continue
		}
		lastLine = line
		// Ожидается, что устройство вернет "0\r" в случае сброшенного счетчика.
		if line == "0\r" {
			break
		}
		// Если устройство не готово, сбрасываем счетчик выстрелов.
		_ = shooter.ResetShotsCount()
		attempt--
		if attempt == 0 {
			return nil, fmt.Errorf("устройство на %s не готово; последняя прочитанная строка: %q", portName, lastLine)
		}
	}

	// Устанавливаем минимальный (незаполненный) режим подробности.
	shooter.SetVerboseFull(false)

	// Запускаем фоновую горутину, которая обновляет счетчик выстрелов каждые 200 мс.
	shooter.wg.Add(1)
	go shooter.updateShotsCountLoop()

	// Быстрое переключение, чтобы инициализировать устройство.
	_ = shooter.SetIsEnable(true)
	_ = shooter.SetIsEnable(false)

	return shooter, nil
}

// discardBuffers очищает данные, которые уже находятся в буфере последовательного порта.
// Примечание: пакет tarm/serial не предоставляет явной функции очистки буфера, поэтому читаем доступные данные.
func (h *HWShooter) discardBuffers() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for {
		h.serialPort.Flush() // Если поддерживается, очищаем буферы.
		h.serialPort.SetReadTimeout(10 * time.Millisecond)
		_, err := h.reader.ReadString('\n')
		if err != nil {
			break
		}
	}
	// Восстанавливаем исходный таймаут.
	h.serialPort.SetReadTimeout(500 * time.Millisecond)
}

// updateShotsCountLoop работает в фоне и обновляет количество выстрелов каждые 200 мс.
func (h *HWShooter) updateShotsCountLoop() {
	defer h.wg.Done()
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// Обновляем количество выстрелов.
			_, _ = h.GetShotsCount() // Ошибки здесь для простоты игнорируются.
		case <-h.stopChan:
			return
		}
	}
}

// write выполняет запись данных в последовательный порт с использованием мьютекса.
func (h *HWShooter) write(data []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.serialPort.Write(data)
	return err
}

// SetVerboseFull включает или выключает детализированный вывод (verbose mode).
func (h *HWShooter) SetVerboseFull(enable bool) error {
	cmd := verboseNone
	if enable {
		cmd = verboseFull
	}
	return h.write(cmd)
}

// SetModeSingle переключает устройство в режим одиночного выстрела.
func (h *HWShooter) SetModeSingle() error {
	return h.write(singleMode)
}

// SetModeExternal переключает устройство в режим внешней синхронизации.
func (h *HWShooter) SetModeExternal() error {
	return h.write(externalMode)
}

// SetModeFreq переключает устройство в режим с фиксированной частотой.
func (h *HWShooter) SetModeFreq() error {
	return h.write(freqMode)
}

// SetPeriod отправляет команду установки периода стрельбы (в микросекундах).
func (h *HWShooter) SetPeriod(periodInUS int) error {
	if periodInUS < h.minPeriod || periodInUS > h.maxPeriod {
		return errors.New("период вне допустимого диапазона")
	}
	// Отправляем команду "P", затем ASCII-представление нового периода.
	if err := h.write(setPeriod); err != nil {
		return err
	}
	return h.write([]byte(fmt.Sprintf("%d", periodInUS)))
}

// SetExternalIN1 устанавливает внешний режим IN1.
func (h *HWShooter) SetExternalIN1() error {
	return h.write(externalIN1)
}

// SetExternalIN2 устанавливает внешний режим IN2.
func (h *HWShooter) SetExternalIN2() error {
	return h.write(externalIN2)
}

// SetIsEnable включает или отключает возможность стрельбы.
func (h *HWShooter) SetIsEnable(enable bool) error {
	cmd := disallowShooting
	if enable {
		cmd = allowShooting
	}
	return h.write(cmd)
}

// GetShotsCount запрашивает и возвращает текущий счетчик выстрелов с устройства.
func (h *HWShooter) GetShotsCount() (int, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	// Очищаем старые данные.
	h.discardBuffers()
	if _, err := h.serialPort.Write(getCount); err != nil {
		return 0, err
	}
	// Читаем строку с ответом.
	line, err := h.reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return 0, err
	}
	h.shotsCount = count
	return count, nil
}

// ResetShotsCount сбрасывает счетчик выстрелов на устройстве.
func (h *HWShooter) ResetShotsCount() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.write(resetCount); err != nil {
		return err
	}
	h.shotsCount = 0
	return nil
}

// Shoot выполняет команду стрельбы, если устройство в режиме одиночного выстрела.
// Отправляет команду выстрела batchCount раз с задержкой, заданной периодом.
func (h *HWShooter) Shoot() error {
	// Проверяем, разрешена ли стрельба.
	if !h.enableShooting {
		return errors.New("стрельба отключена")
	}
	// Если устройство не находится в режиме одиночного выстрела, выходим.
	if h.lastMode != ModeSingleShot {
		return nil
	}
	// Преобразуем период (в микросекундах) в миллисекунды для задержки.
	waitDuration := time.Duration(h.period/1000) * time.Millisecond
	// Отправляем команду выстрела batchCount раз.
	for i := 0; i < h.batchCount; i++ {
		if i > 0 {
			time.Sleep(waitDuration)
		}
		if err := h.write(singleShoot); err != nil {
			return err
		}
	}
	return nil
}

// Close останавливает фоновую горутину обновления и закрывает последовательный порт.
func (h *HWShooter) Close() error {
	// Останавливаем фоновую горутину.
	close(h.stopChan)
	h.wg.Wait()
	return h.serialPort.Close()
}
