package main

import (
	"fmt"
	"github.com/Mansur51-hub/ClubHandler/handler"
	"os"
)

func main() {

        if len(os.Args) < 2 {
		fmt.Println("error: No file path provided")
		return
	}

	path := os.Args[1]
	a, err := handler.HandleInput(path)

	if err != nil {
		fmt.Printf("File read error: %s", err)
		return
	}

	h := handler.NewHandler(a)

	if res, err := h.HandleEvents(); err != nil {
		fmt.Printf("Events handling error: %s", err)
	} else {
		for _, val := range res {
			fmt.Println(val)
		}
	}
}
