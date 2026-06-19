package services

import (
	"sync"
	"time"
)

type Statistic struct {
	LastCall time.Time
	Hit      int
	Key      string
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

func NewStatistics() Statistics {
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

	val, loaded := s.data.Load(maxHitInfos.(*maxHit).key)
	if !loaded {
		return nil
	}

	return val.(*Statistic)
}

func (s *statisticsImpl) Increment(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, loaded := s.data.LoadOrStore(key, &Statistic{LastCall: time.Now(), Hit: 1, Key: key})
	if !loaded {
		return
	}

	val.(*Statistic).Hit++
	val.(*Statistic).LastCall = time.Now()

	s.data.Store(key, val)

	maxHitInfos, loaded := s.data.LoadOrStore("max_hit", &maxHit{key: key})
	if !loaded {
		return
	}

	maxHitVal, _ := s.data.Load(maxHitInfos.(*maxHit).key)
	if val.(*Statistic).Hit > maxHitVal.(*Statistic).Hit {
		s.data.Store("max_hit", &maxHit{key: key})
	}
}
