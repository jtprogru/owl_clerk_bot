package smmemstore

import (
	"errors"
	"github.com/jtprogru/owl_clerk_bot/internal/entities/smentities"
	"github.com/jtprogru/owl_clerk_bot/internal/usecases/app/repo/smrepo"
	"sync"
)

var _ smrepo.SmStore = &States{}

type States struct {
	sync.Mutex
	m map[int]smentities.SM
}

func NewStates() *States {
	return &States{
		m: make(map[int]smentities.SM),
	}
}

func (s States) Create(sm smentities.SM) (int, error) {
	for _, st := range s.m {
		if st.Id == sm.Id {
			return st.Id, nil
		}
	}
	s.m[sm.Id] = sm
	return sm.Id, nil
}

func (s States) Read(ind int) (*smentities.SM, error) {
	s.Lock()
	defer s.Unlock()

	ss, ok := s.m[ind]
	if ok {
		return &ss, nil
	}
	return nil, errors.New("state Not Found")
}

func (s States) Update(smt smentities.SM) (*smentities.SM, error) {
	s.Lock()
	defer s.Unlock()
	s.m[smt.Id] = smt

	return &smt, nil
}

func (s States) Delete(id int) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.m[id]
	if ok {
		delete(s.m, id)
		return nil
	}
	return errors.New("state Not Found")
}
