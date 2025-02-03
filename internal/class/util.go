package class

import (
	"bufio"
	"encoding/binary"
	"io"
)

func readUint8(reader *bufio.Reader) (uint8, error) {
	byte, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return byte, nil
}

func readUint16(reader *bufio.Reader) (uint16, error) {
	data := make([]byte, 2)
	_, err := io.ReadFull(reader, data)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(data), nil
}

func readUint32(reader *bufio.Reader) (uint32, error) {
	data := make([]byte, 4)
	_, err := io.ReadFull(reader, data)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(data), nil
}
