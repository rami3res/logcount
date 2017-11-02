package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/rami3res/logcount/logcount"
)

func main() {
	flag.Parse()

	/*if len(os.Args) == 1 {
		fmt.Println("Файл для чтения не указан")
		os.Exit(1)
	}*/

	//fileNames := flag.Args()
	//filename := fileNames[0]

	filename := "/home/ramieres/logs/access.log.x"

	accLog, err := logcount.ReadLines(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Невозможно прочитать файл: %v: %v\n", filename, err)
		os.Exit(1)
	}

	cs := accLog.SortByIP()

	if len(cs) < *logcount.MaxIPCount {
		*logcount.MaxIPCount = len(cs)
	}

	for i := *logcount.MaxIPCount; i >= 0; i-- {

		fmt.Printf("%v Requests from %s\n", cs[i].Count, cs[i].IP)


		if !*logcount.StatusFlag {

			for statusCode, count := range cs[i].StatusCode {
				fmt.Printf("         Status code %v: %v\n", statusCode, count)
			}

			for statusCode, statusData := range cs[i].StatusData {
				fmt.Printf("         Status data %v: %5.2f kB\n", statusCode, float64(statusData/1024))
			}

			fmt.Println()

		}



	}


	fmt.Printf("\ntotal IP-adresses: %v\n", len(cs))
	fmt.Printf("total Requests: %v\n", len(accLog))
}