package snowflake

import (
	"sync"
	"time"
)

/*
关于节点位长度取值问题(不借秒的情况下)
节点位  支持节点数  每节点最大秒并发
 0     1         2097152
 1     2         1048576
 2     4         524288
 3     8         262144
 4     16        131072
 5     32        65536
 6     64        32768
 7     128       16384
 8     256       8192
 9     512       4096
10     1024      2048
*/

//32位秒 + 17位自增 + 4位机器 = 53位
const (
	timeShift     = 21
	nodeBits      = 5
	sequenceBits  = timeShift - nodeBits
	sequenceShift = nodeBits

	nodeMask     = -1 ^ (-1 << nodeBits)
	sequenceMask = -1 ^ (-1 << sequenceBits)

	timeFormat = "20060102150405"

	//2018-01-04T01:20:00+08:00 //为何取这个时间，想在2018年初找一个好记的时间戳, 这个点挺好，1515+6个0
	epoch = 1515000000 //纪元时间（序号时间起始）
)

func New(node uint) *Snowflake {
	return (&Snowflake{}).Set(node)
}

type Snowflake struct {
	node uint64 //节点

	time     uint64 //时间
	sequence uint64 //序列号

	mu sync.Mutex
}

func (s *Snowflake) Set(node uint) *Snowflake {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.node = uint64(node & nodeMask) //如果超限设为0
	return s
}

func (s *Snowflake) Generate() ID {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := uint64(time.Now().Unix())
	seq := s.sequence
	if now <= s.time {
		if now < s.time {
			now = s.time
		}
		//达到最大值，借秒
		if seq = (seq + 1) & sequenceMask; seq == 0 {
			now++
		}
	} else {
		seq = 0
	}
	s.time = now
	s.sequence = seq

	return ID{
		uts:  now,
		seq:  seq,
		node: s.node,
	}
}
