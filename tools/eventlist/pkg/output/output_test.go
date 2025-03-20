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
	"errors"
	"eventlist/pkg/eval"
	"eventlist/pkg/event"
	"eventlist/pkg/xml/scvd"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestTimeInSecs(t *testing.T) { //nolint:golint,paralleltest
	type args struct {
		time uint64
	}
	tests := []struct {
		name  string
		args  args
		clear bool
		want  float64
	}{
		{"clear", args{77}, true, 4e-8 * 77},
		{"set", args{47}, false, 2.0 * 47},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			if tt.clear {
				TimeFactor = nil
			} else {
				TimeFactor = new(float64)
				*TimeFactor = 2
			}
			if got := TimeInSecs(tt.args.time); got != tt.want {
				t.Errorf("TimeInSecs() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_eventStatistic_init(t *testing.T) {
	t.Parallel()

	type fields struct {
		evFirst  bool
		evStart  bool
		count    int
		start    float64
		tot      float64
		min      float64
		max      float64
		first    float64
		last     float64
		avg      float64
		textB    string
		textMinB string
		textMinE string
		textMaxB string
		textMaxE string
	}
	tests := []struct {
		name   string
		fields fields
		want   eventStatistic
	}{
		{"init", fields{start: 1.234, evFirst: true, evStart: true, count: 123, min: 44, max: 55},
			eventStatistic{start: 1.234, evFirst: false, evStart: false, count: 0, min: math.MaxFloat64, max: 0}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			es := &eventStatistic{
				evFirst:  tt.fields.evFirst,
				evStart:  tt.fields.evStart,
				count:    tt.fields.count,
				start:    tt.fields.start,
				tot:      tt.fields.tot,
				min:      tt.fields.min,
				max:      tt.fields.max,
				first:    tt.fields.first,
				last:     tt.fields.last,
				avg:      tt.fields.avg,
				textB:    tt.fields.textB,
				textMinB: tt.fields.textMinB,
				textMinE: tt.fields.textMinE,
				textMaxB: tt.fields.textMaxB,
				textMaxE: tt.fields.textMaxE,
			}
			es.init()
			if !reflect.DeepEqual(*es, tt.want) {
				t.Errorf("eventStatistic.init() %s = %v, want %v", tt.name, es, tt.want)
			}
		})
	}
}

func Test_eventStatistic_add(t *testing.T) {
	t.Parallel()

	type fields struct {
		evFirst  bool
		evStart  bool
		count    int
		start    float64
		tot      float64
		min      float64
		max      float64
		first    float64
		last     float64
		avg      float64
		textB    string
		textMinB string
		textMinE string
		textMaxB string
		textMaxE string
	}
	type args struct {
		time  float64
		start bool
		text  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   eventStatistic
	}{
		{"start", fields{min: math.MaxFloat64}, args{time: 123, start: true, text: "text"},
			eventStatistic{evStart: true, start: 123, textB: "text", min: math.MaxFloat64}},
		{"start_start", fields{min: math.MaxFloat64, evStart: true}, args{time: 123, start: true, text: "text"},
			eventStatistic{evStart: true, min: math.MaxFloat64}},
		{"!start_!start", fields{min: math.MaxFloat64, evStart: false}, args{time: 123, start: false, text: "text"},
			eventStatistic{evStart: false, min: math.MaxFloat64}},
		{"!start_min", fields{min: math.MaxFloat64, max: 222, evFirst: true, evStart: true, start: 111, textB: "tb"},
			args{time: 123, start: false, text: "text"},
			eventStatistic{evStart: false, start: 111, textB: "tb",
				min: 12, textMinB: "tb", textMinE: "text",
				max:   222,
				first: 0, evFirst: true, last: 12, tot: 12, avg: 12,
				minTime: 111, lastTime: 111, count: 1}},
		{"!start_max", fields{min: 1, max: 0, evFirst: true, evStart: true, start: 111, textB: "tb"},
			args{time: 123, start: false, text: "text"},
			eventStatistic{evStart: false, start: 111, textB: "tb",
				min: 1,
				max: 12, textMaxB: "tb", textMaxE: "text",
				first: 0, evFirst: true, last: 12, tot: 12, avg: 12,
				maxTime: 111, lastTime: 111, count: 1}},
		{"!start_minmax", fields{min: math.MaxFloat64, max: 0, evFirst: true, evStart: true, start: 111, textB: "tb"},
			args{time: 123, start: false, text: "text"},
			eventStatistic{evStart: false, start: 111, textB: "tb",
				min: 12, textMinB: "tb", textMinE: "text",
				max: 12, textMaxB: "tb", textMaxE: "text",
				first: 0, evFirst: true, last: 12, tot: 12, avg: 12,
				minTime: 111, maxTime: 111, lastTime: 111, count: 1}},
		{"!start_first", fields{min: 0, max: 222, evStart: true, start: 111, textB: "tb"},
			args{time: 123, start: false, text: "text"},
			eventStatistic{evStart: false, start: 111, textB: "tb",
				min:   0,
				max:   222,
				first: 12, evFirst: true, last: 12, tot: 12, avg: 12,
				firstTime: 111, lastTime: 111, count: 1}},
		{"!start_minfirst", fields{min: math.MaxFloat64, max: 222, evStart: true, start: 111, textB: "tb"},
			args{time: 123, start: false, text: "text"},
			eventStatistic{evStart: false, start: 111, textB: "tb",
				min: 12, textMinB: "tb", textMinE: "text",
				max:   222,
				first: 12, evFirst: true, last: 12, tot: 12, avg: 12,
				minTime: 111, firstTime: 111, lastTime: 111, count: 1}},
		{"!start_maxfirst", fields{min: 0, max: 0, evStart: true, start: 111, textB: "tb"},
			args{time: 123, start: false, text: "text"},
			eventStatistic{evStart: false, start: 111, textB: "tb",
				min: 0,
				max: 12, textMaxB: "tb", textMaxE: "text",
				first: 12, evFirst: true, last: 12, tot: 12, avg: 12,
				maxTime: 111, firstTime: 111, lastTime: 111, count: 1}},
		{"!start_minmaxfirst", fields{min: math.MaxFloat64, evStart: true, start: 111, textB: "tb"},
			args{time: 123, start: false, text: "text"},
			eventStatistic{evStart: false, start: 111, textB: "tb",
				min: 12, textMinB: "tb", textMinE: "text",
				max: 12, textMaxB: "tb", textMaxE: "text",
				first: 12, evFirst: true, last: 12, tot: 12, avg: 12,
				minTime: 111, maxTime: 111, firstTime: 111, lastTime: 111, count: 1}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			es := &eventStatistic{
				evFirst:  tt.fields.evFirst,
				evStart:  tt.fields.evStart,
				count:    tt.fields.count,
				start:    tt.fields.start,
				tot:      tt.fields.tot,
				min:      tt.fields.min,
				max:      tt.fields.max,
				first:    tt.fields.first,
				last:     tt.fields.last,
				avg:      tt.fields.avg,
				textB:    tt.fields.textB,
				textMinB: tt.fields.textMinB,
				textMinE: tt.fields.textMinE,
				textMaxB: tt.fields.textMaxB,
				textMaxE: tt.fields.textMaxE,
			}
			es.add(tt.args.time, tt.args.start, tt.args.text)
			if !reflect.DeepEqual(*es, tt.want) {
				t.Errorf("eventStatistic.add() %s = %v, want %v", tt.name, *es, tt.want)
			}
		})
	}
}

func Test_eventProperty_init(t *testing.T) {
	t.Parallel()

	type fields struct {
		values [16]eventStatistic
	}
	tests := []struct {
		name   string
		fields fields
		want   eventStatistic
	}{
		{"init", fields{[16]eventStatistic{{start: 1.234, evFirst: true, evStart: true, count: 123, min: 44, max: 55}}},
			eventStatistic{start: 1.234, evFirst: false, evStart: false, count: 0, min: math.MaxFloat64, max: 0}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ep := &eventProperty{
				values: tt.fields.values,
			}
			ep.init()
			if !reflect.DeepEqual(ep.values[0], tt.want) {
				t.Errorf("eventProperty.init() %s = %v, want %v", tt.name, ep.values[0], tt.want)
			}
		})
	}
}

func Test_eventProperty_add(t *testing.T) {
	t.Parallel()

	type fields struct {
		values [16]eventStatistic
	}
	type args struct {
		time  float64
		idx   uint16
		start bool
		text  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantev bool
		wantst float64
	}{
		{"add1", fields{}, args{47.11, 1, true, "text"}, true, 47.11},
		{"addsStop", fields{}, args{47.12, 15, false, "text1"}, false, 0.0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ep := &eventProperty{
				values: tt.fields.values,
			}
			ep.add(tt.args.time, tt.args.idx, tt.args.start, tt.args.text)
			if ep.values[tt.args.idx].evStart != tt.wantev {
				t.Errorf("eventProperty.add() %s = %v, want %v", tt.name,
					ep.values[tt.args.idx].evStart, tt.wantev)
			}
			if ep.values[tt.args.idx].start != tt.wantst {
				t.Errorf("eventProperty.add() %s = %v, want %v", tt.name,
					ep.values[tt.args.idx].start, tt.wantst)
			}
		})
	}
}

func Test_eventProperty_getCount(t *testing.T) {
	t.Parallel()

	type fields struct {
		values [16]eventStatistic
	}
	type args struct {
		idx uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"test0", fields{[16]eventStatistic{{count: 1}, {count: 2}}}, args{0}, 1},
		{"test1", fields{[16]eventStatistic{{count: 1}, {count: 2}}}, args{1}, 2},
		{"test16", fields{[16]eventStatistic{0: {count: 1}, 15: {count: 2}}}, args{16}, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ep := &eventProperty{
				values: tt.fields.values,
			}
			if got := ep.getCount(tt.args.idx); got != tt.want {
				t.Errorf("eventProperty.getCount() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_eventProperty_getAddCount(t *testing.T) {
	t.Parallel()

	type fields struct {
		values [16]eventStatistic
	}
	type args struct {
		idx uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"test0", fields{[16]eventStatistic{{evStart: false}, {evStart: true}}}, args{0}, "  "},
		{"test1", fields{[16]eventStatistic{{evStart: true}, {evStart: true}}}, args{1}, "+1"},
		{"test16", fields{[16]eventStatistic{0: {evStart: true}, 15: {evStart: true}}}, args{16}, "  "},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ep := &eventProperty{
				values: tt.fields.values,
			}
			if got := ep.getAddCount(tt.args.idx); got != tt.want {
				t.Errorf("eventProperty.getAddCount() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_convertUnit(t *testing.T) {
	t.Parallel()

	type args struct {
		v    float64
		unit string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"giga", args{1e9, "x"}, "  1.00000Gx"},
		{"mega", args{1e6, "x"}, "  1.00000Mx"},
		{"kilo", args{1e3, "x"}, "  1.00000kx"},
		{"nix", args{1.0, "x"}, "  1.00000x "},
		{"null", args{0.0, "x"}, "  0.00000x "},
		{"milli", args{1e-3, "x"}, "  1.00000mx"},
		{"micro", args{1e-6, "x"}, "  1.00000µx"},
		{"nano", args{1e-9, "x"}, "  1.00000nx"},
		{"tooSmall", args{1e-12, "x"}, "  0.00000x"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := convertUnit(tt.args.v, tt.args.unit); got != tt.want {
				t.Errorf("convertUnit() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_eventProperty_getTot(t *testing.T) {
	t.Parallel()

	type fields struct {
		values [16]eventStatistic
	}
	type args struct {
		idx uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"test", fields{[16]eventStatistic{{tot: 1.234}}}, args{0}, "  1.23400s "},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ep := &eventProperty{
				values: tt.fields.values,
			}
			if got := ep.getTot(tt.args.idx); got != tt.want {
				t.Errorf("eventProperty.getTot() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_eventProperty_getMin(t *testing.T) {
	t.Parallel()

	type fields struct {
		values [16]eventStatistic
	}
	type args struct {
		idx uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"test", fields{[16]eventStatistic{{min: 1.234}}}, args{0}, "  1.23400s "},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ep := &eventProperty{
				values: tt.fields.values,
			}
			if got := ep.getMin(tt.args.idx); got != tt.want {
				t.Errorf("eventProperty.getMin() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_eventProperty_getMax(t *testing.T) {
	t.Parallel()

	type fields struct {
		values [16]eventStatistic
	}
	type args struct {
		idx uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"test", fields{[16]eventStatistic{{max: 1.234}}}, args{0}, "  1.23400s "},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ep := &eventProperty{
				values: tt.fields.values,
			}
			if got := ep.getMax(tt.args.idx); got != tt.want {
				t.Errorf("eventProperty.getMax() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_eventProperty_getAvg(t *testing.T) {
	t.Parallel()

	type fields struct {
		values [16]eventStatistic
	}
	type args struct {
		idx uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"test", fields{[16]eventStatistic{{avg: 1.234}}}, args{0}, "  0.00000s "},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ep := &eventProperty{
				values: tt.fields.values,
			}
			if got := ep.getAvg(tt.args.idx); got != tt.want {
				t.Errorf("eventProperty.getAvg() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_eventProperty_getFirst(t *testing.T) {
	t.Parallel()

	type fields struct {
		values [16]eventStatistic
	}
	type args struct {
		idx uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"test", fields{[16]eventStatistic{{first: 1.234}}}, args{0}, "  1.23400s "},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ep := &eventProperty{
				values: tt.fields.values,
			}
			if got := ep.getFirst(tt.args.idx); got != tt.want {
				t.Errorf("eventProperty.getFirst() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_eventProperty_getLast(t *testing.T) {
	t.Parallel()

	type fields struct {
		values [16]eventStatistic
	}
	type args struct {
		idx uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"test", fields{[16]eventStatistic{{last: 1.234}}}, args{0}, "  1.23400s "},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ep := &eventProperty{
				values: tt.fields.values,
			}
			if got := ep.getLast(tt.args.idx); got != tt.want {
				t.Errorf("eventProperty.getLast() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestOutput_buildStatistic(t *testing.T) { //nolint:golint,paralleltest
	eds0 := make(scvd.Events)
	eds := make(scvd.Events)
	eds[0xEF00] = scvd.EventType{Brief: "briefbriefbrief", Property: "propertypropertyproperty", Value: "value"}

	tds := make(eval.Typedefs)

	var s1 = "../../testdata/test1.binary"
	var s3 = "../../testdata/test3.binary"
	var s4 = "../../testdata/test4.binary"
	var s6 = "../../testdata/test6.binary"
	var s7 = "../../testdata/test7.binary"

	type fields struct {
		evProps       [4]eventProperty
		columns       []string
		componentSize int
		propertySize  int
	}
	type args struct {
		file     string
		evdefs   scvd.Events
		typedefs eval.Typedefs
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  int
		want2  int
		want3  float64
	}{
		{"test1", fields{[4]eventProperty{}, []string{"Index", "Time (s)", "Component", "Event Property", "Value"}, 0, 0}, args{s1, eds0, tds}, 0, 9, 14, 0.0},
		{"test3", fields{[4]eventProperty{}, []string{"Index", "Time (s)", "Component", "Event Property", "Value"}, 0, 0}, args{s3, eds0, tds}, 1, 9, 14, 0.0},
		{"test4", fields{[4]eventProperty{}, []string{"Index", "Time (s)", "Component", "Event Property", "Value"}, 0, 0}, args{s4, eds0, tds}, 1, 9, 14, 0.5},
		{"test6", fields{[4]eventProperty{}, []string{"Index", "Time (s)", "Component", "Event Property", "Value"}, 0, 0}, args{s6, eds0, tds}, 1, 9, 14, 0.25},
		{"test7a", fields{[4]eventProperty{}, []string{"Index", "Time (s)", "Component", "Event Property", "Value"}, 0, 0}, args{s7, eds0, tds}, 1, 9, 14, 0.25},
		{"test7b", fields{[4]eventProperty{}, []string{"Index", "Time (s)", "Component", "Event Property", "Value"}, 0, 0}, args{s7, eds, tds}, 1, 15, 24, 0.25},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			o := &Output{
				evProps:       tt.fields.evProps,
				columns:       tt.fields.columns,
				componentSize: tt.fields.componentSize,
				propertySize:  tt.fields.propertySize,
			}
			TimeFactor = nil
			var b event.Binary
			in := b.Open(&tt.args.file)
			if got := o.buildStatistic(in, tt.args.evdefs, tt.args.typedefs); got != tt.want {
				t.Errorf("Output.buildStatistic() %s = %v, want %v", tt.name, got, tt.want)
			}
			b.Close()
			if o.componentSize != tt.want1 || o.propertySize != tt.want2 {
				t.Errorf("Output.buildStatistic() %s = %v,%v, want %v,%v", tt.name, o.componentSize, o.propertySize, tt.want1, tt.want2)
			}
			if TimeFactor != nil && *TimeFactor != tt.want3 {
				t.Errorf("Output.buildStatistic() %s = %v, want %v", tt.name, TimeFactor, tt.want3)
			}
		})
	}
}

func TestOutput_printStatistic(t *testing.T) { //nolint:golint,paralleltest
	var b bytes.Buffer

	props0 := [4]eventProperty{}
	props1 := [4]eventProperty{{[16]eventStatistic{{evFirst: true, count: 1, tot: 2, min: 3, max: 4, avg: 5, first: 6, last: 7}}}}

	header := "   Start/Stop event statistic\n" +
		"   --------------------------\n\n" +
		"Event count      total       min         max         average     first       last\n" +
		"----- -----      -----       ---         ---         -------     -----       ----\n"

	line1 := "A(0)      1     2.00000s    3.00000s    4.00000s    5.00000s    6.00000s    7.00000s \n" +
		"      Min: Start: 0.00000000  Stop: 3.00000000 \n" +
		"      Max: Start: 0.00000000  Stop: 4.00000000 \n\n"

	type fields struct {
		evProps       [4]eventProperty
		componentSize int
		propertySize  int
	}
	type args struct {
		out        *bufio.Writer
		eventCount int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"header", fields{props0, 15, 20}, args{nil, 1}, header, false},
		{"line1", fields{props1, 15, 20}, args{nil, 1}, header + line1, false},
	}
	eventsTable := EventsTable{
		Events:     []EventRecord{},
		Statistics: []EventRecordStatistic{},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			tt.args.out = bufio.NewWriter(&b)
			o := &Output{
				evProps:       tt.fields.evProps,
				componentSize: tt.fields.componentSize,
				propertySize:  tt.fields.propertySize,
			}
			if err := o.printStatistic(tt.args.out, tt.args.eventCount, &eventsTable); (err != nil) != tt.wantErr {
				t.Errorf("Output.printStatistic() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.args.out.Flush()
			str, err := b.ReadString('\000')
			if err != nil && !errors.Is(err, io.EOF) {
				t.Errorf("Output.printStatistic() err = %v", err)
			}
			if str != tt.want {
				t.Errorf("Output.printStatistic() %s = %v, want %v", tt.name, str, tt.want)
			}
		})
	}
}

func Test_escapeGen(t *testing.T) {
	t.Parallel()

	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"apo", args{"a\x27-"}, "a\\'-"},
		{"quote", args{"a\x22-"}, "a\\\"-"},
		{"a", args{"a\x07-"}, "a\\a-"},
		{"b", args{"a\x08-"}, "a\\b-"},
		{"esc", args{"a\x1b-"}, "a\\e-"},
		{"f", args{"a\x0c-"}, "a\\f-"},
		{"n", args{"a\x0a-"}, "a\\n-"},
		{"r", args{"a\x0d-"}, "a\\r-"},
		{"t", args{"a\x09-"}, "a\\t-"},
		{"v", args{"a\x0b-"}, "a\\v-"},
		{"?", args{"a\x01-"}, "a\\001-"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := escapeGen(tt.args.s); got != tt.want {
				t.Errorf("escapeGen() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestOutput_printEvents(t *testing.T) { //nolint:golint,paralleltest
	var b bytes.Buffer

	eds := make(scvd.Events)
	eds[0xFE00] = scvd.EventType{Brief: "briefbriefbrief", Property: "propertypropertyproperty", Value: "value"}
	eds[0xFF03] = scvd.EventType{Brief: "briefbriefbrief", Property: "propertypropertyproperty", Value: "value"}

	var s0 = "../../testdata/test0.binary"
	var s1 = "../../testdata/test1.binary"
	var s10 = "../../testdata/test10.binary"
	var s11 = "../../testdata/test11.binary"
	var sNix = "../../testdata/xxxx"

	line1 := "    0 0.00000124 0xFF     0xFF03       val1=0x00000004, val2=0x00000002\n" +
		"    1 0.00000124 0xFE     0xFE00       \"hello wo\"\n"
	line2 := "    0 0.00000124 briefbriefbrief propertypropertyproperty value\n" +
		"    1 0.00000124 briefbriefbrief propertypropertyproperty \"hello wo\"\n"
	line3 := "    0 0.00000124 0xFF     0xFF00       val1=0x00000004, val2=0x00000002\n" +
		"    1 0.00000124 0xFE     0xFE00       \"hello wo\"\n"

	type fields struct {
		evProps       [4]eventProperty
		columns       []string
		componentSize int
		propertySize  int
	}
	type args struct {
		out      *bufio.Writer
		in       *bufio.Reader
		evdefs   scvd.Events
		typedefs eval.Typedefs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		file    *string
		want    string
		wantErr bool
	}{
		{"readErr0", fields{}, args{}, &s0, "", false},
		{"readErr1", fields{}, args{}, &s1, "", true},
		{"read1", fields{}, args{}, &s10, line1, false},
		{"read2", fields{}, args{evdefs: eds}, &s10, line2, false},
		{"read3", fields{}, args{}, &s11, line3, false},
		{"readNix", fields{}, args{}, &sNix, "", false},
	}
	eventsTable := EventsTable{
		Events:     []EventRecord{},
		Statistics: []EventRecordStatistic{},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			tt.args.out = bufio.NewWriter(&b)

			TimeFactor = nil
			var ib event.Binary
			tt.args.in = ib.Open(tt.file)
			o := &Output{
				evProps:       tt.fields.evProps,
				columns:       tt.fields.columns,
				componentSize: tt.fields.componentSize,
				propertySize:  tt.fields.propertySize,
			}
			if err := o.printEvents(tt.args.out, tt.args.in, tt.args.evdefs, tt.args.typedefs, &eventsTable); (err != nil) != tt.wantErr {
				t.Errorf("Output.printEvents() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			tt.args.out.Flush()
			str, err := b.ReadString('\000')
			if err != nil && !errors.Is(err, io.EOF) {
				t.Errorf("Output.printStatistic() err = %v", err)
			}
			if str != tt.want {
				t.Errorf("Output.printStatistic() %s = %v, want %v", tt.name, str, tt.want)
			}
		})
	}
}

func TestOutput_printHeader(t *testing.T) { //nolint:golint,paralleltest
	var b bytes.Buffer

	type fields struct {
		evProps       [4]eventProperty
		columns       []string
		componentSize int
		propertySize  int
	}
	type args struct {
		out *bufio.Writer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want1  string
		want2  string
	}{
		{"test", fields{columns: []string{"a", "b", "c", "d", "e"}, componentSize: 15, propertySize: 20}, args{}, "c", "d"},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			tt.args.out = bufio.NewWriter(&b)
			o := &Output{
				evProps:       tt.fields.evProps,
				columns:       tt.fields.columns,
				componentSize: tt.fields.componentSize,
				propertySize:  tt.fields.propertySize,
			}
			if err := o.printHeader(tt.args.out); err != nil {
				t.Errorf("printHeader() err = %v", err)
			}
			tt.args.out.Flush()
			str, err := b.ReadString('\000')
			if err != nil && !errors.Is(err, io.EOF) {
				t.Errorf("printHeader() err = %v, want %v", err, tt.want1)
			}
			out := fmt.Sprintf("%5s %-10s %*s %*s %s", "a", "b",
				-o.componentSize, tt.want1, -o.propertySize, tt.want2, "e")
			i := strings.IndexByte(str, '\n')
			if i > 0 {
				j := (i+1)*2 + 1
				if j >= len(str) {
					t.Errorf("printHeader() = %v, want %v", str[i+1:], out)
				} else {
					str = str[j:]
					i = strings.IndexByte(str, '\n')
					if i > 0 {
						str = str[:i]
					}
				}
			}
			if i < 0 {
				t.Errorf("printHeader() = %v, want %v", str, out)
			}
			if str != out {
				t.Errorf("printHeader() = %v, want %v", str, out)
			}
		})
	}
}

func TestOutput_print(t *testing.T) { //nolint:golint,paralleltest
	var b bytes.Buffer

	//	var e0 = "../../testdata/test.xml"
	var s10 = "../../testdata/test10.binary"
	var s11 = "../../testdata/nix.binary"

	line1 := "   Detailed event list\n" +
		"   -------------------\n\n" +
		"Index Time (s)   Component Event Property Value\n" +
		"----- --------   --------- -------------- -----\n" +
		"    0 7.75000000 0xFF      0xFF03         val1=0x00000004, val2=0x00000002\n" +
		"    1 7.75000000 0xFE      0xFE00         \"hello wo\"\n\n" +
		"   Start/Stop event statistic\n" +
		"   --------------------------\n\n" +
		"Event count      total       min         max         average     first       last\n" +
		"----- -----      -----       ---         ---         -------     -----       ----\n"

	line2 := "   Start/Stop event statistic\n" +
		"   --------------------------\n\n" +
		"Event count      total       min         max         average     first       last\n" +
		"----- -----      -----       ---         ---         -------     -----       ----\n\n" +
		"   Detailed event list\n" +
		"   -------------------\n\n" +
		"Index Time (s)   Component Event Property Value\n" +
		"----- --------   --------- -------------- -----\n" +
		"    0 7.75000000 0xFF      0xFF03         val1=0x00000004, val2=0x00000002\n" +
		"    1 7.75000000 0xFE      0xFE00         \"hello wo\"\n"

	type fields struct {
		evProps       [4]eventProperty
		columns       []string
		componentSize int
		propertySize  int
	}
	type args struct {
		out           *bufio.Writer
		eventFile     *string
		evdefs        scvd.Events
		typedefs      eval.Typedefs
		statBegin     bool
		showStatistic bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"stat empty", fields{}, args{eventFile: nil}, line1, true},
		{"stat wrong", fields{}, args{eventFile: &s11}, line1, true},
		{"statEnd", fields{}, args{eventFile: &s10}, line1, false},
		{"statBegin", fields{}, args{eventFile: &s10, statBegin: true}, line2, false},
	}
	eventsTable := EventsTable{
		Events:     []EventRecord{},
		Statistics: []EventRecordStatistic{},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			tt.args.out = bufio.NewWriter(&b)

			TimeFactor = nil
			o := &Output{
				evProps:       tt.fields.evProps,
				columns:       tt.fields.columns,
				componentSize: tt.fields.componentSize,
				propertySize:  tt.fields.propertySize,
			}
			if err := o.print(tt.args.out, tt.args.eventFile, tt.args.evdefs, tt.args.typedefs, tt.args.statBegin, tt.args.showStatistic, &eventsTable); (err != nil) != tt.wantErr {
				t.Errorf("Output.print() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.args.out.Flush()
			str, err := b.ReadString('\000')
			if err != nil && !errors.Is(err, io.EOF) {
				t.Errorf("Output.print() err = %v", err)
			}
			if !errors.Is(err, io.EOF) && str != tt.want {
				t.Errorf("Output.print() %s = %v, want %v", tt.name, str, tt.want)
			}
		})
	}
}

func TestPrint(t *testing.T) { //nolint:golint,paralleltest
	o1 := "testOutput.out"

	var s10 = "../../testdata/test10.binary"

	lines1 := [...]string{
		"   Detailed event list\n",
		"   -------------------\n",
		"\n",
		"Index Time (s)   Component Event Property Value\n",
		"----- --------   --------- -------------- -----\n",
		"    0 7.75000000 0xFF      0xFF03         val1=0x00000004, val2=0x00000002\n",
		"    1 7.75000000 0xFE      0xFE00         \"hello wo\"\n",
		"\n",
		"   Start/Stop event statistic\n",
		"   --------------------------\n",
		"\n",
		"Event count      total       min         max         average     first       last\n",
		"----- -----      -----       ---         ---         -------     -----       ----\n"}

	type args struct {
		filename      *string
		eventFile     *string
		evdefs        scvd.Events
		typedefs      eval.Typedefs
		statBegin     bool
		showStatistic bool
	}
	formatType := "txt"
	level := ""
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test", args{filename: &o1, eventFile: &s10}, false},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			TimeFactor = nil
			defer os.Remove(*tt.args.filename)
			if err := Print(tt.args.filename, &formatType, &level, tt.args.eventFile, tt.args.evdefs, tt.args.typedefs, tt.args.statBegin, tt.args.showStatistic); (err != nil) != tt.wantErr {
				t.Errorf("Print() error = %v, wantErr %v", err, tt.wantErr)
			}
			file, err := os.Open(*tt.args.filename)
			if err != nil {
				t.Errorf("Print() error = %v, output file not created", err)
			}
			if file != nil {
				defer file.Close()
				in := bufio.NewReader(file)
				var l string
				end := false
				for _, l = range lines1 {
					line, err := in.ReadString('\n')
					if errors.Is(err, io.EOF) {
						end = true
						break // end of lines reached
					}
					if line != l {
						t.Errorf("Print() %s = %v, want %v", tt.name, line, l)
					}
				}
				line, err := in.ReadString('\n')
				if errors.Is(err, io.EOF) {
					end = true
				} else {
					t.Errorf("Print() %s = %v, want EOF", tt.name, line)
				}
				if !end {
					t.Errorf("Print() %s = EOF, want %v", tt.name, l)
				}
			}
		})
	}
}

func TestPrintJSON(t *testing.T) { //nolint:golint,paralleltest
	o1 := "testOutput.json"

	var s10 = "../../testdata/test10.binary"

	lines1 := [...]string{
		"{\"events\":[{\"index\":0,\"time\":7.75,\"component\":\"0xFF\",\"eventProperty\":\"0xFF03\",\"value\":\"val1=0x00000004, val2=0x00000002\"},{\"index\":1,\"time\":7.75,\"component\":\"0xFE\",\"eventProperty\":\"0xFE00\",\"value\":\"hello wo\"}],\"statistics\":[]}",
	}

	type args struct {
		filename      *string
		eventFile     *string
		evdefs        scvd.Events
		typedefs      eval.Typedefs
		statBegin     bool
		showStatistic bool
	}
	formatType := "json"
	level := ""
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{filename: &o1, eventFile: &s10}, false},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			TimeFactor = nil
			defer os.Remove(*tt.args.filename)
			if err := Print(tt.args.filename, &formatType, &level, tt.args.eventFile, tt.args.evdefs, tt.args.typedefs, tt.args.statBegin, tt.args.showStatistic); (err != nil) != tt.wantErr {
				t.Errorf("Print() error = %v, wantErr %v", err, tt.wantErr)
			}
			file, err := os.Open(*tt.args.filename)
			if err != nil {
				t.Errorf("Print() error = %v, output file not created", err)
			}
			if file != nil {
				defer file.Close()
				in := bufio.NewReader(file)
				var l string
				end := false
				for _, l = range lines1 {
					line, _ := in.ReadString('\n')
					if line != l {
						t.Errorf("Print() %s = %v, want %v", tt.name, line, l)
					}
				}
				line, err := in.ReadString('\n')
				if errors.Is(err, io.EOF) {
					end = true
				} else {
					t.Errorf("Print() %s = %v, want EOF", tt.name, line)
				}
				if !end {
					t.Errorf("Print() %s = EOF, want %v", tt.name, l)
				}
			}
		})
	}
}

func TestPrintXML(t *testing.T) { //nolint:golint,paralleltest
	o1 := "testOutput.xml"

	var s10 = "../../testdata/test10.binary"

	lines1 := [...]string{
		"<EventsTable><events><index>0</index><time>7.75</time><component>0xFF</component><eventProperty>0xFF03</eventProperty><value>val1=0x00000004, val2=0x00000002</value></events><events><index>1</index><time>7.75</time><component>0xFE</component><eventProperty>0xFE00</eventProperty><value>hello wo</value></events></EventsTable>",
	}

	type args struct {
		filename      *string
		eventFile     *string
		evdefs        scvd.Events
		typedefs      eval.Typedefs
		statBegin     bool
		showStatistic bool
	}
	formatType := "xml"
	level := ""
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test", args{filename: &o1, eventFile: &s10}, false},
	}
	for _, tt := range tests { //nolint:golint,paralleltest
		t.Run(tt.name, func(t *testing.T) {
			TimeFactor = nil
			defer os.Remove(*tt.args.filename)
			if err := Print(tt.args.filename, &formatType, &level, tt.args.eventFile, tt.args.evdefs, tt.args.typedefs, tt.args.statBegin, tt.args.showStatistic); (err != nil) != tt.wantErr {
				t.Errorf("Print() error = %v, wantErr %v", err, tt.wantErr)
			}
			file, err := os.Open(*tt.args.filename)
			if err != nil {
				t.Errorf("Print() error = %v, output file not created", err)
			}
			if file != nil {
				defer file.Close()
				in := bufio.NewReader(file)
				var l string
				end := false
				for _, l = range lines1 {
					line, _ := in.ReadString('\n')
					if line != l {
						t.Errorf("Print() %s = %v, want %v", tt.name, line, l)
					}
				}
				line, err := in.ReadString('\n')
				if errors.Is(err, io.EOF) {
					end = true
				} else {
					t.Errorf("Print() %s = %v, want EOF", tt.name, line)
				}
				if !end {
					t.Errorf("Print() %s = EOF, want %v", tt.name, l)
				}
			}
		})
	}
}
