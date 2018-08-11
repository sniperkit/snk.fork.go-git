/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package storer

// Storer is a basic storer for encoded objects and references.
type Storer interface {
	EncodedObjectStorer
	ReferenceStorer
}

// Initializer should be implemented by storers that require to perform any
// operation when creating a new repository (i.e. git init).
type Initializer interface {
	// Init performs initialization of the storer and returns the error, if
	// any.
	Init() error
}
