package util

import (
	"database/sql"
)

func CommitOrRollback(tx *sql.Tx) error {
	err := recover()
	if err != nil {
		errRollback := tx.Rollback()
		return errRollback
	} else {
		errorCommit := tx.Commit()
		return errorCommit
	}
}
