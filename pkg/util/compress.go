package main

import (
    "fmt"
    "github.com/pierrec/lz4/v4"
    "bytes"
    "io"
    "os"
    "strings"
)

var fileContent = `CompressBlock compresses the source buffer starting at soffet into the destination one.
This is the fast version of LZ4 compression and also the default one.
The size of the compressed data is returned. If it is 0 and no error, then the data is incompressible.
An error is returned if the destination buffer is too small.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
aaaaaaaaaaaa`

var compressedCode = `000000020000000100000001001573797374656d5f7669727475616c5f736368656d6100096b6579737061636573000d6b657973706163655f6e616d65000d000000020000000c73797374656d5f76696577730000001573797374656d5f7669727475616c5f736368656d61c2061a20`

func main () {
    // Compress and uncompress an input string.
    s := "hello world"
    r := strings.NewReader(s)

    // The pipe will uncompress the data from the writer.
    pr, pw := io.Pipe()
    zw := lz4.NewWriter(pw)
    zr := lz4.NewReader(pr)

    go func() {
        // Compress the input string.
        _, _ = io.Copy(zw, r)
        _ = zw.Close() // Make sure the writer is closed
        _ = pw.Close() // Terminate the pipe
    }()

    _, _ = io.Copy(os.Stdout, zr)

    compressed, _ := compress([]byte(fileContent))
    fmt.Printf("Compressed data is %d and %d \n", len(compressed), len(fileContent))
    fmt.Printf("The compressed value is %x \n", compressed)

    testCode, err := decompress([]byte(compressedCode))

    if err != nil {
        fmt.Printf("The error is %v \n", err)
        return
    }
    fmt.Printf("The uncompressed value is %s", testCode)
}

func compress(in []byte) ([]byte, error) {
    r := bytes.NewReader(in)
    w := &bytes.Buffer{}
    zw := lz4.NewWriter(w)
    _, err := io.Copy(zw, r)
    if err != nil {
        return nil, err
    }
    // Closing is *very* important
    if err := zw.Close(); err != nil {
        return nil, err
    }
    return w.Bytes(), nil

//    bToCompress := []byte(toCompress)
//    compressed := make([]byte, len(bToCompress))
//
//    //compress
//    l, err := lz4.CompressBlock(bToCompress, compressed , 0)
//    if err != nil {
//        panic(err)
//    }
//    return compressed
//    // fmt.Println("compressed Data:", string(compressed[:l]))
}

func decompress(in []byte) ([]byte, error) {
    r := bytes.NewReader(in)
    w := &bytes.Buffer{}
    zr := lz4.NewReader(r)
    _, err := io.Copy(w, zr)
    if err != nil {
        return nil, err
    }
    return w.Bytes(), nil
}

// Output:
// hello world



