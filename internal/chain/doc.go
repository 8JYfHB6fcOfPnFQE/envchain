// Package chain implements ordered chaining of environment variable contexts
// for the envchain tool.
//
// A Chain holds an ordered list of context names. When Resolve is called,
// it iterates through each context in order, loading its associated EnvSet
// and merging the results. Contexts later in the chain take precedence over
// earlier ones when the same key appears in multiple sets.
//
// Typical usage:
//
//	reg := context.NewRegistry()
//	// ... register contexts ...
//
//	c := chain.New(reg)
//	c.Add("base")
//	c.Add("production")
//
//	resolved, err := c.Resolve()
//	if err != nil {
//		log.Fatal(err)
//	}
package chain
