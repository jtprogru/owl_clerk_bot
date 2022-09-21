package smrepo

import (
	"fmt"
	"github.com/jtprogru/owl_clerk_bot/internal/entities/smentities"
)

type SmStore interface {
	Create(smentities.SM) (int, error)
	Read(int) (*smentities.SM, error)
	Update(smentities.SM) (*smentities.SM, error)
	Delete(u int) error
}
type States struct {
	smStore SmStore
}

func NewSmStore(smstore SmStore) *States {
	return &States{
		smStore: smstore,
	}
}

func (st *States) Create(sm smentities.SM) (int, error) {
	smt, err := st.smStore.Create(sm)
	if err != nil {
		return -100, fmt.Errorf("error create State: %w", err)
	}
	return smt, nil
}

func (st *States) Read(id int) (*smentities.SM, error) {
	smt, err := st.smStore.Read(id)
	if err != nil {
		return nil, fmt.Errorf("error read State: %w", err)
	}
	return smt, nil
}

func (st *States) Update(smt *smentities.SM) (*smentities.SM, error) {
	smt, err := st.smStore.Update(*smt)
	if err != nil {
		return nil, fmt.Errorf("error update State: %w", err)
	}
	return smt, nil
}

func (st *States) Delete(u int) error {
	err := st.smStore.Delete(u)
	if err != nil {
		return fmt.Errorf("error delete State: %w", err)
	}
	return nil
}
