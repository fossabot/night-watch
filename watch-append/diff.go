package watch_append

import (
	"syscall"
	"go.uber.org/zap"
	"path/filepath"
)

type DiffResult struct {
	Speed float64
	TotalSize int64
	Count  int64
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
	a      States
	Result DiffResult
}

func TimeRange(a, b syscall.Timespec) float64 {
	return float64(a.Sec-b.Sec) + float64(a.Nsec-b.Nsec)/1e9
}




func NewDiff(a,o States) Diff{
	return Diff{
		a:a,
		o:o,
		Result:DiffResult{
			Speed: 0,
			TotalSize: 0,
			Count: 0,
		},
	}
}

func (d *Diff) diff(asf State) {
	d.getRate(asf)
	return
}

// Deprecated
func (d *Diff) getRate(asf State) float64 {
	osf, ok := d.o.States[asf.INode]

	// 1# NEW FILE
	// if file is new
	// means no rate
	if !ok {
		zap.S().Infow("asf is new file", "source", asf.Source, "size", asf.Size)
		speed := float64(asf.Size) / TimeRange(d.a.RecordAt, d.o.RecordAt)
		d.Result.add(asf.Size, speed)
		return speed
	}

	// modify change
	mtRange := TimeRange(asf.ModifyAt, osf.ModifyAt)

	// 2# ROTATE
	// only once rotate can control
	if asf.INode != osf.INode {
		size := d.getRotateAppendSize(asf, osf)
		speed := float64(size) / mtRange
		d.Result.add(size, speed)
		return speed
	}

	// 3# File No Change or File
	if (asf.Size-osf.Size) <= 0 || mtRange <= 0 {
		if (asf.Size - osf.Size) < 0 {
			zap.S().Infow("3# Size Reduce", "source", asf.Source, "asf_size", asf.Size, "osf_size", osf.Size)
		}
		if mtRange < 0 {
			zap.S().Infow("3# mtRange <0", "source", asf.Source, "omt", asf.ModifyAt, "osf_size", osf.ModifyAt)
		}
		d.Result.add(0,0)
		return float64(0)
	}


	// 4# File Append
	speed := float64(asf.Size - osf.Size) / mtRange
	d.Result.add(asf.Size - osf.Size, speed)
	return speed
}

func (d *Diff)  getRotateAppendSize(asf, osf State) int64{
	files, _ := filepath.Glob(asf.Source + "*")
	for _, path := range  files {
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
	for _, asf := range d.a.States {
		d.diff(asf)
	}
}
