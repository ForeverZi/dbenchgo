package task

import (
	"context"
	"dbenchgo/utils"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const createTable = `
CREATE TABLE IF NOT EXISTS dbenchgo.insertbench (
    id  bigint not null,
    c1 varchar(64) not null,
    c2 varchar(64) not null,
    c3 varchar(64) not null,
    i1 bigint not null,
    i2 bigint not null,
    i3 bigint not null,
    primary key(id),
    index idx_i1_c1 (i1, c1)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
`

func NewInsertBench(parallel int, recordCount int) *InsertBench {
	task := &InsertBench{
		BaseTask: BaseTask{
			Driver:   "mysql",
			Parallel: parallel,
		},
		recordCount: recordCount,
		idChan:      make(chan int64, parallel*3),
		tokenBucket: utils.NewTokenBucket(parallel, true),
		result:      NewOpResult(),
	}
	return task
}

type InsertBench struct {
	BaseTask
	recordCount int
	completed   atomic.Int64
	idChan      chan int64
	tokenBucket *utils.TokenBucket
	result      *OpResult
}

func (task *InsertBench) SetUp(ctx context.Context) (err error) {
	_, err = task.DB().ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS dbenchgo")
	if err != nil {
		return
	}
	_, err = task.DB().ExecContext(ctx, "DROP TABLE dbenchgo.insertbench")
	if err != nil {
		return
	}
	_, err = task.DB().ExecContext(ctx, createTable)
	return
}

func (task *InsertBench) Run(ctx context.Context) (statusCh <-chan Status, err error) {
	sendCh := make(chan Status, task.Parallel*3)
	statusCh = sendCh
	go func() {
		var wg sync.WaitGroup
		for i := 0; i < task.Parallel; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for id := range task.idChan {
					err = task.doInsert(context.Background(), id)
					if err != nil {
						log.Print("can't insert id:%v due to %v\n", id, err)
					}
					completed := task.completed.Add(1)
					status := Status{
						Total:     task.recordCount,
						Completed: int(completed),
						Done:      completed >= int64(task.recordCount),
						Error:     err,
					}
					// status will not block the process
					select {
					case sendCh <- status:
					default:
					}
				}
			}()
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(task.idChan)
			for i := 0; i < task.recordCount; i++ {
				task.idChan <- int64(i)
			}
		}()
		wg.Wait()
		close(sendCh)
	}()
	return
}

func (task *InsertBench) doInsert(ctx context.Context, id int64) (err error) {
	c1 := utils.GenerateRandomString(30)
	c2 := utils.GenerateRandomString(30)
	c3 := utils.GenerateRandomString(30)
	i1 := rand.Int63()
	i2 := rand.Int63()
	i3 := rand.Int63()
	begin := time.Now()
	_, err = task.DB().ExecContext(ctx, "INSERT INTO dbenchgo.insertbench(id, c1,c2,c3,i1,i2,i3) values(?,?,?,?,?,?,?)",
		id, c1, c2, c3, i1, i2, i3)
	task.result.Record(time.Since(begin), err == nil)
	return
}

func (task *InsertBench) CleanUp(ctx context.Context) (err error) {
	return
}

func (task *InsertBench) CollectResult(ctx context.Context) (result Result, err error) {
	result = task.result
	return
}
