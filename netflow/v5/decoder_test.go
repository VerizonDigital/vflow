//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    decoder_test.go
//: details: netflow v5 decoder tests
//: author:  Christopher Noel
//: date:    12/10/2018ls
//:
//: Licensed under the Apache License, Version 2.0 (the "License");
//: you may not use this file except in compliance with the License.
//: You may obtain a copy of the License at
//:
//:     http://www.apache.org/licenses/LICENSE-2.0
//:
//: Unless required by applicable law or agreed to in writing, software
//: distributed under the License is distributed on an "AS IS" BASIS,
//: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//: See the License for the specific language governing permissions and
//: limitations under the License.
//: ----------------------------------------------------------------------------

package netflow5

import (
	"bytes"
	"fmt"
	"net"
	"testing"
)

var TestV5FlowPacket = []byte{0x00, 0x05, 0x00, 0x1d, 0x03, 0x11, 0x5d, 0xd8, 0x5c, 0x0e, 0xd7, 0xa5, 0x00, 0x00, 0x00, 0x00, 0x34, 0x16, 0x41, 0xa6, 0x00, 0x00, 0x03, 0xe8, 0x7d, 0xee, 0x2e, 0x30, 0x72, 0x17, 0xec, 0x60, 0x72, 0x17, 0x03, 0xe7, 0x03, 0x17, 0x03, 0x31, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x06, 0xac, 0x03, 0x10, 0x55, 0xa1, 0x03, 0x10, 0xcf, 0x30, 0xc0, 0x51, 0x01, 0xbb, 0x00, 0x10, 0x06, 0x00, 0x12, 0xa3, 0xda, 0xde, 0x14, 0x16, 0x00, 0x00, 0x7d, 0xee, 0x2e, 0x30, 0x72, 0x17, 0xec, 0x60, 0x72, 0x17, 0x03, 0xe7, 0x03, 0x17, 0x03, 0x31, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x01, 0xb9, 0x03, 0x10, 0xaf, 0x71, 0x03, 0x10, 0xaf, 0x71, 0xc0, 0x51, 0x01, 0xbb, 0x00, 0x18, 0x06, 0x00, 0x12, 0xa3, 0xda, 0xde, 0x14, 0x16, 0x00, 0x00, 0xd2, 0x05, 0x35, 0x30, 0x67, 0x16, 0xc8, 0xd2, 0x7a, 0x38, 0x76, 0x9d, 0x02, 0x34, 0x03, 0x22, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x05, 0xdc, 0x03, 0x10, 0x9b, 0xa8, 0x03, 0x10, 0x9b, 0xa8, 0x00, 0x50, 0xdb, 0x2c, 0x00, 0x10, 0x06, 0x00, 0xda, 0xde, 0x34, 0x17, 0x18, 0x17, 0x00, 0x00, 0x68, 0x10, 0x3c, 0x30, 0x72, 0x17, 0xfe, 0x48, 0x72, 0x17, 0x03, 0xe7, 0x02, 0x26, 0x03, 0x31, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x05, 0xa7, 0x03, 0x10, 0x63, 0x41, 0x03, 0x10, 0x63, 0x41, 0x00, 0x50, 0xdf, 0x2a, 0x00, 0x18, 0x06, 0x00, 0x34, 0x17, 0xda, 0xde, 0x14, 0x17, 0x00, 0x00, 0x6f, 0xa1, 0x40, 0x30, 0x72, 0x17, 0xf1, 0x30, 0x72, 0x17, 0x03, 0xe7, 0x03, 0x22, 0x03, 0x31, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x28, 0x03, 0x10, 0xb0, 0x67, 0x03, 0x10, 0xb0, 0x67, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x12, 0xe5, 0xda, 0xde, 0x0d, 0x18, 0x00, 0x00, 0x17, 0x34, 0x46, 0x30, 0x72, 0x17, 0xdf, 0x67, 0x72, 0x17, 0x03, 0xe7, 0x02, 0x26, 0x03, 0x31, 0x00, 0x00, 0x00, 0x0d, 0x00, 0x00, 0x4b, 0xc4, 0x03, 0x10, 0x67, 0x41, 0x03, 0x10, 0x6e, 0xe0, 0x01, 0xbb, 0x4a, 0x41, 0x00, 0x10, 0x06, 0x00, 0x51, 0xcc, 0xda, 0xde, 0x18, 0x16, 0x00, 0x00, 0x68, 0x10, 0x4f, 0x30, 0x72, 0x17, 0xe1, 0x2b, 0x72, 0x17, 0x03, 0xe7, 0x02, 0x26, 0x03, 0x31, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x28, 0x03, 0x10, 0x4e, 0x19, 0x03, 0x10, 0x4e, 0x19, 0x01, 0xbb, 0xd0, 0xb2, 0x00, 0x10, 0x06, 0x00, 0x34, 0x17, 0xda, 0xde, 0x14, 0x17, 0x00, 0x00, 0x72, 0x17, 0x63, 0x30, 0xcc, 0x5d, 0x8d, 0x7b, 0x7a, 0x38, 0x76, 0x9d, 0x02, 0x34, 0x03, 0x22, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x05, 0xd4, 0x03, 0x10, 0x6e, 0x57, 0x03, 0x10, 0x6e, 0x57, 0xf8, 0x23, 0x01, 0xbb, 0x00, 0x10, 0x06, 0x00, 0xda, 0xde, 0x5b, 0x38, 0x16, 0x11, 0x00, 0x00, 0x72, 0x17, 0x6d, 0x30, 0x9d, 0xf0, 0x08, 0x13, 0x7a, 0x38, 0x76, 0x9d, 0x02, 0x34, 0x03, 0x22, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x68, 0x03, 0x10, 0x45, 0x54, 0x03, 0x10, 0x8b, 0x9f, 0xbb, 0x26, 0x01, 0xbb, 0x00, 0x10, 0x06, 0x00, 0xda, 0xde, 0x80, 0xa6, 0x16, 0x18, 0x00, 0x00, 0x34, 0x6d, 0x70, 0x30, 0x72, 0x17, 0x1a, 0x05, 0x72, 0x17, 0x03, 0xfb, 0x02, 0x26, 0x02, 0x34, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x05, 0xdc, 0x03, 0x11, 0x02, 0x7e, 0x03, 0x11, 0x02, 0x7e, 0x01, 0xbb, 0xf7, 0xff, 0x00, 0x10, 0x06, 0x00, 0x1f, 0x8b, 0xda, 0xde, 0x0c, 0x1f, 0x00, 0x00, 0x34, 0x6d, 0x70, 0x30, 0x72, 0x17, 0xd8, 0x0e, 0x72, 0x17, 0x03, 0xe7, 0x02, 0x26, 0x03, 0x31, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x28, 0x03, 0x11, 0x0f, 0xdd, 0x03, 0x11, 0x0f, 0xdd, 0x01, 0xbb, 0xcb, 0xd5, 0x00, 0x10, 0x06, 0x00, 0x1f, 0x8b, 0xda, 0xde, 0x0c, 0x17, 0x00, 0x00, 0x34, 0x6d, 0x70, 0x30, 0x72, 0x17, 0xe9, 0x56, 0x72, 0x17, 0x03, 0xe7, 0x02, 0x26, 0x03, 0x31, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x05, 0xd4, 0x03, 0x10, 0xa8, 0x7a, 0x03, 0x10, 0xa8, 0x7a, 0x01, 0xbb, 0xfc, 0x8d, 0x00, 0x10, 0x06, 0x00, 0x1f, 0x8b, 0xda, 0xde, 0x0c, 0x16, 0x00, 0x00, 0x34, 0x6d, 0x70, 0x30, 0x72, 0x17, 0xf1, 0x6c, 0x72, 0x17, 0x03, 0xe7, 0x02, 0x26, 0x03, 0x31, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x93, 0x03, 0x10, 0x70, 0x27, 0x03, 0x10, 0x70, 0x27, 0x01, 0xbb, 0xca, 0xcc, 0x00, 0x18, 0x06, 0x00, 0x1f, 0x8b, 0xda, 0xde, 0x0c, 0x18, 0x00, 0x00, 0x34, 0x6d, 0x70, 0x30, 0x72, 0x17, 0x64, 0x79, 0x72, 0x17, 0x03, 0xfb, 0x02, 0x26, 0x02, 0x34, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x93, 0x03, 0x10, 0x68, 0x0d, 0x03, 0x10, 0x68, 0x0d, 0x01, 0xbb, 0xc8, 0x0b, 0x00, 0x18, 0x06, 0x00, 0x1f, 0x8b, 0xda, 0xde, 0x0c, 0x16, 0x00, 0x00, 0x72, 0x17, 0x79, 0x30, 0xb0, 0x09, 0x4a, 0x05, 0x7a, 0x38, 0x76, 0x9d, 0x03, 0x31, 0x03, 0x22, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0xba, 0x80, 0x03, 0x10, 0x3b, 0x89, 0x03, 0x11, 0x0f, 0x6f, 0xf0, 0xdc, 0xe6, 0x42, 0x00, 0x10, 0x06, 0x38, 0xda, 0xde, 0x61, 0x6c, 0x18, 0x10, 0x00, 0x00, 0x72, 0x17, 0x79, 0x30, 0x63, 0x49, 0xbf, 0xb2, 0x7a, 0x38, 0x76, 0x9d, 0x03, 0x31, 0x03, 0x22, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x2e, 0xa0, 0x03, 0x10, 0x3f, 0x35, 0x03, 0x11, 0x11, 0x3c, 0xc4, 0xf9, 0xe6, 0x42, 0x00, 0x10, 0x06, 0x38, 0xda, 0xde, 0x1b, 0x6a, 0x18, 0x0f, 0x00, 0x00, 0x72, 0x17, 0x79, 0x30, 0x56, 0x9e, 0xe3, 0xbb, 0x7a, 0x38, 0x76, 0x9d, 0x03, 0x31, 0x03, 0x22, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x05, 0xd4, 0x03, 0x11, 0x0d, 0xed, 0x03, 0x11, 0x0d, 0xed, 0xea, 0x28, 0x61, 0xe2, 0x00, 0x10, 0x06, 0x38, 0xda, 0xde, 0x0b, 0x28, 0x18, 0x0b, 0x00, 0x00, 0x72, 0x17, 0x7b, 0x30, 0x34, 0x5f, 0x83, 0x10, 0x7a, 0x38, 0x76, 0x9d, 0x03, 0x31, 0x03, 0x22, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x28, 0x03, 0x11, 0x05, 0x5b, 0x03, 0x11, 0x05, 0x5b, 0xf5, 0xb4, 0x01, 0xbb, 0x00, 0x10, 0x06, 0x00, 0xda, 0xde, 0x40, 0x7d, 0x18, 0x18, 0x00, 0x00, 0x72, 0x17, 0x8a, 0x30, 0x9d, 0xf0, 0x08, 0x13, 0x7a, 0x38, 0x76, 0x9d, 0x03, 0x31, 0x03, 0x22, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x34, 0x03, 0x10, 0xef, 0xda, 0x03, 0x10, 0xef, 0xda, 0xc4, 0x8a, 0x01, 0xbb, 0x00, 0x10, 0x06, 0x00, 0xda, 0xde, 0x80, 0xa6, 0x18, 0x18, 0x00, 0x00, 0x72, 0x17, 0x8a, 0x30, 0x9d, 0xf0, 0x08, 0x13, 0x7a, 0x38, 0x76, 0x9d, 0x03, 0x31, 0x03, 0x22, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x34, 0x03, 0x10, 0xd7, 0x97, 0x03, 0x10, 0xd7, 0x97, 0xea, 0x8a, 0x01, 0xbb, 0x00, 0x10, 0x06, 0x00, 0xda, 0xde, 0x80, 0xa6, 0x18, 0x18, 0x00, 0x00, 0x72, 0x17, 0x8e, 0x30, 0x34, 0x6d, 0x70, 0x2a, 0x2b, 0xf3, 0x15, 0x17, 0x02, 0xff, 0x02, 0x26, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x28, 0x03, 0x10, 0x74, 0x70, 0x03, 0x10, 0x74, 0x70, 0xe1, 0xc5, 0x01, 0xbb, 0x00, 0x10, 0x06, 0x00, 0xda, 0xde, 0x1f, 0x8b, 0x16, 0x0c, 0x00, 0x00, 0x72, 0x17, 0x8e, 0x30, 0x77, 0x09, 0x9a, 0x2d, 0x2b, 0xf3, 0x15, 0x1b, 0x02, 0xff, 0x02, 0x26, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x28, 0x03, 0x10, 0x5f, 0x14, 0x03, 0x10, 0x5f, 0x14, 0xe4, 0x64, 0x13, 0xe2, 0x00, 0x10, 0x06, 0x00, 0xda, 0xde, 0xe5, 0x3b, 0x16, 0x12, 0x00, 0x00, 0x72, 0x17, 0x8e, 0x30, 0x34, 0x72, 0x9e, 0x32, 0x2b, 0xf3, 0x15, 0x17, 0x02, 0xff, 0x02, 0x26, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x05, 0xc8, 0x03, 0x11, 0x14, 0xb4, 0x03, 0x11, 0x14, 0xb4, 0xc7, 0x31, 0x01, 0xbb, 0x00, 0x18, 0x06, 0x00, 0xda, 0xde, 0x1f, 0x8b, 0x16, 0x0e, 0x00, 0x00, 0x72, 0x17, 0x8e, 0x30, 0x23, 0xba, 0xc2, 0x3a, 0x7a, 0x38, 0x76, 0x9d, 0x02, 0xff, 0x03, 0x22, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x10, 0xa4, 0x03, 0x10, 0x57, 0x58, 0x03, 0x10, 0xa2, 0xf7, 0xc9, 0xa4, 0x01, 0xbb, 0x00, 0x10, 0x06, 0x00, 0xda, 0xde, 0x3b, 0x41, 0x16, 0x10, 0x00, 0x00, 0x72, 0x17, 0x8f, 0x30, 0x23, 0xbd, 0x11, 0x92, 0x7a, 0x38, 0x76, 0x9d, 0x02, 0xff, 0x03, 0x22, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x28, 0x03, 0x10, 0xf1, 0x3e, 0x03, 0x10, 0xf1, 0x3e, 0xe8, 0xf3, 0x01, 0xbb, 0x00, 0x10, 0x06, 0x00, 0xda, 0xde, 0x3b, 0x41, 0x16, 0x13, 0x00, 0x00, 0x72, 0x17, 0x8f, 0x30, 0x28, 0x64, 0x92, 0xb2, 0x2b, 0xf3, 0x15, 0x17, 0x02, 0xff, 0x02, 0x26, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x05, 0x78, 0x03, 0x10, 0x91, 0x3d, 0x03, 0x10, 0x91, 0x3d, 0xe4, 0x62, 0x01, 0xbb, 0x00, 0x10, 0x06, 0x00, 0xda, 0xde, 0x1f, 0x8b, 0x16, 0x0a, 0x00, 0x00, 0xd2, 0x37, 0x8f, 0x30, 0x6f, 0x41, 0xe6, 0x64, 0x72, 0x17, 0x03, 0xfb, 0x03, 0x17, 0x02, 0x34, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x05, 0xdc, 0x03, 0x10, 0xfe, 0x45, 0x03, 0x10, 0xfe, 0x45, 0x67, 0x2b, 0x00, 0x19, 0x00, 0x10, 0x06, 0x00, 0x12, 0x28, 0xda, 0xde, 0x18, 0x1b, 0x00, 0x00, 0x72, 0x17, 0x96, 0x30, 0x4a, 0x7d, 0x18, 0x6c, 0x7a, 0x38, 0x76, 0x9d, 0x03, 0x31, 0x03, 0x22, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x34, 0x03, 0x10, 0x3f, 0xf7, 0x03, 0x10, 0x3f, 0xf7, 0xf9, 0x48, 0x03, 0xe1, 0x00, 0x10, 0x06, 0x00, 0xda, 0xde, 0x3b, 0x41, 0x17, 0x18, 0x00, 0x00, 0x72, 0x17, 0x96, 0x30, 0x4a, 0x7d, 0x18, 0x6c, 0x7a, 0x38, 0x76, 0x9d, 0x03, 0x31, 0x03, 0x22, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x34, 0x03, 0x10, 0x3e, 0xa0, 0x03, 0x10, 0x3e, 0xa0, 0xf7, 0x56, 0x03, 0xe1, 0x00, 0x10, 0x06, 0x00, 0xda, 0xde, 0x3b, 0x41, 0x17, 0x18, 0x00, 0x00, 0x00, 0x00, 0xd4, 0x05, 0x00, 0x00}

func TestV5HeaderDecode(t *testing.T) {
	ip := net.ParseIP("114.23.3.231")
	rawBody := TestV5FlowPacket
	t.Log(fmt.Printf("Test Packet contains %v bytes\n", len(rawBody)))
	d := NewDecoder(ip, rawBody)
	if msg, err := d.Decode(); err != nil {
		t.Error(fmt.Printf("expected a message but got an error: %v\n", err))
	} else if msg == nil {
		t.Error("Expected a message but got nothing")
	} else {
		buf := new(bytes.Buffer)
		msg.JSONMarshal(buf)
		t.Log(fmt.Printf("Got a message \n%v\n", buf.String()))
	}
}
