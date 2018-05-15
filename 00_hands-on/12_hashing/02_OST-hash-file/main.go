package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const (
	ChunkSize = 65536 // 64k
)

func main() {
	file, err := os.Open("../Supernatural.S13E21.720p.HDTV.x264-SVA[eztv].mkv")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	st, _ := file.Stat()
	hsh, _ := HashFile(file)
	fmt.Printf("Hash found: %x\n", hsh)
	fmt.Println("Size file: ", st.Size())
}

// Generate an OSDB hash for an *os.File.
func HashFile(file *os.File) (hash uint64, err error) {
	fi, err := file.Stat()
	if err != nil {
		return
	}
	if fi.Size() < ChunkSize {
		return 0, fmt.Errorf("File is too small")
	}

	// Read head and tail blocks.
	buf := make([]byte, ChunkSize*2)
	err = readChunk(file, 0, buf[:ChunkSize])
	if err != nil {
		return
	}
	err = readChunk(file, fi.Size()-ChunkSize, buf[ChunkSize:])
	if err != nil {
		return
	}

	// Convert to uint64, and sum.
	var nums [(ChunkSize * 2) / 8]uint64
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &nums)
	if err != nil {
		return 0, err
	}
	for _, num := range nums {
		hash += num
	}

	return hash + uint64(fi.Size()), nil
}

// Read a chunk of a file at `offset` so as to fill `buf`.
func readChunk(file *os.File, offset int64, buf []byte) (err error) {
	n, err := file.ReadAt(buf, offset)
	if err != nil {
		return
	}
	if n != ChunkSize {
		return fmt.Errorf("Invalid read %v", n)
	}
	return
}
