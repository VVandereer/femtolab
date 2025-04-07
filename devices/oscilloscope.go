package devices
import (
    "fmt"
    "bufio"
    "errors"
    "strconv"
    "string"
    "time"
    "sync"
    "github.com/tarm/serial"

type Oscilloscope struct {
    serialPort      *serial.Port
    numChannels     int

    }