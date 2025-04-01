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

package scvd

import (
	"encoding/xml"
	"eventlist/pkg/eval"
	"os"
	"strconv"
	"strings"
)

type Value string
type ID string
type Endian string // B, b, L, l

type ComponentsType struct {
	Name      string `xml:"name,attr"`
	Shortname string `xml:"shortname,attr"`
	Version   string `xml:"version,attr"`
}

type EnumType struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
	Info  string `xml:"info,attr"`
}

type MemberType struct {
	Name   string     `xml:"name,attr"`
	Type   string     `xml:"type,attr"`
	Offset string     `xml:"offset,attr"`
	Size   uint       `xml:"size,attr"`
	Info   string     `xml:"info,attr"`
	Endian string     `xml:"endian,attr"`
	Enums  []EnumType `xml:"enum"`
}

type VarType struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
	Type  string `xml:"type,attr"`
	Size  uint   `xml:"size,attr"`
	Info  string `xml:"info,attr"`
}

type TypedefType struct {
	Name    string       `xml:"name,attr"`
	Size    uint         `xml:"size,attr"`
	Const   bool         `xml:"const,attr"`
	Info    string       `xml:"info,attr"`
	Endian  string       `xml:"endian,attr"`
	Import  string       `xml:"import,attr"`
	Members []MemberType `xml:"member"`
	Vars    []VarType    `xml:"var"`
}

type TypedefsType struct {
	Typedef []TypedefType `xml:"typedef"`
}

type PrintType struct {
	Cond     string `xml:"cond,attr"`
	Value    Value  `xml:"value,attr"`
	Property string `xml:"property,attr"`
	Alert    string `xml:"alert,attr"`
	Bold     string `xml:"bold,attr"`
}

type EventType struct {
	ID       ID          `xml:"id,attr"`    // limits to 16 bit
	Level    string      `xml:"level,attr"` // Enum: Error, API, Op, Detail
	Name     string      `xml:"name,attr"`
	Brief    string      `xml:"brief,attr"`
	Val1     string      `xml:"val1,attr"`
	Val2     string      `xml:"val2,attr"`
	Val3     string      `xml:"val3,attr"`
	Val4     string      `xml:"val4,attr"`
	Val5     string      `xml:"val5,attr"`
	Val6     string      `xml:"val6,attr"`
	Value    Value       `xml:"value,attr"`
	Property string      `xml:"property,attr"`
	Info     string      `xml:"info,attr"`
	Doc      string      `xml:"doc,attr"`
	Alert    string      `xml:"alert,attr"` // expression resolves to bool
	Bold     string      `xml:"bold,attr"`  // expression resolves to bool
	State    string      `xml:"state,attr"`
	Handle   string      `xml:"handle,attr"`
	HName    string      `xml:"hname,attr"`
	Tracking string      `xml:"tracking,attr"`
	Reset    bool        `xml:"reset,attr"`
	Prints   []PrintType `xml:"print"`
}

type State struct {
	Name string `xml:"name,attr"`
	Plot string `xml:"plot,attr"`
	Bold bool   `xml:"bold,attr"`
}

type ComponentType struct {
	Name   string  `xml:"name,attr"`
	Brief  string  `xml:"brief,attr"`
	No     string  `xml:"no,attr"`
	Prefix string  `xml:"prefix,attr"`
	Info   string  `xml:"info,attr"`
	States []State `xml:"state"`
}

type GroupType struct {
	Name      string          `xml:"name,attr"`
	Component []ComponentType `xml:"component"`
}

type EventsType struct {
	Group  GroupType   `xml:"group"`
	Events []EventType `xml:"event"`
}

type ComponentViewer struct {
	Component     ComponentsType `xml:"component"`
	Typedefs      TypedefsType   `xml:"typedefs"`
	Events        EventsType     `xml:"events"`
	SchemaVersion string         `xml:"schemaVersion,attr"`
}

type IDType uint16
type Events map[IDType]EventType

// getFromFile reads an XML file specified by the given filename and decodes its content
// into the ComponentViewer receiver. It returns an error if the file cannot be read or
// if the XML decoding fails.
//
// Parameters:
//   - name: A pointer to a string containing the filename of the XML file to be read.
//
// Returns:
//   - error: An error if the file cannot be read or if the XML decoding fails, otherwise nil.
func (viewer *ComponentViewer) getFromFile(name *string) error {
	data, err := os.ReadFile(*name)
	if err == nil {
		d := xml.NewDecoder(strings.NewReader(string(data)))
		err = d.Decode(&viewer)
	}
	return err
}

// getOne reads and processes event and typedef data from a specified file.
// It populates the provided Events and Typedefs structures with the extracted data.
//
// Parameters:
//   - filename: A pointer to the name of the file to read from.
//   - events: An Events structure to be populated with event data.
//   - typedefs: A Typedefs structure to be populated with typedef data.
//
// Returns:
//   - error: An error if any issues occur during file reading or data processing, otherwise nil.
func getOne(filename *string, events Events, typedefs eval.Typedefs) error {
	var viewer ComponentViewer
	var err error
	if err = viewer.getFromFile(filename); err == nil {
		// create a components map indexed by "no" to speed up things
		components := make(map[uint8]*ComponentType)
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
			id, err := eval.GetIdValue(string(event.ID), typedefs)
			if err != nil {
				return err // cannot decode IdValue
			}
			if components[uint8(id>>8)] != nil {
				event.Brief = components[uint8(id>>8)].Brief
			}
			events[IDType(id)] = event
		}
		// extract enums from typedefs
		for _, typedef := range viewer.Typedefs.Typedef {
			if len(typedef.Members) > 0 {
				members := make(map[string]eval.Member)
				for _, member := range typedef.Members {
					var mem eval.Member
					if len(member.Enums) > 0 {
						mem.Enums = make(map[int64]string)
						for _, enum := range member.Enums {
							var enu int64
							if enu, err = eval.GetValue(enum.Value, typedefs); err != nil {
								return err
							}
							mem.Enums[enu] = enum.Name
						}
					}
					mem.IType = eval.ITypes[member.Type]
					mem.Offset = member.Offset
					members[member.Name] = mem
				}
				if len(members) > 0 {
					typedefs[typedef.Name] = eval.ITypedef{Size: uint32(typedef.Size), BigEndian: typedef.Endian == "B" || typedef.Endian == "b", Members: members}
				}
			}
		}
	}
	return err
}

// Get processes a list of SCVD files and populates the provided events and typedefs.
//
// Parameters:
//
//	scvdFiles - A pointer to a slice of strings, where each string is a path to an SCVD file.
//	events - An Events object to be populated with data from the SCVD files.
//	typedefs - A Typedefs object to be populated with data from the SCVD files.
//
// Returns:
//
//	An error if any of the SCVD files could not be processed, otherwise nil.
func Get(scvdFiles *[]string, events Events, typedefs eval.Typedefs) error {
	if scvdFiles != nil {
		for _, scvdFile := range *scvdFiles {
			if err := getOne(&scvdFile, events, typedefs); err != nil {
				return err
			}
		}
	}
	return nil
}
