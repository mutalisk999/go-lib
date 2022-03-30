package serialization

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

func PackBool(writer io.Writer, argBool bool) error {
	boolByte := uint8(0)
	if argBool {
		boolByte = uint8(1)
	}
	bytes := []byte{boolByte}
	_, err := writer.Write(bytes)
	return err
}

func PackInt8(writer io.Writer, argInt8 int8) error {
	_, err := writer.Write([]byte{uint8(argInt8)})
	return err
}

func PackUint8(writer io.Writer, argUint8 uint8) error {
	_, err := writer.Write([]byte{argUint8})
	return err
}

func PackInt16(writer io.Writer, argInt16 int16) error {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, uint16(argInt16))
	_, err := writer.Write(bytes)
	return err
}

func PackUint16(writer io.Writer, argUint16 uint16) error {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, argUint16)
	_, err := writer.Write(bytes)
	return err
}

func PackInt32(writer io.Writer, argInt32 int32) error {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(argInt32))
	_, err := writer.Write(bytes)
	return err
}

func PackUint32(writer io.Writer, argUint32 uint32) error {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, argUint32)
	_, err := writer.Write(bytes)
	return err
}

func PackInt64(writer io.Writer, argInt64 int64) error {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(argInt64))
	_, err := writer.Write(bytes)
	return err
}

func PackUint64(writer io.Writer, argUint64 uint64) error {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, argUint64)
	_, err := writer.Write(bytes)
	return err
}

func PackInt(writer io.Writer, argInt int) error {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(argInt))
	_, err := writer.Write(bytes)
	return err
}

func PackUint(writer io.Writer, argUint uint) error {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(argUint))
	_, err := writer.Write(bytes)
	return err
}

func PackFloat32(writer io.Writer, argFloat32 float32) error {
	argByteU32 := math.Float32bits(argFloat32)
	return PackUint32(writer, argByteU32)
}

func PackFloat64(writer io.Writer, argFloat64 float64) error {
	argByteU64 := math.Float64bits(argFloat64)
	return PackUint64(writer, argByteU64)
}

func PackString(writer io.Writer, argString string) error {
	strLen := uint32(len(argString))
	err := PackUint32(writer, strLen)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(argString))
	return err
}

func Pack(writer io.Writer, argPack interface{}) error {
	typeKind := reflect.TypeOf(argPack).Kind()
	typeValue := reflect.ValueOf(argPack)

	var err error = nil
	switch typeKind {
	case reflect.Bool:
		return PackBool(writer, typeValue.Bool())
	case reflect.Int8:
		return PackInt8(writer, int8(typeValue.Int()))
	case reflect.Uint8:
		return PackUint8(writer, uint8(typeValue.Uint()))
	case reflect.Int16:
		return PackInt16(writer, int16(typeValue.Int()))
	case reflect.Uint16:
		return PackUint16(writer, uint16(typeValue.Uint()))
	case reflect.Int32:
		return PackInt32(writer, int32(typeValue.Int()))
	case reflect.Uint32:
		return PackUint32(writer, uint32(typeValue.Uint()))
	case reflect.Int64:
		return PackInt64(writer, typeValue.Int())
	case reflect.Uint64:
		return PackUint64(writer, typeValue.Uint())
	case reflect.Int:
		return PackInt(writer, int(typeValue.Int()))
	case reflect.Uint:
		return PackUint(writer, uint(typeValue.Uint()))
	case reflect.Float32:
		return PackFloat32(writer, float32(typeValue.Float()))
	case reflect.Float64:
		return PackFloat64(writer, typeValue.Float())
	case reflect.String:
		return PackString(writer, typeValue.String())
	case reflect.Array:
		for i := 0; i < typeValue.Len(); i++ {
			err = Pack(writer, typeValue.Index(i).Interface())
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Slice:
		sliceLen := uint32(typeValue.Len())
		err = PackUint32(writer, sliceLen)
		if err != nil {
			return err
		}
		for i := 0; i < typeValue.Len(); i++ {
			err = Pack(writer, typeValue.Index(i).Interface())
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		mapLen := uint32(typeValue.Len())
		err = PackUint32(writer, mapLen)
		if err != nil {
			return err
		}
		keys := typeValue.MapKeys()
		for _, key := range keys {
			err = Pack(writer, key.Interface())
			if err != nil {
				return err
			}
			err = Pack(writer, typeValue.MapIndex(key).Interface())
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Struct:
		methodPack := typeValue.MethodByName("Pack")
		errValue := methodPack.Call([]reflect.Value{reflect.ValueOf(writer)})
		err := errValue[0].Interface()
		if err != nil {
			return err.(error)
		}
		return nil
	default:
		err = errors.New(fmt.Sprintf("Not Support Pack Type: %s", typeKind.String()))
		return err
	}
}

func UnPackBool(reader io.Reader) (bool, error) {
	var bytes [1]byte
	_, err := reader.Read(bytes[0:1])
	if err != nil {
		return false, err
	}
	if bytes[0] == uint8(1) {
		return true, nil
	} else if bytes[0] == uint8(0) {
		return false, nil
	} else {
		return false, errors.New("UnPackBool: Unexpected byte")
	}
}

func UnPackInt8(reader io.Reader) (int8, error) {
	var bytes [1]byte
	_, err := reader.Read(bytes[0:1])
	if err != nil {
		return 0, err
	}
	return int8(bytes[0]), nil
}

func UnPackUint8(reader io.Reader) (uint8, error) {
	var bytes [1]byte
	_, err := reader.Read(bytes[0:1])
	if err != nil {
		return 0, err
	}
	return bytes[0], nil
}

func UnPackInt16(reader io.Reader) (int16, error) {
	var bytes [2]byte
	_, err := reader.Read(bytes[0:2])
	if err != nil {
		return 0, err
	}
	return int16(binary.LittleEndian.Uint16(bytes[:])), nil
}

func UnPackUint16(reader io.Reader) (uint16, error) {
	var bytes [2]byte
	_, err := reader.Read(bytes[0:2])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(bytes[:]), nil
}

func UnPackInt32(reader io.Reader) (int32, error) {
	var bytes [4]byte
	_, err := reader.Read(bytes[0:4])
	if err != nil {
		return 0, err
	}
	return int32(binary.LittleEndian.Uint32(bytes[:])), nil
}

func UnPackUint32(reader io.Reader) (uint32, error) {
	var bytes [4]byte
	_, err := reader.Read(bytes[0:4])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(bytes[:]), nil
}

func UnPackInt64(reader io.Reader) (int64, error) {
	var bytes [8]byte
	_, err := reader.Read(bytes[0:8])
	if err != nil {
		return 0, err
	}
	return int64(binary.LittleEndian.Uint64(bytes[:])), nil
}

func UnPackUint64(reader io.Reader) (uint64, error) {
	var bytes [8]byte
	_, err := reader.Read(bytes[0:8])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(bytes[:]), nil
}

func UnPackInt(reader io.Reader) (int, error) {
	var bytes [8]byte
	_, err := reader.Read(bytes[0:8])
	if err != nil {
		return 0, err
	}
	return int(binary.LittleEndian.Uint64(bytes[:])), nil
}

func UnPackUint(reader io.Reader) (uint, error) {
	var bytes [8]byte
	_, err := reader.Read(bytes[0:8])
	if err != nil {
		return 0, err
	}
	return uint(binary.LittleEndian.Uint64(bytes[:])), nil
}

func UnPackFloat32(reader io.Reader) (float32, error) {
	var bytes [4]byte
	_, err := reader.Read(bytes[0:4])
	if err != nil {
		return 0.0, err
	}
	bits := binary.LittleEndian.Uint32(bytes[:])
	return math.Float32frombits(bits), nil
}

func UnPackFloat64(reader io.Reader) (float64, error) {
	var bytes [8]byte
	_, err := reader.Read(bytes[0:8])
	if err != nil {
		return 0.0, err
	}
	bits := binary.LittleEndian.Uint64(bytes[:])
	return math.Float64frombits(bits), nil
}

func UnPackString(reader io.Reader) (string, error) {
	strLen, err := UnPackUint32(reader)
	if err != nil {
		return "", err
	}
	bytes := make([]byte, strLen)
	_, err = reader.Read(bytes[0:strLen])
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func UnPack(reader io.Reader, ty reflect.Type) (interface{}, error) {
	var err error = nil
	switch ty.Kind() {
	case reflect.Bool:
		return UnPackBool(reader)
	case reflect.Int8:
		return UnPackInt8(reader)
	case reflect.Uint8:
		return UnPackUint8(reader)
	case reflect.Int16:
		return UnPackInt16(reader)
	case reflect.Uint16:
		return UnPackUint16(reader)
	case reflect.Int32:
		return UnPackInt32(reader)
	case reflect.Uint32:
		return UnPackUint32(reader)
	case reflect.Int64:
		return UnPackInt64(reader)
	case reflect.Uint64:
		return UnPackUint64(reader)
	case reflect.Int:
		return UnPackInt(reader)
	case reflect.Uint:
		return UnPackUint(reader)
	case reflect.Float32:
		return UnPackFloat32(reader)
	case reflect.Float64:
		return UnPackFloat64(reader)
	case reflect.String:
		return UnPackString(reader)
	case reflect.Array:
		arrayLen := ty.Len()
		array := reflect.New(ty).Elem()
		for i := 0; i < arrayLen; i++ {
			val, err := UnPack(reader, ty.Elem())
			if err != nil {
				return nil, err
			}
			array.Index(i).Set(reflect.ValueOf(val))
		}
		return array.Interface(), nil
	case reflect.Slice:
		sliceLen, err := UnPackUint32(reader)
		if err != nil {
			return nil, err
		}
		slice := reflect.MakeSlice(ty, 0, int(sliceLen))
		for i := 0; i < int(sliceLen); i++ {
			val, err := UnPack(reader, ty.Elem())
			if err != nil {
				return nil, err
			}
			slice = reflect.Append(slice, reflect.ValueOf(val))
		}
		return slice.Interface(), nil
	case reflect.Map:
		mapLen, err := UnPackUint32(reader)
		if err != nil {
			return nil, err
		}
		m := reflect.MakeMap(ty)
		for i := 0; i < int(mapLen); i++ {
			k, err := UnPack(reader, ty.Key())
			if err != nil {
				return nil, err
			}
			v, err := UnPack(reader, ty.Elem())
			if err != nil {
				return nil, err
			}
			m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
		}
		return m.Interface(), nil
	case reflect.Struct:
		refObj := reflect.New(ty)
		methodUnPack := refObj.MethodByName("UnPack")

		retValue := methodUnPack.Call([]reflect.Value{reflect.ValueOf(reader)})
		err := retValue[1].Interface()
		if err != nil {
			return nil, err.(error)
		}
		return retValue[0].Interface(), nil
	default:
		err = errors.New(fmt.Sprintf("Not Support UnPack Type: %s", ty.String()))
		return nil, err
	}
}
