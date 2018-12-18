package watch_append

import (
	"syscall"
	"go.uber.org/zap"
	"path/filepath"
)

type DiffResult struct {
	Speed     float64
	TotalSize int64
	Count     int64
}

func (dr *DiffResult) add(size int64, speed float64) {
	dr.Speed += speed
	dr.TotalSize += size
	dr.Count += 1
}

type Diff struct {
	// oldStates
	o States
	// actualStates
	a       States
	Pattern string
	Result  DiffResult
	Metric  *WatchMetric
}

func TimeRange(a, b syscall.Timespec) float64 {
	return float64(a.Sec-b.Sec) + float64(a.Nsec-b.Nsec)/1e9
}

func NewDiff(a, o States, pattern string, metric *WatchMetric) Diff {
	return Diff{
		a:       a,
		o:       o,
		Pattern: pattern,
		Result: DiffResult{
			Speed:     0,
			TotalSize: 0,
			Count:     0,
		},
		Metric: metric,
	}
}

func (d *Diff) diff(asf *State, osf *State) {

	// 1# if asf == nil, osf ==nil.
	// emm is a error ..
	if asf == nil && osf == nil {
		zap.S().Error("0# asf and osf is nil")
		d.Metric.DiffNullError.Add()
		return
	}

	// 2# NEW FILE
	// if file is new
	// means no rate
	if asf != nil && osf == nil {
		speed := float64(asf.Size) / TimeRange(d.a.RecordAt, d.o.RecordAt)
		d.Result.add(asf.Size, speed)
		d.Metric.DiffNewFile.Add()
		return
	}

	// 3# ROTATE
	// only once rotate can control
	if asf == nil && osf != nil {
		// todo will simple to check
		isOld := true
		for _, a := range d.a.States {
			if a.Source == osf.Source {
				isOld = false
				break
			}
		}

		if isOld {
			zap.S().Infow("Ignore old state,",
				"osf.source", osf.Source,
				"pattern", d.Pattern)
			d.Result.add(0, 0)
			d.Metric.DiffRotateIgnore.Add()
			return
		}

		state, err := FoundFileByInode(osf.Source, osf.INode)
		if err != nil {
			zap.S().Info(err.Error())
			d.Result.add(0, 0)
			d.Metric.DiffRotateNotExists.Add()
			return
		}
		size := state.Size - osf.Size
		speed := float64(size) / TimeRange(state.ModifyAt, osf.ModifyAt)
		d.Result.add(size, speed)
		d.Metric.DiffRotate.Add()
		return
	}

	mtRange := TimeRange(asf.ModifyAt, osf.ModifyAt)

	// 4# File No Change or File
	if (asf.Size-osf.Size) <= 0 || mtRange <= 0 {
		if (asf.Size - osf.Size) < 0 {
			zap.S().Infow("3# Size Reduce", "source", asf.Source, "asf_size", asf.Size, "osf_size", osf.Size)
		}
		if mtRange < 0 {
			zap.S().Infow("3# mtRange <0", "source", asf.Source, "omt", asf.ModifyAt, "osf_size", osf.ModifyAt)
		}
		d.Result.add(0, 0)
		d.Metric.DiffNoChange.Add()
		return
	}

	// 5# File Append
	speed := float64(asf.Size-osf.Size) / mtRange
	d.Result.add(asf.Size-osf.Size, speed)
	d.Metric.DiffAppend.Add()
	return
}

func (d *Diff) getRotateAppendSize(asf, osf State) int64 {
	files, _ := filepath.Glob(asf.Source + "*")
	for _, path := range files {
		var stat syscall.Stat_t
		if err := syscall.Stat(path, &stat); err != nil {
			zap.S().Infow("	scan get file stat failed",
				"file_path:", path,
				"err", err,
			)
			continue
		}
		if stat.Ino == osf.INode {
			return asf.Size + stat.Size - osf.Size
		}
	}
	zap.S().Infow("2# Can't found inode",
		"osf.inode", osf.INode,
		"osf.size", osf.Size,
		"osf.mtime", osf.ModifyAt,
		"asf.mtime", asf.ModifyAt,
	)
	return asf.Size
}

func (d *Diff) Diff() {
	// actualStates Set
	asSet := map[uint64]bool{}
	for k := range d.a.States {
		asSet[k] = true
	}

	// oldStates Set
	osSet := map[uint64]bool{}
	for k := range d.o.States {
		osSet[k] = true
	}

	asSet_union_osSet := map[uint64]bool{}
	// in actualStatesoldStates Set, not in oldStates Set
	in_asSet_not_osSet := map[uint64]bool{}
	not_asSet_in_osSet := map[uint64]bool{}
	for k := range asSet {
		if _, ok := osSet[k]; ok {
			asSet_union_osSet[k] = true
		} else {
			in_asSet_not_osSet[k] = true
		}
	}
	for k := range osSet {
		if _, ok := asSet[k]; !ok {
			not_asSet_in_osSet[k] = true
		}
	}

	d.Metric.AsCount = len(asSet)
	d.Metric.OsCount = len(osSet)
	d.Metric.AsUnionOs = len(asSet_union_osSet)

	for k := range asSet_union_osSet {
		osf := d.o.States[k]
		asf := d.a.States[k]
		d.diff(&asf, &osf)
	}

	for k := range in_asSet_not_osSet {
		asf := d.a.States[k]
		d.diff(&asf, nil)
	}

	for k := range not_asSet_in_osSet {
		osf := d.o.States[k]
		d.diff(nil, &osf)
	}

}
