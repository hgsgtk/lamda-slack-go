package lambdahandler

import (
	"database/sql"

	"github.com/pkg/errors"
)

// KPI 新規獲得ユーザー数
func CountUser(db *sql.DB) (int, error) {
	q := `SELECT count(*) from users`

	stmt, err := db.Prepare(q)
	if err != nil {
		return 0, errors.Wrap(err, "preparing statement error")
	}

	var userNum int
	if err := stmt.QueryRow().Scan(&userNum); err != nil {
		return 0, errors.Wrap(err, "querying and scanning database error")
	}
	return userNum, nil
}

// KPI 累計実績金額
func GetTotalOrder(db *sql.DB) (int, error) {
	q := `SELECT sum(amount) from orders`

	stmt, err := db.Prepare(q)
	if err != nil {
		return 0, errors.Wrap(err, "preparing statement error")
	}

	var totalAmount int
	if err := stmt.QueryRow().Scan(&totalAmount); err != nil {
		return 0, errors.Wrap(err, "scanning database error")
	}

	return totalAmount, nil
}
