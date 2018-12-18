package watch_append

import "time"

type WatchMetric struct {
	GlobMatchCount int
	ExcludeCount   int
	AsCount        int
	OsCount        int
	AsUnionOs      int
	DeepFindCount  int

	StartAt int64
	Spent   int64

	DiffNewFile  MetricInt
	DiffRotate   MetricInt
	DiffNoChange MetricInt
	DiffAppend   MetricInt


	DiffNullError       MetricInt
	DiffRotateIgnore    MetricInt
	DiffRotateDeepIgnore    MetricInt
	DiffRotateNotExists MetricInt


}

func NewWatchMetric() WatchMetric {
	return WatchMetric{
		GlobMatchCount: 0,
		ExcludeCount:   0,
	}
}

func (w *WatchMetric) Start() {
	w.StartAt = time.Now().UnixNano()
}

func (w *WatchMetric) End() {
	w.Spent = time.Now().UnixNano() - w.StartAt
}

type MetricInt struct {
	count int
}

func (m *MetricInt) Add() {
	m.count += 1
}
func (m *MetricInt) Get() int {
	return m.count
}
