package net

import (
	"net"
	"reflect"

	"github.com/segmentio/objconv"
)

// EncodeTCPAddr encodes the net.TCPAddr value in v as a string using e.
func EncodeTCPAddr(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.TCPAddr)
	return e.Emitter.EmitString(a.String())
}

// EncodeUDPAddr encodes the net.UDPAddr value in v as a string using e.
func EncodeUDPAddr(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.UDPAddr)
	return e.Emitter.EmitString(a.String())
}

// EncodeUnixAddr encodes the net.UnixAddr value in v as a string using e.
func EncodeUnixAddr(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.UnixAddr)
	return e.Emitter.EmitString(a.String())
}

// EncodeIPAddr encodes the net.IPAddr value in v as a string using e.
func EncodeIPAddr(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.IPAddr)
	return e.Emitter.EmitString(a.String())
}

// EncodeIP encodes the net.IP value in v as a string using e.
func EncodeIP(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(net.IP)
	return e.Emitter.EmitString(a.String())
}
