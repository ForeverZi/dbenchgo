package main

import (
	"context"
	_ "dbenchgo/conf"
	_ "dbenchgo/driver"
	"dbenchgo/task"
	"fmt"
)

func main() {
	benchTask := task.NewInsertBench(30, 100000)
	ctx := context.Background()
	err := benchTask.SetUp(ctx)
	if err != nil {
		panic(err)
	}
	statusCh, err := benchTask.Run(ctx)
	if err != nil {
		panic(err)
	}
	for status := range statusCh {
		fmt.Print("\r")
		fmt.Printf("total:%v,completed:%v", status.Total, status.Completed)
	}
	fmt.Println("")
	result, err := benchTask.CollectResult(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Summary())
}
