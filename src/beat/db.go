package beat

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	_ "github.com/go-sql-driver/mysql"
)

// Beat 구조체 정의
type BeatDTO struct {
	ID           int
	Key          string
	BeatType     string
	RegTs        time.Time
	PresignedUrl string
}

func getDB() *sql.DB {
	database_host := os.Getenv("database_host")
	database_user := os.Getenv("database_user")
	database_password := os.Getenv("database_password")
	database_name := os.Getenv("database_name")
	source := fmt.Sprintf("%s:%s@tcp(%s)/%s", database_user, database_password, database_host, database_name)
	db, err := sql.Open("mysql", source+"?parseTime=true")
	fmt.Println("source is ... " + source)
	if err != nil {
		fmt.Println(err)
	}
	return db
}

func CreateBeat(fileKey string, beatType string) {
	db := getDB()

	beat := BeatDTO{
		Key:      fileKey,
		BeatType: beatType,
	}

	stmt, err := db.Prepare("INSERT INTO `tenseconds`.`beat`(`key`, `beat_type`, `reg_ts`) VALUES(?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(beat.Key, beat.BeatType, time.Now())
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}

func GetAllBeats() ([]BeatDTO, error) {
	db := getDB()
	rows, err := db.Query("SELECT `id`, `key`, beat_type, reg_ts FROM tenseconds.beat")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	beats := []BeatDTO{}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile("tenseconds"),
		config.WithRegion("ap-northeast-2"),
	)
	client := s3.NewFromConfig(cfg)
	bucket := os.Getenv("aws_s3_bucket")

	presignClient := s3.NewPresignClient(client)
	presigner := Presigner{PresignClient: presignClient}

	for rows.Next() {
		beat := BeatDTO{}
		err := rows.Scan(&beat.ID, &beat.Key, &beat.BeatType, &beat.RegTs)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		presignedGetRequest, err := presigner.GetObject(bucket, beat.Key, 60)
		presignedURL := presignedGetRequest.URL
		beat.PresignedUrl = presignedURL
		beats = append(beats, beat)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return beats, nil
}

func DeleteBeatById(id string) {
	db := getDB()

	stmt, err := db.Prepare("DELETE FROM `tenseconds`.`beat` where `id` = ?")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}
