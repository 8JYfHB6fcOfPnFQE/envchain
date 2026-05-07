// Package freezer provides point-in-time immutable captures of environment
// variable maps. A Freezer stores labelled snapshots and returns defensive
// copies on retrieval, ensuring that downstream consumers cannot accidentally
// modify a frozen frame.
//
// Typical usage:
//
//	f := freezer.New()
//
//	// Capture the resolved environment before deployment.
//	if err := f.Freeze("pre-deploy", resolvedEnv); err != nil {
//		log.Fatal(err)
//	}
//
//	// Later, retrieve the frozen state for auditing or rollback comparison.
//	env, err := f.Thaw("pre-deploy")
//	if err != nil {
//		log.Fatal(err)
//	}
package freezer
