package groot

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"time"
)

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

// unzip_root_buffer implements the ROOT unzip algorithm
func unzip_root_buffer(src []byte) (buf []byte, err error) {
	const HDRSIZE = 9
	const DEFLATE = 8

	buf = make([]byte, 0)

	// check header

	if len(src) < HDRSIZE {
		return buf, fmt.Errorf("groot.utils.unzip: too small source")
	}

	if src[0] != byte('Z') || src[1] != byte('L') || src[2] != DEFLATE {
		return buf, fmt.Errorf("groot.utils.unzip: error in header: %v",
			src[:3])
	}

	rbuf := src[HDRSIZE:]
	dec, err := zlib.NewReader(bytes.NewBuffer(rbuf))
	if err != nil {
		return []byte{}, err
	}
	buf, err = ioutil.ReadAll(dec)
	if err != nil {
		return []byte{}, err
	}
	err = dec.Close()
	if err != nil {
		return []byte{}, err
	}
	return buf, err
}

// EOF
