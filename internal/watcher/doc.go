// Package watcher provides drift detection for environment variable snapshots.
//
// A Watcher periodically resolves the current environment state and compares
// it against a recorded baseline snapshot using the differ package. When
// changes are detected — added, removed, or modified keys — a user-supplied
// DriftHandler callback is invoked with the full Diff result.
//
// Example usage:
//
//	w, err := watcher.New(
//		30*time.Second,
//		baselineSnapshot,
//		resolveFn,
//		func(d *differ.Diff) {
//			log.Printf("drift detected: %d changes", len(d.Changes()))
//		},
//	)
//	if err != nil { ... }
//	w.Start()
//	defer w.Stop()
package watcher
