package main

// File-level constants.
const fInt int = 42
const fString string = "file"

// Using _ on file level is not supported.
// var _ = fInt
// var _ = fString

// Typedefed constant group.
type HttpStatus int

const (
	StatusOK       HttpStatus = 200
	StatusNotFound HttpStatus = 404
	StatusError    HttpStatus = 500
	statusSecret   HttpStatus = 999
)

// Regular constant group.
type ServerState string

const (
	StateIdle      ServerState = "idle"
	StateConnected ServerState = "connected"
	StateError     ServerState = "error"
)

// Iota is not supported.
// const (
// 	Sunday = iota
// 	Monday
// 	Tuesday
// )

func main() {
	{
		// Local constants.
		const lInt = 500000000
		_ = lInt
		const lFloat = 3e20 / lInt
		_ = lFloat
		const lString = "local"
		_ = lString
	}
	{
		// Using constants in expressions.
		status := StatusOK
		_ = status != StatusNotFound

		secret := statusSecret
		_ = secret > StatusOK

		state := StateConnected
		_ = state == StateIdle
	}
	{
		// Using _ on file level is not supported,
		// so silence the unused file-level constants here.
		_ = fInt
		_ = fString
	}
}
