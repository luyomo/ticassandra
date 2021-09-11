package main

import (
    "github.com/luyomo/ticql/pkg/tidb"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "fmt"
)

func main() {
    db, err := sql.Open("mysql", "cqluser:cqluser@tcp(192.168.1.108:4000)/test")
    if err != nil {
        panic (err) 
        return
    }
    defer db.Close()
    //output := tidb.QueryAnything(db, "SELECT col01, col02 FROM test01" )
    //output := tidb.QueryAnything(db, "select version, TABLE_SCHEMA, CREATE_TIME from information_schema.tables limit 1 " )
    //output := tidb.QueryAnything(db, "select version, TABLE_SCHEMA, CREATE_TIME from information_schema.tables where table_name = 'test03' limit 1 " )
    //output := tidb.QueryAnything(db, "select host, db as keyspace_name, user from mysql.db limit 1 " )
    output := tidb.QueryAnything(db, "select * from information_schema.columns limit 1" )
    fmt.Printf("The data is <%s>\n", output)
}
