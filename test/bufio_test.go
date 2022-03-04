package test

import (
	"bufio"
	"fmt"
	"testing"
)

type Writer int

func (*Writer) Write(p []byte) (n int, err error) {
	fmt.Printf("%q\n", p)
	return len(p), nil
}

type R struct{}

func (r *R) Read(p []byte) (n int, err error) {
	fmt.Println("Read")
	str := "12345678901234567890"
	copy(p, str)
	return len(str), nil
}

func TestBufio(t *testing.T) {
	// fmt.Println("Unbuffered I/O")
	// w := new(Writer)
	// w.Write([]byte{'a'})
	// w.Write([]byte{'b'})
	// w.Write([]byte{'c'})
	// w.Write([]byte{'d'})

	// fmt.Println("Buffered I/O")
	// bw := bufio.NewWriterSize(w, 3)
	// bw.Write([]byte{'a'})
	// bw.Write([]byte{'b'})
	// bw.Write([]byte{'c'})
	// bw.Write([]byte{'d'})
	// err := bw.Flush()
	// if err != nil {
	// 	panic(err)
	// }

	// {
	// 	w := new(Writer)
	// 	bw := bufio.NewWriterSize(w, 10)
	// 	bw.Write([]byte("abcd"))
	// 	bw.Flush()
	// }

	{
		r := new(R)
		br := bufio.NewReaderSize(r, 100)
		buf1 := make([]byte, 2)
		_, err := br.Read(buf1[:])
		if err != nil {
			panic(err)
		}
		buf2 := make([]byte, 4)
		_, err = br.Read(buf2[:])
		if err != nil {
			panic(err)
		}

		buf3 := make([]byte, 2)
		_, err = br.Read(buf3[:])
		if err != nil {
			panic(err)
		}
	}
}
