package websocket

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"bytes"
)

type TopicData struct {
	feedType           string
	exchange           string
	symbol             string
	tSymbol            string
	multiplier         int
	precision          int
	precisionValue     int
	jsonArray          interface{}
	fieldDataArray     [100]interface{}
	updatedFieldsArray [100]interface{}
}

func NewTopicData(feedType string) *TopicData {
	t := &TopicData{
		feedType:       feedType,
		multiplier:     1,
		precision:      2,
		precisionValue: 100,
	}
	t.fieldDataArray[StringIndex["NAME"]] = feedType
	return t
}

func (t *TopicData) getKey() string {
	return fmt.Sprintf("%s|%s", t.exchange, t.symbol)
}

func (t *TopicData) setLongValues(indexVal int, value int64) {
	if t.fieldDataArray[indexVal] != value && value != int64(TrashVal) {
		t.fieldDataArray[indexVal] = value
		t.updatedFieldsArray[indexVal] = true
	}
}

func (t *TopicData) prepareCommonData() {
	t.updatedFieldsArray[StringIndex["NAME"]] = true
	t.updatedFieldsArray[StringIndex["EXCHG"]] = true
	t.updatedFieldsArray[StringIndex["SYMBOL"]] = true
}

func (t *TopicData) setStringValues(e int, d string) {
	if e == StringIndex["SYMBOL"] {
		t.symbol = d
		t.fieldDataArray[StringIndex["SYMBOL"]] = d
	} else if e == StringIndex["EXCHG"] {
		t.exchange = d
		t.fieldDataArray[StringIndex["EXCHG"]] = d
	} else if e == StringIndex["TSYMBOL"] {
		t.tSymbol = d
		t.fieldDataArray[StringIndex["TSYMBOL"]] = d
		t.updatedFieldsArray[StringIndex["TSYMBOL"]] = true
	}
}

type DepthTopicData struct {
	*TopicData
}

func NewDepthTopicData() *DepthTopicData {
	t := &DepthTopicData{
		TopicData: NewTopicData(TopicTypes["DEPTH"]),
	}
	t.updatedFieldsArray = [100]interface{}{}
	t.multiplier = 0
	t.precision = 0
	t.precisionValue = 0
	return t
}

func (t *DepthTopicData) setMultiplierAndPrec() {
	if t.updatedFieldsArray[DEPTH_INDEX["PRECISION"]] != nil {
		t.precision = t.fieldDataArray[DEPTH_INDEX["PRECISION"]].(int)
		t.precisionValue = int(pow(10, t.precision))
	}
	if t.updatedFieldsArray[DEPTH_INDEX["MULTIPLIER"]] != nil {
		t.multiplier = t.fieldDataArray[DEPTH_INDEX["MULTIPLIER"]].(int)
	}
}

func pow(x, y int) int {
	result := 1
	for i := 0; i < y; i++ {
		result *= x
	}
	return result
}

func (t *DepthTopicData) prepareData(reqType interface{}) map[string]interface{} {
	t.prepareCommonData()
	jsonRes := make(map[string]interface{})
	for d, c := range DepthMapping {
		e := t.fieldDataArray[d]
		if t.updatedFieldsArray[d] != nil && e != nil && c != nil {
			switch c.Type {
			case FieldTypes["FLOAT32"]:
				e = round(float64(e.(int))/float64(t.multiplier*t.precisionValue), t.precision)
			case FieldTypes["DATE"]:
				e = getFormatDate(int64(e.(int)))
			}
			jsonRes[c.Name] = fmt.Sprintf("%v", e)
		}
	}
	t.updatedFieldsArray = [100]interface{}{}
	if reqType != nil {
		jsonRes["request_type"] = reqType
	}
	return jsonRes
}

func round(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func getAcknowledgementReq(a int32) []byte {
	buffer := NewByteData(11)
	buffer.markStartOfMsg()
	buffer.appendByte(byte(BinRespTypes["ACK_TYPE"]))
	buffer.appendByte(1)
	buffer.appendByte(1)
	buffer.appendShort(4)
	buffer.appendInt(a)
	buffer.markEndOfMsg()
	return buffer.getBytes()
}

func prepareConnectionRequest(a string) []byte {
	userIDLen := len(a)
	src := "JS_API"
	srcLen := len(src)
	buffer := make([]byte, userIDLen+srcLen+10)
	buffer[0] = byte(BinRespTypes["CONNECTION_TYPE"])
	buffer[1] = 2
	buffer[2] = 1
	binary.BigEndian.PutUint16(buffer[3:5], uint16(userIDLen))
	copy(buffer[5:], a)
	buffer[5+userIDLen] = 2
	binary.BigEndian.PutUint16(buffer[6+userIDLen:], uint16(srcLen))
	copy(buffer[8+userIDLen:], src)
	buffer[8+userIDLen+srcLen] = byte(BinRespTypes["END_OF_MSG"])
	return buffer
}

func prepareConnectionRequest2(a, c string) []byte {
	src := "JS_API"
	srcLen := len(src)
	jwtLen := len(a)
	redisLen := len(c)
	buffer := NewByteData(srcLen + jwtLen + redisLen + 13)
	buffer.markStartOfMsg()
	buffer.appendByte(byte(BinRespTypes["CONNECTION_TYPE"]))
	buffer.appendByte(3)
	buffer.appendByte(1)
	buffer.appendShort(int16(jwtLen))
	buffer.appendString(a)
	buffer.appendByte(2)
	buffer.appendShort(int16(redisLen))
	buffer.appendString(c)
	buffer.appendByte(3)
	buffer.appendShort(int16(srcLen))
	buffer.appendString(src)
	buffer.markEndOfMsg()
	return buffer.getBytes()
}

func isScripOK(a string) bool {
	scripsCount := len(bytes.Split([]byte(a), []byte("&")))
	if scripsCount > MaxScrips {
		fmt.Println("Maximum scrips allowed per request is", MaxScrips)
		return false
	}
	return true
}

func prepareSubsUnSubsRequest(scrips, subscribeType, scripPrefix string, channelNum int) []byte {
	if !isScripOK(scrips) {
		return nil
	}
	dataArr := getScripByteArray(scrips, scripPrefix)
	buffer := NewByteData(len(dataArr) + 11)
	buffer.markStartOfMsg()
	buffer.appendByte(byte(subscribeType[0]))
	buffer.appendByte(2)
	buffer.appendByte(1)
	buffer.appendShort(int16(len(dataArr)))
	buffer.appendByteArr(dataArr, len(dataArr))
	buffer.appendByte(2)
	buffer.appendShort(1)
	buffer.appendByte(byte(channelNum))
	buffer.markEndOfMsg()
	return buffer.getBytes()
}

func prepareSnapshotRequest(a, c, d string) []byte {
	if !isScripOK(a) {
		return nil
	}
	dataArr := getScripByteArray(a, d)
	buffer := NewByteData(len(dataArr) + 7)
	buffer.markStartOfMsg()
	buffer.appendByte(byte(c[0]))
	buffer.appendByte(1)
	buffer.appendByte(2)
	buffer.appendShort(int16(len(dataArr)))
	buffer.appendByteArr(dataArr, len(dataArr))
	buffer.markEndOfMsg()
	return buffer.getBytes()
}

func prepareChannelRequest(c int, a []int) []byte {
	buffer := make([]byte, 15)
	buffer[0] = byte(c)
	buffer[1] = 1
	buffer[2] = 1
	binary.BigEndian.PutUint16(buffer[3:5], 8)
	var int1, int2 int
	for _, d := range a {
		switch {
		case 0 < d && d <= 32:
			int1 |= 1 << d
		case 32 < d && d <= 64:
			int2 |= 1 << d
		default:
			fmt.Println("Error: Channel values must be in this range  [ val > 0 && val < 65 ]")
		}
	}
	binary.BigEndian.PutUint32(buffer[5:9], uint32(int2))
	binary.BigEndian.PutUint32(buffer[9:13], uint32(int1))
	return buffer
}

func prepareThrottlingIntervalRequest(a int) []byte {
	buffer := make([]byte, 11)
	buffer[0] = byte(BinRespTypes["THROTTLING_TYPE"])
	buffer[1] = 1
	buffer[2] = 1
	binary.BigEndian.PutUint16(buffer[3:5], 4)
	binary.BigEndian.PutUint32(buffer[5:9], uint32(a))
	return buffer
}

func getScripByteArray(c, a string) []byte {
	if c[len(c)-1] == '&' {
		c = c[:len(c)-1]
	}
	scripArray := bytes.Split([]byte(c), []byte("&"))
	scripsCount := len(scripArray)
	dataLen := 0
	for i := range scripArray {
		scripArray[i] = append([]byte(a+"|"), scripArray[i]...)
		dataLen += len(scripArray[i]) + 1
	}
	bytes := make([]byte, dataLen+2)
	pos := 0
	bytes[pos] = byte((scripsCount >> 8) & 255)
	pos++
	bytes[pos] = byte(scripsCount & 255)
	pos++
	for _, currScrip := range scripArray {
		scripLen := len(currScrip)
		bytes[pos] = byte(scripLen & 255)
		pos++
		copy(bytes[pos:], currScrip)
		pos += scripLen
	}
	return bytes
}

func getOpcChainSubsRequest(d string, e int64, a, c, f byte) []byte {
	opcKeyLen := len(d)
	buffer := make([]byte, opcKeyLen+30)
	pos := 0
	buffer[pos] = byte(BinRespTypes["OPC_SUBSCRIBE"])
	pos++
	buffer[pos] = 5
	pos++
	buffer[pos] = 1
	pos++
	buffer[pos] = byte(opcKeyLen >> 8 & 255)
	pos++
	buffer[pos] = byte(opcKeyLen & 255)
	pos++
	copy(buffer[pos:], d)
	pos += opcKeyLen
	buffer[pos] = 2
	pos++
	binary.BigEndian.PutUint16(buffer[pos:pos+2], 8)
	pos += 2
	binary.BigEndian.PutUint64(buffer[pos:pos+8], uint64(e))
	pos += 8
	buffer[pos] = 3
	pos++
	binary.BigEndian.PutUint16(buffer[pos:pos+2], 1)
	pos += 2
	buffer[pos] = a
	pos++
	buffer[pos] = 4
	pos++
	binary.BigEndian.PutUint16(buffer[pos:pos+2], 1)
	pos += 2
	buffer[pos] = c
	pos++
	buffer[pos] = 5
	pos++
	binary.BigEndian.PutUint16(buffer[pos:pos+2], 1)
	pos += 2
	buffer[pos] = f
	return buffer
}

func sendJSONArrResp(a interface{}) string {
	jsonArrRes := []interface{}{a}
	jsonData, _ := json.Marshal(jsonArrRes)
	return string(jsonData)
}

func buf2long(a []byte) int64 {
	val := int64(0)
	leng := len(a)
	for i := 0; i < leng; i++ {
		j := leng - 1 - i
		val += int64(a[j]) << (i * 8)
	}
	return val
}

func buf2string(a []byte) string {
	return string(a)
}
