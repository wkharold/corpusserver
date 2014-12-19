package catalog

import (
	"fmt"
	"os"
	"path"

	"github.com/cznic/ql"
)

const (
	CorpusServerDir = "/var/lib/corpusserver"
)

type MbxMsg struct {
	ID      int64
	Mailbox string
	Folder  string
	Msgfile string
}

var (
	schema ql.List
	insmsg ql.List
)

func init() {
	_, err := os.Open(CorpusServerDir)
	if err != nil {
		switch err.(*os.PathError).Err.Error() {
		case "no such file or directory":
			if os.Mkdir(CorpusServerDir, os.ModePerm) != nil {
				panic(fmt.Sprintf("Can't create server directory [%v]", err))
			}
		default:
			panic(fmt.Sprintf("Unexpected error [%v]", err))
		}
	}

	schema = ql.MustSchema((*MbxMsg)(nil), "", nil)
	insmsg = ql.MustCompile(`
		BEGIN TRANSACTION;
			INSERT INTO MbxMsg VALUES($1, $2, $3);
		COMMIT;`)

	db, err := ql.OpenFile(path.Join(CorpusServerDir, "catalog.db"), &ql.Options{true, nil, nil})
	if err != nil {
		panic(fmt.Sprintf("Can't open catalog database [%v]", err))
	}
	if _, _, err = db.Execute(ql.NewRWCtx(), schema); err != nil {
		panic(fmt.Sprintf("Can't create schema [%v]", err))
	}
	db.Close()
}

func Cataloger(mc chan MbxMsg, done chan struct{}) {
	db, err := ql.OpenFile(path.Join(CorpusServerDir, "catalog.db"), &ql.Options{true, nil, nil})
	if err != nil {
		panic(fmt.Sprintf("Can't open catalog database [%v]", err))
	}
	defer db.Close()

loop:
	for {
		select {
		case msg := <-mc:
			fmt.Println(msg)
			if _, _, err := db.Execute(ql.NewRWCtx(), insmsg, ql.MustMarshal(&msg)...); err != nil {
				panic(fmt.Sprintf("Message insert failed [%v]", err))
			}
		case <-done:
			break loop
		}
	}
}
