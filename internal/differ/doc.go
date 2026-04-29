// Package differ compares two environment variable maps and produces a
// structured diff describing which keys were added, removed, or changed.
//
// Usage:
//
//	d := differ.New()
//	result, err := d.Compare(baseEnv, nextEnv)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, entry := range result.Entries() {
//		fmt.Printf("%s %s\n", entry.Kind, entry.Key)
//	}
//
// DiffKind values are: Added, Removed, Changed.
// Entries are always returned in lexicographic key order for deterministic output.
package differ
