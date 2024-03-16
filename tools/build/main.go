package main 

import (
  "github.com/ngn13/venom/builder"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Venom build tool")
		fmt.Println("=======================================================")
		fmt.Println("Applies config and mapping and builds agent from source")
		fmt.Printf("Source will be searched in the specified directory\n\n")
		fmt.Printf("Usage: %s <debug> <config> <output> <directory>\n", os.Args[0])
		return
	}

	var (
    ctx builder.Ctx
		err error
	)

	ctx.Dir = os.Args[4]
	ctx.Out = os.Args[3]
	ctx.Key = []byte("")

	ctx.Debug = os.Args[1] == "yes" || os.Args[1] == "true"
	ctx.Config, err = os.ReadFile(os.Args[2])

	if err != nil {
		fmt.Printf("Failed to load config: %s\n", err.Error())
		os.Exit(1)
	}

	err = builder.Run(&ctx)
	if err != nil {
		fmt.Printf("Build failed: %s\n", err.Error())
	}
}
