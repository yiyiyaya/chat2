package usage

import (

	"database/sql"
	"io/ioutil"
	"path/filepath"
	"github.com/asiainfoLDP/datahub_commons/log"
	"bytes"
)
type DatabaseJoin_Base struct {
    TableCreateSqlFile   string
}

func (join *DatabaseJoin_Base) CreateTables(db *sql.DB) error{
    if join.TableCreateSqlFile == "" {
        return nil
    }
    data, err := ioutil.ReadFile(filepath.Join("db", join.TableCreateSqlFile))
    if err != nil {
        log.DefaultlLogger().Errorf("file name read err")
        return err
    }
    sqls := bytes.SplitAfter(data, []byte("DEFAULT CHARSET=UTF8;"))
    sqls = sqls[:len(sqls)-1]
    for _, sql := range sqls{
      _, err :=  db.Exec(string(sql))
      if err != nil {
          return err
      }
    }
    return nil
}