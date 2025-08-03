package section

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type section1 struct {
	length                    uint32
	sectionNumber             uint8
	originatingCenter         uint16
	originatingSubcenter      uint16
	masterTablesVersion       uint8
	localTablesVersion        uint8
	referenceTimeSignificance uint8
	year                      uint16
	month                     uint8
	day                       uint8
	hour                      uint8
	minute                    uint8
	second                    uint8
	productionStatus          uint8
	productType               uint8
	reserved                  []byte
}

var _ Section1 = (*section1)(nil)

func (s *section1) Length() uint32 {
	return s.length
}

func (s *section1) SectionNumber() uint8 {
	return s.sectionNumber
}

func (s *section1) OriginatingCenter() uint16 {
	return s.originatingCenter
}

func (s *section1) OriginatingSubcenter() uint16 {
	return s.originatingSubcenter
}

func (s *section1) MasterTablesVersion() uint8 {
	return s.masterTablesVersion
}

func (s *section1) LocalTablesVersion() uint8 {
	return s.localTablesVersion
}

func (s *section1) ReferenceTimeSignificance() uint8 {
	return s.referenceTimeSignificance
}

func (s *section1) Year() uint16 {
	return s.year
}

func (s *section1) Month() uint8 {
	return s.month
}

func (s *section1) Day() uint8 {
	return s.day
}

func (s *section1) Hour() uint8 {
	return s.hour
}

func (s *section1) Minute() uint8 {
	return s.minute
}

func (s *section1) Second() uint8 {
	return s.second
}

func (s *section1) ProductionStatus() uint8 {
	return s.productionStatus
}

func (s *section1) DataType() uint8 {
	return s.productType
}

func NewSection1FromBytes(data []byte, keepReserved bool) (Section1, error) {
	if len(data) < 21 {
		return nil, fmt.Errorf("section1: data too short")
	}

	br := bytes.NewReader(data)

	var s section1
	var err error
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.length))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.sectionNumber))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.originatingCenter))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.originatingSubcenter))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.masterTablesVersion))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.localTablesVersion))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.referenceTimeSignificance))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.year))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.month))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.day))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.hour))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.minute))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.second))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.productionStatus))
	err = errors.Join(err, binary.Read(br, binary.BigEndian, &s.productType))

	reservedN := s.length - 21
	reservedWriter := io.Discard
	if keepReserved {
		s.reserved = make([]byte, reservedN)
		reservedWriter = bytes.NewBuffer(s.reserved)
	}
	if _, err := io.CopyN(reservedWriter, br, int64(reservedN)); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &s, nil
}
