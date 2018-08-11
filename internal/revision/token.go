/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package revision

// token represents a entity extracted from string parsing
type token int

const (
	eof token = iota

	aslash
	asterisk
	at
	caret
	cbrace
	colon
	control
	dot
	emark
	minus
	number
	obrace
	obracket
	qmark
	slash
	space
	tilde
	tokenError
	word
)
