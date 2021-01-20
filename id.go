package snowflake

import (
	"fmt"
	"strconv"
	"time"
)

var locCST = time.FixedZone("CST", 28800)

type ID struct {
	uts  uint64
	seq  uint64
	node uint64
}

func (s ID) IsZero() bool {
	return s.uts <= 0
}

func (s ID) Time() time.Time {
	if s.uts > 0 {
		return time.Unix(int64(s.uts), 0)
	}
	return time.Time{}
}

func (s ID) Node() uint {
	return uint(s.node)
}

func (s ID) Seq() uint64 {
	return s.seq
}

func (s ID) Uint() uint64 {
	return (s.uts-epoch)<<timeShift | s.seq<<sequenceShift | s.node
}

func (s ID) Int() int64 {
	return int64(s.Uint())
}

func (s ID) String() string {
	return s.FormatHuman(nil)
}

func (s ID) FormatHuman(loc *time.Location) string {
	if loc == nil {
		loc = locCST
	}
	return fmt.Sprintf("%s%05d%02d",
		time.Unix(int64(s.uts), 0).In(loc).Format(timeFormat),
		s.seq,
		s.node)
}

func (s ID) FormatRadix(radix int) string {
	if radix < 2 || radix > 36 {
		radix = 36
	}
	return strconv.FormatUint(s.Uint(), radix)
}

func FromInt(iid int64) ID {
	return FromUint(uint64(iid))
}

func FromUint(iid uint64) ID {
	var id ID
	id.uts = iid>>timeShift + epoch
	id.seq = iid>>sequenceShift - iid>>timeShift<<sequenceBits
	id.node = iid - iid>>sequenceShift<<sequenceShift
	return id
}

func ParseRadix(sid string, radix int) ID {
	if radix < 2 {
		radix = 36
	}
	i, _ := strconv.ParseUint(sid, radix, 64)
	return FromUint(i)
}

func ParseHuman(hid string, loc *time.Location) ID {
	if loc == nil {
		loc = locCST
	}

	var id ID
	tim, _ := time.ParseInLocation(timeFormat, hid[:14], loc)
	id.uts = uint64(tim.Unix())
	id.seq, _ = strconv.ParseUint(hid[14:19], 10, 64)
	id.node, _ = strconv.ParseUint(hid[19:], 10, 64)
	return id
}
