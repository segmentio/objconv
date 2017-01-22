package net

import (
	"net"
	"reflect"

	"github.com/segmentio/objconv"
)

func init() {
	objconv.Install(reflect.TypeOf(net.TCPAddr{}), objconv.Adapter{
		Encode: EncodeTCPAddr,
		Decode: DecodeTCPAddr,
	})

	objconv.Install(reflect.TypeOf(net.UDPAddr{}), objconv.Adapter{
		Encode: EncodeUDPAddr,
		Decode: DecodeUDPAddr,
	})

	objconv.Install(reflect.TypeOf(net.UnixAddr{}), objconv.Adapter{
		Encode: EncodeUnixAddr,
		Decode: DecodeUnixAddr,
	})

	objconv.Install(reflect.TypeOf(net.IPAddr{}), objconv.Adapter{
		Encode: EncodeIPAddr,
		Decode: DecodeIPAddr,
	})

	objconv.Install(reflect.TypeOf(net.IP(nil)), objconv.Adapter{
		Encode: EncodeIP,
		Decode: DecodeIP,
	})
}
