package websocket

type ByteData struct {
	pos        int
	bytes      []byte
	startOfMsg int
}

func NewByteData(c int) *ByteData {
	return &ByteData{
		pos:        0,
		bytes:      make([]byte, c),
		startOfMsg: 0,
	}
}

func (bd *ByteData) length() int {
	return len(bd.bytes)
}

func (bd *ByteData) markStartOfMsg() {
	bd.startOfMsg = bd.pos
	bd.pos += 2
}

func (bd *ByteData) markEndOfMsg() {
	length := bd.pos - bd.startOfMsg - 2
	bd.bytes[bd.startOfMsg] = byte((length >> 8) & 255)
	bd.bytes[bd.startOfMsg+1] = byte(length & 255)
}

func (bd *ByteData) clear() {
	bd.pos = 0
}

func (bd *ByteData) getPosition() int {
	return bd.pos
}

func (bd *ByteData) getBytes() []byte {
	return bd.bytes
}

func (bd *ByteData) appendByte(d byte) {
	bd.bytes[bd.pos] = d
	bd.pos++
}

func (bd *ByteData) appendByteAtPos(e int, d byte) {
	bd.bytes[e] = d
}

func (bd *ByteData) appendChar(d byte) {
	bd.bytes[bd.pos] = d
	bd.pos++
}

func (bd *ByteData) appendCharAtPos(e int, d byte) {
	bd.bytes[e] = d
}

func (bd *ByteData) appendShort(d int16) {
	bd.bytes[bd.pos] = byte((d >> 8) & 255)
	bd.pos++
	bd.bytes[bd.pos] = byte(d & 255)
	bd.pos++
}

func (bd *ByteData) appendInt(d int32) {
	bd.bytes[bd.pos] = byte((d >> 24) & 255)
	bd.pos++
	bd.bytes[bd.pos] = byte((d >> 16) & 255)
	bd.pos++
	bd.bytes[bd.pos] = byte((d >> 8) & 255)
	bd.pos++
	bd.bytes[bd.pos] = byte(d & 255)
	bd.pos++
}

func (bd *ByteData) appendLong(d int64) {
	bd.bytes[bd.pos] = byte((d >> 56) & 255)
	bd.pos++
	bd.bytes[bd.pos] = byte((d >> 48) & 255)
	bd.pos++
	bd.bytes[bd.pos] = byte((d >> 40) & 255)
	bd.pos++
	bd.bytes[bd.pos] = byte((d >> 32) & 255)
	bd.pos++
	bd.bytes[bd.pos] = byte((d >> 24) & 255)
	bd.pos++
	bd.bytes[bd.pos] = byte((d >> 16) & 255)
	bd.pos++
	bd.bytes[bd.pos] = byte((d >> 8) & 255)
	bd.pos++
	bd.bytes[bd.pos] = byte(d & 255)
	bd.pos++
}

func (bd *ByteData) appendLongAsBigInt(e int64) {
	d := e
	bd.bytes = append(bd.bytes, byte((d>>56)&255))
	bd.bytes = append(bd.bytes, byte((d>>48)&255))
	bd.bytes = append(bd.bytes, byte((d>>40)&255))
	bd.bytes = append(bd.bytes, byte((d>>32)&255))
	bd.bytes = append(bd.bytes, byte((d>>24)&255))
	bd.bytes = append(bd.bytes, byte((d>>16)&255))
	bd.bytes = append(bd.bytes, byte((d>>8)&255))
	bd.bytes = append(bd.bytes, byte(d&255))
}

func (bd *ByteData) appendString(d string) {
	for i := 0; i < len(d); i++ {
		bd.bytes[bd.pos] = d[i]
		bd.pos++
	}
}

func (bd *ByteData) appendByteArray(d []byte) {
	for i := 0; i < len(d); i++ {
		bd.bytes[bd.pos] = d[i]
		bd.pos++
	}
}

func (bd *ByteData) appendByteArr(e []byte, d int) {
	for i := 0; i < d; i++ {
		bd.bytes[bd.pos] = e[i]
		bd.pos++
	}
}
