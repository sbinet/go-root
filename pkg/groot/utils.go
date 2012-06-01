package groot

import (
	"encoding/binary"
	"io"
	"os"
	"time"
)

type breader struct {
	order binary.ByteOrder
}

func (b breader) skip(r io.Seeker, nbytes int64) {
	_, err := r.Seek(nbytes, os.SEEK_CUR)
	if err != nil {
		panic(err)
	}
}

func (b breader) ntoi2(r io.Reader) (o int16) {
	err := binary.Read(r, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b breader) ntoi4(r io.Reader) (o int32) {
	err := binary.Read(r, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b breader) ntoi8(r io.Reader) (o int64) {
	err := binary.Read(r, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b breader) ntobyte(r io.Reader) (o byte) {
	err := binary.Read(r, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b breader) ntou2(r io.Reader) (o uint16) {
	err := binary.Read(r, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b breader) ntou4(r io.Reader) (o uint32) {
	err := binary.Read(r, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b breader) ntou8(r io.Reader) (o uint64) {
	err := binary.Read(r, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b breader) ntof(r io.Reader) (o float32) {
	err := binary.Read(r, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b breader) ntod(r io.Reader) (o float64) {
	err := binary.Read(r, b.order, &o)
	if err != nil {
		panic(err)
	}
	return
}

func (b breader) readArrayF(r io.Reader) (o []float32) {
	n := int(b.ntou4(r))
	o = make([]float32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntof(r)
	}
	return
}

func (b breader) readArrayD(r io.Reader) (o []float64) {
	n := int(b.ntou4(r))
	o = make([]float64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntod(r)
	}
	return
}

func (b breader) readArrayS(r io.Reader) (o []int16) {
	n := int(b.ntou4(r))
	o = make([]int16, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi2(r)
	}
	return
}

func (b breader) readArrayI(r io.Reader) (o []int32) {
	n := int(b.ntou4(r))
	o = make([]int32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi4(r)
	}
	return
}

func (b breader) readArrayL(r io.Reader) (o []int64) {
	n := int(b.ntou4(r))
	o = make([]int64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi8(r)
	}
	return
}

func (b breader) readArrayC(r io.Reader) (o []byte) {
	n := int(b.ntou4(r))
	o = make([]byte, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntobyte(r)
	}
	return
}

func (b breader) readStaticArray(r io.Reader) (o []uint32) {
	n := int(b.ntou4(r))
	o = make([]uint32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntou4(r)
	}
	return
}

func (b breader) readFastArrayF(r io.Reader, n int) (o []float32) {
	o = make([]float32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntof(r)
	}
	return
}

func (b breader) readFastArrayD(r io.Reader, n int) (o []float64) {
	o = make([]float64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntod(r)
	}
	return
}

func (b breader) readFastArrayS(r io.Reader, n int) (o []int16) {
	o = make([]int16, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi2(r)
	}
	return
}

func (b breader) readFastArrayI(r io.Reader, n int) (o []int32) {
	o = make([]int32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi4(r)
	}
	return
}

func (b breader) readFastArrayL(r io.Reader, n int) (o []int64) {
	o = make([]int64, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntoi8(r)
	}
	return
}

func (b breader) readFastArrayC(r io.Reader, n int) (o []byte) {
	o = make([]byte, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntobyte(r)
	}
	return
}

func (b breader) readFastArrayTString(r io.Reader, n int) (o []string) {
	o = make([]string, n)
	for i := 0; i < n; i++ {
		o[i] = b.readTString(r)
	}
	return
}

func (b breader) readFastArray(r io.Reader, n int) (o []uint32) {
	o = make([]uint32, n)
	for i := 0; i < n; i++ {
		o[i] = b.ntou4(r)
	}
	return
}

func (b breader) readTString(r io.Reader) string {
	n := int(b.ntobyte(r))
	if n == 255 {
		// large string
		n = int(b.ntou4(r))
	}
	if n == 0 {
		return ""
	}
	v := b.ntobyte(r)
	if v == 0 {
		return ""
	}
	o := make([]byte, n)
	o[0] = v
	err := binary.Read(r, b.order, o[1:])
	if err != nil {
		panic(err)
	}
	return string(o)
}

//FIXME
// readBasicPointerElem
// readBasicPointer

func (b breader) readString(r io.Reader, max int) string {
	o := []byte{}
	n := 0
	var v byte
	for {
		v = b.ntobyte(r)
		if v == 0 {
			break
		}
		n += 1
		if max > 0 && n >= max {
			break
		}
		o = append(o, v)
	}
	return string(o)
}

//FIXME
// readObjectAny
// readTList
// readTObjArray
// readTClonesArray
// readTCollection
// readTHashList
// readTNamed
// readTCanvas

//FIXME
// getStreamer(f TFile, name string) Streamer

// datime2time converts a uint32 holding a ROOT's TDatime into a time.Time
func datime2time(d uint32) time.Time {

	// ROOT's TDatime begins in January 1995...
	var year uint32 = (d >> 26) + 1995
	var month uint32 = (d << 6) >> 28
	var day uint32 = (d << 10) >> 27
	var hour uint32 = (d << 15) >> 27
	var min uint32 = (d << 20) >> 26
	var sec uint32 = (d << 26) >> 26
	nsec := 0
	return time.Date(int(year), time.Month(month), int(day),
		int(hour), int(min), int(sec), nsec, time.UTC)
}

// EOF
