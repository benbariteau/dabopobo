package lib

import (
	"database/sql"
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

	statement := "UPDATE channel_karma SET " + opToColumn[mutation.op] + "=" + opToColumn[mutation.op] + " + 1 WHERE identifier = ? AND channel = ?"
	stmt, err := txn.Prepare(statement)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(mutation.identifier, channel)
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

	stmt, err = txn.Prepare("INSERT INTO channel_karma (identifier, plusplus, minusminus, plusminus, channel) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	result, err = stmt.Exec(mutation.identifier, plusplus, minusminus, plusminus, channel)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

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
