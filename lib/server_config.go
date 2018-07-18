package lib

import (
	"database/sql"
	"time"
)

type sqliteBackend struct {
	db *sql.DB
}

var opToColumn = map[string]string{
	"++": "plusplus",
	"--": "minusminus",
	"+-": "plusminus",
}

func (s sqliteBackend) incr(key string) (err error) {
	identifier, operation := key[:len(key)-2], key[len(key)-2:]

	txn, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func(txn *sql.Tx) {
		if err != nil {
			txn.Rollback()
		}
	}(txn)

	statement := "UPDATE karma SET " + opToColumn[operation] + "=" + opToColumn[operation] + " + 1 WHERE identifier = ?"
	stmt, err := txn.Prepare(statement)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(identifier)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	// successful update
	if rows == 1 {
		txn.Commit()
		return nil
	}

	plusplus := 0
	minusminus := 0
	plusminus := 0
	switch operation {
	case "++":
		plusplus += 1
	case "--":
		minusminus += 1
	case "+-":
		plusminus += 1
	}

	stmt, err = txn.Prepare("INSERT INTO karma (identifier, plusplus, minusminus, plusminus) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	result, err = stmt.Exec(identifier, plusplus, minusminus, plusminus)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s sqliteBackend) addChannelKarma(mutation karmaMutation, channel string) (err error) {
	txn, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func(txn *sql.Tx) {
		if err != nil {
			txn.Rollback()
		}
	}(txn)

	dateBucket := time.Now().Format("2006-01-02")

	statement := "UPDATE channel_karma SET " + opToColumn[mutation.op] + "=" + opToColumn[mutation.op] + " + 1 WHERE identifier = ? AND channel = ? AND date_bucket = ?"
	stmt, err := txn.Prepare(statement)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(mutation.identifier, channel, dateBucket)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	// successful update
	if rows == 1 {
		txn.Commit()
		return nil
	}

	plusplus := 0
	minusminus := 0
	plusminus := 0
	switch mutation.op {
	case "++":
		plusplus += 1
	case "--":
		minusminus += 1
	case "+-":
		plusminus += 1
	}

	stmt, err = txn.Prepare("INSERT INTO channel_karma (identifier, channel, date_bucket, plusplus, minusminus, plusminus) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	result, err = stmt.Exec(
		mutation.identifier,
		channel,
		dateBucket,
		plusplus,
		minusminus,
		plusminus,
	)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

// TODO use channel_karma for this
func (s sqliteBackend) getInt(key string) (value int) {
	identifier, operation := key[:len(key)-2], key[len(key)-2:]

	stmt, err := s.db.Prepare("SELECT " + opToColumn[operation] + " FROM karma WHERE identifier = ?")
	if err != nil {
		panic(err)
	}
	err = stmt.QueryRow(identifier).Scan(&value)
	if err != nil {
		return 0
	}
	return value
}

type karmaWithId struct {
	identifier string
	karma      karmaSet
}

func (s sqliteBackend) getChannelTop3(channel string, minimumDate string) (karmaList []karmaWithId) {
	stmt, err := s.db.Prepare("SELECT identifier, SUM(plusplus) as plusplus_sum, SUM(minusminus) as minusminus_sum, SUM(plusminus) FROM channel_karma WHERE channel = ? AND date_bucket > ? GROUP BY identifier ORDER BY (plusplus_sum - minusminus_sum) DESC LIMIT 3")
	if err != nil {
		panic(err)
	}
	rows, err := stmt.Query(channel, minimumDate)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var karma karmaWithId
		if err := rows.Scan(&karma.identifier, &karma.karma.plusplus, &karma.karma.minusminus, &karma.karma.plusminus); err != nil {
			panic(err)
		}
		karmaList = append(karmaList, karma)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	return karmaList
}
