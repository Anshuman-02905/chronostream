package dlq

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Anshuman-02905/chronostream/internal/event"
	"github.com/Anshuman-02905/chronostream/internal/monotime"
)

type FileDlq struct {
	file *os.File
	mu   sync.Mutex
	ts   monotime.TimeSource
}

func NewFileDlq(directory string, instanceID string, ts monotime.TimeSource) (*FileDlq, error) {
	//Create a directory if it does not exist
	//Open a file like `dql-node-1-2026-03-19.json` in append mode
	//Check if directory exists if not exists then create
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return nil, err
	}
	//Generate file name with date
	now := ts.Now()
	dateStr := now.Format("2006-01-02")

	fileName := fmt.Sprintf(
		"dlq-%s-%s.json",
		instanceID,
		dateStr,
	)

	fullPath := filepath.Join(directory, fileName)
	file, err := os.OpenFile(
		fullPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)

	if err != nil {
		return nil, err
	}
	return &FileDlq{
		file: file,
		ts:   ts,
	}, nil

}

func (fq *FileDlq) Writebatch(ctx context.Context, events []event.Event) error {
	fq.mu.Lock()
	defer fq.mu.Unlock()
	//JSON Marshal each event
	//  Append new line "\n"
	//  Write to f.file
	// f.file.Sync to gurantee persistence"
	batchData := []byte{}
	for _, ev := range events {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		//Marshal event to JSON
		data, err := json.Marshal(ev)
		if err != nil {
			return fmt.Errorf("dlq amrshal failed :%w ", err)
		}
		data = append(data, '\n')
		batchData = append(batchData, data...)
	}
	//write to File
	_, err := fq.file.Write(batchData)
	if err != nil {
		return fmt.Errorf("write to dlq failed %v", err)
	}
	err = fq.file.Sync()
	if err != nil {
		return fmt.Errorf("dlq sync failed :%v", err)
	}

	return nil
}

func (fq *FileDlq) Close(ctx context.Context) error {
	fq.mu.Lock()
	defer fq.mu.Unlock()

	if fq.file != nil {
		return fq.file.Close()
	}
	return nil
}
