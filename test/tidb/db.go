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
    output := tidb.QueryAnything(db, "SELECT col01, col02 FROM test01" )
    fmt.Printf("The data is <%s>\n", output)
}
