# InfiniMaze : an infinite, persistent, procedurally generated, explorable maze

Built off the great work by @itchyny here: https://github.com/itchyny/maze

## Usage
The `infiniMaze` command without arguments starts a console version of InfiniMaze in the Color mode with a size that matches the console window.
```sh
infiniMaze
```

To run the web verison run the below command and visit [http://localhost:3000/](http://localhost:3000/)
```sh
infiniMaze --web
```

You can get a list of command line options with the below command:
```sh
$ infiniMaze --h
NAME:
   infiniMaze - InfiniMaze is an infinite, persistent, procedurally generated, explorable maze

USAGE:
   infiniMaze [global options] [arguments...]

VERSION:
   0.0.2

AUTHOR:
   Static-Flow

GLOBAL OPTIONS:
   --width value     The width of the infiniMaze
   --height value    The height of the infiniMaze
   --format default  Output format, default or `ascii`
   --seed value      The random seed
   --help            Shows the help of the command (default: false)
   --scale value     Scale output image (default: 0)
   --web             Runs the web version of InfiniMaze (default: false)
   --version, -v     print the version (default: false)
```

## Installation

### Build from source
```bash
go get -u github.com/Static-Flow/infiniMaze/infiniMaze
```

## Bug Tracker
Report bug at [Issues・Static-Flow/infiniMaze - GitHub](https://github.com/Static-Flow/infiniMaze/issues).

## Author
Static-Flow (https://github.com/Static-Flow)

## License
This software is released under the MIT License, see LICENSE.

## Special thanks
Special thanks to the [termbox-go](https://github.com/nsf/termbox-go) library.
