// Package watcher provides drift detection for environment variable snapshots.
//
// A Watcher periodically resolves the current environment state and compares
// it against a recorded baseline snapshot using the differ package. When
// changes are detected — added, removed, or modified keys — a user-supplied
// DriftHandler callback is invoked with the full Diff result.
//
// Lifecycle
//
// Call Start to begin periodic polling. The watcher runs in a background
// goroutine and the caller retains full control over shutdown via Stop, which
// blocks until the background goroutine has exited. It is safe to call Stop
// multiple times.
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
