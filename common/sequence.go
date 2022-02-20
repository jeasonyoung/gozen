package common

import (
	"errors"
	"sync"
	"time"
)

// Sequence 序号接口
type Sequence interface {
	// NextId 生成序列号
	NextId() (uint64, error)
}

const (
	CEpoch        = 1474802888000
	CWorkerIdBits = 10 //Num of WorkerId Bits
	//CSenquenceBits = 12 //Num of Sequece Bits

	CWorkerIdShift  = 12
	CTimeStampShift = 22

	CSequenceMask = 0xfff //equal as getSequenceMask()
	CMaxWorker    = 0x3ff //equal as getMaxWorkerId()
)

// SnowFlake 雪花数算法
type SnowFlake struct {
	//机器ID(二进制五位)
	workerId int64
	//最新时间戳
	lastTimeStamp int64
	//代表1毫秒内生成多个id的最新序号 2进制12位, 4095个
	sequence int64
	//机房ID(二进制五位)
	maxWorkerId int64
	//同步锁
	lock *sync.Mutex
}

// NewSnowFlake 雪花数算法实例
func NewSnowFlake(workerId uint64) (Sequence, error) {
	sf := new(SnowFlake)
	sf.maxWorkerId = -1 ^ -1<<CWorkerIdBits
	if int64(workerId) > sf.maxWorkerId || workerId < 0 {
		return nil, errors.New("worker not fit")
	}
	sf.workerId = int64(workerId)
	sf.lastTimeStamp = -1
	sf.sequence = 0
	sf.lock = new(sync.Mutex)
	return sf, nil
}

func (sf *SnowFlake) timeGen() int64 {
	return time.Now().UnixNano() / 1000 / 1000
}

func (sf *SnowFlake) timeReGen(last int64) int64 {
	ts := time.Now().UnixNano() / 1000 / 1000
	for {
		if ts <= last {
			ts = sf.timeGen()
		} else {
			break
		}
	}
	return ts
}

func (sf *SnowFlake) NextId() (uint64, error) {
	sf.lock.Lock()
	defer sf.lock.Unlock()
	ts := sf.timeGen()
	if ts == sf.lastTimeStamp {
		sf.sequence = (sf.sequence + 1) & CSequenceMask
		if sf.sequence == 0 {
			ts = sf.timeReGen(ts)
		}
	} else {
		sf.sequence = 0
	}
	if ts < sf.lastTimeStamp {
		return 0, errors.New("clock moved backwards,Refuse gen id")
	}
	sf.lastTimeStamp = ts
	ts = (ts-CEpoch)<<CTimeStampShift | sf.workerId<<CWorkerIdShift | sf.sequence
	return uint64(ts), nil
}

func ParseSnowFlake(id uint64) (t time.Time, ts int64, workerId int64, seq int64) {
	seq = int64(id) & CSequenceMask
	workerId = (int64(id) >> CWorkerIdShift) & CMaxWorker
	ts = (int64(id) >> CTimeStampShift) + CEpoch
	t = time.Unix(ts/1000, (ts%1000)*1000000)
	return
}
