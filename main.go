package main

import (
	// "fmt"
	// "go_test"
	"fmt"
	"go_test2/aco"
	// "go_test2/plot"
)


func main(){
	// plot.Run()
	task := aco.NewDownloadTask(20, 300, 10)
	// task := aco.NewDownloadTask(50, 500, 5)
	// task := aco.NewDownloadTask(node_count:50, chunck_count:500, concurrent_download:10)
	task.Run()
	fmt.Printf("AAAA")
	// fmt.Println("DONE")
	// util.New()
	// util.New2()
}