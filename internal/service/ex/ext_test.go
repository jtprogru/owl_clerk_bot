package ex

import (
	"github.com/jtprogru/owl_clerk_bot/internal/transport/tg"
	"testing"
)

type MockAnswer struct {
	msg string
	kb  []string
}

func (ma MockAnswer) GetMessage() string {
	return ma.msg
}

func (ma MockAnswer) GetKeyboard() []string {
	return ma.kb
}

func TestButtonSelect(t *testing.T) {
	tests := []struct {
		name    string
		p       tg.Profile
		m       MockAnswer
		nx      []int
		btn     []string
		cs      int
		want    int  // Ожидаемый результат
		wantErr bool // Ожидаем ошибку?
	}{
		{
			name: "+Return Current state",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "State1",
			},
			cs:      2,
			nx:      []int{},
			btn:     []string{"State1"},
			want:    2,
			wantErr: false,
		},
		{
			name: "-Return Current state",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "State1",
			},
			cs:      2,
			nx:      []int{},
			btn:     []string{"State1"},
			want:    0,
			wantErr: true,
		},
		{
			name: "+ Return [0]",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "State1",
			},
			cs:      2,
			nx:      []int{1},
			btn:     []string{"State1"},
			want:    1,
			wantErr: false,
		},
		{
			name: "- Return [0]",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "State5",
			},
			cs:      2,
			nx:      []int{5},
			btn:     []string{"State5"},
			want:    1,
			wantErr: true,
		},
		{
			name: "+ Return [index]",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "State6",
			},
			cs:      4,
			nx:      []int{5, 6, 1, 2},
			btn:     []string{"State5", "State6", "State1", "State2"},
			want:    6,
			wantErr: false,
		},
		{
			name: "- Return [index]",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "State6",
			},
			cs:      4,
			nx:      []int{5, 6, 1, 2},
			btn:     []string{"State5", "State6", "State1", "State2"},
			want:    2,
			wantErr: true,
		},
		{
			name: "- Answer != Button",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "State66",
			},
			cs:      4,
			nx:      []int{5, 6, 1, 2},
			btn:     []string{"State5", "State6", "State1", "State2"},
			want:    2,
			wantErr: true,
		},
		{
			name: "+ Answer != Button",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "State66",
			},
			cs:      4,
			nx:      []int{4, 5, 6, 1, 2},
			btn:     []string{"State4", "State5", "State6", "State1", "State2"},
			want:    4,
			wantErr: false,
		},
		{
			name: "+ != Len",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "State6",
			},
			cs:      4,
			nx:      []int{5, 6, 1},
			btn:     []string{"State5", "State6", "State1", "State2"},
			want:    4,
			wantErr: false,
		},
		{
			name: "- != Len",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "State6",
			},
			cs:      4,
			nx:      []int{5, 6, 1},
			btn:     []string{"State5", "State6", "State1", "State2"},
			want:    2,
			wantErr: true,
		},
		{
			name:    "- != Len",
			p:       tg.Profile{},
			m:       MockAnswer{},
			cs:      4,
			nx:      []int{5, 6, 1},
			btn:     []string{"State5", "State6", "State1", "State2"},
			want:    2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		{
			t.Run(tt.name, func(t *testing.T) {
				got := ButtonSelect(tt.nx, tt.btn, tt.cs)
				res := got(tt.p, tt.m)
				if (res != tt.want) != tt.wantErr {
					t.Errorf("ButtonSelect() result = %d, want result %v", res, tt.want)
				}
			})
		}
	}
}

func TestIntSelect(t *testing.T) {
	tests := []struct {
		name    string
		p       tg.Profile
		m       MockAnswer
		nx      []int
		btn     []string
		cs      int
		want    int  // Ожидаемый результат
		wantErr bool // Ожидаем ошибку?
	}{
		{
			name: "+Odd",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "3",
			},
			cs:      2,
			nx:      []int{},
			btn:     []string{"State1"},
			want:    0,
			wantErr: false,
		},
		{
			name: "-Odd",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "3",
			},
			cs:      2,
			nx:      []int{},
			btn:     []string{"State1"},
			want:    1,
			wantErr: true,
		},
		{
			name: "Odd return CS",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "3,5",
			},
			cs:      2,
			nx:      []int{},
			btn:     []string{"State1"},
			want:    2,
			wantErr: false,
		},
		{
			name: "+Even",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "4",
			},
			cs:      5,
			nx:      []int{},
			btn:     []string{"State1"},
			want:    1,
			wantErr: false,
		},
		{
			name: "-Even",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "4",
			},
			cs:      5,
			nx:      []int{},
			btn:     []string{"State1"},
			want:    5,
			wantErr: true,
		},
		{
			name: "+String return CS",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "STRING",
			},
			cs:      5,
			nx:      []int{},
			btn:     []string{"State1"},
			want:    5,
			wantErr: false,
		},
		{
			name: "-String Return CS",
			p:    tg.Profile{},
			m: MockAnswer{
				msg: "STRING",
			},
			cs:      5,
			nx:      []int{},
			btn:     []string{"State1"},
			want:    2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		{
			t.Run(tt.name, func(t *testing.T) {
				res := IntSelect(tt.nx, tt.btn, tt.cs)(tt.p, tt.m)
				//res := got(tt.p, tt.m)
				if (res != tt.want) != tt.wantErr {
					t.Errorf("ButtonSelect() result = %d, want result %v", res, tt.want)
				}
			})
		}
	}
}
