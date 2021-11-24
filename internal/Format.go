package internal

// Format is the printing format of the infiniMaze
type Format struct {
	Wall      string
	Path      string
	Cursor    string
	ExitDown  string
	ExitRight string
	ExitLeft  string
	ExitUp    string
}

// Ascii format
var Ascii = &Format{
	Wall:      "##",
	Path:      "  ",
	Cursor:    "@@",
	ExitDown:  "VV",
	ExitRight: ">>",
	ExitLeft:  "<<",
	ExitUp:    "^^",
}

// Default Color format
var Color = &Format{
	Wall:      "  ",
	Path:      "\x1b[7m  \x1b[0m",
	ExitDown:  "\x1b[42;1mVV\x1b[0m",
	ExitRight: "\x1b[42;1m>>\x1b[0m",
	ExitLeft:  "\x1b[42;1m<<\x1b[0m",
	ExitUp:    "\x1b[42;1m^^\x1b[0m",
	Cursor:    "\x1b[43;1m@@\x1b[0m",
}
