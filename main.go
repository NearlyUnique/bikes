package main

import (
	"fmt"
	"os"
)

func main() {
	act := getAction(os.Args)
	var err error
	switch act {
	case "init":
		err = createIndex()
	case "find":
		if len(os.Args) < 3 {
			fmt.Println("missing find param")
			os.Exit(1)
		}
		err = find(os.Args[2:])
	case "show":
		err = show(os.Args[2:])
	}

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func getAction(args []string) string {
	if len(args) < 2 {
		return ""
	}
	return args[1]
}
