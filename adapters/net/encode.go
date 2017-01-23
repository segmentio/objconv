package net

import (
	"net"
	"reflect"

	"github.com/segmentio/objconv"
)

func encodeTCPAddr(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.TCPAddr)
	return e.Emitter.EmitString(a.String())
}

func encodeUDPAddr(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.UDPAddr)
	return e.Emitter.EmitString(a.String())
}

func encodeUnixAddr(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.UnixAddr)
	return e.Emitter.EmitString(a.String())
}

func encodeIPAddr(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.IPAddr)
	return e.Emitter.EmitString(a.String())
}

func encodeIP(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.IP)
	return e.Emitter.EmitString(a.String())
}
