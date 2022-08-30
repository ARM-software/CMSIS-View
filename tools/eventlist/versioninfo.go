// Auto-generated file by goversioninfo. Do not edit.
package main

import (
	"encoding/json"

	"github.com/josephspurrier/goversioninfo"
)

func unmarshalGoVersionInfo(b []byte) goversioninfo.VersionInfo {
	vi := goversioninfo.VersionInfo{}
	json.Unmarshal(b, &vi)
	return vi
}

var versionInfo = unmarshalGoVersionInfo([]byte(`{
	"FixedFileInfo":{
		"FileVersion": {
			"Major": 1,
			"Minor": 0,
			"Patch": 9,
			"Build": 0
		},
		"ProductVersion": {
			"Major": 1,
			"Minor": 0,
			"Patch": 9,
			"Build": 0
		},
		"FileFlagsMask": "3f",
		"FileFlags": "",
		"FileOS": "040004",
		"FileType": "01",
		"FileSubType": "00"
	},
	"StringFileInfo":{
		"Comments": "",
		"CompanyName": "ARM Limited",
		"FileDescription": "",
		"FileVersion": "1.0.9.0",
		"InternalName": "eventlist",
		"LegalCopyright": "Copyright (C) 2022 ARM Limited or its Affiliates. All rights reserved.",
		"LegalTrademarks": "",
		"OriginalFilename": "eventlist",
		"PrivateBuild": "",
		"ProductName": "eventlist",
		"ProductVersion": "1.0.9.0",
		"SpecialBuild": ""
	},
	"VarFileInfo":{
		"Translation": {
			"LangID": 1033,
			"CharsetID": 1200
		}
	}
}`))
