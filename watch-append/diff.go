package watch_append

import (
	"syscall"
	"go.uber.org/zap"
	"path/filepath"
)


type Diff struct {
	// oldStates
	o States
	// actualStates
	a States
}

func TimeRange(a, b syscall.Timespec) float64 {
	return float64(a.Sec-b.Sec) + float64(a.Nsec-b.Nsec)/1e9
}




func NewDiff(a,o States) Diff{
	return Diff{
		a:a,
		o:o,
	}
}



func (d *Diff) getRate(asf State) float64 {
	osf, ok := d.o.States[asf.INode]

	// 1# NEW FILE
	// if file is new
	// means no rate
	if !ok {
		zap.S().Infow("asf is new file", "source", asf.Source, "size", asf.Size)
		return float64(asf.Size) / TimeRange(d.a.RecordAt, d.o.RecordAt)
	}

	// modify change
	mtRange := TimeRange(asf.ModifyAt, osf.ModifyAt)

	// 2# ROTATE
	// only once rotate can control
	if asf.INode != osf.INode {
		size := d.getRotateAppendSize(asf, osf)
		return size / mtRange
	}

	// 3# File No Change or File
	if (asf.Size-osf.Size) <= 0 || mtRange <= 0 {
		if (asf.Size - osf.Size) < 0 {
			zap.S().Infow("3# Size Reduce", "source", asf.Source, "asf_size", asf.Size, "osf_size", osf.Size)
		}
		if mtRange < 0 {
			zap.S().Infow("3# mtRange <0", "source", asf.Source, "omt", asf.ModifyAt, "osf_size", osf.ModifyAt)
		}
		return float64(0)
	}


	// 4# File Append
	return float64(asf.Size - osf.Size) / mtRange
}

func (d *Diff)  getRotateAppendSize(asf, osf State) float64{
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
			return float64(asf.Size + stat.Size - osf.Size)
		}
	}
	zap.S().Infow("2# Can't found inode",
		"osf.inode", osf.INode,
		"osf.size", osf.Size,
		"osf.mtime", osf.ModifyAt,
		"asf.mtime", asf.ModifyAt,
	)
	return float64(asf.Size)
}

func (d *Diff) GetTotalRate() float64{
	totalRate := float64(0)

	for _, asf := range d.a.States {
		totalRate = totalRate + d.getRate(asf)
	}
	return totalRate
}
