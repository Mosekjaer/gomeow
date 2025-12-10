package memory

import (
	"encoding/binary"
	"fmt"
	"math"

	"gomeow/process"
	"gomeow/vec"
)

// ReadInt8 reads a signed 8-bit integer
func ReadInt8(p *process.Process, address uintptr) (int8, error) {
	buf := make([]byte, 1)
	if err := Read(p, address, buf); err != nil {
		return 0, err
	}
	return int8(buf[0]), nil
}

// ReadInt16 reads a signed 16-bit integer
func ReadInt16(p *process.Process, address uintptr) (int16, error) {
	buf := make([]byte, 2)
	if err := Read(p, address, buf); err != nil {
		return 0, err
	}
	return int16(binary.LittleEndian.Uint16(buf)), nil
}

// ReadInt32 reads a signed 32-bit integer
func ReadInt32(p *process.Process, address uintptr) (int32, error) {
	buf := make([]byte, 4)
	if err := Read(p, address, buf); err != nil {
		return 0, err
	}
	return int32(binary.LittleEndian.Uint32(buf)), nil
}

// ReadInt64 reads a signed 64-bit integer
func ReadInt64(p *process.Process, address uintptr) (int64, error) {
	buf := make([]byte, 8)
	if err := Read(p, address, buf); err != nil {
		return 0, err
	}
	return int64(binary.LittleEndian.Uint64(buf)), nil
}

// ReadUint8 reads an unsigned 8-bit integer
func ReadUint8(p *process.Process, address uintptr) (uint8, error) {
	buf := make([]byte, 1)
	if err := Read(p, address, buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}

// ReadUint16 reads an unsigned 16-bit integer
func ReadUint16(p *process.Process, address uintptr) (uint16, error) {
	buf := make([]byte, 2)
	if err := Read(p, address, buf); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf), nil
}

// ReadUint32 reads an unsigned 32-bit integer
func ReadUint32(p *process.Process, address uintptr) (uint32, error) {
	buf := make([]byte, 4)
	if err := Read(p, address, buf); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf), nil
}

// ReadUint64 reads an unsigned 64-bit integer
func ReadUint64(p *process.Process, address uintptr) (uint64, error) {
	buf := make([]byte, 8)
	if err := Read(p, address, buf); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf), nil
}

// ReadFloat32 reads a 32-bit float
func ReadFloat32(p *process.Process, address uintptr) (float32, error) {
	buf := make([]byte, 4)
	if err := Read(p, address, buf); err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint32(buf)
	return math.Float32frombits(bits), nil
}

// ReadFloat64 reads a 64-bit float
func ReadFloat64(p *process.Process, address uintptr) (float64, error) {
	buf := make([]byte, 8)
	if err := Read(p, address, buf); err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(bits), nil
}

// ReadBool reads a boolean value
func ReadBool(p *process.Process, address uintptr) (bool, error) {
	buf := make([]byte, 1)
	if err := Read(p, address, buf); err != nil {
		return false, err
	}
	return buf[0] != 0, nil
}

// ReadString reads a null-terminated string
func ReadString(p *process.Process, address uintptr, maxLen int) (string, error) {
	if maxLen <= 0 {
		maxLen = 64
	}
	buf := make([]byte, maxLen)
	if err := Read(p, address, buf); err != nil {
		return "", err
	}
	// Find null terminator
	for i, b := range buf {
		if b == 0 {
			return string(buf[:i]), nil
		}
	}
	return string(buf), nil
}

// ReadBytes reads a byte slice of the given length
func ReadBytes(p *process.Process, address uintptr, size int) ([]byte, error) {
	buf := make([]byte, size)
	if err := Read(p, address, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// ReadVec2 reads a 2D vector (two float32s)
func ReadVec2(p *process.Process, address uintptr) (vec.Vec2, error) {
	buf := make([]byte, 8)
	if err := Read(p, address, buf); err != nil {
		return vec.Vec2{}, err
	}
	return vec.Vec2{
		X: math.Float32frombits(binary.LittleEndian.Uint32(buf[0:4])),
		Y: math.Float32frombits(binary.LittleEndian.Uint32(buf[4:8])),
	}, nil
}

// ReadVec3 reads a 3D vector (three float32s)
func ReadVec3(p *process.Process, address uintptr) (vec.Vec3, error) {
	buf := make([]byte, 12)
	if err := Read(p, address, buf); err != nil {
		return vec.Vec3{}, err
	}
	return vec.Vec3{
		X: math.Float32frombits(binary.LittleEndian.Uint32(buf[0:4])),
		Y: math.Float32frombits(binary.LittleEndian.Uint32(buf[4:8])),
		Z: math.Float32frombits(binary.LittleEndian.Uint32(buf[8:12])),
	}, nil
}

// WriteInt8 writes a signed 8-bit integer
func WriteInt8(p *process.Process, address uintptr, value int8) error {
	return Write(p, address, []byte{byte(value)})
}

// WriteInt16 writes a signed 16-bit integer
func WriteInt16(p *process.Process, address uintptr, value int16) error {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(value))
	return Write(p, address, buf)
}

// WriteInt32 writes a signed 32-bit integer
func WriteInt32(p *process.Process, address uintptr, value int32) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(value))
	return Write(p, address, buf)
}

// WriteInt64 writes a signed 64-bit integer
func WriteInt64(p *process.Process, address uintptr, value int64) error {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(value))
	return Write(p, address, buf)
}

// WriteUint8 writes an unsigned 8-bit integer
func WriteUint8(p *process.Process, address uintptr, value uint8) error {
	return Write(p, address, []byte{value})
}

// WriteUint16 writes an unsigned 16-bit integer
func WriteUint16(p *process.Process, address uintptr, value uint16) error {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, value)
	return Write(p, address, buf)
}

// WriteUint32 writes an unsigned 32-bit integer
func WriteUint32(p *process.Process, address uintptr, value uint32) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, value)
	return Write(p, address, buf)
}

// WriteUint64 writes an unsigned 64-bit integer
func WriteUint64(p *process.Process, address uintptr, value uint64) error {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, value)
	return Write(p, address, buf)
}

// WriteFloat32 writes a 32-bit float
func WriteFloat32(p *process.Process, address uintptr, value float32) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, math.Float32bits(value))
	return Write(p, address, buf)
}

// WriteFloat64 writes a 64-bit float
func WriteFloat64(p *process.Process, address uintptr, value float64) error {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(value))
	return Write(p, address, buf)
}

// WriteBool writes a boolean value
func WriteBool(p *process.Process, address uintptr, value bool) error {
	var b byte
	if value {
		b = 1
	}
	return Write(p, address, []byte{b})
}

// WriteString writes a null-terminated string
func WriteString(p *process.Process, address uintptr, value string) error {
	buf := append([]byte(value), 0)
	return Write(p, address, buf)
}

// WriteVec2 writes a 2D vector
func WriteVec2(p *process.Process, address uintptr, value vec.Vec2) error {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint32(buf[0:4], math.Float32bits(value.X))
	binary.LittleEndian.PutUint32(buf[4:8], math.Float32bits(value.Y))
	return Write(p, address, buf)
}

// WriteVec3 writes a 3D vector
func WriteVec3(p *process.Process, address uintptr, value vec.Vec3) error {
	buf := make([]byte, 12)
	binary.LittleEndian.PutUint32(buf[0:4], math.Float32bits(value.X))
	binary.LittleEndian.PutUint32(buf[4:8], math.Float32bits(value.Y))
	binary.LittleEndian.PutUint32(buf[8:12], math.Float32bits(value.Z))
	return Write(p, address, buf)
}

// PointerChain32 follows a chain of 32-bit pointers and returns the final address
func PointerChain32(p *process.Process, base uintptr, offsets ...uintptr) (uintptr, error) {
	if len(offsets) == 0 {
		return base, nil
	}

	addr, err := ReadUint32(p, base)
	if err != nil {
		return 0, fmt.Errorf("failed to read base pointer: %v", err)
	}

	current := uintptr(addr)
	for i := 0; i < len(offsets)-1; i++ {
		addr, err := ReadUint32(p, current+offsets[i])
		if err != nil {
			return 0, fmt.Errorf("failed to read pointer at offset %d: %v", i, err)
		}
		current = uintptr(addr)
	}

	return current + offsets[len(offsets)-1], nil
}

// PointerChain64 follows a chain of 64-bit pointers and returns the final address
func PointerChain64(p *process.Process, base uintptr, offsets ...uintptr) (uintptr, error) {
	if len(offsets) == 0 {
		return base, nil
	}

	addr, err := ReadUint64(p, base)
	if err != nil {
		return 0, fmt.Errorf("failed to read base pointer: %v", err)
	}

	current := uintptr(addr)
	for i := 0; i < len(offsets)-1; i++ {
		addr, err := ReadUint64(p, current+offsets[i])
		if err != nil {
			return 0, fmt.Errorf("failed to read pointer at offset %d: %v", i, err)
		}
		current = uintptr(addr)
	}

	return current + offsets[len(offsets)-1], nil
}
