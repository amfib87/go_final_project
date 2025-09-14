package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var Db *sql.DB
var FileDbEnv string

func Init(dbFile string) error {

	const schema = `CREATE TABLE scheduler (
    			id INTEGER PRIMARY KEY AUTOINCREMENT,
    			date CHAR(8) NOT NULL DEFAULT "",
    			title VARCHAR(256) NOT NULL DEFAULT "",
				comment TEXT NOT NULL DEFAULT "",
				repeat VARCHAR(128) NOT NULL DEFAULT "" );

			CREATE INDEX ind_date ON scheduler (date);`

	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	Db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = Db.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
