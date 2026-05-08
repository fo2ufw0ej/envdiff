// Package snapshot provides functionality for saving and loading comparison
// result snapshots to disk, and comparing two snapshots to detect changes
// in key statuses over time.
//
// A snapshot captures the full set of comparator.Result entries along with
// a timestamp and an optional label, serialised as JSON.
//
// Use Compare to produce a list of Delta values describing which keys changed
// status (or were newly added) between an older and a newer snapshot.
package snapshot
