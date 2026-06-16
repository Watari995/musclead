// Package main は outbox consumer Lambda。
// SQS で受け取ったメッセージを、 DynamoDB で重複排除しつつ SES でメール送信する (ADR 0020)。
// VPC 外で動かす想定 (DynamoDB / SES / SQS は公開 API、 NAT 不要 = ゼロ円)。
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sestypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type emailTemplate struct {
	subject string
	body    string
}

var emailTemplates = map[string]emailTemplate{
	"PaymentSucceeded": {
		subject: "サブスクリプションの申し込みに成功しました",
		body:    "お申し込みありがとうございます。Proプランが有効になりました。\nこのメールは送信専用です。",
	},
	// 将来ここに追加していく
}

// message は relay (payment.PublishMessage) が SQS に送る JSON 契約。
// 別モジュール (internal) なので型は共有できず、 ここで同じ形を定義する (契約は同期させること)。
type message struct {
	EventID string `json:"EventID"`
	Type    string `json:"Type"`
	Email   string `json:"Email"`
}

type consumer struct {
	ddb       *dynamodb.Client
	ses       *sesv2.Client
	tableName string // 重複排除テーブル
	fromAddr  string // 送信元 (no-reply@musclead.com)
	ttlSecond int64  // dedup レコードの TTL (秒)
}

// Handle は SQS バッチを受け取り、 1 件ずつ処理する。
func (c *consumer) Handle(ctx context.Context, e events.SQSEvent) error {
	for _, record := range e.Records {
		var msg message
		if err := json.Unmarshal([]byte(record.Body), &msg); err != nil {
			// 壊れたメッセージは再試行しても直らないので飛ばす
			slog.Error("consumer: invalid message", "err", err, "body", record.Body)
			continue
		}
		if err := c.processOne(ctx, msg); err != nil {
			// error を返すとそのメッセージは SQS に戻り再試行される (at-least-once)
			return err
		}
	}
	return nil
}

// processOne は 1 メッセージを冪等に処理する。
func (c *consumer) processOne(ctx context.Context, msg message) error {
	_, err := c.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &c.tableName,
		Item: map[string]ddbtypes.AttributeValue{
			"event_id":   &ddbtypes.AttributeValueMemberS{Value: msg.EventID},
			"expires_at": &ddbtypes.AttributeValueMemberN{Value: strconv.FormatInt(time.Now().Unix()+c.ttlSecond, 10)},
		},
		ConditionExpression: aws.String("attribute_not_exists(event_id)"),
	})
	if err != nil {
		var condErr *ddbtypes.ConditionalCheckFailedException
		if errors.As(err, &condErr) {
			return nil // condition expression が失敗した = すでに処理済みなので、何もしない
		}
		return err // それ以外のエラーはそのまま返す (SQS 再試行に委ねる)
	}
	emailTmpl, ok := emailTemplates[msg.Type]
	if !ok {
		slog.Info("consumer: no template, skip", "type", msg.Type)
		return nil
	}
	// 問題なく書けたらSESでメール送信 (送信元は c.fromAddr、宛先は msg.Email、本文は msg.Type で出し分け)
	_, err = c.ses.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: &c.fromAddr,
		Destination: &sestypes.Destination{
			ToAddresses: []string{msg.Email},
		},
		Content: &sestypes.EmailContent{
			Simple: &sestypes.Message{
				Subject: &sestypes.Content{Data: aws.String(emailTmpl.subject)},
				Body: &sestypes.Body{
					Text: &sestypes.Content{Data: aws.String(emailTmpl.body)},
				},
			},
		},
	})
	if err != nil {
		// 再度送信し直せるようにレコードを削除する
		_, _ = c.ddb.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: &c.tableName,
			Key: map[string]ddbtypes.AttributeValue{
				"event_id": &ddbtypes.AttributeValueMemberS{Value: msg.EventID},
			},
		})
		return err
	}
	return nil
}

func main() {
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		slog.Error("consumer: load aws config", "err", err)
		os.Exit(1)
	}
	c := &consumer{
		ddb:       dynamodb.NewFromConfig(awsCfg),
		ses:       sesv2.NewFromConfig(awsCfg),
		tableName: os.Getenv("DEDUP_TABLE_NAME"),
		fromAddr:  os.Getenv("SES_FROM_ADDRESS"),
		ttlSecond: 3600, // 1h
	}
	lambda.Start(c.Handle)
}
