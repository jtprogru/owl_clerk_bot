package smrepo

import (
	"fmt"
	"github.com/jtprogru/owl_clerk_bot/internal/entities/smentiti"
)

type SmStore interface {
	CreateState(smentiti.SM) (int, error)
	ReadState(int) (*smentiti.SM, error)
	UpdateState(smentiti.SM) (*smentiti.SM, error)
	DeleteState(u int) error
}

type States struct {
	smStore SmStore
}

func NewSmRepo(smstore SmStore) *States {
	return &States{
		smStore: smstore,
	}
}

func (st *States) CreateState(sm smentiti.SM) (int, error) {
	smt, err := st.smStore.CreateState(sm)
	if err != nil {
		return -100, fmt.Errorf("error create State: %w", err)
	}
	return smt, nil
}

func (st *States) ReadState(id int) (*smentiti.SM, error) {
	smt, err := st.smStore.ReadState(id)
	if err != nil {
		return nil, fmt.Errorf("error read State: %w", err)
	}
	return smt, nil
}

func (st *States) UpdateState(sm smentiti.SM) (*smentiti.SM, error) {
	smt, err := st.smStore.UpdateState(sm)
	if err != nil {
		return nil, fmt.Errorf("error update State: %w", err)
	}
	return smt, nil
}

func (st *States) DeleteState(u int) error {
	err := st.smStore.DeleteState(u)
	if err != nil {
		return fmt.Errorf("error delete State: %w", err)
	}
	return nil
}
