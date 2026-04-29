// Package exporter provides functionality for serializing environment variable
// sets into various output formats suitable for shell consumption or
// configuration management tooling.
//
// Supported formats:
//
//   - dotenv  — KEY=VALUE pairs, one per line (compatible with .env files)
//   - export  — Shell-ready `export KEY=VALUE` statements
//   - json    — A flat JSON object mapping keys to string values
//
// Example usage:
//
//	e, err := exporter.New(exporter.FormatDotenv, os.Stdout)
//	if err != nil {
//		log.Fatal(err)
//	}
//	e.Write(map[string]string{"APP_ENV": "production", "PORT": "8080"})
package exporter
