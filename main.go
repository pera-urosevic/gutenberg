package main

import (
	"os"

	"somnusalis.org/gutenberg/system"
	"somnusalis.org/gutenberg/templates/stamps"
)

func main() {
	system.Log("Gutenberg")

	if len(os.Args) < 2 {
		system.Log("Usage: gutenberg <project> [args...]")
		return
	}

	template := os.Args[1]
	switch template {
	case "stamps":
		stamps.Stamps()
	}
}
