/*
 * Copyright (c) 2022-2023 Arm Limited. All rights reserved.
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

package scvd

import (
	"encoding/xml"
	"errors"
	"eventlist/pkg/eval"
	"os"
	"strconv"
	"strings"
)

type Value string
type ID string

type Component struct {
	Name      string `xml:"name,attr"`
	Shortname string `xml:"shortname,attr"`
	Version   string `xml:"version,attr"`
}

type Enum struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
	Info  string `xml:"info,attr"`
}

type Member struct {
	Name   string `xml:"name,attr"`
	Type   string `xml:"type,attr"`
	Offset string `xml:"offset,attr"`
	Info   string `xml:"info,attr"`
	Enums  []Enum `xml:"enum"`
}

type Var struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Info  string `xml:"info,attr"`
	Enums []Enum `xml:"enum"`
}

type Typedef struct {
	Name    string   `xml:"name,attr"`
	Info    string   `xml:"info,attr"`
	Size    string   `xml:"size,attr"`
	Members []Member `xml:"member"`
	Vars    []Var    `xml:"var"`
}

type Typedefs struct {
	Typedef []Typedef `xml:"typedef"`
}

type State struct {
	Name string `xml:"name,attr"`
	Plot string `xml:"plot,attr"`
}

type Event struct {
	ID       ID     `xml:"id,attr"`
	Level    string `xml:"level,attr"`
	Property string `xml:"property,attr"`
	Tracking string `xml:"tracking,attr"`
	State    string `xml:"state,attr"`
	Handle   string `xml:"handle,attr"`
	HName    string `xml:"hname,attr"`
	Value    Value  `xml:"value,attr"`
	Info     string `xml:"info,attr"`
	Val1     string `xml:"val1,attr"`
	Val2     string `xml:"val2,attr"`
	Val3     string `xml:"val3,attr"`
	Val4     string `xml:"val4,attr"`
	Brief    string
}

type GroupComponent struct {
	Name   string  `xml:"name,attr"`
	Brief  string  `xml:"brief,attr"`
	No     string  `xml:"no,attr"`
	Prefix string  `xml:"prefix,attr"`
	Info   string  `xml:"info,attr"`
	States []State `xml:"state"`
}

type Group struct {
	Name      string           `xml:"name,attr"`
	Component []GroupComponent `xml:"component"`
}

type Events struct {
	Group  Group   `xml:"group"`
	Events []Event `xml:"event"`
}

type ComponentViewer struct {
	Component Component `xml:"component"`
	Typedefs  Typedefs  `xml:"typedefs"`
	Events    Events    `xml:"events"`
}

func (viewer *ComponentViewer) getFromFile(name *string) error {
	data, err := os.ReadFile(*name)
	if err == nil {
		d := xml.NewDecoder(strings.NewReader(string(data)))
		err = d.Decode(&viewer)
	}
	return err
}

// get the enum value with calculation
func (enum *Enum) getInfo() (int16, error) {
	n, err := eval.Eval(&enum.Value, nil)
	if err != nil && !errors.Is(err, eval.ErrEof) {
		return 0, err
	}
	return int16(n.GetInt64()), nil
}

func (id *ID) getIdValue() (uint16, error) { //nolint:golint,revive
	sid := string(*id)
	n, err := eval.Eval(&sid, nil)
	if err != nil && !errors.Is(err, eval.ErrEof) {
		return 0, err
	}
	return uint16(n.GetInt64()), nil
}

type ScvdData struct {
	Events   map[uint16]Event
}

func (s *ScvdData) GetOne(filename *string) error {
	var viewer ComponentViewer
	var err error
	if err = viewer.getFromFile(filename); err == nil {
		// create a components map indexed by "no" to speed up things
		components := make(map[uint8]*GroupComponent)
		for _, component := range viewer.Events.Group.Component {
			var no int64
			no, err = strconv.ParseInt(component.No, 0, 0)
			if err != nil {
				break
			}
			components[uint8(no)] = &component
		}
		if err != nil {
			return err // cannot decode component number
		}
		for _, event := range viewer.Events.Events {
			id, err := event.ID.getIdValue()
			if err != nil {
				return err // cannot decode IdValue
			}
			if components[uint8(id>>8)] != nil {
				event.Brief = components[uint8(id>>8)].Brief
			}
			s.Events[id] = event
		}
		// extract enums from typedefs
		for _, typedef := range viewer.Typedefs.Typedef {
			if len(typedef.Members) > 0 {
				members := make(map[string]eval.TdMember)
				for _, member := range typedef.Members {
					t := members[member.Name]
					var off int64
					off, err = strconv.ParseInt(member.Offset, 0, 0)
					if err != nil {
						return err
					}
					t.Offset = int32(off)
					ty := eval.ITypes[member.Type]
					t.Type = ty
					members[member.Name] = t
					for _, enum := range member.Enums {
						var en int16
						if en, err = enum.getInfo(); err != nil {
							return err
						}
						t := members[member.Name]
						t.Enum = make(map[int16]string)
						t.Enum[en] = enum.Name
						members[member.Name] = t
					}
				}
				if len(members) > 0 {
					eval.Typedefs[typedef.Name] = members
				}
			}
		}
	}
	return err
}

// returns the events and typedef map
func (s *ScvdData) Get(scvdFiles *[]string) error {
	if scvdFiles != nil {
		s.Events = make(map[uint16]Event)
		eval.Typedefs = make(map[string]map[string]eval.TdMember)
		for _, scvdFile := range *scvdFiles {
			if err := s.GetOne(&scvdFile); err != nil {
				return err
			}
		}
	}
	return nil
}
