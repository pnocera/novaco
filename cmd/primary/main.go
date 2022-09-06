package main

import (
	"flag"

	"github.com/pnocera/novaco/internal/service"
)

func main() {
	mode := flag.String("mode", "primary", "primary or secondary or dev")

	_ = service.StartNew(*mode)
}
