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
	n, err := eval.Eval(&enum.Value)
	if err != nil && !errors.Is(err, eval.ErrEof) {
		return 0, err
	}
	return int16(n.GetInt()), nil
}

func (id *ID) getIdValue() (uint16, error) { //nolint:golint,revive
	sid := string(*id)
	n, err := eval.Eval(&sid)
	if err != nil && !errors.Is(err, eval.ErrEof) {
		return 0, err
	}
	return uint16(n.GetInt()), nil
}

func getOne(filename *string, events map[uint16]Event,
	typedefs map[string]map[string]map[int16]string) error {
	var viewer ComponentViewer
	var err error
	if err = viewer.getFromFile(filename); err == nil {
		// create a components map indexed by "no" to speed up things
		components := make(map[uint8]*GroupComponent)
		for _, component := range viewer.Events.Group.Component {
			var no uint64
			no, err = strconv.ParseUint(component.No, 0, 8)
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
			events[id] = event
		}
		// extract enums from typedefs
		for _, typedef := range viewer.Typedefs.Typedef {
			if len(typedef.Members) > 0 {
				members := make(map[string]map[int16]string)
				for _, member := range typedef.Members {
					if len(member.Enums) > 0 {
						enums := make(map[int16]string)
						for _, enum := range member.Enums {
							var en int16
							if en, err = enum.getInfo(); err != nil {
								return err
							}
							enums[en] = enum.Name
						}
						members[member.Name] = enums
					}
				}
				if len(members) > 0 {
					typedefs[typedef.Name] = members
				}
			}
		}
	}
	return err
}

// returns the events and typedef map
func Get(scvdFiles *[]string, events map[uint16]Event,
	typedefs map[string]map[string]map[int16]string) error {
	if scvdFiles != nil {
		for _, scvdFile := range *scvdFiles {
			if err := getOne(&scvdFile, events, typedefs); err != nil {
				return err
			}
		}
	}
	return nil
}
