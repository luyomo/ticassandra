package server

// https://cassandra.apache.org/_/native_protocol.html
const (
    OPCODE_ERROR          = 0x00
    OPCODE_STARTUP        = 0x01
    OPCODE_READY          = 0x02
    OPCODE_AUTHENTICATE   = 0x03
    OPCODE_OPTIONS        = 0x05
    OPCODE_SUPPORTED      = 0x06
    OPCODE_QUERY          = 0x07
    OPCODE_RESULT         = 0x08
    OPCODE_PREPARE        = 0x09
    OPCODE_EXECUTE        = 0x0A
    OPCODE_REGISTER       = 0x0B
    OPCODE_EVENT          = 0x0C
    OPCODE_BATCH          = 0x0D
    OPCODE_AUTH_CHALLENGE = 0x0E
    OPCODE_AUTH_RESPONSE  = 0x0F
    OPCODE_AUTH_SUCCESS   = 0x10
)



const (
    FRAME_DTYPE_INT             = "int"              // A 4 bytes integer
    FRAME_DTYPE_LONG            = "long"             // A 8 bytes integer
    FRAME_DTYPE_BYTE            = "byte"             // A 1 byte unsigned integer
    FRAME_DTYPE_SHORT           = "short"            // A 2 bytes unsigned integer
    FRAME_DTYPE_STRING          = "string"           // A "short" n, followed by n bytes representing an UTF-8 string.
    FRAME_DTYPE_LONG_STRING     = "long string"      // An "int" n, followed by n bytes representing an UTF-8 string.
    FRAME_DTYPE_UUID            = "uuid"             // A 16 bytes long uuid.
    FRAME_DTYPE_STRING_LIST     = "string list"      // A "short" n, followed by n "string".
    FRAME_DTYPE_BYTES           = "bytes"            // A "int" n, followed by n bytes if n >= 0. If n < 0,
                                                     // no byte should follow and the value represented is `null`.
    FRAME_DTYPE_VALUE           = "value"            // A "int" n, followed by n bytes if n >= 0.
                                                     // If n == -1 no byte should follow and the value represented is `null`.
                                                     // If n == -2 no byte should follow and the value represented is
                                                     // `not set` not resulting in any change to the existing value.
                                                     // n < -2 is an invalid value and results in an error.
    FRAME_DTYPE_SHORT_BYTES     = "short bytes"      // A "short" n, followed by n bytes if n >= 0.
    FRAME_DTYPE_UNISGHED_VINT   = "unsigned vint"    // An unsigned variable length integer. A vint is encoded with the most significant byte (MSB) first.
                                                     // The most significant byte will contains the information about how many extra bytes need to be read
                                                     // as well as the most significant bits of the integer.
                                                     // The number of extra bytes to read is encoded as 1 bits on the left side.
                                                     // For example, if we need to read 2 more bytes the first byte will start with 110
                                                     // (e.g. 256 000 will be encoded on 3 bytes as "110"00011 11101000 00000000)
                                                     // If the encoded integer is 8 bytes long the vint will be encoded on 9 bytes and the first
                                                     // byte will be: 11111111
   FRAME_DTYPE_VINT             = "vint"             // A signed variable length integer. This is encoded using zig-zag encoding and then sent
                                                     // like an "unsigned vint". Zig-zag encoding converts numbers as follows:
                                                     // 0 = 0, -1 = 1, 1 = 2, -2 = 3, 2 = 4, -3 = 5, 3 = 6 and so forth.
                                                     // The purpose is to send small negative values as small unsigned values, so that we save bytes on the wire.
                                                     // To encode a value n use "(n >> 31) ^ (n << 1)" for 32 bit values, and "(n >> 63) ^ (n << 1)"
                                                     // for 64 bit values where "^" is the xor operation, "<<" is the left shift operation and ">>" is
                                                     // the arithemtic right shift operation (highest-order bit is replicated).
                                                     // Decode with "(n >> 1) ^ -(n & 1)".
    FRAME_DTYPE_OPTION          = "option"           // A pair of <id><value> where <id> is a "short" representing
                                                     // the option id and <value> depends on that option (and can be
                                                     // of size 0). The supported id (and the corresponding <value>)
                                                     // will be described when this is used.
    FRAME_DTYPE_OPTION_LIST     = "option list"      // A "short" n, followed by n "option".
    FRAME_DTYPE_INET            = "inet"             // An address (ip and port) to a node. It consists of one
                                                     // "byte" n, that represents the address size, followed by n
                                                     // "byte" representing the IP address (in practice n can only be
                                                     // either 4 (IPv4) or 16 (IPv6)), following by one "int"
                                                     // representing the port.
    FRAME_DTYPE_INETADDR        = "inetaddr"         // An IP address (without a port) to a node. It consists of one
                                                     // "byte" n, that represents the address size, followed by n
                                                     // "byte" representing the IP address.
    FRAME_DTYPE_CONSISTENCY     = "consistency"      // A consistency level specification. This is a "short"
                                                     //  representing a consistency level with the following
                                                     //  correspondance:
                                                     //    0x0000    ANY
                                                     //    0x0001    ONE
                                                     //    0x0002    TWO
                                                     //    0x0003    THREE
                                                     //    0x0004    QUORUM
                                                     //    0x0005    ALL
                                                     //    0x0006    LOCAL_QUORUM
                                                     //    0x0007    EACH_QUORUM
                                                     //    0x0008    SERIAL
                                                     //    0x0009    LOCAL_SERIAL
                                                     //    0x000A    LOCAL_ONE
    FRAME_DTYPE_STRING_MAP      = "string map"       // A "short" n, followed by n pair <k><v> where <k> and <v> are "string".
    FRAME_DTYPE_STRING_MULTIMAP = "string multimap"  // A "short" n, followed by n pair <k><v> where <k> is a "string" and <v> is a "string list".
    FRAME_DTYPE_BYTES_MAP       = "bytes map"        // A "short" n, followed by n pair <k><v> where <k> is a
                                                     // "string" and <v> is a "bytes".
)

var OPCODEMAP = map[byte]string{
    OPCODE_SUPPORTED : FRAME_DTYPE_STRING_MULTIMAP ,
    OPCODE_STARTUP   : FRAME_DTYPE_STRING_MAP,
}
