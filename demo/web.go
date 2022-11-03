package demo

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello "+os.Getenv("USERNAME"))
	})

	err := http.ListenAndServe(":3000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Print("server closed\n")
	} else if err != nil {
		fmt.Print("error starting server: \n", err)
		os.Exit(1)
	}
}
