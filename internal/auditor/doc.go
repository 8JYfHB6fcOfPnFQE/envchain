// Package auditor provides audit logging for resolved environment snapshots
// within envchain. It records which context was resolved, which keys were
// present, and an optional human-readable note for each resolution event.
//
// Audit entries are stored in-memory and can be retrieved or cleared at any
// time. The Auditor is safe to use within a single goroutine; external
// synchronisation is required for concurrent access.
//
// Example usage:
//
//	a := auditor.New()
//	err := a.Record("production", snap, "deployment triggered")
//	for _, entry := range a.Entries() {
//		fmt.Println(entry.Timestamp, entry.ContextName, entry.Note)
//	}
package auditor
