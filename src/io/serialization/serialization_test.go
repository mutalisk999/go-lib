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
	Pack(writer, c.b)
	return nil
}

func (c TestType) UnPack(reader io.Reader) (TestType, error) {
	var it interface{}
	it, _ = UnPack(reader, c.a)
	c.a = it.(int)
	it, _ = UnPack(reader, c.b)
	c.b = it.(string)
	return c, nil
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

	var it interface{}
	var i int
	it, _ = UnPack(file, i)
	i = it.(int)
	fmt.Println(i)

	var b byte
	it, _ = UnPack(file, b)
	b = it.(byte)
	fmt.Println(b)

	var i16 int16
	it, _ = UnPack(file, i16)
	i16 = it.(int16)
	fmt.Println(i16)

	var i64 int64
	it, _ = UnPack(file, i64)
	i64 = it.(int64)
	fmt.Println(i64)

	var ni int
	it, _ = UnPack(file, ni)
	ni = it.(int)
	fmt.Println(ni)

	var ni8 int8
	it, _ = UnPack(file, ni8)
	ni8 = it.(int8)
	fmt.Println(ni8)

	var ni16 int16
	it, _ = UnPack(file, ni16)
	ni16 = it.(int16)
	fmt.Println(ni16)

	var ni64 int64
	it, _ = UnPack(file, ni64)
	ni64 = it.(int64)
	fmt.Println(ni64)

	var f32 float32
	it, _ = UnPack(file, f32)
	f32 = it.(float32)
	fmt.Println(f32)

	var f64 float64
	it, _ = UnPack(file, f64)
	f64 = it.(float64)
	fmt.Println(f64)

	a := []int{1}
	it, _ = UnPack(file, a)
	fmt.Println(it)

	m := []map[string]int{{"a1": 1}}
	it, _ = UnPack(file, m)
	fmt.Println(it)

	c := TestType{a: 0, b: ""}
	it, _ = UnPack(file, c)
	c = it.(TestType)
	fmt.Println(c)
}
