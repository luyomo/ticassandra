package main 

import (
    "fmt"
//    "time"
    "log"
    "strconv"
     "database/sql"
//     "encoding/binary"
//    "reflect"
      _ "github.com/go-sql-driver/mysql"
     "github.com/didi/gendry/scanner"
)

const (
    CUSTOM_TYPE = 0x0000
    AsciiType = 0x0001
    LongType = 0x0002
    BytesType = 0x0003
    BooleanType = 0x0004
    CounterColumnType = 0x0005
    DecimalType = 0x0006
    DoubleType = 0x0007
    FloatType = 0x0008
    Int32Type = 0x0009
    UTF8Type = 0x000A
    DateType = 0x000B
    UUIDType = 0x000C
    VarcharType = 0x000D
    IntegerType = 0x000E
    TimeUUIDType = 0x000F
    InetAddressType = 0x0010
    SimpleDateType = 0x0011
    TimeType = 0x0012
    ShortType = 0x0013
    ByteType = 0x0014
    DurationType = 0x0015
    ListType = 0x0020
    MapType = 0x0021
    SetType = 0x0022
    UserType = 0x0030
    TupleType = 0x0031
)

const (
    HEADER_CODE = 0x00000002
    OPTION = 0x00000001
)


type ColumnType struct {
    ColumnName      string
    TiDataType      string
    Length          int
    TiPrecision     int
    CDataTypeCode   uint16
}

var MapDataType = map[string]uint16{
    "INT" : 0x000E,
}

var MapDataTypeLength = map[string]uint16{
    "INT" : 0x04,
}

func queryAnything(db *sql.DB, query string) (string) {
    var output string
    keyspace := "dummykeyspace"
    tablename := "dummytablename"
    rows, err := db.Query(query)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    columnNames, err := rows.Columns()
    if err != nil {
        panic(err)
    }

    header := fmt.Sprintf("%08x%08x%08x%04x%x%04x%x", HEADER_CODE, OPTION, len(columnNames), len(keyspace), keyspace, len(tablename), tablename)
    output = header
    fmt.Printf("header: <%s>\n", header)


    columnTypes := make ([]ColumnType, len(columnNames))
    //fmt.Printf("The column names %#v\n", columnNames)
    for _idx, _columnName := range columnNames {
        columnTypes[_idx].ColumnName = _columnName
    }
    fieldType, err := rows.ColumnTypes()
    if err != nil {
        panic(err)
    }
    for _idx, schema  := range fieldType {
        fmt.Printf("The schema is <%s>\n", schema.DatabaseTypeName() )
        columnTypes[_idx].TiDataType     = schema.DatabaseTypeName()
        columnTypes[_idx].CDataTypeCode  = MapDataType[schema.DatabaseTypeName()]
    }
    fmt.Printf("The data is <%#v>\n", columnTypes)

    for _, columnMetaData := range columnTypes {
        columnMeta := fmt.Sprintf("%04x%x%04x", len(columnMetaData.ColumnName), columnMetaData.ColumnName, columnMetaData.CDataTypeCode) 
        fmt.Printf("column meta data: <%s>\n", columnMeta)
        output += columnMeta
    }

    result,err := scanner.ScanMap(rows)
    fmt.Printf("The result is <%v> %d \n", result, len(result))
    body := fmt.Sprintf("%08x", len(result))
    for _,record := range result {
        for _, columnMetaData := range columnTypes {
            if value, ok := (record[columnMetaData.ColumnName]).([]byte); ok {
                //fmt.Printf("length is -> %d \n", int(value ))
                //fmt.Printf("The value is <%d>, %x \n", binary.BigEndian.Uint16(value), 1)
                //fmt.Printf("The value is <%d>, %x \n", uint16(value), 1)
                __intValue, _err := strconv.Atoi(string(value))
                if _err != nil {
                    panic("Error")
                }
                fmt.Printf("The value here is <%d>\n", __intValue)
                body += fmt.Sprintf("%08x%08x",MapDataTypeLength[columnMetaData.TiDataType] , __intValue)
            } else {
                fmt.Printf("It's wrong.")
            }
            //body += fmt.Sprintf("%08x%d",MapDataTypeLength[columnMetaData.TiDataType] , int(record[columnMetaData.ColumnName] ))
            //fmt.Printf("The body is : <%s>", body)
            //fmt.Printf("Data size is <%d> value: <%#v> \n" , MapDataTypeLength[columnMetaData.TiDataType] , binary.BigEndian.Uint32([]byte(record[columnMetaData.ColumnName])))
            //fmt.Printf("Data size is <%08x> value: <%08x> \n" , MapDataTypeLength[columnMetaData.TiDataType] , value )
            //columnMeta := fmt.Sprintf("%04x%x%04x", len(columnMetaData.ColumnName), columnMetaData.ColumnName, columnMetaData.CDataTypeCode) 
            //fmt.Printf("column meta data: <%s>\n", columnMeta)
            //output += columnMeta
        }
    }
    fmt.Printf("body: %s\n", body)

    return output + body
}

func main() {
    db, err := sql.Open("mysql", "cqluser:cqluser@tcp(192.168.1.108:4000)/test")
    if err != nil {
        panic (err) 
        return
    }
    defer db.Close()
    output := queryAnything(db, "SELECT col01, col02 FROM test01" )
    fmt.Printf("The data is <%s>\n", output)
}
