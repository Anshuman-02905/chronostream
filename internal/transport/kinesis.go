package transport

import (
	"context"
	"encoding/json"
	"fmt"

	cfg "github.com/Anshuman-02905/chronostream/internal/config"
	"github.com/Anshuman-02905/chronostream/internal/event"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"github.com/sirupsen/logrus"
)

// Option addition logger metric partition strategy
// Acceptance criterion Struct Compiles Transport interface
type AwsKinesisTransport struct {
	client     *kinesis.Client
	streamName string
	partition  string
}

func NewAwsKinesisTransport(ctx context.Context, kcfg cfg.Config) (*AwsKinesisTransport, error) {
	if !kcfg.Kinesis.Enabled {
		return nil, fmt.Errorf("Kinesis transport diabled")
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(kcfg.Kinesis.Region),
	)

	if err != nil {
		return nil, err
	}
	client := kinesis.NewFromConfig(awsCfg)
	logrus.Infof("Instance ID is. %v", kcfg.Instance.ID)
	return &AwsKinesisTransport{
		client:     client,
		streamName: kcfg.Kinesis.StreamName,
		partition:  kcfg.Instance.ID,
	}, nil

}

func (k *AwsKinesisTransport) Send(ctx context.Context, e event.Event) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = k.client.PutRecord(
		ctx,
		&kinesis.PutRecordInput{
			StreamName:   aws.String(k.streamName),
			Data:         data,
			PartitionKey: aws.String(e.UserID), // Route by UserID so one user = one shard
		},
	)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"event_id":  e.ID,
		"stream":    k.streamName,
		"partition": k.partition,
	}).Debug("Event pushed to Kinesis")
	return err
}

func (k *AwsKinesisTransport) Close(ctx context.Context) error {
	return nil
}

func (k *AwsKinesisTransport) SendBatch(ctx context.Context, events []event.Event) error {
	if len(events) == 0 {
		return nil
	}
	records := make([]types.PutRecordsRequestEntry, 0, len(events))

	for _, e := range events {
		data, err := json.Marshal(e)
		if err != nil {
			return err
		}
		record := types.PutRecordsRequestEntry{
			Data:         data,
			PartitionKey: aws.String(k.partition),
		}
		records = append(records, record)
	}
	output, err := k.client.PutRecords(
		ctx,
		&kinesis.PutRecordsInput{
			StreamName: aws.String(k.streamName),
			Records:    records,
		},
	)
	if err != nil {
		return err
	}
	//It only returns count need to integrate a struct to map records back to original events

	if output.FailedRecordCount != nil && *output.FailedRecordCount > 0 {
		var failedIDs []string
		for i, res := range output.Records {
			if res.ErrorCode != nil {
				failedIDs = append(failedIDs, events[i].ID)
			}
		}
		return fmt.Errorf(
			"kinesis batch failed: %d records failed. IDs: %v",
			*output.FailedRecordCount,
			failedIDs,
		)
	}

	logrus.WithFields(logrus.Fields{
		"count":  len(events),
		"stream": k.streamName,
	}).Info("Batch pushed to Kinesis")

	return nil
}
