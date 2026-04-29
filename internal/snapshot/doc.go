// Package snapshot provides point-in-time capture and comparison of
// environment variable sets within envchain.
//
// A Snapshot records the state of a resolved environment at a specific
// moment, allowing operators to detect configuration drift between
// deployments or across chain resolution runs.
//
// Basic usage:
//
//	// Capture two states
//	before, err := snapshot.New("before-deploy", resolvedBefore)
//	after, err := snapshot.New("after-deploy", resolvedAfter)
//
//	// Compare them
//	diff := before.Compare(after)
//	if diff.HasChanges() {
//		// inspect diff.Added, diff.Removed, diff.Changed
//	}
package snapshot
