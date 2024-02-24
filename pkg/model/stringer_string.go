package model

type StringerString string

func (s StringerString) String() string {
	return string(s)
}
