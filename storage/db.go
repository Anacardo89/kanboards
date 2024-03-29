package storage

import (
	"database/sql"
	"os"

	"github.com/Anacardo89/kanboards/fsops"
	"github.com/Anacardo89/kanboards/logger"
	_ "github.com/mattn/go-sqlite3"
)

const (
	ErrCreatSQLstmt string = "Error creating SQL statement:"
	ErrExecSQLstmt  string = "Error executing SQL statement:"
	ErrSQLrowScan   string = "Error scanning rows:"
)

var DB *sql.DB

func DBExists() bool {
	_, err := os.Open(fsops.DBPath)
	return err == nil
}

func OpenDB() {
	var err error
	DB, err = sql.Open("sqlite3", fsops.DBPath)
	if err != nil {
		logger.Error.Fatal("Cannot establish DB connection:", err)
		err = nil
	}
}

func CreateDBTables() {
	file, err := os.OpenFile(fsops.DBPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		logger.Error.Fatal("Cannot open DB file:", err)
	}
	defer file.Close()

	OpenDB()

	createProjects, err := DB.Prepare(CreateTableProjects)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = createProjects.Exec()
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}

	createBoards, err := DB.Prepare(CreateTableBoards)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = createBoards.Exec()
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}

	createLabels, err := DB.Prepare(CreateTableLabels)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = createLabels.Exec()
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}

	createCards, err := DB.Prepare(CreateTableCards)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = createCards.Exec()
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}

	createCardLabels, err := DB.Prepare(CreateTableCardLabels)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = createCardLabels.Exec()
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}

	createCheckItems, err := DB.Prepare(CreateTableCheckItems)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = createCheckItems.Exec()
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
}
