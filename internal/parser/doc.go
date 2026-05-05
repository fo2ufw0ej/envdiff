// Package parser provides utilities for reading and parsing .env files
// into key-value maps.
//
// Supported syntax:
//
//	# This is a comment
//	KEY=value
//	KEY="quoted value"
//	KEY='single quoted'
//
// Blank lines and lines beginning with '#' are ignored.
// Keys must be non-empty strings. Values may be optionally wrapped
// in single or double quotes, which will be stripped.
//
// Example usage:
//
//	env, err := parser.ParseFile(".env.production")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(env["DATABASE_URL"])
package parser
