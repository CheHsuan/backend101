package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type balanceRecord struct {
	ID          int64
	UserID      int64
	Balance     int64
	Deleted     bool
	CreatedTime time.Time
	UpdatedTime time.Time
}

const (
	userID = 123456
)

func main() {
	db, err := initDB()
	if err != nil {
		panic(err.Error())
	}

	switch os.Getenv("APP_TYPE") {
	case "blocking":
		// init record
		id, err := initRecord(db)
		if err != nil {
			panic(err.Error())
		}

		tx := db.Begin()
		log.Println("start blocking transaction")
		recd := balanceRecord{}
		if err := tx.Where("id = ?", id).First(&recd).Error; err != nil {
			tx.Rollback()
		}
		now := time.Now()
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&recd).Update("balance", recd.Balance-10).Update("updated_time", now).Error; err != nil {
			tx.Rollback()
		}

		// sleep 3 min
		// simulate unexpected errors or bugs that stuck the transaction.
		time.Sleep(3 * 60 * time.Second)

		if err := tx.Commit().Error; err != nil {
			log.Println("blocking app failed to commit: ", err)
		}
	case "normal":
		if err := db.Transaction(func(tx *gorm.DB) error {
			log.Println("start normal transaction")
			recd := balanceRecord{}
			if err := tx.Where("user_id = ?", userID).First(&recd).Error; err != nil {
				return err
			}
			now := time.Now()
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&recd).Update("balance", recd.Balance-10).Update("updated_time", now).Error; err != nil {
				return err
			}

			return nil
		}); err != nil {
			log.Println("normal app failed to commit: ", err)
		}
	}

	log.Println("app finished")
}

func initRecord(db *gorm.DB) (int64, error) {
	now := time.Now()
	recd := balanceRecord{UserID: userID, Balance: 100, CreatedTime: now, UpdatedTime: now}
	if result := db.Create(&recd); result.Error != nil {
		return 0, result.Error
	}

	return recd.ID, nil
}

func initDB() (*gorm.DB, error) {
	var db *sql.DB
	if err := backoff.Retry(func() error {
		var err error
		db, err = sql.Open("mysql", "user:pass@tcp(mysql:3306)/testdb?parseTime=true")
		return err
	}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5)); err != nil {
		return nil, err
	}

	return gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{SkipDefaultTransaction: true})
}
