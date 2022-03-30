package serialization

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

type TestStruct struct {
	a int
	b []string
	c map[string]int
}

func (g TestStruct) Pack(writer io.Writer) error {
	_ = Pack(writer, g.a)
	_ = Pack(writer, g.b)
	_ = Pack(writer, g.c)
	return nil
}

func (g TestStruct) UnPack(reader io.Reader) (TestStruct, error) {
	a, _ := UnPack(reader, reflect.TypeOf(g.a))
	b, _ := UnPack(reader, reflect.TypeOf(g.b))
	c, _ := UnPack(reader, reflect.TypeOf(g.c))
	g.a = a.(int)
	g.b = b.([]string)
	g.c = c.(map[string]int)
	return g, nil
}

func TestSerialize(t *testing.T) {
	file, err := os.Create("test_pack.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_ = Pack(file, int(0x12345678))
	_ = Pack(file, true)
	_ = Pack(file, int16(0x1234))
	_ = Pack(file, int64(0x1234567887654321))
	_ = Pack(file, int(-1))
	_ = Pack(file, int8(-1))
	_ = Pack(file, int16(-1))
	_ = Pack(file, int64(-1))
	_ = Pack(file, float32(1.1))
	_ = Pack(file, float64(2.2))

	a := [5]int{1, 2, 3, 4, 5}
	s := []string{"11", "22", "33", "44", "55"}
	m := [3]map[string]int{{"a1": 1, "b2": 10, "c3": 2}, {}, {}}
	mm := []map[string]int{{"a11": 11, "b22": 1010, "c33": 22}, {}, {}, {}, {"a111": 111, "b222": 101010, "c333": 222}}
	_ = Pack(file, a)
	_ = Pack(file, s)
	_ = Pack(file, m)
	_ = Pack(file, mm)

	c := TestStruct{a: 0x12345678, b: []string{"abcdefg", "hijklmn"}, c: map[string]int{"a": 1, "b": 2}}
	_ = Pack(file, c)
}

func TestUnSerialize(t *testing.T) {
	file, err := os.Open("test_pack.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var it interface{}
	var i int
	it, _ = UnPack(file, reflect.TypeOf(i))
	i = it.(int)
	fmt.Printf("%x\n", i)

	var b bool
	it, _ = UnPack(file, reflect.TypeOf(b))
	b = it.(bool)
	fmt.Println(b)

	var i16 int16
	it, _ = UnPack(file, reflect.TypeOf(i16))
	i16 = it.(int16)
	fmt.Printf("%x\n", i16)

	var i64 int64
	it, _ = UnPack(file, reflect.TypeOf(i64))
	i64 = it.(int64)
	fmt.Printf("%x\n", i64)

	var ni int
	it, _ = UnPack(file, reflect.TypeOf(ni))
	ni = it.(int)
	fmt.Println(ni)

	var ni8 int8
	it, _ = UnPack(file, reflect.TypeOf(ni8))
	ni8 = it.(int8)
	fmt.Println(ni8)

	var ni16 int16
	it, _ = UnPack(file, reflect.TypeOf(ni16))
	ni16 = it.(int16)
	fmt.Println(ni16)

	var ni64 int64
	it, _ = UnPack(file, reflect.TypeOf(ni64))
	ni64 = it.(int64)
	fmt.Println(ni64)

	var f32 float32
	it, _ = UnPack(file, reflect.TypeOf(f32))
	f32 = it.(float32)
	fmt.Println(f32)

	var f64 float64
	it, _ = UnPack(file, reflect.TypeOf(f64))
	f64 = it.(float64)
	fmt.Println(f64)

	var a [5]int
	it, _ = UnPack(file, reflect.TypeOf(a))
	fmt.Println(it)

	var s []string
	it, _ = UnPack(file, reflect.TypeOf(s))
	fmt.Println(it)

	var m [3]map[string]int
	it, _ = UnPack(file, reflect.TypeOf(m))
	fmt.Println(it)

	var mm []map[string]int
	it, _ = UnPack(file, reflect.TypeOf(mm))
	fmt.Println(it)

	c := TestStruct{}
	it, _ = UnPack(file, reflect.TypeOf(c))
	c = it.(TestStruct)
	fmt.Println(c)
}
