package tidb

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
     "github.com/araddon/dateparse"
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
    "INT"      : 0x000E,
    "BIGINT"   : 0x0002,
    "VARCHAR"  : 0x000D,
    "TEXT"     : 0x000D,
    "CHAR"     : 0x000D,
    "DATETIME" : 0x000B,
}

var MapDataTypeLength = map[string]uint16{
    "INT"      : 0x04,
    "BIGINT"   : 0x08,
    "DATETIME" : 0x08,
}

func QueryAnything(db *sql.DB, query string) (string) {
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
        // fmt.Printf("The schema is <%s>\n", schema.DatabaseTypeName() )
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
            if record[columnMetaData.ColumnName] == nil {
                switch  columnMetaData.TiDataType {
                    case
                        "INT",
                        "BIGINT":
                            if columnMetaData.TiDataType == "INT" {
                                body += fmt.Sprintf("%08x%08x",MapDataTypeLength[columnMetaData.TiDataType], 0x00)
                            }else if columnMetaData.TiDataType == "BIGINT" {
                                body += fmt.Sprintf("%08x%016x",MapDataTypeLength[columnMetaData.TiDataType], 0x00)
                            }
                    case "VARCHAR",
                         "CHAR",
                         "TEXT":
                       body += fmt.Sprintf("%08x",0x00)
                    case "DATETIME":
                        body += fmt.Sprintf("%08x%016x",MapDataTypeLength[columnMetaData.TiDataType] , 0 )
                }
                continue
            }

            fmt.Printf("The data is %s: <%#v> \n", columnMetaData.ColumnName  , record[columnMetaData.ColumnName])
            if byteValue, ok := (record[columnMetaData.ColumnName]).([]byte); ok {
                switch  columnMetaData.TiDataType {
                    case
                        "INT",
                        "BIGINT":
                           //if value, ok := (record[columnMetaData.ColumnName]).([]byte); ok {
                               fmt.Printf("The column [%s] is %s \n" , columnMetaData.TiDataType, columnMetaData.ColumnName )
                               //__intValue, _err := strconv.Atoi(string(value))
                               __intValue, _err := strconv.Atoi(string(byteValue))
                               if _err != nil {
                                   panic("Error")
                               }
                               fmt.Printf("The value here is <%d>\n", __intValue)
                               if columnMetaData.TiDataType == "INT" {
                                   body += fmt.Sprintf("%08x%08x",MapDataTypeLength[columnMetaData.TiDataType] , __intValue)
                               }else if columnMetaData.TiDataType == "BIGINT" {
                                   body += fmt.Sprintf("%08x%016x",MapDataTypeLength[columnMetaData.TiDataType] , __intValue)
                               }
                           //} else {
                           //    fmt.Printf("It's wrong.")
                          // }
                    case "VARCHAR",
                         "CHAR",
                         "TEXT":
                       fmt.Printf("The column [VARCHAR] is %s \n" , columnMetaData.ColumnName )
                       fmt.Printf("The string is <%#v>\n", string(byteValue))
                       body += fmt.Sprintf("%08x%x",len(byteValue), byteValue)
                       //body += fmt.Sprintf("%08x%08x",MapDataTypeLength[columnMetaData.TiDataType] , __intValue)
                    case "DATETIME":
                        fmt.Printf("The column [DATETIME] is %s \n" , columnMetaData.ColumnName )
                        fmt.Printf("The date time is <%s>\n", byteValue )
                        theDateTime, err := dateparse.ParseLocal(string(byteValue) )
                        if err != nil {
                            panic(err.Error())
                        }
                        fmt.Printf("The date time is <%d>\n",  theDateTime.Unix())
                        body += fmt.Sprintf("%08x%016x",MapDataTypeLength[columnMetaData.TiDataType] , theDateTime.Unix() )
                }
            } else {
                panic(fmt.Sprintf("Failed to parse the value: name: <%s>, value: <%#v>", columnMetaData.ColumnName,  record[columnMetaData.ColumnName]))
            }
        }
    }
    fmt.Printf("body: %s\n", body)

    return output + body
}

