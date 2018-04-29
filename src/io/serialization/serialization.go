package serialization

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"unsafe"
)

func PackByte(writer io.Writer, argByte byte) {
	bytes := []byte{argByte}
	writer.Write(bytes)
}

func PackShort(writer io.Writer, argShort int16) {
	bytes := []byte{byte((argShort >> 8) & 0xFF), byte((argShort & 0x00FF) >> 0)}
	writer.Write(bytes)
}

func PackUShort(writer io.Writer, argUShort uint16) {
	bytes := []byte{byte((argUShort & 0xFF00) >> 8), byte((argUShort & 0x00FF) >> 0)}
	writer.Write(bytes)
}

func PackInt(writer io.Writer, argInt int) {
	bytes := []byte{
		byte((argInt >> 24) & 0xFF), byte((argInt & 0x00FF0000) >> 16),
		byte((argInt & 0x0000FF00) >> 8), byte((argInt & 0x000000FF) >> 0)}
	writer.Write(bytes)
}

func PackUInt(writer io.Writer, argUInt uint) {
	bytes := []byte{
		byte((argUInt & 0xFF000000) >> 24), byte((argUInt & 0x00FF0000) >> 16),
		byte((argUInt & 0x0000FF00) >> 8), byte((argUInt & 0x000000FF) >> 0)}
	writer.Write(bytes)
}

func PackLong(writer io.Writer, argLong int64) {
	bytes := []byte{
		byte((argLong >> 56) & 0xFF), byte((argLong & 0x00FF000000000000) >> 48),
		byte((argLong & 0x0000FF0000000000) >> 40), byte((argLong & 0x000000FF00000000) >> 32),
		byte((argLong & 0x00000000FF000000) >> 24), byte((argLong & 0x0000000000FF0000) >> 16),
		byte((argLong & 0x000000000000FF00) >> 8), byte((argLong & 0x00000000000000FF) >> 0)}
	writer.Write(bytes)
}

func PackULong(writer io.Writer, argULong uint64) {
	bytes := []byte{
		byte((argULong & 0xFF00000000000000) >> 56), byte((argULong & 0x00FF000000000000) >> 48),
		byte((argULong & 0x0000FF0000000000) >> 40), byte((argULong & 0x000000FF00000000) >> 32),
		byte((argULong & 0x00000000FF000000) >> 24), byte((argULong & 0x0000000000FF0000) >> 16),
		byte((argULong & 0x000000000000FF00) >> 8), byte((argULong & 0x00000000000000FF) >> 0)}
	writer.Write(bytes)
}

func PackString(writer io.Writer, argString string) {
	PackInt(writer, len(argString))
	bytes := []byte(argString)
	writer.Write(bytes)
}

func PackFloat(writer io.Writer, argFloat float32) {
	unsafePtr := uintptr(unsafe.Pointer(&argFloat))
	bytes := []byte{
		*(*byte)(unsafe.Pointer(unsafePtr)),
		*(*byte)(unsafe.Pointer(unsafePtr + uintptr(1))),
		*(*byte)(unsafe.Pointer(unsafePtr + uintptr(2))),
		*(*byte)(unsafe.Pointer(unsafePtr + uintptr(3)))}
	writer.Write(bytes)
}

func PackDouble(writer io.Writer, argDouble float64) {
	unsafePtr := uintptr(unsafe.Pointer(&argDouble))
	bytes := []byte{
		*(*byte)(unsafe.Pointer(unsafePtr)),
		*(*byte)(unsafe.Pointer(unsafePtr + uintptr(1))),
		*(*byte)(unsafe.Pointer(unsafePtr + uintptr(2))),
		*(*byte)(unsafe.Pointer(unsafePtr + uintptr(3))),
		*(*byte)(unsafe.Pointer(unsafePtr + uintptr(4))),
		*(*byte)(unsafe.Pointer(unsafePtr + uintptr(5))),
		*(*byte)(unsafe.Pointer(unsafePtr + uintptr(6))),
		*(*byte)(unsafe.Pointer(unsafePtr + uintptr(7)))}
	writer.Write(bytes)
}

func Pack(writer io.Writer, argPack interface{}) error {
	typeStr := reflect.TypeOf(argPack).String()

	if typeStr == "byte" {
		argByte := argPack.(byte)
		PackByte(writer, argByte)
	} else if typeStr == "int8" {
		argByte := byte(argPack.(int8))
		PackByte(writer, argByte)
	} else if typeStr == "uint8" {
		argByte := byte(argPack.(uint8))
		PackByte(writer, argByte)
	} else if typeStr == "int16" {
		argShort := argPack.(int16)
		PackShort(writer, argShort)
	} else if typeStr == "uint16" {
		argUShort := argPack.(uint16)
		PackUShort(writer, argUShort)
	} else if typeStr == "int" {
		argInt := argPack.(int)
		PackInt(writer, argInt)
	} else if typeStr == "uint" {
		argUInt := argPack.(uint)
		PackUInt(writer, argUInt)
	} else if typeStr == "int32" {
		argInt := int(argPack.(int32))
		PackInt(writer, argInt)
	} else if typeStr == "uint32" {
		argUInt := uint(argPack.(uint32))
		PackUInt(writer, argUInt)
	} else if typeStr == "int64" {
		argLong := argPack.(int64)
		PackLong(writer, argLong)
	} else if typeStr == "uint64" {
		argULong := argPack.(uint64)
		PackULong(writer, argULong)
	} else if typeStr == "string" {
		argString := argPack.(string)
		PackString(writer, argString)
	} else if typeStr == "float32" {
		argFloat := argPack.(float32)
		PackFloat(writer, argFloat)
	} else if typeStr == "float64" {
		argDouble := argPack.(float64)
		PackDouble(writer, argDouble)
	} else if typeStr[0] == '[' {
		argArray := reflect.ValueOf(argPack).Convert(reflect.TypeOf(argPack))
		PackInt(writer, argArray.Len())
		for i := 0; i < argArray.Len(); i++ {
			err := Pack(writer, argArray.Index(i).Interface())
			if err != nil {
				return errors.New("fail to pack element of array")
			}
		}
	} else if strings.Contains(typeStr, "map[") && typeStr[0:4] == "map[" {
		argMap := reflect.ValueOf(argPack).Convert(reflect.TypeOf(argPack))
		keys := argMap.MapKeys()
		PackInt(writer, len(keys))
		for i := 0; i < len(keys); i++ {
			err := Pack(writer, keys[i].Interface())
			if err != nil {
				return errors.New("fail to pack key of map")
			}
			err = Pack(writer, argMap.MapIndex(keys[i]).Interface())
			if err != nil {
				return errors.New("fail to pack value of map")
			}
		}
	} else {
		argObj := reflect.ValueOf(argPack).Convert(reflect.TypeOf(argPack))
		methodPack := argObj.MethodByName("Pack")

		errValue := methodPack.Call([]reflect.Value{reflect.ValueOf(writer).Convert(reflect.TypeOf(writer))})
		err := errValue[0].Interface().(error)
		if err != nil {
			return errors.New("fail to pack object of " + typeStr)
		}
	}

	return nil
}

func UnPackByte(reader io.Reader) byte {
	var bytes [1]byte
	reader.Read(bytes[0:1])
	return bytes[0]
}

func UnPackChar(reader io.Reader) int8 {
	var bytes [1]byte
	reader.Read(bytes[0:1])
	return int8(bytes[0])
}

func UnPackUChar(reader io.Reader) uint8 {
	var bytes [1]byte
	reader.Read(bytes[0:1])
	return uint8(bytes[0])
}

func UnPackShort(reader io.Reader) int16 {
	var bytes [2]byte
	reader.Read(bytes[0:2])
	return int16(int16(bytes[0])<<8 | int16(bytes[1]))
}

func UnPackUShort(reader io.Reader) uint16 {
	var bytes [2]byte
	reader.Read(bytes[0:2])
	return uint16(uint16(bytes[0])<<8 | uint16(bytes[1]))
}

func UnPackInt(reader io.Reader) int {
	var bytes [4]byte
	reader.Read(bytes[0:4])
	return int(int(bytes[0])<<24 | int(bytes[1])<<16 | int(bytes[2])<<8 | int(bytes[3]))
}

func UnPackUInt(reader io.Reader) uint {
	var bytes [4]byte
	reader.Read(bytes[0:4])
	return uint(uint(bytes[0])<<24 | uint(bytes[1])<<16 | uint(bytes[2])<<8 | uint(bytes[3]))
}

func UnPackLong(reader io.Reader) int64 {
	var bytes [8]byte
	reader.Read(bytes[0:8])
	return int64(int64(bytes[0])<<56 | int64(bytes[1])<<48 | int64(bytes[2])<<40 | int64(bytes[3])<<32 | int64(bytes[4])<<24 | int64(bytes[5])<<16 | int64(bytes[6])<<8 | int64(bytes[7]))
}

func UnPackULong(reader io.Reader) uint64 {
	var bytes [8]byte
	reader.Read(bytes[0:8])
	return uint64(uint64(bytes[0])<<56 | uint64(bytes[1])<<48 | uint64(bytes[2])<<40 | uint64(bytes[3])<<32 | uint64(bytes[4])<<24 | uint64(bytes[5])<<16 | uint64(bytes[6])<<8 | uint64(bytes[7]))
}

func UnPackString(reader io.Reader) string {
	strLength := UnPackInt(reader)
	bytes := make([]byte, strLength)
	reader.Read(bytes[0:strLength])
	return string(bytes)
}

func UnPackFloat(reader io.Reader) float32 {
	var bytes [4]byte
	reader.Read(bytes[0:4])
	var valFloat float32
	unsafePtr := uintptr(unsafe.Pointer(&valFloat))
	*(*byte)(unsafe.Pointer(unsafePtr)) = bytes[0]
	*(*byte)(unsafe.Pointer(unsafePtr + uintptr(1))) = bytes[1]
	*(*byte)(unsafe.Pointer(unsafePtr + uintptr(2))) = bytes[2]
	*(*byte)(unsafe.Pointer(unsafePtr + uintptr(3))) = bytes[3]
	return valFloat
}

func UnPackDouble(reader io.Reader) float64 {
	var bytes [8]byte
	reader.Read(bytes[0:8])
	var valDouble float64
	unsafePtr := uintptr(unsafe.Pointer(&valDouble))
	*(*byte)(unsafe.Pointer(unsafePtr)) = bytes[0]
	*(*byte)(unsafe.Pointer(unsafePtr + uintptr(1))) = bytes[1]
	*(*byte)(unsafe.Pointer(unsafePtr + uintptr(2))) = bytes[2]
	*(*byte)(unsafe.Pointer(unsafePtr + uintptr(3))) = bytes[3]
	*(*byte)(unsafe.Pointer(unsafePtr + uintptr(4))) = bytes[4]
	*(*byte)(unsafe.Pointer(unsafePtr + uintptr(5))) = bytes[5]
	*(*byte)(unsafe.Pointer(unsafePtr + uintptr(6))) = bytes[6]
	*(*byte)(unsafe.Pointer(unsafePtr + uintptr(7))) = bytes[7]
	return valDouble
}

func UnPack(reader io.Reader, argExample interface{}) (interface{}, error) {
	typeStr := reflect.TypeOf(argExample).String()

	if typeStr == "byte" {
		argByte := UnPackByte(reader)
		return argByte, nil
	} else if typeStr == "int8" {
		argChar := UnPackChar(reader)
		return argChar, nil
	} else if typeStr == "uint8" {
		argUChar := UnPackUChar(reader)
		return argUChar, nil
	} else if typeStr == "int16" {
		argShort := UnPackShort(reader)
		return argShort, nil
	} else if typeStr == "uint16" {
		argUShort := UnPackUShort(reader)
		return argUShort, nil
	} else if typeStr == "int" {
		argInt := UnPackInt(reader)
		return argInt, nil
	} else if typeStr == "uint" {
		argUInt := UnPackUInt(reader)
		return argUInt, nil
	} else if typeStr == "int32" {
		argInt := UnPackInt(reader)
		return int32(argInt), nil
	} else if typeStr == "uint32" {
		argUInt := UnPackUInt(reader)
		return uint32(argUInt), nil
	} else if typeStr == "int64" {
		argLong := UnPackLong(reader)
		return argLong, nil
	} else if typeStr == "uint64" {
		argULong := UnPackULong(reader)
		return argULong, nil
	} else if typeStr == "string" {
		argString := UnPackString(reader)
		return argString, nil
	} else if typeStr == "float32" {
		argFloat := UnPackFloat(reader)
		return argFloat, nil
	} else if typeStr == "float64" {
		argDouble := UnPackDouble(reader)
		return argDouble, nil
	} else if typeStr[0] == '[' {
		var argArray []interface{}
		argArrayExample := reflect.ValueOf(argExample).Convert(reflect.TypeOf(argExample))
		if argArrayExample.Len() == 0 {
			return nil, errors.New("invalid example with nil array")
		}
		arraySize := UnPackInt(reader)
		for i := 0; i < arraySize; i++ {
			argElement, err := UnPack(reader, argArrayExample.Index(0).Interface())
			if err != nil {
				return nil, err
			} else {
				argArray = append(argArray, argElement)
			}
		}
		return argArray, nil
	} else if strings.Contains(typeStr, "map[") && typeStr[0:4] == "map[" {
		argMap := make(map[interface{}]interface{})
		argMapExample := reflect.ValueOf(argExample).Convert(reflect.TypeOf(argExample))
		if argMapExample.Len() == 0 {
			return nil, errors.New("invalid example with nil map")
		}
		exampleKeys := argMapExample.MapKeys()
		mapSize := UnPackInt(reader)
		for i := 0; i < mapSize; i++ {
			argKey, err := UnPack(reader, exampleKeys[0].Interface())
			if err != nil {
				return nil, err
			}
			argValue, err := UnPack(reader, argMapExample.MapIndex(exampleKeys[0]).Interface())
			if err != nil {
				return nil, err
			}
			argMap[argKey] = argValue
		}
		return argMap, nil
	} else {
		refObj := reflect.ValueOf(argExample).Convert(reflect.TypeOf(argExample))
		methodUnPack := refObj.MethodByName("UnPack")

		retValue := methodUnPack.Call([]reflect.Value{reflect.ValueOf(reader).Convert(reflect.TypeOf(reader))})
		err := retValue[1].Interface()
		if err != nil {
			return nil, errors.New("fail to pack object of " + typeStr)
		}
		return retValue[0].Interface(), nil
	}

	return nil, nil
}
