// Package comparator provides side-by-side comparison of two environment
// variable maps, identifying keys that match, conflict, or exist exclusively
// in one of the two environments.
//
// It is useful for auditing configuration drift between deployment contexts
// such as staging vs. production, or local vs. CI.
//
// Basic usage:
//
//	c, err := comparator.New("staging", "production")
//	if err != nil { ... }
//
//	result, err := c.Compare(stagingEnv, productionEnv)
//	if err != nil { ... }
//
//	fmt.Println("Conflicts:", result.Conflicts)
//	fmt.Println("Only in staging:", result.LeftOnly)
//	fmt.Println("Only in production:", result.RightOnly)
package comparator
