package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

// Message 구조체 정의
type Message struct {
	ID          int    `json:"id"`
	Data        string `json:"data"`
	RandomValue int    `json:"random_value"`
}

func main() {
	// Kafka 설정
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0 // Kafka 버전 (Strimzi와 호환되는 Kafka 버전)
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // "earliest"와 동일

	brokers := []string{"study-cluster-kafka-bootstrap.kafka.svc.cluster.local:9092"}
	topic := "log"
	groupID := "my-group"

	// Consumer Group 생성
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	// 종료 신호 처리
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	// Consumer Group 핸들러
	handler := &ConsumerGroupHandler{}

	// 메시지 소비 루프
	go func() {
		for {
			if err := consumerGroup.Consume(context.Background(), []string{topic}, handler); err != nil {
				log.Fatalf("Error consuming: %v", err)
			}
		}
	}()

	fmt.Println("Kafka Consumer started. Press Ctrl+C to stop.")
	<-sigchan
	fmt.Println("Shutting down consumer...")
}

// ConsumerGroupHandler는 sarama.ConsumerGroupHandler를 구현합니다.
type ConsumerGroupHandler struct{}

// Setup은 Consumer Group이 초기화될 때 호출됩니다.
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup은 Consumer Group이 종료될 때 호출됩니다.
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim은 메시지를 처리합니다.
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var message Message
		if err := json.Unmarshal(msg.Value, &message); err != nil {
			log.Printf("Failed to parse message: %s\n", err)
		} else {
			fmt.Printf("Received message:\n  Partition: %d\n  Offset: %d\n  ID: %d\n  Data: %s\n  RandomValue: %d\n",
				msg.Partition, msg.Offset, message.ID, message.Data, message.RandomValue)
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
