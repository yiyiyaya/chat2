package api

import (
    "os"
	"fmt"
    "github.com/asiainfoLDP/datahub_commons/log"
    //"log"
	//"unicode/utf8"
    "database/sql"
	"sync"
	"time"
    "github.com/asiainfoLDP/chat/usage"
)
var (
    dbMutex sync.Mutex
    dbInstance *sql.DB
)
func setDB(db *sql.DB){
    dbMutex.Lock()
    dbInstance = db
    defer dbMutex.Unlock()

}
func getDB() *sql.DB {
    dbMutex.Lock()
	defer dbMutex.Unlock()
	return dbInstance
}
func connectDB(){
    DB_ADDR := os.Getenv("MYSQL_ADDR")
    DB_PORT := os.Getenv("MySQL_PORT")
    DB_DATABASE := os.Getenv("MYSQL_DATABASE")
    DB_USER := os.Getenv("MYSQL_USER")
    DB_PASSWORD := os.Getenv("MYSQL_PASSWORD")
    DB_URL := fmt.Sprintf(`%s:%s@tcp(%s:port)/%s?charset=utf8&&parseTime=true`, DB_USER, DB_PASSWORD, DB_ADDR, DB_PORT, DB_DATABASE)
    log.DefaultlLogger().Info("connect to", DB_URL)
    db, err := sql.Open("mysql", DB_URL)
    if err == nil {
        err = db.Ping()
    }
    if err != nil {
        log.DefaultlLogger().Errorf("err: %s\n", err)
        if db != nil {
            db.Close()
        }
        return
    }
    //保存实例
    setDB(db)
}
func initDB()  {
    for i := 0; i < 3 ; i++ {
        connectDB()
        if getDB() == nil {
            select{
                case <- time.After(time.Second * 10) :
                    continue
            }
        } else {
            break
        }
    }
    
    dbjoin := usage.NewDatabaseCreate_0()
    err := dbjoin.CreateTables(getDB())
    if err != nil {
        log.DefaultlLogger().Errorf("create table failed err")
        return 
    }
    go  updateDB()
}
func updateDB () {
    ticker := time.Tick(5 * time.Second)
    for range ticker{
        db := getDB()
        if db == nil {
            connectDB()
        } else if err := db.Ping(); err != nil {
            log.DefaultlLogger().Errorf("db ping err: %s\n", err)
            db.Close()
            connectDB()
        }

    }
}
