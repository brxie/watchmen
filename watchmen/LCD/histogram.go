package LCD

import (
    "time"
    "math"
    "strconv"
	"fmt"
)

type histogram struct {
    lcdWrapper      *LcdWrapper
    position        byte
    height          byte
    view            *view
    dataCollector   *dataCollector
}

type dataCollector struct {
    stakes          []int
    lastStakeId     int
    startTime       time.Time
}

type view struct {
    barWidth        int
    sepSize         int
    leftMargin      int
}

func NewHistogram(lcd *LcdWrapper, position, height byte) *histogram {
    return &histogram {
        lcdWrapper:     lcd,
        position:       position,
        height:         height,
        view:           &view {
            barWidth:       10,
            sepSize:        3,
            leftMargin:     20,
        },
        dataCollector:  &dataCollector {
            startTime:  time.Now(),
            stakes:     make([]int, 8),
        },
    }
}

func (dc *dataCollector) fill(count int) {
    dc.stakes[len(dc.stakes) - 1]++
}

func (dc *dataCollector) updateSeries() {
    elapsed := time.Since(dc.startTime.Round(time.Second)).Hours()
    timeUnits := 1
    stakeId := int(int(elapsed) / timeUnits) 
    if stakeId > dc.lastStakeId {
        seriesNo := stakeId - dc.lastStakeId 
        if seriesNo > len(dc.stakes) {
            seriesNo = len(dc.stakes)
        }
        dc.stakes = append(dc.stakes[seriesNo:], make([]int, seriesNo)[:]...)
    }
    dc.lastStakeId = stakeId
}

func (h *histogram) drawBar(size, position int) {
    var bar int64
    // let's glue bits one by one to create binary 
    // representation of the bar piece
    for i := 0; i < size; i++ {
        bar += int64(math.Pow(2, float64(i)))
    }

    // split big int into 8-bit chunks
    barSlices := make([]byte, h.height)
    for i := 0; i <= int(h.height - 1); i++ {
        slice := byte(bar >> (uint(i) * 8) & 0xFF)
        slice = h.reverseBits(slice)
        barSlices[i] = slice
    }  
    
    for key, slice := range barSlices {
        pos := int(h.height + h.position - byte(key)) * h.lcdWrapper.Width
        pos += (position * h.view.barWidth) + (position * h.view.sepSize)
        pos += h.view.barWidth + h.view.leftMargin
        for i := 1; i <= h.view.barWidth; i++ {
            h.lcdWrapper.Dev.WriteData(slice, pos - (i - 1))
        }
    }
}

func (h *histogram) reverseBits(b byte) byte {
    // return bits in reversed order - 00000111 -> 11100000 
    b = (b & 0xF0) >> 4 | (b & 0x0F) << 4;
    b = (b & 0xCC) >> 2 | (b & 0x33) << 2;
    b = (b & 0xAA) >> 1 | (b & 0x55) << 1;
    return b
}

func (h *histogram) plot() {
    h.dataCollector.updateSeries()
    max := h.max(h.dataCollector.stakes)
    for key, val := range h.dataCollector.stakes {
        if max == 0 {
            break
        }

        point := float32(h.height * 8.0) / float32(max)
        h.drawBar(int(float32(val) * point), key)
    }
    h.drawBoundVal()
    h.drawScale()
    h.lcdWrapper.Dev.Display()
}

func (h *histogram) drawBoundVal() {
    // draw max and min values next to histogram
    max := h.max(h.dataCollector.stakes)
    h.lcdWrapper.DisplayString("   ", int(h.position + 1) * h.lcdWrapper.Width)
    h.lcdWrapper.DisplayString("0", (int(h.position) + int(h.height)) * h.lcdWrapper.Width)
    h.lcdWrapper.DisplayString(strconv.Itoa(max), int(h.position + 1) * h.lcdWrapper.Width)
}

func (h *histogram) drawScale() {
    // draw scale under the histogram
    starPos := (int(h.height + 1 + h.position) * h.lcdWrapper.Width) + h.view.leftMargin
    maxPos := (int(h.height + 1 + h.position) * h.lcdWrapper.Width) + h.lcdWrapper.Width
    voids := 2
    labelWidth := 12
    for i := 0;; i++ {
        pos := starPos + (i * (h.view.barWidth + h.view.sepSize) * voids)
        if pos > maxPos - labelWidth {
            break
        } 
        h.lcdWrapper.DisplayString(h.getBeginHour(i * voids), pos)
    }
    h.lcdWrapper.Dev.Display()
}

func (h *histogram) getBeginHour(offsetHours int) string {
    // returns the hout of the first bar on the histogram with the given shift
    hoursNr := time.Duration(len(h.dataCollector.stakes))
    beginTime := time.Now().Add((time.Hour * hoursNr) * - 1)
    hour := beginTime.Add(time.Hour * time.Duration(offsetHours + 1)).Hour()
    return fmt.Sprintf("%02d", hour)
}

func (h *histogram) max(arr []int) int {
    var n, max int
    for v := range arr {
    if v > n {
      n = v
      max = n
    }
  }
  return max
}