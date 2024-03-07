package server

import (
    "net"
    "fmt"
    "bytes"
    "encoding/binary"
//    "errors"
)


type Header struct {
    version  byte    // One byte: 0x03 / 0x83
    flag     uint16  // Two bytes: flag
    streamId int     // One byte: 
    opCode   byte    // It's only for request from client
    length   uint32  // Four bytes: 
}

func (h *Header) ToString() string {
    return fmt.Sprintf("version: %x , flag:, %d, stream id: %d, op code: %x, length: %d", h.version, h.flag, h.streamId, h.opCode, h.length )
}

type ACPacket struct {
    header Header
    buffer []byte
    pos    int
    length int
}

func NewACPacket(conn net.Conn) (*ACPacket, error) {
    msg := make([]byte, 4096)
    readLen, err := conn.Read(msg)
    if err != nil {
        return nil, err
    }

    return &ACPacket{
        buffer: msg,
        length: readLen,
        pos   : 0,
    }, nil
}

// Read the header into Header
func (a *ACPacket) DecodeHeader() error {
    fmt.Printf("The buffer: %#v \n\n", bytes.Trim(a.buffer, "\x00") )
    a.header.version  = a.buffer[a.pos]
    a.pos += 1
    a.header.flag     = binary.BigEndian.Uint16(a.buffer[a.pos:a.pos+2])
    a.pos += 2
    a.header.streamId = int(a.buffer[a.pos])
    a.pos += 1
    a.header.opCode   = a.buffer[a.pos]
    a.pos += 1
    a.header.length   = binary.BigEndian.Uint32(a.buffer[a.pos:a.pos+4])
    a.pos += 4
    return nil
}

type FrameBody interface {
    Decode([]byte) error
    Encode() ([]byte, error)
}


func (a *ACPacket) DecodeBody() (*FrameBody, error) {
    frameDataType := OPCODEMAP[a.header.opCode]
    if frameDataType == "" {
        fmt.Printf("Does not find data type \n")
        return nil, nil
    } 

    switch frameDataType {
        case FRAME_DTYPE_STRING_MAP:
            fmt.Printf("Starting to parse data string map data type \n")
            stringMap := NewStringMap(&a.header)
            err := stringMap.Decode(a.buffer[a.pos:a.header.length])
	    fmt.Printf("--------------------  <%#v> \n", stringMap)
            if err != nil {
                return nil, nil
            }
	    // fmt.Printf("The message body: %#v \n", "To add")
        default:
            fmt.Printf(fmt.Sprintf("Not found data type, %s ", frameDataType))
        return nil, nil
//        return nil, errors.New("Unexpected OP Code")
    }
    return nil, nil
}

func (a *ACPacket) GetBuffer() []byte {
    return a.buffer
}

func (a *ACPacket) GetLength() int {
    return a.length
}

func (a *ACPacket) GetHeader() *Header {
    return &a.header
}

// ---------- String Map ----------
type StringMap struct {
    Header
    values map[string]string
    keys []string
}

func NewStringMap(header *Header) *StringMap {
    return &StringMap{
        Header: *header,
    values: make(map[string]string),
    }
}

func(m *StringMap)Decode(buffer []byte)error {
    // 01. Read number of key value
    pos := 0
    numKV := int(binary.BigEndian.Uint16(buffer[pos:pos+2]))

    pos += 2
    for idx:= 0; idx < numKV; idx++ {
	    // 01. Read Key
	    key, err := decodeString(buffer, &pos)
	    if err != nil {
		    return nil
	    }
	    // 02. Read Value
	    value, err := decodeString(buffer, &pos)
	    if err != nil {
		    return nil
	    }

	    m.keys = append(m.keys, key)
	    m.values[key] = value
    }
    return nil
}

func(m *StringMap)Encode() ([]byte, error){
    return nil, nil
} 

func decodeString(buffer []byte, pos *int) (ret string, err error){
    lenStr := int(binary.BigEndian.Uint16(buffer[*pos:*pos+2]))
    *pos = *pos + 2

    ret = string(buffer[*pos:*pos+lenStr])
    *pos += lenStr

    return
} 
