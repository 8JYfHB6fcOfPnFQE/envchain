// Package linter provides a rule-based linting engine for environment variable
// maps. It allows callers to define custom rules or use the built-in
// DefaultRules to detect common issues such as empty values, lowercase keys,
// and values with extraneous whitespace.
//
// Usage:
//
//	l, err := linter.New(linter.DefaultRules())
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	issues, err := l.Lint(env)
//	for _, issue := range issues {
//		fmt.Println(issue)
//	}
//
// Custom rules can be appended to DefaultRules or supplied independently.
// Each Rule must have a unique non-blank Name and a non-nil Check function
// that returns true when the rule is violated for a given key/value pair.
package linter
