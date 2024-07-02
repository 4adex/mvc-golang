package main
import (
	"fmt"
	"github.com/4adex/mvc-golang/pkg/api"
)

func main() {
	fmt.Println("Started the API server")
	api.Start()
}