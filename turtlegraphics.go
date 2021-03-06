package lsystem

import (
	"bufio"
	"image"
	"image/png"
	"math"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
)

type Vector struct {
	X, Y float64
}

type TGFunc func(tg *TurtleGraphics, arg1 int)

type TGCall struct {
	Func TGFunc
	Arg  int
}

type TGState struct {
	P Vector
	A float64
}

// TurtleGraphics rules.
//
// Maps the L-System state string to some code to be executed
type TurtleGraphicsRules struct {
	graphics map[rune]TGCall
}

func NewTurtleGraphicsRules() *TurtleGraphicsRules {
	return &TurtleGraphicsRules{make(map[rune]TGCall)}
}

func (tg *TurtleGraphicsRules) Add(let rune, fn TGFunc, arg int) *TurtleGraphicsRules {
	tg.graphics[let] = TGCall{fn, arg}
	return tg
}

// Turtle Graphics system.
//
// Allows you to draw an LSystem onto an image buffer,
// and dump it to a file.
type TurtleGraphics struct {
	Image *image.RGBA
	Gc    *draw2dimg.GraphicContext

	rules *TurtleGraphicsRules

	ReductionScale     float64
	InitPosX, InitPosY float64
	InitAngle          float64

	// Current position
	Pos Vector

	// Current Angle
	CA float64

	// TODO:
	// add stack functions
	// add better way to insert variables to functions
	// Maps alphabet letter to turtlegraphics call
	//TurtleGraphics map[rune]string

	// TurtleGraphics call to function map
	//TurtleFns map[string]TGFunc

	// TurtleGraphics state
	//PosStack, AngleStack Stack
	Stack []TGState
}

// Draw the LSystem to the image buffer
func (tg *TurtleGraphics) Draw(ls *LSystem) {
	str := ls.State()

	tg.Pos = Vector{tg.InitPosX, tg.InitPosY}
	tg.CA = tg.InitAngle
	tg.Gc.BeginPath()
	tg.Gc.MoveTo(tg.InitPosX, tg.InitPosY)

	for _, c := range str {
		if call, ok := tg.rules.graphics[c]; ok {
			call.Func(tg, call.Arg)
		}
	}

	tg.Gc.FillStroke()
}

// Calculate new line position given angle and distance
func (tg *TurtleGraphics) getnewpos(angle, distance float64) Vector {
	d2r := math.Pi / 180
	op := math.Sin(angle*d2r) * distance
	ad := math.Cos(angle*d2r) * distance

	newp := Vector{op, ad}

	newp.X += tg.Pos.X
	newp.Y += tg.Pos.Y

	return newp
}

func Push(tg *TurtleGraphics, unused int) {
	st := TGState{tg.Pos, tg.CA}
	tg.Stack = append(tg.Stack, st)
}

func Pop(tg *TurtleGraphics, unused int) {
	if tg.Stack == nil || len(tg.Stack) == 0 {
		return
	}

	st := tg.Stack[len(tg.Stack)-1]
	tg.Stack = tg.Stack[:len(tg.Stack)-1]

	tg.Pos, tg.CA = st.P, st.A
	tg.Gc.MoveTo(tg.Pos.X, tg.Pos.Y)
}

// Draw function.
//
// Draw a line from where we are to 10 px in front of us at the current angle.
func DrawFwd(tg *TurtleGraphics, px int) {
	newp := tg.getnewpos(tg.CA, float64(px))

	//fmt.Printf("OK draw (%d px %f angle) from %v to %v\n", px, tg.CA, tg.Pos, newp)

	tg.Gc.LineTo(newp.X, newp.Y)

	tg.Pos = newp
}

// Turn function.
//
// Turn a number of degrees.
func Turn(tg *TurtleGraphics, deg int) {
	tg.CA += float64(deg)
}

// Write this TurtleGraphics instance to a png file
func (tg *TurtleGraphics) SavePNG(filePath string) error {
	f, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer f.Close()

	b := bufio.NewWriter(f)

	err = png.Encode(b, tg.Image)

	if err != nil {
		return err
	}

	err = b.Flush()

	if err != nil {
		return err
	}

	return nil
}

// Make a new TurtleGraphics System.
//
// A ruleset needs to be passed in, translating the
// LSystem state into drawing function calls.
func NewTurtleGraphics(width, height int, rules *TurtleGraphicsRules) *TurtleGraphics {
	tg := &TurtleGraphics{
		rules: rules,

		InitPosX:  float64(width / 2),
		InitPosY:  float64(height / 2),
		InitAngle: 90.0,
	}

	tg.Image = image.NewRGBA(image.Rect(0, 0, width, height))
	tg.Gc = draw2dimg.NewGraphicContext(tg.Image)
	tg.Gc.Clear()

	return tg
}
