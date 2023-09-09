package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	sequenceBits  uint8 = 12
	workerBits    uint8 = 10
	timeStampBits uint8 = 41

	MaxSequence int64 = -1 ^ (-1 << sequenceBits)
	MaxWorkerId int64 = -1 ^ (-1 << workerBits)

	workerShiftBits = sequenceBits
	timeShiftBits   = sequenceBits + workerBits

	customEpoch int64 = 1609459200000 // Timestamp for 1st January 2021 in milliseconds
)

type Snowflake struct {
	lock          sync.Mutex
	workerId      int64
	sequence      int64
	lastTimeStamp int64
	startEpoch    int64
}

func NewSnowflake(workerId int64) (*Snowflake, error) {
	if workerId < 0 || workerId > MaxWorkerId {
		return nil, errors.New("workerId out of range")
	}

	return &Snowflake{
		workerId:      workerId,
		sequence:      0,
		lastTimeStamp: 0,
		startEpoch:    customEpoch,
	}, nil
}

func (sf *Snowflake) getCurrentTime() int64 {
	return (time.Now().UnixNano()/int64(time.Millisecond) - sf.startEpoch)
}

func (sf *Snowflake) GenerateId() (int64, error) {
	sf.lock.Lock()
	defer sf.lock.Unlock()

	now := sf.getCurrentTime()

	if now < int64(sf.lastTimeStamp) {
		return -1, errors.New("living backwards in time")
	}

	if now == int64(sf.lastTimeStamp) {

		sf.sequence = (sf.sequence + 1) & MaxSequence

		if sf.sequence == 0 {
			for now <= sf.lastTimeStamp {
				now = sf.getCurrentTime()
			}
		}

	} else {
		sf.sequence = 0
	}

	sf.lastTimeStamp = now

	id := sf.sequence | (sf.workerId << int64(workerShiftBits)) | (now << int64(timeShiftBits))

	return id, nil
}
