package lsystem

import (
	"fmt"
	"testing"
)

type ruletest struct {
	vars  *Variables
	cons  *Constants
	rules *Rules

	start    string
	steps    int
	expected string
}

var (
	expandtests = []ruletest{
		// Sierpinski Triangle
		ruletest{
			&Variables{[]rune{'A', 'B'}},
			&Constants{[]rune{'-', '+'}},
			&Rules{map[rune]string{'A': "B-A-B", 'B': "A+B+A"}},
			"A",
			4,
			"A+B+A-B-A-B-A+B+A+B-A-B+A+B+A+B-A-B+A+B+A-B-A-B-A+B+A-B-A-B+A+B+A+B-A-B-A+B+A-B-A-B-A+B+A-B-A-B+A+B+A+B-A-B-A+B+A-B-A-B-A+B+A+B-A-B+A+B+A+B-A-B+A+B+A-B-A-B-A+B+A",
		},
	}
)

func TestExpandRules(t *testing.T) {
	for _, test := range expandtests {
		ls := NewLSystem(test.start, test.vars, test.cons, test.rules)

		ls.Run(test.steps)

		state := ls.State()
		if state != test.expected {
			t.Errorf("LSystem state is %q, expected %q", state, test.expected)
		}
	}
}

// This will dump files names testN.png for each test case
func TestDrawStuff(t *testing.T) {
	for i, test := range expandtests {
		ls := NewLSystem(test.start, test.vars, test.cons, test.rules)
		ls.Run(test.steps)

		rules := NewTurtleGraphicsRules()
		tg := NewTurtleGraphics(640, 480, rules)
		rules.Add('A', DrawFwd, 10)
		rules.Add('B', DrawFwd, 10)
		rules.Add('-', Turn, -60)
		rules.Add('+', Turn, 60)

		tg.Draw(ls)

		tg.SavePNG(fmt.Sprintf("test%d.png", i))

    fmt.Printf("test %d dumped to test%d.png\n", i, i)
	}

}
