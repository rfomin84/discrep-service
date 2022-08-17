package balance_history

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	balance_history "github.com/rfomin84/discrep-service/internal/services/balance_history/domain"
	"github.com/rfomin84/discrep-service/pkg/logger"
	"github.com/rfomin84/discrep-service/pkg/store/mysql_client"
	"github.com/spf13/viper"
	"log"
	"time"
)

type BalanceHistoryStorage struct {
	conn *sql.DB
}

func NewBalanceHistoryStorage(cfg *viper.Viper) *BalanceHistoryStorage {
	conn, err := mysql_client.NewMysqlClient(
		cfg.GetString("MYSQL_HOST"),
		cfg.GetString("MYSQL_PORT"),
		cfg.GetString("MYSQL_USERNAME"),
		cfg.GetString("MYSQL_PASSWORD"),
		cfg.GetString("MYSQL_DATABASE"),
	)
	if err != nil {
		logger.Error("error connect : " + err.Error())
		log.Fatal("error")
	}

	err = conn.Ping()
	if err != nil {
		logger.Error("error ping : " + err.Error())
		log.Fatal("error")
	}

	return &BalanceHistoryStorage{
		conn: conn,
	}
}

func (storage *BalanceHistoryStorage) Save(data []balance_history.BalanceHistory) {
	ctx := context.Background()
	// шаг 1 — объявляем транзакцию
	tx, err := storage.conn.Begin()
	if err != nil {
		logger.Error(err.Error())
	}
	// шаг 1.1 — если возникает ошибка, откатываем изменения
	defer tx.Rollback()

	// шаг 2 — готовим инструкцию
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO balance_history(feed_id, date, cost, approved) VALUES(?,?,?,?)")
	if err != nil {
		logger.Error("ERROR 1" + err.Error())
	}

	// шаг 2.1 — не забываем закрыть инструкцию, когда она больше не нужна
	defer stmt.Close()

	for _, row := range data {
		// шаг 3 — указываем, что каждое видео будет добавлено в транзакцию
		if _, err := stmt.ExecContext(ctx, row.FeedId, row.Date, row.Cost, 0); err != nil {
			logger.Error("ERROR 2 : " + err.Error())
			//log.Fatal(res)
		}

	}
	// шаг 4 — сохраняем изменения
	if err = tx.Commit(); err != nil {
		logger.Error("ERROR 3" + err.Error())
	}

	logger.Info(fmt.Sprintf("result save balance history : %v", nil))
}

func (storage *BalanceHistoryStorage) DeleteNotApproveStatistics(start, end time.Time) {
	ctx := context.Background()

	approved := 0
	_, err := storage.conn.ExecContext(ctx, "DELETE FROM balance_history WHERE date BETWEEN ? AND ? and approved = ?;",
		start,
		end,
		approved,
	)

	if err != nil {
		logger.Error(err.Error())
		log.Fatal(err.Error())
	}
}

func (storage *BalanceHistoryStorage) GetReserveFeedBalance() ([]balance_history.ReservedBalance, error) {
	result := make([]balance_history.ReservedBalance, 0)
	ctx := context.Background()
	rows, err := storage.conn.QueryContext(
		ctx,
		"SELECT feed_id, SUM(cost) FROM balance_history WHERE approved = ? GROUP BY feed_id",
		0,
	)
	if err != nil {
		return nil, err
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var reserveBalance balance_history.ReservedBalance
		err = rows.Scan(&reserveBalance.FeedId, &reserveBalance.Cost)
		if err != nil {
			return nil, err
		}

		result = append(result, reserveBalance)
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}
