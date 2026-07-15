// Package json implements a small, streaming JSON reader and writer.
//
// Unlike Go's encoding/json, So has no reflection, so there is no
// Marshal/Unmarshal over arbitrary structs. Instead this package exposes a
// token-level API: a [Decoder] that pulls one validated token at a time, and
// an [Encoder] that writes tokens to an [io.Writer] while inserting commas and
// colons automatically.
//
// The Decoder reads either a complete in-memory document ([NewDecoder]) or a
// stream pulled from an [io.Reader] ([NewReader]). The token API is the same
// for both. The Encoder buffers its output, so call [Encoder.Flush] when the
// document is complete.
//
// Both the Decoder and Encoder reject invalid JSON syntax and non-UTF-8 strings.
//
// # Limitations
//
// A document holds a single root value, and a Decoder reads exactly one
// document. [Decoder.Next] fails with [ErrSyntax] on anything but whitespace
// past the root value.
//
// A streaming Decoder reads ahead in chunks, so the bytes that follow the
// document end up in its buffer, and no method gives them back. A second
// Decoder over the same reader does not help: those bytes are already read.
// Formats that put several values on one stream, such as newline-delimited
// JSON, are not supported.
//
// Reading to the end (for d.Next() {}) makes the Decoder read one byte past the
// document to check that nothing follows the root value. On a reader that stays
// open with nothing left to send, such as a socket, that read blocks.
package json

import "solod.dev/so/errors"

// MaxDepth is the deepest nesting of objects and arrays the Encoder and Decoder
// accept. Going deeper yields [ErrDepth]. The decoder tracks nesting in a fixed
// array rather than by recursion, so this is a structural limit, not a guard
// against stack overflow.
const MaxDepth = 128

var (
	ErrDepth     = errors.New("json: nesting too deep")
	ErrKind      = errors.New("json: wrong token kind")
	ErrNonFinite = errors.New("json: non-finite number")
	ErrSyntax    = errors.New("json: invalid syntax")
	ErrTooLong   = errors.New("json: token too long")
	ErrValue     = errors.New("json: invalid token value")
)
