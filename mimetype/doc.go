// Package mimetype provides a global registry that maps mime types to encoders
// and decoders.
//
// Packages providing implementation of parsers and emitters for specific
// formats should register their codecs to the mimetype package within an init
// function, so programs loading these packages also get the ability to use the
// mimetype package as a side effect.
package mimetype
