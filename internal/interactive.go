package internal

import (
	"unicode"

	"github.com/nsf/termbox-go"
)

type keyDir struct {
	key  termbox.Key
	char rune
	dir  int
}

var keyDirs = []*keyDir{
	{termbox.KeyArrowUp, 'k', Up},
	{termbox.KeyArrowDown, 'j', Down},
	{termbox.KeyArrowLeft, 'h', Left},
	{termbox.KeyArrowRight, 'l', Right},
}

func interactive(maze *InfiniMaze, format *Format) {
	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()
	strwriter := make(chan string)
	go printTermbox(strwriter)
	maze.CurrentMaze.Write(strwriter, format)
loop:
	for {
		select {
		case event := <-events:
			if event.Type == termbox.EventKey {

				for _, keydir := range keyDirs {
					if event.Key == keydir.key || event.Ch == keydir.char {
						maze.Move(keydir.dir)
						maze.CurrentMaze.Write(strwriter, format)
						continue loop
					}
				}
				if event.Ch == 'q' || event.Ch == 'Q' || event.Key == termbox.KeyCtrlC || event.Key == termbox.KeyCtrlD {
					break loop
				}
			}
		}
	}
}

func printTermbox(strwriter chan string) {
	x, y := 1, 0
	for {
		str := <-strwriter
		switch str {
		case "\u0000":
			_ = termbox.Flush()
			x, y = 1, 0
		default:
			printString(str, &x, &y)
		}
	}
}

func printString(str string, x *int, y *int) {
	attr, skip, d0, d1, d := false, false, '0', '0', false
	fg, bg := termbox.ColorDefault, termbox.ColorDefault
	for _, c := range str {
		if c == '\n' {
			*x, *y = (*x)+1, 0
		} else if c == '\x1b' || attr && c == '[' {
			attr = true
		} else if attr && unicode.IsDigit(c) {
			if !skip {
				if d {
					d1 = c
				} else {
					d0, d = c, true
				}
			}
		} else if attr && c == ';' {
			skip = true
		} else if attr && c == 'm' {
			if d0 == '7' && d1 == '0' {
				fg, bg = termbox.AttrReverse, termbox.AttrReverse
			} else if d0 == '3' {
				fg, bg = termbox.Attribute(uint64(d1-'0'+1)), termbox.ColorDefault
			} else if d0 == '4' {
				fg, bg = termbox.ColorDefault, termbox.Attribute(uint64(d1-'0'+1))
			} else {
				fg, bg = termbox.ColorDefault, termbox.ColorDefault
			}
			attr, skip, d0, d1, d = false, false, '0', '0', false
		} else {
			termbox.SetCell(*y, *x, c, fg, bg)
			*y = *y + 1
		}
	}
}
