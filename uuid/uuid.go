package uuid // import "github.com/webnice/log/v2/uuid"

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

type UUID [16]byte

var (
	hardwareAddr   []byte
	clockSeq       uint32
	testingModeOne bool // set true from test program
	testingModeTwo bool // set true from test program
)

const (
	VariantNCSCompat = 0
	VariantIETF      = 2
	VariantMicrosoft = 6
	VariantFuture    = 7
)

func init() { initialize() }

func initialize() {
	if interfaces, err := net.Interfaces(); err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagLoopback == 0 && len(i.HardwareAddr) > 0 {
				hardwareAddr = i.HardwareAddr
				break
			}
		}
	}
	if testingModeOne {
		hardwareAddr = nil
	}
	if hardwareAddr == nil {
		// If we failed to obtain the MAC address of the current computer,
		// we will use a randomly generated 6 byte sequence instead and set
		// the multicast bit as recommended in RFC 4122.
		hardwareAddr = make([]byte, 6)
		_, err := io.ReadFull(rand.Reader, hardwareAddr)
		if testingModeTwo {
			err = errors.New("io.ReadFull() error residing from rand.Reader")
		}
		if err != nil {
			panic(err)
		}
		hardwareAddr[0] = hardwareAddr[0] | 0x01
	}
	// initialize the clock sequence with a random number
	var clockSeqRand [2]byte
	_, _ = io.ReadFull(rand.Reader, clockSeqRand[:])
	clockSeq = uint32(clockSeqRand[1])<<8 | uint32(clockSeqRand[0])
}

// ParseUUID parses a 32 digit hexadecimal number (that might contain hypens)
// representing an UUID.
func ParseUUID(input string) (u UUID, err error) {
	var (
		j int
		r rune
	)

	for _, r = range input {
		switch {
		case r == '-' && j&1 == 0:
			continue
		case r >= '0' && r <= '9' && j < 32:
			u[j/2] |= byte(r-'0') << uint(4-j&1*4)
		case r >= 'a' && r <= 'f' && j < 32:
			u[j/2] |= byte(r-'a'+10) << uint(4-j&1*4)
		case r >= 'A' && r <= 'F' && j < 32:
			u[j/2] |= byte(r-'A'+10) << uint(4-j&1*4)
		default:
			return UUID{}, fmt.Errorf("invalid UUID %q", input)
		}
		j += 1
	}
	if j != 32 {
		return UUID{}, fmt.Errorf("invalid UUID %q", input)
	}

	return
}

// UUIDFromBytes converts a raw byte slice to an UUID.
func UUIDFromBytes(input []byte) (u UUID, err error) {
	const err16byte = `UUIDs must be exactly 16 bytes long`

	if len(input) != 16 {
		return u, errors.New(err16byte)
	}
	copy(u[:], input)

	return
}

// RandomUUID generates a totally random UUID (version 4) as described in
// RFC 4122.
func RandomUUID() (u UUID, err error) {
	_, err = io.ReadFull(rand.Reader, u[:])
	if testingModeOne {
		err = errors.New("reader closed")
	}
	if err != nil {
		return
	}
	u[6] &= 0x0F // clear version
	u[6] |= 0x40 // set version to 4 (random uuid)
	u[8] &= 0x3F // clear variant
	u[8] |= 0x80 // set to IETF variant

	return
}

var timeBase = time.Date(1582, time.October, 15, 0, 0, 0, 0, time.UTC).Unix()

// TimeUUID generates a new time based UUID (version 1) using the current
// time as the timestamp.
func TimeUUID() UUID {
	return UUIDFromTime(time.Now())
}

// UUIDFromTime generates a new time based UUID (version 1) as described in
// RFC 4122. This UUID contains the MAC address of the node that generated
// the UUID, the given timestamp and a sequence number.
func UUIDFromTime(aTime time.Time) (u UUID) {
	var (
		utcTime time.Time
		t       uint64
		clock   uint32
	)

	utcTime = aTime.In(time.UTC)
	t = uint64(utcTime.Unix()-timeBase)*10000000 + uint64(utcTime.Nanosecond()/100)
	u[0], u[1], u[2], u[3] = byte(t>>24), byte(t>>16), byte(t>>8), byte(t)
	u[4], u[5] = byte(t>>40), byte(t>>32)
	u[6], u[7] = byte(t>>56)&0x0F, byte(t>>48)
	//
	clock = atomic.AddUint32(&clockSeq, 1)
	u[8] = byte(clock >> 8)
	u[9] = byte(clock)
	//
	copy(u[10:], hardwareAddr)
	//
	u[6] |= 0x10 // set version to 1 (time based uuid)
	u[8] &= 0x3F // clear variant
	u[8] |= 0x80 // set to IETF variant

	return u
}

// String returns the UUID in it's canonical form, a 32 digit hexadecimal
// number in the form of xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u UUID) String() (ret string) {
	const hexString = `0123456789abcdef`
	var (
		offsets [16]int
		r       []byte
		n       int
		b       byte
	)

	offsets = [...]int{0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34}
	r = make([]byte, 36)
	for n, b = range u {
		r[offsets[n]] = hexString[b>>4]
		r[offsets[n]+1] = hexString[b&0xF]
	}
	r[8], r[13], r[18], r[23] = '-', '-', '-', '-'
	ret = string(r)

	return

}

// Bytes returns the raw byte slice for this UUID. A UUID is always 128 bits
// (16 bytes) long.
func (u UUID) Bytes() []byte { return u[:] }

// Variant returns the variant of this UUID. This package will only generate
// UUIDs in the IETF variant.
func (u UUID) Variant() int {
	var x = u[8]

	if x&0x80 == 0 {
		return VariantNCSCompat
	}
	if x&0x40 == 0 {
		return VariantIETF
	}
	if x&0x20 == 0 {
		return VariantMicrosoft
	}

	return VariantFuture
}

// Version extracts the version of this UUID variant. The RFC 4122 describes
// five kinds of UUIDs.
func (u UUID) Version() int {
	return int(u[6] & 0xF0 >> 4)
}

// Node extracts the MAC address of the node who generated this UUID. It will
// return nil if the UUID is not a time based UUID (version 1).
func (u UUID) Node() []byte {
	if u.Version() != 1 {
		return nil
	}
	return u[10:]
}

// Timestamp extracts the timestamp information from a time based UUID
// (version 1).
func (u UUID) Timestamp() int64 {
	if u.Version() != 1 {
		return 0
	}
	return int64(uint64(u[0])<<24|uint64(u[1])<<16|
		uint64(u[2])<<8|uint64(u[3])) +
		int64(uint64(u[4])<<40|uint64(u[5])<<32) +
		int64(uint64(u[6]&0x0F)<<56|uint64(u[7])<<48)
}

// Time is like Timestamp, except that it returns a time.Time.
func (u UUID) Time() (ret time.Time) {
	var (
		t, sec, nsec int64
	)
	if u.Version() != 1 {
		return time.Time{}
	}
	t = u.Timestamp()
	sec = t / 1e7
	nsec = (t % 1e7) * 100
	ret = time.Unix(sec+timeBase, nsec).UTC()

	return
}

// Marshaling for JSON
func (u UUID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.String() + `"`), nil
}

// Unmarshalling for JSON
func (u *UUID) UnmarshalJSON(data []byte) (err error) {
	var (
		str    string
		parsed UUID
	)

	str = strings.Trim(string(data), `"`)
	if len(str) > 36 {
		err = fmt.Errorf("invalid JSON UUID %s", str)
		return
	}
	if parsed, err = ParseUUID(str); err == nil {
		copy(u[:], parsed[:])
	}

	return
}
