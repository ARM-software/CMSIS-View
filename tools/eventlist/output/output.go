/*
 * Copyright (c) 2022 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the License); you may
 * not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an AS IS BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package output

import (
	"bufio"
	"errors"
	"eventlist/eval"
	"eventlist/event"
	"eventlist/xml/scvd"
	"fmt"
	"math"
	"os"
)

var errNoEvents = errors.New("cannot open event file")

var TimeFactor *float64

func TimeInSecs(time uint64) float64 {
	if TimeFactor == nil {
		return 4e-8 * float64(time) // default
	}
	return *TimeFactor * float64(time)
}

type eventStatistic struct {
	evFirst   bool // true if not first time appeared
	evStart   bool // true if started, false if stopped
	count     int
	start     float64
	tot       float64
	min       float64
	max       float64
	first     float64
	last      float64
	avg       float64
	minTime   float64
	maxTime   float64
	firstTime float64
	lastTime  float64
	textB     string
	textMinB  string
	textMinE  string
	textMaxB  string
	textMaxE  string
}

func (es *eventStatistic) init() {
	es.evFirst = false
	es.evStart = false
	es.count = 0
	es.min = math.MaxFloat64
	es.max = 0
	es.first = 0
	es.last = 0
	es.avg = 0
	es.minTime = 0
	es.maxTime = 0
	es.firstTime = 0
	es.lastTime = 0
}

func (es *eventStatistic) add(time float64, start bool, text string) {
	if start {
		if es.evStart {
			return // ignore start event, was not stopped yet
		}
		es.evStart = true
		es.start = time
		es.textB = text
	} else {
		if !es.evStart {
			return // ignore already stopped events
		}
		es.evStart = false
		diff := time - es.start
		if diff < es.min {
			es.min = diff
			es.minTime = es.start
			es.textMinB = es.textB
			es.textMinE = text
		}
		if diff > es.max {
			es.max = diff
			es.maxTime = es.start
			es.textMaxB = es.textB
			es.textMaxE = text
		}
		if !es.evFirst {
			es.first = diff
			es.firstTime = es.start
			es.evFirst = true
		}
		es.last = diff
		es.lastTime = es.start
		es.tot += diff
		es.avg += diff
		es.count++
	}
}

type eventProperty struct {
	values [16]eventStatistic
}

func (ep *eventProperty) init() {
	for i := uint16(0); i < uint16(len(ep.values)); i++ {
		ep.values[i].init()
	}
}

func (ep *eventProperty) add(time float64, idx uint16, start bool, text string) {
	if idx == 15 && !start { // stop 15 means stop all
		for i := uint16(0); i < uint16(len(ep.values)); i++ {
			ep.values[i].add(time, start, text)
		}
	} else {
		ep.values[idx].add(time, start, text)
	}
}

func (ep *eventProperty) getCount(idx uint16) int {
	if int(idx) >= len(ep.values) {
		return 0
	}
	return ep.values[idx].count
}

func (ep *eventProperty) getAddCount(idx uint16) string {
	if int(idx) < len(ep.values) && ep.values[idx].evStart {
		return "+1"
	}
	return "  "
}

func convertUnit(v float64, unit string) string { //nolint:golint,unparam
	switch {
	case v >= 1e9:
		unit = "G" + unit
		v /= 1e9
	case v >= 1e6:
		unit = "M" + unit
		v /= 1e6
	case v >= 1e3:
		unit = "k" + unit
		v /= 1e3
	case v >= 1 || v == 0.0:
		unit = unit + " "
	case v >= 1e-3:
		unit = "m" + unit
		v *= 1e3
	case v >= 1e-6:
		unit = "µ" + unit
		v *= 1e6
	case v >= 1e-9:
		unit = "n" + unit
		v *= 1e9
	}
	return fmt.Sprintf("%9.5f%s", v, unit)
}

func (ep *eventProperty) getTot(idx uint16) string {
	return convertUnit(ep.values[idx].tot, "s")
}

func (ep *eventProperty) getMin(idx uint16) string {
	return convertUnit(ep.values[idx].min, "s")
}

func (ep *eventProperty) getMax(idx uint16) string {
	return convertUnit(ep.values[idx].max, "s")
}

func (ep *eventProperty) getAvg(idx uint16) string {
	if ep.values[idx].count != 0 {
		return convertUnit(ep.values[idx].avg/float64(ep.values[idx].count), "s")
	}
	return convertUnit(0, "s")
}

func (ep *eventProperty) getFirst(idx uint16) string {
	return convertUnit(ep.values[idx].first, "s")
}

func (ep *eventProperty) getLast(idx uint16) string {
	return convertUnit(ep.values[idx].last, "s")
}

type Output struct {
	evProps       [4]eventProperty
	columns       []string
	componentSize int
	propertySize  int
}

func (o *Output) buildStatistic(in *bufio.Reader, evdefs map[uint16]scvd.Event,
	typedefs map[string]map[string]scvd.TdMember) int {
	o.componentSize = len(o.columns[2]) // use minimum width of header
	o.propertySize = len(o.columns[3])
	for i := uint16(0); i < uint16(len(o.evProps)); i++ {
		o.evProps[i].init()
	}
	var beforeClockEvent float64
	var lastClockEvent uint64
	var eventCount int
	for {
		var ev event.Data
		if err := ev.Read(in); err != nil {
			if errors.Is(err, eval.ErrEof) {
				break
			}
			fmt.Println(err)
			return 0
		}
		eventCount++
		var evdef scvd.Event
		var ok bool
		var rep string
		if evdef, ok = evdefs[ev.Info.ID]; ok {
			if len(evdef.Brief) > o.componentSize {
				o.componentSize = len(evdef.Brief)
			}
			if len(evdef.Property) > o.propertySize {
				o.propertySize = len(evdef.Property)
			}
			class, _, _, _ := ev.Info.SplitID()
			switch class {
			case 0xEF:
				rep, _ = ev.EvalLine(evdef, typedefs)
			}
		}
		class, group, idx, start := ev.Info.SplitID()
		switch class {
		case 0xEF:
			if !ok { // rep not yet built up because of wrong or missing SCVD files
				rep = ev.GetValuesAsString()
			}
			o.evProps[group].add(beforeClockEvent+TimeInSecs(ev.Time-lastClockEvent), idx, start, rep)
		case 0xFF:
			switch ev.Info.ID {
			case 0xFF00: // EventRecorderInitialize
				if ev.Value2 != 0 {
					beforeClockEvent = TimeInSecs(ev.Time)
					lastClockEvent = ev.Time
					if TimeFactor == nil {
						TimeFactor = new(float64)
					}
					*TimeFactor = 1.0 / float64(ev.Value2)
				}
			case 0xFF03: // EventRecorderClock
				if ev.Value1 != 0 {
					beforeClockEvent = TimeInSecs(ev.Time - lastClockEvent)
					lastClockEvent = ev.Time
					if TimeFactor == nil {
						TimeFactor = new(float64)
					}
					*TimeFactor = 1.0 / float64(ev.Value1)
				}
			}
		}
	}
	return eventCount
}

func (o *Output) printStatistic(out *bufio.Writer, eventCount int) error {
	var err error

	if out != nil && eventCount > 0 {
		if _, err = out.WriteString("   Start/Stop event statistic\n"); err != nil {
			return err
		}
		if _, err = out.WriteString("   --------------------------\n\n"); err != nil {
			return err
		}
		if _, err = out.WriteString("Event count      total       min         max         average     first       last\n"); err != nil {
			return err
		}
		if _, err = out.WriteString("----- -----      -----       ---         ---         -------     -----       ----\n"); err != nil {
			return err
		}
		for i := uint16(0); i < uint16(len(o.evProps)); i++ {
			for j := uint16(0); j < uint16(len(o.evProps[i].values)); j++ {
				if o.evProps[i].values[j].evFirst {
					_, err = fmt.Fprintf(out, "%c(%d)", byte(i+'A'), j)
					if err == nil && j < 10 {
						err = out.WriteByte(' ')
					}
					if err != nil {
						return err
					}
					_, err = fmt.Fprintf(out, " %5d%s %s %s %s %s %s %s\n",
						o.evProps[i].getCount(j),
						o.evProps[i].getAddCount(j),
						o.evProps[i].getTot(j),
						o.evProps[i].getMin(j),
						o.evProps[i].getMax(j),
						o.evProps[i].getAvg(j),
						o.evProps[i].getFirst(j),
						o.evProps[i].getLast(j))
					if err != nil {
						return err
					}
					_, err = fmt.Fprintf(out, "      Min: Start: %.8f %s Stop: %.8f %s\n",
						o.evProps[i].values[j].minTime,
						o.evProps[i].values[j].textMinB,
						o.evProps[i].values[j].minTime+o.evProps[i].values[j].min,
						o.evProps[i].values[j].textMinE)
					if err != nil {
						return err
					}
					_, err = fmt.Fprintf(out, "      Max: Start: %.8f %s Stop: %.8f %s\n\n",
						o.evProps[i].values[j].maxTime,
						o.evProps[i].values[j].textMaxB,
						o.evProps[i].values[j].maxTime+o.evProps[i].values[j].max,
						o.evProps[i].values[j].textMaxE)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return err
}

func escapeGen(s string) string {
	var t string
	for _, c := range s {
		switch c {
		case '\'':
			t += "\\'"
		case '"':
			t += "\\\""
		case '\a':
			t += "\\a"
		case '\b':
			t += "\\b"
		case '\x1b':
			t += "\\e"
		case '\f':
			t += "\\f"
		case '\n':
			t += "\\n"
		case '\r':
			t += "\\r"
		case '\t':
			t += "\\t"
		case '\v':
			t += "\\v"
		default:
			if c < ' ' {
				t += fmt.Sprintf("\\%03o", byte(c))
			} else {
				t += string(c)
			}
		}
	}
	return t
}

func (o *Output) printEvents(out *bufio.Writer, in *bufio.Reader, evdefs map[uint16]scvd.Event,
	typedefs map[string]map[string]scvd.TdMember) error {
	if out == nil || in == nil {
		return nil
	}
	var err error
	no := 0
	var beforeClockEvent float64
	var lastClockEvent uint64
	for {
		var ev event.Data
		if err = ev.Read(in); err != nil {
			if errors.Is(err, eval.ErrEof) {
				err = nil
				break // end of event data reached
			}
			fmt.Println(err)
		}
		if err != nil {
			break
		}
		if ev.Info.ID == 0xFF00 { // EventRecorderInitialize
			if ev.Value2 != 0 {
				beforeClockEvent = TimeInSecs(ev.Time)
				lastClockEvent = ev.Time
				if TimeFactor == nil {
					TimeFactor = new(float64)
				}
				*TimeFactor = 1.0 / float64(ev.Value2)
			}
		}
		if ev.Info.ID == 0xFF03 { // Clock event
			if ev.Value1 != 0 {
				beforeClockEvent = TimeInSecs(ev.Time - lastClockEvent)
				lastClockEvent = ev.Time
				if TimeFactor == nil {
					TimeFactor = new(float64)
				}
				*TimeFactor = 1.0 / float64(ev.Value1)
			}
		}
		var rep string
		if evdef, ok := evdefs[ev.Info.ID]; ok {
			if ev.Info.ID == 0xFE00 && ev.Data != nil { // special case stdout
				s := escapeGen(string(*ev.Data))
				_, err = fmt.Fprintf(out, "%5d %.8f %*s %*s \"%s\"\n",
					no, beforeClockEvent+TimeInSecs(ev.Time-lastClockEvent),
					-o.componentSize, evdef.Brief, -o.propertySize, evdef.Property, s)
			} else {
				rep, err = ev.EvalLine(evdef, typedefs)
				if err == nil {
					_, err = fmt.Fprintf(out, "%5d %.8f %*s %*s %s\n",
						no, beforeClockEvent+TimeInSecs(ev.Time-lastClockEvent),
						-o.componentSize, evdef.Brief, -o.propertySize, evdef.Property, rep)
				}
			}
		} else {
			if ev.Info.ID == 0xFE00 && ev.Data != nil { // special case stdout
				s := escapeGen(string(*ev.Data))
				_, err = fmt.Fprintf(out, "%5d %.8f 0x%02X%*s 0x%04X%*s \"%s\"\n",
					no, beforeClockEvent+TimeInSecs(ev.Time-lastClockEvent),
					uint8(ev.Info.ID>>8), -(o.componentSize - 4), "",
					ev.Info.ID, -(o.propertySize - 6), "", s)
			} else {
				rep = ev.GetValuesAsString()
				_, err = fmt.Fprintf(out, "%5d %.8f 0x%02X%*s 0x%04X%*s %s\n",
					no, beforeClockEvent+TimeInSecs(ev.Time-lastClockEvent),
					uint8(ev.Info.ID>>8), -(o.componentSize - 4), "",
					ev.Info.ID, -(o.propertySize - 6), "", rep)
			}
		}
		if err != nil {
			break
		}
		no++
	}
	return err
}

func (o *Output) printHeader(out *bufio.Writer) error {
	var err error
	if _, err = out.WriteString("   Detailed event list\n"); err != nil {
		return err
	}
	if _, err = out.WriteString("   -------------------\n\n"); err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "%5s %-10s %*s %*s %s\n", o.columns[0], o.columns[1],
		-o.componentSize, o.columns[2], -o.propertySize, o.columns[3], o.columns[4])
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "----- --------   %*s %*s -----\n",
		-o.componentSize, "---------", -o.propertySize, "--------------")
	return err
}

func (o *Output) print(out *bufio.Writer, eventFile *string, evdefs map[uint16]scvd.Event,
	typedefs map[string]map[string]scvd.TdMember, statBegin bool, showStatistic bool) error {
	var b event.Binary
	var err error
	var eventCount int

	o.columns = []string{"Index", "Time (s)", "Component", "Event Property", "Value"}

	if eventFile == nil {
		return errNoEvents
	}
	in := b.Open(eventFile)
	if in != nil {
		eventCount = o.buildStatistic(in, evdefs, typedefs)
		err = b.Close()
	} else {
		err = errNoEvents
	}

	if err == nil && statBegin {
		err = o.printStatistic(out, eventCount)
		if err == nil && !showStatistic {
			_, err = out.WriteString("\n")
		}
	}

	if err == nil && !showStatistic {
		err = o.printHeader(out)
		if err == nil {
			in = b.Open(eventFile)
			if in != nil {
				err = o.printEvents(out, in, evdefs, typedefs)
				if err != nil {
					_ = b.Close()
				} else {
					err = b.Close()
				}
			} else {
				err = errNoEvents // cannot happen because eventFile already was read
			}
		}
	}

	if err == nil && !statBegin {
		if !showStatistic {
			_, err = out.WriteString("\n")
		}
		if err == nil {
			err = o.printStatistic(out, eventCount)
		}
	}

	if err == nil {
		err = out.Flush()
	}
	return err
}

func Print(filename *string, eventFile *string, evdefs map[uint16]scvd.Event,
	typedefs map[string]map[string]scvd.TdMember, statBegin bool, showStatistic bool) error {
	var file *os.File
	var err error
	var o Output

	if TimeFactor == nil {
		TimeFactor = new(float64)
	}
	if *TimeFactor == 0.0 {
		*TimeFactor = 4e-8
	}

	if filename != nil && len(*filename) != 0 {
		if file, err = os.Create(*filename); err != nil {
			return err
		}
		defer file.Close()
	} else {
		file = os.Stdout
	}

	out := bufio.NewWriter(file)
	err = o.print(out, eventFile, evdefs, typedefs, statBegin, showStatistic)
	if err == nil {
		err = out.Flush()
	} else {
		_ = out.Flush()
	}
	return err
}
