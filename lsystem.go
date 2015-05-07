package lsystem

import (
	"bytes"
)

// Variables
type Variables struct {
	v []rune
}

func (v *Variables) Add(neu rune) *Variables {
	v.v = append(v.v, neu)
	return v
}

// Constants
type Constants struct {
	c []rune
}

func (c *Constants) Add(neu rune) *Constants {
	c.c = append(c.c, neu)
	return c
}

// L-System rules.
// Format like:
// "A": "A-B-A"
type Rules struct {
	rules map[rune]string
}

func NewRules() *Rules {
	return &Rules{make(map[rune]string)}
}

func (r *Rules) Add(tok rune, repl string) *Rules {
	r.rules[tok] = repl
	return r
}

type LSystem struct {
	state *bytes.Buffer

	vars  *Variables
	cons  *Constants
	rules *Rules
}

func (ls *LSystem) State() string {
	return ls.state.String()
}

func (ls *LSystem) SetState(s string) string {
	old := ls.state.String()
	ls.state = bytes.NewBufferString(s)
	return old
}

func (ls *LSystem) Run(steps int) {
	now := &bytes.Buffer{}

	ls.state.WriteTo(now)

	for i := 0; i < steps; i++ {
		next := &bytes.Buffer{}

		r, _, err := now.ReadRune()
		for err == nil {
			if exp, ok := ls.rules.rules[r]; ok {
				next.WriteString(exp)
			} else {
				next.WriteRune(r)
			}
			r, _, err = now.ReadRune()
		}

		next.WriteTo(now)
	}

	now.WriteTo(ls.state)
}

func NewLSystem(start string, vars *Variables, cons *Constants, rules *Rules) *LSystem {
	ls := &LSystem{
		state: bytes.NewBufferString(start),
		vars:  vars,
		cons:  cons,
		rules: rules,
	}

	return ls
}
