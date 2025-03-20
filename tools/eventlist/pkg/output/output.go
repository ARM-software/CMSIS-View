/*
 * Copyright (c) 2022-2025 Arm Limited. All rights reserved.
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
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"eventlist/pkg/eval"
	"eventlist/pkg/event"
	"eventlist/pkg/xml/scvd"
	"fmt"
	"math"
	"os"
)

var errNoEvents = errors.New("cannot open event file")

var TimeFactor *float64
var FormatType = "txt"
var Level = ""

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

type EventRecord struct {
	Index         int     `json:"index" xml:"index"`
	Time          float64 `json:"time" xml:"time"`
	Component     string  `json:"component" xml:"component"`
	EventProperty string  `json:"eventProperty" xml:"eventProperty"`
	Value         string  `json:"value" xml:"value"`
}

type EventRecordStatistic struct {
	Event       string  `json:"event" xml:"event"`
	Count       int     `json:"count" xml:"count"`
	AddCount    string  `json:"addCount" xml:"addCount"`
	Start       string  `json:"start" xml:"start"`
	MinStopTime float64 `json:"minStopTime" xml:"minStopTime"`
	MaxStopTime float64 `json:"maxStopTime" xml:"maxStopTime"`
	Total       string  `json:"total" xml:"total"`
	Min         string  `json:"min" xml:"min"`
	Max         string  `json:"max" xml:"max"`
	First       string  `json:"first" xml:"first"`
	Last        string  `json:"last" xml:"last"`
	Avg         string  `json:"avg" xml:"avg"`
	MinTime     float64 `json:"minTime" xml:"minTime"`
	MaxTime     float64 `json:"maxTime" xml:"maxTime"`
	FirstTime   string  `json:"firstTime" xml:"firstTime"`
	LastTime    string  `json:"lastTime" xml:"lastTime"`
	TextB       string  `json:"textB" xml:"textB"`
	TextMinB    string  `json:"textMinB" xml:"textMinB"`
	TextMinE    string  `json:"textMinE" xml:"textMinE"`
	TextMaxB    string  `json:"textMaxB" xml:"textMaxB"`
	TextMaxE    string  `json:"textMaxE" xml:"textMaxE"`
}

type EventsTable struct {
	Events     []EventRecord          `json:"events" xml:"events"`
	Statistics []EventRecordStatistic `json:"statistics" xml:"statistics"`
}

// init initializes the eventStatistic struct by setting default values for its fields.
// It sets boolean fields evFirst and evStart to false, count to 0, min to the maximum float64 value,
// max to 0, and other numeric fields (first, last, avg, minTime, maxTime, firstTime, lastTime) to 0.
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

// add records the start and stop events with their respective times and texts,
// and updates the event statistics accordingly.
//
// Parameters:
//   - time: The time at which the event occurred.
//   - start: A boolean indicating whether the event is a start event (true) or a stop event (false).
//   - text: A string containing additional information about the event.
//
// The function updates the following statistics:
//   - Minimum event duration and its start time and texts.
//   - Maximum event duration and its start time and texts.
//   - First event duration and its start time.
//   - Last event duration and its start time.
//   - Total duration of all events.
//   - Average duration of all events.
//   - Count of events.
//
// If a start event is recorded while another start event is still active, it is ignored.
// Similarly, if a stop event is recorded without a corresponding start event, it is ignored.
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

// init initializes all the values in the eventProperty's values slice by calling
// their respective init methods. It iterates over each element in the slice and
// invokes the init method on each one.
func (ep *eventProperty) init() {
	for i := uint16(0); i < uint16(len(ep.values)); i++ {
		ep.values[i].init()
	}
}

// add adds an event to the eventProperty. If the idx is 15 and start is false,
// it stops all events by iterating through all values and calling their add method.
// Otherwise, it adds the event to the specific index in the values slice.
//
// Parameters:
//   - time: The time at which the event occurs.
//   - idx: The index of the event in the values slice.
//   - start: A boolean indicating whether the event is a start event.
//   - text: A string containing additional information about the event.
func (ep *eventProperty) add(time float64, idx uint16, start bool, text string) {
	if idx == 15 && !start { // stop 15 means stop all
		for i := uint16(0); i < uint16(len(ep.values)); i++ {
			ep.values[i].add(time, start, text)
		}
	} else {
		ep.values[idx].add(time, start, text)
	}
}

// getCount returns the count of events at the specified index.
// If the index is out of range, it returns 0.
//
// Parameters:
//
//	idx - the index of the event property to retrieve the count for.
//
// Returns:
//
//	The count of events at the specified index, or 0 if the index is out of range.
func (ep *eventProperty) getCount(idx uint16) int {
	if int(idx) >= len(ep.values) {
		return 0
	}
	return ep.values[idx].count
}

// getAddCount returns a string indicating whether an event should be added.
// If the event at the given index has started, it returns "+1". Otherwise,
//
//	it returns "  " as place holder.
//
// Parameters:
//
//	idx (uint16): The index of the event in the values slice.
//
// Returns:
//
//	string: "+1" if the event at the given index has started, otherwise "  ".
func (ep *eventProperty) getAddCount(idx uint16) string {
	if int(idx) < len(ep.values) && ep.values[idx].evStart {
		return "+1"
	}
	return "  "
}

// convertUnit converts a given float64 value `v` to a string representation
// with an appropriate unit prefix (G, M, k, m, µ, n) based on its magnitude.
// The `unit` parameter is the base unit to which the prefix will be added.
// The function returns a formatted string with the value adjusted to the
// appropriate magnitude and the corresponding unit prefix.
//
// Parameters:
//   - v: The float64 value to be converted.
//   - unit: The base unit as a string.
//
// Returns:
//
//	A string representing the value `v` with the appropriate unit prefix.
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

// getTot retrieves the total time (tot) value from the eventProperty at the specified index,
// converts it to a string representation in seconds, and returns it.
//
// Parameters:
//
//	idx (uint16): The index of the eventProperty value to retrieve.
//
// Returns:
//
//	string: The total time value converted to a string in seconds.
func (ep *eventProperty) getTot(idx uint16) string {
	return convertUnit(ep.values[idx].tot, "s")
}

// getMin returns the minimum value from the eventProperty values at the specified index,
// converted to a string representation with the unit "s" (seconds).
//
// Parameters:
//
//	idx (uint16): The index of the value to retrieve.
//
// Returns:
//
//	string: The minimum value at the specified index, converted to a string with the unit "s".
func (ep *eventProperty) getMin(idx uint16) string {
	return convertUnit(ep.values[idx].min, "s")
}

// getMax retrieves the maximum value from the eventProperty's values at the specified index,
// converts it to a string with the appropriate unit, and returns it.
//
// Parameters:
//
//	idx - The index of the value to retrieve.
//
// Returns:
//
//	A string representing the maximum value at the specified index, converted to the appropriate unit.
func (ep *eventProperty) getMax(idx uint16) string {
	return convertUnit(ep.values[idx].max, "s")
}

// getAvg calculates the average value for the event property at the specified index.
// It returns the average value as a string with the appropriate unit.
//
// Parameters:
//
//	idx (uint16): The index of the event property value to calculate the average for.
//
// Returns:
//
//	string: The average value as a string with the unit "s". If the count is zero, it returns "0s".
func (ep *eventProperty) getAvg(idx uint16) string {
	if ep.values[idx].count != 0 {
		return convertUnit(ep.values[idx].avg/float64(ep.values[idx].count), "s")
	}
	return convertUnit(0, "s")
}

// getFirst retrieves the first value from the eventProperty at the specified index,
// converts it to a string representation in seconds, and returns it.
//
// Parameters:
//
//	idx (uint16): The index of the value to retrieve.
//
// Returns:
//
//	string: The string representation of the first value in seconds.
func (ep *eventProperty) getFirst(idx uint16) string {
	return convertUnit(ep.values[idx].first, "s")
}

// getLast retrieves the last value from the eventProperty at the specified index,
// converts it to a string representation in seconds, and returns it.
//
// Parameters:
//
//	idx - The index of the value to retrieve.
//
// Returns:
//
//	A string representation of the last value at the specified index, converted to seconds.
func (ep *eventProperty) getLast(idx uint16) string {
	return convertUnit(ep.values[idx].last, "s")
}

type Output struct {
	evProps       [4]eventProperty
	columns       []string
	componentSize int
	propertySize  int
}

// TODO: escape-sequeces for Color and Bold

// buildStatistic processes events from the provided bufio.Reader and updates
// the Output structure with statistics about the events. It reads events,
// evaluates their properties, and updates the component and property sizes
// based on the event definitions. It also handles specific event classes and
// updates the event properties accordingly.
//
// Parameters:
//   - in: a bufio.Reader from which events are read.
//   - evdefs: a map of event definitions (scvd.Events).
//   - typedefs: a map of type definitions (eval.Typedefs).
//
// Returns:
//
//	The total number of events processed.
func (o *Output) buildStatistic(in *bufio.Reader, evdefs scvd.Events, typedefs eval.Typedefs) int {
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
		var evdef scvd.EventType
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

// conditionalWrite writes formatted data to the provided bufio.Writer
// if the global FormatType is set to "txt". It uses fmt.Fprintf to
// format the data according to the specified format string and arguments.
//
// Parameters:
//   - out: A pointer to a bufio.Writer where the formatted data will be written.
//   - format: A format string as described in the fmt package documentation.
//   - a: Variadic arguments to be formatted according to the format string.
//
// Returns:
//   - err: An error value if writing to the bufio.Writer fails, otherwise nil.
func conditionalWrite(out *bufio.Writer, format string, a ...any) (err error) {
	if FormatType == "txt" {
		_, err = fmt.Fprintf(out, format, a...)
		return err
	}
	return nil
}

// printStatistic writes the event statistics to the provided bufio.Writer.
// It includes details such as event count, total, min, max, average, first, and last occurrences.
// The statistics are formatted and written in a tabular form.
//
// Parameters:
//
//	out - A pointer to a bufio.Writer where the statistics will be written.
//	eventCount - The total number of events to be processed.
//	eventTable - A pointer to an EventsTable where the statistics will be stored.
//
// Returns:
//
//	An error if any write operation fails, otherwise nil.
func (o *Output) printStatistic(out *bufio.Writer, eventCount int, eventTable *EventsTable) error {
	var err error

	if out != nil && eventCount > 0 {
		if err = conditionalWrite(out, "   Start/Stop event statistic\n"); err != nil {
			return err
		}
		if err = conditionalWrite(out, "   --------------------------\n\n"); err != nil {
			return err
		}
		if err = conditionalWrite(out, "Event count      total       min         max         average     first       last\n"); err != nil {
			return err
		}
		if err = conditionalWrite(out, "----- -----      -----       ---         ---         -------     -----       ----\n"); err != nil {
			return err
		}
		for i := uint16(0); i < uint16(len(o.evProps)); i++ {
			for j := uint16(0); j < uint16(len(o.evProps[i].values)); j++ {
				if o.evProps[i].values[j].evFirst {
					eventStat := EventRecordStatistic{
						Event:       fmt.Sprintf("%c(%d)", byte(i+'A'), j),
						AddCount:    o.evProps[i].getAddCount(j),
						Count:       o.evProps[i].getCount(j),
						Total:       o.evProps[i].getTot(j),
						Min:         o.evProps[i].getMin(j),
						Max:         o.evProps[i].getMax(j),
						Avg:         o.evProps[i].getAvg(j),
						First:       o.evProps[i].getFirst(j),
						Last:        o.evProps[i].getLast(j),
						MinTime:     o.evProps[i].values[j].minTime,
						TextMinB:    o.evProps[i].values[j].textMinB,
						TextMinE:    o.evProps[i].values[j].textMinE,
						MinStopTime: o.evProps[i].values[j].minTime + o.evProps[i].values[j].min,
						MaxStopTime: o.evProps[i].values[j].maxTime + o.evProps[i].values[j].max,
						MaxTime:     o.evProps[i].values[j].maxTime,
						TextMaxB:    o.evProps[i].values[j].textMaxB,
						TextMaxE:    o.evProps[i].values[j].textMaxE,
					}
					err = conditionalWrite(out, eventStat.Event)
					if err == nil && j < 10 {
						err = conditionalWrite(out, " ")
					}
					if err != nil {
						return err
					}
					err = conditionalWrite(out, " %5d%s %s %s %s %s %s %s\n",
						eventStat.Count,
						eventStat.AddCount,
						eventStat.Total,
						eventStat.Min,
						eventStat.Max,
						eventStat.Avg,
						eventStat.First,
						eventStat.Last)
					if err != nil {
						return err
					}
					err = conditionalWrite(out, "      Min: Start: %.8f %s Stop: %.8f %s\n",
						eventStat.MinTime,
						eventStat.TextMinB,
						eventStat.MinStopTime,
						eventStat.TextMinE)
					if err != nil {
						return err
					}
					err = conditionalWrite(out, "      Max: Start: %.8f %s Stop: %.8f %s\n\n",
						eventStat.MaxTime,
						eventStat.TextMaxB,
						eventStat.MaxStopTime,
						eventStat.TextMaxE)
					if err != nil {
						return err
					}
					eventTable.Statistics = append(eventTable.Statistics, eventStat)
				}
			}
		}
	}
	return err
}

// escapeGen takes a string as input and returns a new string with certain characters
// escaped using backslashes. The function handles the following characters:
// - Single quote ('): escaped as \'
// - Double quote ("): escaped as \"
// - Alert/bell (\a): escaped as \a
// - Backspace (\b): escaped as \b
// - Escape (\x1b): escaped as \e
// - Form feed (\f): escaped as \f
// - Newline (\n): escaped as \n
// - Carriage return (\r): escaped as \r
// - Horizontal tab (\t): escaped as \t
// - Vertical tab (\v): escaped as \v
// For any other control characters (ASCII values less than 32), the function
// escapes them using their octal representation (e.g., \000).
// All other characters are left unchanged.
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

// printEvents reads events from the input buffer, processes them, and writes the formatted events to the output buffer.
// It also updates the event table with the processed events.
//
// Parameters:
//   - out: A buffered writer to which the formatted events will be written.
//   - in: A buffered reader from which the events will be read.
//   - evdefs: A map of event definitions used to interpret the events.
//   - typedefs: A map of type definitions used for evaluating event data.
//   - eventTable: A table to store the processed events.
//
// Returns:
//   - error: An error if any occurs during reading, processing, or writing events.
//
// The function processes special events such as EventRecorderInitialize (ID 0xFF00) and Clock event (ID 0xFF03)
// to adjust the time factor. It filters events based on the specified level and formats them accordingly.
// The function handles special cases like stdout events (ID 0xFE00) differently.
func (o *Output) printEvents(out *bufio.Writer, in *bufio.Reader, evdefs scvd.Events, typedefs eval.Typedefs, eventTable *EventsTable) error {
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
		eventRecord := EventRecord{
			Index: no,
			Time:  beforeClockEvent + TimeInSecs(ev.Time-lastClockEvent),
		}
		var rep string
		if evdef, ok := evdefs[ev.Info.ID]; ok {
			// Filter events by level
			if Level == "" || evdef.Level == Level {
				eventRecord.Component = evdef.Brief
				eventRecord.EventProperty = evdef.Property
				if ev.Info.ID == 0xFE00 && ev.Data != nil { // special case stdout
					s := escapeGen(string(*ev.Data))
					eventRecord.Value = s
					err = conditionalWrite(out, "%5d %.8f %*s %*s \"%s\"\n",
						eventRecord.Index, eventRecord.Time, -o.componentSize,
						eventRecord.Component, -o.propertySize, eventRecord.EventProperty, eventRecord.Value)
				} else {
					rep, err = ev.EvalLine(evdef, typedefs)
					if err == nil {
						eventRecord.Value = rep
						err = conditionalWrite(out, "%5d %.8f %*s %*s %s\n",
							eventRecord.Index, eventRecord.Time, -o.componentSize,
							eventRecord.Component, -o.propertySize, eventRecord.EventProperty, eventRecord.Value)
					}
				}
			}
		} else {
			eventRecord.Component = fmt.Sprintf("0x%02X%*s", uint8(ev.Info.ID>>8), 0, "")
			eventRecord.EventProperty = fmt.Sprintf("0x%04X%*s", ev.Info.ID, 0, "")
			if ev.Info.ID == 0xFE00 && ev.Data != nil { // special case stdout
				s := escapeGen(string(*ev.Data))
				eventRecord.Value = s
				err = conditionalWrite(out, "%5d %.8f 0x%02X%*s 0x%04X%*s \"%s\"\n",
					eventRecord.Index, eventRecord.Time,
					uint8(ev.Info.ID>>8), -(o.componentSize - 4), "",
					ev.Info.ID, -(o.propertySize - 6), "", eventRecord.Value)
			} else {
				rep = ev.GetValuesAsString()
				eventRecord.Value = rep
				err = conditionalWrite(out, "%5d %.8f 0x%02X%*s 0x%04X%*s %s\n",
					eventRecord.Index, eventRecord.Time,
					uint8(ev.Info.ID>>8), -(o.componentSize - 4), "",
					ev.Info.ID, -(o.propertySize - 6), "", eventRecord.Value)
			}
		}
		eventTable.Events = append(eventTable.Events, eventRecord)
		if err != nil {
			break
		}
		no++
	}
	return err
}

// printHeader writes the header section of the detailed event list to the provided bufio.Writer.
// It includes the title, a separator line, and column headers formatted according to the Output struct's settings.
//
// Parameters:
//
//	out (*bufio.Writer): The buffered writer to which the header will be written.
//
// Returns:
//
//	error: An error if any write operation fails, otherwise nil.
func (o *Output) printHeader(out *bufio.Writer) error {
	var err error
	if err = conditionalWrite(out, "   Detailed event list\n"); err != nil {
		return err
	}
	if err = conditionalWrite(out, "   -------------------\n\n"); err != nil {
		return err
	}
	err = conditionalWrite(out, "%5s %-10s %*s %*s %s\n", o.columns[0], o.columns[1],
		-o.componentSize, o.columns[2], -o.propertySize, o.columns[3], o.columns[4])
	if err != nil {
		return err
	}
	err = conditionalWrite(out, "----- --------   %*s %*s -----\n",
		-o.componentSize, "---------", -o.propertySize, "--------------")
	return err
}

// print generates and writes the output for the given event file and definitions.
// It processes the events, builds statistics, and prints the event details and statistics
// based on the provided flags.
//
// Parameters:
//   - out: A buffered writer to write the output.
//   - eventFile: A pointer to the event file path string.
//   - evdefs: Event definitions.
//   - typedefs: Type definitions.
//   - statBegin: A flag indicating whether to print statistics at the beginning.
//   - showStatistic: A flag indicating whether to show statistics.
//   - eventsTable: A pointer to the events table.
//
// Returns:
//   - An error if any operation fails, otherwise nil.
func (o *Output) print(out *bufio.Writer, eventFile *string, evdefs scvd.Events,
	typedefs eval.Typedefs, statBegin bool, showStatistic bool, eventsTable *EventsTable) error {
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
		err = o.printStatistic(out, eventCount, eventsTable)
		if err == nil && !showStatistic {
			err = conditionalWrite(out, "\n")
		}
	}

	if err == nil && !showStatistic {
		err = o.printHeader(out)
		if err == nil {
			in = b.Open(eventFile)
			if in != nil {
				err = o.printEvents(out, in, evdefs, typedefs, eventsTable)
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
			err = conditionalWrite(out, "\n")
		}
		if err == nil {
			err = o.printStatistic(out, eventCount, eventsTable)
		}
	}
	if err == nil {
		err = out.Flush()
	}
	return err
}

// Print generates and writes event data to a specified file or standard output in a given format.
// It supports XML and JSON formats and can include statistics if specified.
//
// Parameters:
//   - filename: Pointer to the name of the file where the output will be written. If nil or empty, output is written to stdout.
//   - formatType: Pointer to the format type ("xml" or "json"). If nil or invalid, default format is used.
//   - level: Pointer to the level of detail for the output. If nil or empty, default level is used.
//   - eventFile: Pointer to the event file name.
//   - evdefs: Event definitions.
//   - typedefs: Type definitions.
//   - statBegin: Boolean flag indicating whether to include statistics at the beginning.
//   - showStatistic: Boolean flag indicating whether to show statistics.
//
// Returns:
//   - error: An error if the file could not be created or written to, or if there was an error during the output generation.
func Print(filename *string, formatType *string, level *string, eventFile *string, evdefs scvd.Events,
	typedefs eval.Typedefs, statBegin bool, showStatistic bool) error {
	var file *os.File
	var err error
	var o Output

	eventsTable := EventsTable{
		Events:     []EventRecord{},
		Statistics: []EventRecordStatistic{},
	}

	if TimeFactor == nil {
		TimeFactor = new(float64)
	}
	if *TimeFactor == 0.0 {
		*TimeFactor = 4e-8
	}
	if formatType != nil {
		if *formatType == "xml" || *formatType == "json" {
			FormatType = *formatType
		}
	}
	if level != nil && *level != "" {
		Level = *level
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
	err = o.print(out, eventFile, evdefs, typedefs, statBegin, showStatistic, &eventsTable)
	if err == nil {
		if FormatType == "json" {
			output, err := json.Marshal(eventsTable)
			if err == nil {
				buf := bytes.NewBuffer(output)
				_, err = fmt.Fprint(out, buf)
				if err == nil {
					out.Flush()
				}
			}
		} else if FormatType == "xml" {
			output, err := xml.Marshal(eventsTable)
			if err == nil {
				buf := bytes.NewBuffer(output)
				_, err = fmt.Fprint(out, buf)
				if err == nil {
					out.Flush()
				}
			}
		} else {
			err = out.Flush()
		}
	} else {
		_ = out.Flush()
	}
	return err
}
