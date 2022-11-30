package storage

import (
	"time"
)

// SetDedupInterval sets the deduplication interval, which is applied to raw samples during data ingestion and querying.
//
// De-duplication is disabled if dedupInterval is 0.
//
// This function must be called before initializing the storage.
func SetDedupInterval(dedupInterval time.Duration) {
	globalDedupInterval = dedupInterval.Milliseconds()
}

// GetDedupInterval returns the dedup interval in milliseconds, which has been set via SetDedupInterval.
func GetDedupInterval() int64 {
	return globalDedupInterval
}

var globalDedupInterval int64

func isDedupEnabled() bool {
	return globalDedupInterval > 0
}

// DeduplicateSamples removes samples from src* if they are closer to each other than dedupInterval in millseconds.
func DeduplicateSamples(srcTimestamps []int64, srcValues []float64, dedupInterval int64) ([]int64, []float64) {
	if !needsDedup(srcTimestamps, dedupInterval) {
		// Fast path - nothing to deduplicate
		return srcTimestamps, srcValues
	}
	tsNext := srcTimestamps[0] + dedupInterval - 1
	tsNext -= tsNext % dedupInterval
	dstTimestamps := srcTimestamps[:0]
	dstValues := srcValues[:0]

	var prevTs int64
	var value, prevValue float64
	for i, ts := range srcTimestamps[1:] {
		value = srcValues[i]
		if i > 0 && srcTimestamps[i] == srcTimestamps[i-1] && value < prevValue {
			// prefer biggest value on timestamp conflict
			value = prevValue
		}
		prevValue = value
		prevTs = srcTimestamps[i]
		if ts <= tsNext {
			continue
		}
		dstTimestamps = append(dstTimestamps, srcTimestamps[i])
		dstValues = append(dstValues, value)
		tsNext += dedupInterval
		if tsNext < ts {
			tsNext = ts + dedupInterval - 1
			tsNext -= tsNext % dedupInterval
		}
	}

	ts := srcTimestamps[len(srcTimestamps)-1]
	v := srcValues[len(srcValues)-1]
	dstTimestamps = append(dstTimestamps, ts)
	if ts == prevTs && v < value {
		v = value
	}
	dstValues = append(dstValues, v)
	return dstTimestamps, dstValues
}

func deduplicateSamplesDuringMerge(srcTimestamps, srcValues []int64, dedupInterval int64) ([]int64, []int64) {
	if !needsDedup(srcTimestamps, dedupInterval) {
		// Fast path - nothing to deduplicate
		return srcTimestamps, srcValues
	}
	tsNext := srcTimestamps[0] + dedupInterval - 1
	tsNext -= tsNext % dedupInterval
	dstTimestamps := srcTimestamps[:0]
	dstValues := srcValues[:0]

	var prevTs int64
	var value, prevValue int64
	for i, ts := range srcTimestamps[1:] {
		value = srcValues[i]
		if i > 0 && srcTimestamps[i] == srcTimestamps[i-1] && value < prevValue {
			// prefer biggest value on timestamp conflict
			value = prevValue
		}
		prevValue = value
		if ts <= tsNext {
			continue
		}
		dstTimestamps = append(dstTimestamps, srcTimestamps[i])
		dstValues = append(dstValues, value)
		tsNext += dedupInterval
		if tsNext < ts {
			tsNext = ts + dedupInterval - 1
			tsNext -= tsNext % dedupInterval
		}
	}
	ts := srcTimestamps[len(srcTimestamps)-1]
	v := srcValues[len(srcValues)-1]
	dstTimestamps = append(dstTimestamps, ts)
	if ts == prevTs && v < value {
		v = value
	}
	dstValues = append(dstValues, v)
	return dstTimestamps, dstValues
}

func needsDedup(timestamps []int64, dedupInterval int64) bool {
	if len(timestamps) < 2 || dedupInterval <= 0 {
		return false
	}
	tsNext := timestamps[0] + dedupInterval - 1
	tsNext -= tsNext % dedupInterval
	for _, ts := range timestamps[1:] {
		if ts <= tsNext {
			return true
		}
		tsNext += dedupInterval
		if tsNext < ts {
			tsNext = ts + dedupInterval - 1
			tsNext -= tsNext % dedupInterval
		}
	}
	return false
}
