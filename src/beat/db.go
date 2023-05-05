package beat

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Beat 구조체 정의
type BeatDTO struct {
	ID       int
	Key      string
	BeatType string
}

func getDB() *sql.DB {
	database_user := os.Getenv("database_user")
	database_password := os.Getenv("database_password")
	database_name := os.Getenv("database_name")
	source := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", database_user, database_password, database_name)
	db, err := sql.Open("mysql", source)
	if err != nil {
		panic(err)
	}
	return db
}

func CreateBeat(fileKey string, beatType string) {
	db := getDB()

	beat := BeatDTO{
		Key:      fileKey,
		BeatType: beatType,
	}

	stmt, err := db.Prepare("INSERT INTO `tenseconds`.`beat`(`key`, `beat_type`) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(beat.Key, beat.BeatType)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllBeats() ([]BeatDTO, error) {
	db := getDB()
	rows, err := db.Query("SELECT `id`, `key`, beat_type FROM tenseconds.beat")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beats := []BeatDTO{}
	for rows.Next() {
		beat := BeatDTO{}
		err := rows.Scan(&beat.ID, &beat.Key, &beat.BeatType)
		if err != nil {
			return nil, err
		}
		beats = append(beats, beat)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return beats, nil
}
