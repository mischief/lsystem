package lsystem

import (
	"bytes"

//	"strconv"
//	"strings"
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

func (ls *LSystem) Draw() {
	//var commands []string

	//ls.Gc.MoveTo(ls.InitPosX, ls.InitPosY)

	//x, y := ls.InitPosX, ls.InitPosY

	//ang := ls.InitAngle

	/*

		for _, ch := range ls.WorkingSet {
			fmt.Printf("%c %s\n", ch, ls.TurtleGraphics[ch])
			commands = append(commands, ls.TurtleGraphics[ch])
		}

		for _, cmd := range commands {
			spl := strings.Split(cmd, " ")
			if spl != nil && spl[0] == "" {
				ls.TurtleFns[cmd](ls, 0)
			} else {
				arg, err := strconv.Atoi(spl[1])
				if err != nil {
					fmt.Printf("conversion error: %s: %s\n", spl[1], err)
					break
				}
				ls.TurtleFns[spl[0]](ls, arg)
			}
		}
	*/
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
