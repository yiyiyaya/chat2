package usage

import (

)
type DatabaseJoin_0 struct {
    DatabaseJoin_Base
}
func NewDatabaseCreate_0() *DatabaseJoin_0 {
    updater := &DatabaseJoin_0{}
    updater.TableCreateSqlFile = "initdb.sql"
    return updater
}