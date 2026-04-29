// Package loader provides utilities for loading environment variable
// definitions from .env files into key-value maps.
//
// A .env file is a plain-text file where each non-blank, non-comment line
// defines a single environment variable in KEY=VALUE format. Values may
// optionally be wrapped in single or double quotes, which are stripped
// during parsing.
//
// Example .env file:
//
//	# Database configuration
//	DB_HOST=localhost
//	DB_PORT=5432
//	DB_PASSWORD="s3cr3t"
//
// Usage:
//
//	l := loader.New(".env")
//	vars, err := l.Load()
//	if err != nil {
//		log.Fatal(err)
//	}
package loader
