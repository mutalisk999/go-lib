package serialization

import (
	"fmt"
	"io"
	"os"
	"testing"
)

type TestType struct {
	a int
	b string
}

func (c TestType) Pack(writer io.Writer) error {
	Pack(writer, c.a)
	PackString(writer, c.b)
	return nil
}

func TestSerialize(t *testing.T) {
	file, err := os.Create("test_pack.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	Pack(file, int(0x12345678))
	Pack(file, byte(0x01))
	Pack(file, int16(0x1234))
	Pack(file, int64(0x1234567887654321))
	Pack(file, int(-1))
	Pack(file, int8(-1))
	Pack(file, int16(-1))
	Pack(file, int64(-1))
	Pack(file, float32(1.1))
	Pack(file, float64(2.2))

	a := [5]int{1, 2, 3, 4, 5}
	m := [3]map[string]int{{"a1": 1, "b2": 10, "c3": 2}, {}, {}}
	Pack(file, a)
	Pack(file, m)

	c := TestType{a: 0x12345678, b: "abcdefg"}
	c.Pack(file)
}

func TestUnSerialize(t *testing.T) {
	file, err := os.Open("test_pack.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var a byte
	UnPack(file, &a)
}
