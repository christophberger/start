package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Args)

	yes := flag.Bool("yes", false, "Say yes")
	fmt.Println(*yes)
}
