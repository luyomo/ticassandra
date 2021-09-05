package server

import (
	"fmt"
//	"io"
	"net"
    "encoding/binary"
    "regexp"
    "encoding/hex"
//    "hex"
//    "reflect"
)

// TCPServer struct
type TCPServer struct {
	Bind string
	Port int
}

// Start TCPServer
func (s *TCPServer) Start() {
	fmt.Printf("started tcp echo server... ... \n")
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Bind, s.Port))
	defer ln.Close()
	if err != nil {
		panic(err)
	}
	for {
        //message := "0000000000000000000"
        var message []byte
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
        fmt.Println("------------------")
		go func(conn net.Conn ) {
            sessionVersion := 0
			defer conn.Close()
            for {
                tableName := ""
                msg := make([]byte, 256)
                //msg := make([]byte, 1024)
                fmt.Println("\n\n\nStarting to collect data ********** ********** ")
                _, err = conn.Read(msg)
                fmt.Println("Gathered the info from remote ")
                if err != nil {
                    fmt.Println("Errorr")
                    panic(err)
                }
                // fmt.Println(msg)
                fmt.Printf("%x\n", msg)
                version := int(msg[0])
                op := int(msg[4])

                if sessionVersion != 0 && version != sessionVersion {
                    version = int(msg[6])
                    op = int(msg[10])
                    dataLength := binary.BigEndian.Uint32(msg[11:15])
                    streamID := binary.BigEndian.Uint16(msg[8:10])
                    fmt.Printf("The version is %x and the opt is %x length: <%d> streamID: <%d> \n", version, op, dataLength, streamID)
                    // data := string(msg[16:17+dataLength])
                    if op == 0x0b {
                        data := string(msg[16:35])
                        data02 := string(msg[36:50])
                        data03 := string(msg[51:64])
                        //data := string(msg[16:])
                        fmt.Printf("data: <%s>\n", data)
                        fmt.Printf("data: <%s>\n", data02)
                        fmt.Printf("data: <%s>\n", data03)
                    }

                    if op == 0x07 {
                        data := string(msg[16:17+dataLength-8])
                        fmt.Printf("data: <%s>\n", data)
                        re01 := regexp.MustCompile(".*FROM (.*?) ")
                        re02 := regexp.MustCompile(".*FROM (.*)$")
                        match01 := re01.FindStringSubmatch(data)
                        if len(match01) > 1 {
                            tableName = match01[1]
                        } else {
                            match02 := re02.FindStringSubmatch(data)
                            if len(match02) > 1 {
                                tableName = match02[1]
                            }
                        }
                        fmt.Printf("version: <%x>, opt: <%x>, tableName: <%s> \n", version, op, tableName )
                    }
                } else {
                    fmt.Printf("The version is %x and the opt is %x \n", version, op )
                }

                if  version == 66 &&  op == 5 {
                    message=returnUnsupport01()
                    conn.Write(message)
                    break
                }
                if  version == 65 &&  op == 5 {
                    message=returnUnsupported02()
                    conn.Write(message)
                    break
                }
                if version == 5 && op == 5 {
                    sessionVersion = 5
                    message=returnServerMeta()
                    conn.Write(message)
                }
                if version == 5 && op == 1 {
                    message=returnStartup()
                    conn.Write(message)
                }
                if version == 5 && op == 0x0b {
                    message=returnMsg01()
                    conn.Write(message)
                }
                // fmt.Printf("version: <%x>, opt: <%x>, tableName: <%x> vs <%x> \n", version, op, tableName, "system.peers_v2" )
                if op == 0x07 && tableName == "system.peers_v2" {
                    fmt.Println("peers_v2: data sending")
                    message=returnPeerV2()
                    conn.Write(message)
                }
                if op == 0x07 && tableName == "system.local" {
                    fmt.Println("local: data sending")
                    message=returnLocal()
                    conn.Write(message)
                }
                if op == 0x07 && tableName == "system_schema.keyspaces" {
                    fmt.Println("keyspace: data sending")
                    message=returnKeyspaces()
                    conn.Write(message)
                }
                if op == 0x07 && tableName == "system_schema.tables" {
                    fmt.Println("tables: data sending")
                    message=returnTables()
                    conn.Write(message)
                }
            }
		}(conn)
        fmt.Println(">>>>>>>>>>>")
	}
}

func returnUnsupport01() []byte{
    decodedByteArray, _ :=  hex.DecodeString("8500000000000000680000000a0062496e76616c6964206f7220756e737570706f727465642070726f746f636f6c2076657273696f6e20283636293b20737570706f727465642076657273696f6e73206172652028332f76332c20342f76342c20352f76352c20362f76362d6265746129")
    return decodedByteArray
}

func returnUnsupported02() []byte{
    decodedByteArray, _ :=  hex.DecodeString("8500000000000000680000000a0062496e76616c6964206f7220756e737570706f727465642070726f746f636f6c2076657273696f6e20283636293b20737570706f727465642076657273696f6e73206172652028332f76332c20342f76342c20352f76352c20362f76362d6265746129")
    return decodedByteArray
}

func returnServerMeta() []byte{
    decodedByteArray, _ :=  hex.DecodeString("8500000006000000660003001150524f544f434f4c5f56455253494f4e5300040004332f76330004342f76340004352f76350009362f76362d62657461000b434f4d5052455353494f4e00020006736e6170707900036c7a34000b43514c5f56455253494f4e00010005332e342e35")
    return decodedByteArray
}

func returnStartup() []byte{
    decodedByteArray, _ :=  hex.DecodeString("850000010200000000")
    return decodedByteArray
}

func returnMsg01() []byte{
    decodedByteArray, _ :=  hex.DecodeString("090002a4c8c185000002020000000044fb95d4")
    return decodedByteArray
}

func returnPeerV2() []byte{
    decodedByteArray, _ :=  hex.DecodeString("d6000277bd3f8500000308000000cd00000002000000010000000c000673797374656d000870656572735f763200047065657200100009706565725f706f72740009000b646174615f63656e746572000d0007686f73745f6964000c000e6e61746976655f616464726573730010000b6e61746976655f706f72740009000c7072656665727265645f69700010000e7072656665727265645f706f7274000900047261636b000d000f72656c656173655f76657273696f6e000d000e736368656d615f76657273696f6e000c0006746f6b656e730022000d000000007e715c6a")
    return decodedByteArray
}

func returnLocal() []byte{
    decodedByteArray, _ :=  hex.DecodeString("2804028ed4bb85000004080000041f000000020000000100000014000673797374656d00056c6f63616c00036b6579000d000c626f6f747374726170706564000d001162726f6164636173745f616464726573730010000e62726f6164636173745f706f72740009000c636c75737465725f6e616d65000d000b63716c5f76657273696f6e000d000b646174615f63656e746572000d0011676f737369705f67656e65726174696f6e00090007686f73745f6964000c000e6c697374656e5f616464726573730010000b6c697374656e5f706f7274000900176e61746976655f70726f746f636f6c5f76657273696f6e000d000b706172746974696f6e6572000d00047261636b000d000f72656c656173655f76657273696f6e000d000b7270635f61646472657373001000087270635f706f72740009000e736368656d615f76657273696f6e000c0006746f6b656e730022000d000c7472756e63617465645f61740021000c000300000001000000056c6f63616c00000009434f4d504c4554454400000004c0a8016a0000000400001b580000000c5465737420436c757374657200000005332e342e350000000b6461746163656e74657231000000046131f7b9000000100d36697ecd054a0b8b36f8e8fc1d884200000004c0a8016a0000000400001b5800000001350000002b6f72672e6170616368652e63617373616e6472612e6468742e4d75726d757233506172746974696f6e6572000000057261636b3100000005342e302e3000000004000000000000000400002352000000102207c2a9f5983971986b2926e09e239d0000017900000010000000142d31343137323236323136393233393636343138000000142d32323631363032383634343135313238313734000000142d33333532383939393735313037333430343334000000142d34343438393431383735333432323738373233000000142d36323237353736313736313635383635383532000000142d37323035373231363238363831393739393237000000142d38383234373434393234333739303233323235000000133130393435343236323933313836343936333300000013323131323339343033343036313037363634380000001132333034323031323937383631363832300000001332393330303635393636363733343237343433000000133430303538323138343738393930303639383200000013343734313733373534303630393930303632330000001335383231313130333438303537353032323831000000133730363436333033353337333830393432353300000013383033313939343539363831363632393735320000005c0000000200000010176c39cdb93d33a5a2188eb06a56f66e000000140000017bab2f7afe0000001c0000017bab2fa72c00000010618f817b005f3678b8a453f3930b8e86000000140000017bab2f7afe0000001c0000017bab2fa4282301fbf6")
    return decodedByteArray
}

// 05
func returnKeyspaces() []byte{
    decodedByteArray, _ :=  hex.DecodeString("6202021aa9a5850000050800000259000000020000000100000003000d73797374656d5f736368656d6100096b6579737061636573000d6b657973706163655f6e616d65000d000e64757261626c655f7772697465730004000b7265706c69636174696f6e0021000d000d000000050000000b73797374656d5f617574680000000101000000570000000200000005636c6173730000002b6f72672e6170616368652e63617373616e6472612e6c6f6361746f722e53696d706c655374726174656779000000127265706c69636174696f6e5f666163746f7200000001310000000d73797374656d5f736368656d6100000001010000003b0000000100000005636c6173730000002a6f72672e6170616368652e63617373616e6472612e6c6f6361746f722e4c6f63616c53747261746567790000001273797374656d5f64697374726962757465640000000101000000570000000200000005636c6173730000002b6f72672e6170616368652e63617373616e6472612e6c6f6361746f722e53696d706c655374726174656779000000127265706c69636174696f6e5f666163746f7200000001330000000673797374656d00000001010000003b0000000100000005636c6173730000002a6f72672e6170616368652e63617373616e6472612e6c6f6361746f722e4c6f63616c53747261746567790000000d73797374656d5f7472616365730000000101000000570000000200000005636c6173730000002b6f72672e6170616368652e63617373616e6472612e6c6f6361746f722e53696d706c655374726174656779000000127265706c69636174696f6e5f66616374 6f7200000001329b3170e0")
    return decodedByteArray
}

// 06
func returnTables() []byte{
    decodedByteArray, _ :=  hex.DecodeString("")
    return decodedByteArray
}

// 07
func returnColumns() []byte{
    decodedByteArray, _ :=  hex.DecodeString("")
    return decodedByteArray
}

// 08
func returnTypes() []byte{
    decodedByteArray, _ :=  hex.DecodeString("")
    return decodedByteArray
}

// 09
func returnFunctions() []byte{
    decodedByteArray, _ :=  hex.DecodeString("")
    return decodedByteArray
}

// 10
func returnAggregates() []byte{
    decodedByteArray, _ :=  hex.DecodeString("")
    return decodedByteArray
}

// 11
func returnTriggers() []byte{
    decodedByteArray, _ :=  hex.DecodeString("")
    return decodedByteArray
}

// 12
func returnIndexes() []byte{
    decodedByteArray, _ :=  hex.DecodeString("")
    return decodedByteArray
}

// 13
func returnViews() []byte{
    decodedByteArray, _ :=  hex.DecodeString("")
    return decodedByteArray
}
