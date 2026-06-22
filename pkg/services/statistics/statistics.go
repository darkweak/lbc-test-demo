package statistics

import (
	"sync"
	"time"
)

type Statistic struct {
	LastCall time.Time `json:"last_call"`
	Hit      int       `json:"hit"`
	Key      string    `json:"key"`
}

type maxHit struct {
	key string
}

type Statistics interface {
	GetMostRecent() *Statistic
	Increment(key string)
}

type statisticsImpl struct {
	data sync.Map
	mu   sync.RWMutex
}

var _ Statistics = (*statisticsImpl)(nil)

func NewStatistics() *statisticsImpl {
	return &statisticsImpl{
		data: sync.Map{},
		mu:   sync.RWMutex{},
	}
}

func (s *statisticsImpl) GetMostRecent() *Statistic {
	s.mu.RLock()
	defer s.mu.RUnlock()

	maxHitInfos, loaded := s.data.Load("max_hit")
	if !loaded {
		return nil
	}

	maxhit, ok := maxHitInfos.(*maxHit)
	if !ok {
		return nil
	}

	val, loaded := s.data.Load(maxhit.key)
	if !loaded {
		return nil
	}

	value, ok := val.(*Statistic)
	if !ok {
		return nil
	}

	return value
}

func (s *statisticsImpl) Increment(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, loaded := s.data.LoadOrStore(key, &Statistic{LastCall: time.Now(), Hit: 1, Key: key})
	if !loaded {
		return
	}

	stats, ok := val.(*Statistic)
	if !ok {
		return
	}

	stats.Hit++
	stats.LastCall = time.Now()

	s.data.Store(key, val)

	maxHitInfos, loaded := s.data.LoadOrStore("max_hit", &maxHit{key: key})
	if !loaded {
		return
	}

	maxhit, ok := maxHitInfos.(*maxHit)
	if !ok {
		return
	}

	maxHitVal, _ := s.data.Load(maxhit.key)

	mhv, ok := maxHitVal.(*Statistic)
	if !ok {
		return
	}

	if stats.Hit > mhv.Hit {
		s.data.Store("max_hit", &maxHit{key: key})
	}
}
