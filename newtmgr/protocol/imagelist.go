/*
 Copyright 2015 Runtime Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package protocol

import (
	"encoding/binary"
	"fmt"
)

type ImageList struct {
	Images []string
}

const (
	IMGMGR_NMGR_OP_LIST   = 0
	IMGMGR_NMGR_OP_UPLOAD = 1
	IMGMGR_NMGR_OP_BOOT   = 2
)

func NewImageList() (*ImageList, error) {
	s := &ImageList{}
	s.Images = []string{}
	return s, nil
}

func ImageVersStr(major uint8, minor uint8, revision uint16, buildNum uint32) string {
	if major == 0xff && minor == 0xff && revision == 0xffff &&
		buildNum == 0xffffffff {
		return "Not set"
	} else {
		versStr := fmt.Sprintf("%d.%d.%d.%d", major, minor, revision, buildNum)
		return versStr
	}
}

func (i *ImageList) EncodeWriteRequest() (*NmgrReq, error) {
	nmr, err := NewNmgrReq()
	if err != nil {
		return nil, err
	}

	nmr.Op = NMGR_OP_READ
	nmr.Flags = 0
	nmr.Group = NMGR_GROUP_ID_IMAGE
	nmr.Id = IMGMGR_NMGR_OP_LIST
	nmr.Len = 0

	return nmr, nil
}

func DecodeImageListResponse(data []byte) (*ImageList, error) {
	i := &ImageList{}

	for len(data) >= 8 {
		major := uint8(data[0])
		minor := uint8(data[1])
		revision := binary.BigEndian.Uint16(data[2:4])
		buildNum := binary.BigEndian.Uint32(data[4:8])
		data = data[8:]

		versStr := ImageVersStr(major, minor, revision, buildNum)
		i.Images = append(i.Images, versStr)
	}

	return i, nil
}
