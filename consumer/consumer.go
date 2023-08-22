package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Jamesbarford/video-meta/server/database"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const MIN_COMMIT_COUNT = 10

type VideoMessage struct {
	VideoId int `json:"videoId"`
}

/* Don't blow up collect the error in logging */
func deleteVideoFromMetaData(db *database.DbConnection, videoId int) {
	query := "DELETE FROM meta_data WHERE video_id = ?"
	rows, err := db.PrepareStmt(query, videoId)
	defer rows.Close()

	if err != nil {
		log.Printf("Error deleting all meta_data for video with videoId %d: %v", videoId, err)
	}
}

/**
 * Very simple kafka consumer for listening to video delete events which would
 * happen from a `video-service` producing to the kafka topic
 */
func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "host1:9092,host2:9092",
		"group.id":          "kafka-delete-consumer",
		"auto.offset.reset": "smallest",
	})

	if err != nil {
		panic("Failed to connect to create consumer " + err.Error())
	}

	msgCount := 0
	run := 1
	db, err := database.NewDbConnection(database.DbConfigFromEnvironment())

	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	defer consumer.Close()

	if err := consumer.Subscribe("video-delete", nil); err != nil {
		panic("Failed to subscribe to topic: " + err.Error())
	}

	for run == 1 {
		ev := consumer.Poll(100)

		switch e := ev.(type) {
		case kafka.AssignedPartitions:
			log.Printf("%% %v\n", e)
			consumer.Assign(e.Partitions)

		case kafka.RevokedPartitions:
			log.Printf("%% %v\n", e)
			consumer.Unassign()

		case *kafka.Message:
			msgCount += 1
			if msgCount%MIN_COMMIT_COUNT == 0 {
				consumer.Commit()
			}

			var videoMsg VideoMessage
			if err := json.Unmarshal(e.Value, &videoMsg); err != nil {
				log.Printf("Failed to decode message: %s\n", err)
				continue
			}
			deleteVideoFromMetaData(db, videoMsg.VideoId)

		case kafka.PartitionEOF:
			log.Printf("%% Reached %v\n", e)
			break
		case kafka.Error:
			log.Panicf("%% Error: %v\n", e)
			run = 0
		default:
			fmt.Printf("Ignored %v\n", e)
		}
	}
}
