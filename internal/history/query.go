package history

import "time"

// Filter holds optional criteria for querying history records.
type Filter struct {
	JobName string
	Since   time.Time
	OnlyFailed bool
	Limit   int
}

// Query returns records matching the given filter criteria.
// Results are returned in reverse-chronological order (newest first).
func (s *Store) Query(f Filter) []Record {
	s.mu.Lock()
	all := make([]Record, len(s.records))
	copy(all, s.records)
	s.mu.Unlock()

	// Reverse for newest-first ordering.
	for i, j := 0, len(all)-1; i < j; i, j = i+1, j-1 {
		all[i], all[j] = all[j], all[i]
	}

	var out []Record
	for _, r := range all {
		if f.JobName != "" && r.JobName != f.JobName {
			continue
		}
		if !f.Since.IsZero() && r.StartedAt.Before(f.Since) {
			continue
		}
		if f.OnlyFailed && r.Success {
			continue
		}
		out = append(out, r)
		if f.Limit > 0 && len(out) >= f.Limit {
			break
		}
	}
	return out
}

// LastRun returns the most recent record for the given job name, or nil.
func (s *Store) LastRun(jobName string) *Record {
	results := s.Query(Filter{JobName: jobName, Limit: 1})
	if len(results) == 0 {
		return nil
	}
	return &results[0]
}
