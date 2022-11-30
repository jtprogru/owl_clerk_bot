package smmemstore

import (
	"errors"
	"github.com/jtprogru/owl_clerk_bot/internal/entities/smentiti"

	"github.com/jtprogru/owl_clerk_bot/internal/usecases/app/repo/smrepo"
	"sync"
)

var _ smrepo.SmStore = &States{}

type States struct {
	sync.Mutex
	m map[int]smentiti.SM
}

func NewMemStore() *States {
	return &States{
		m: make(map[int]smentiti.SM),
	}
}

func (s *States) CreateState(sm smentiti.SM) (int, error) {
	s.Lock()
	defer s.Unlock()
	for _, st := range s.m {
		if st.Id == sm.Id {
			return st.Id, nil
		}
	}
	s.m[sm.Id] = sm
	return sm.Id, nil
}

func (s *States) ReadState(ind int) (*smentiti.SM, error) {
	s.Lock()
	defer s.Unlock()

	ss, ok := s.m[ind]
	if ok {
		return &ss, nil
	}
	return nil, errors.New("state Not Found")
}

func (s *States) UpdateState(smt smentiti.SM) (*smentiti.SM, error) {
	s.Lock()
	defer s.Unlock()
	s.m[smt.Id] = smt

	return &smt, nil
}

func (s *States) DeleteState(id int) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.m[id]
	if ok {
		delete(s.m, id)
		return nil
	}
	return errors.New("state Not Found")
}
