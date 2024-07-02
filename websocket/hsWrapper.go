package websocket

import (
	"encoding/json"
	"strings"
	"fmt"
	"encoding/binary"
)

type HSWrapper struct {
	counter int
	ackNum  int
}

func NewHSWrapper() *HSWrapper {
	return &HSWrapper{
		counter: 0,
		ackNum:  0,
	}
}

func (h *HSWrapper) getNewTopicData(c string) interface{} {
	parts := strings.Split(c, "|")
	feedType := parts[0]

	switch feedType {
	case TopicTypes["SCRIP"]:
		return NewScripTopicData()
	case TopicTypes["INDEX"]:
		return NewIndexTopicData()
	case TopicTypes["DEPTH"]:
		return NewDepthTopicData()
	default:
		return nil
	}
}

func (h *HSWrapper) getStatus(c []byte, d int) string {
	status := BinRespStat["NOT_OK"]
	dx := int64(d)
	fieldCount := buf2long(c[dx : dx+1])
	dx++
	if fieldCount > 0 {
		dx++
		fieldLength := buf2long(c[dx : dx+2])
		dx += 2
		status = buf2string(c[dx : dx+fieldLength])
		dx += fieldLength
	}
	return status
}

func (h *HSWrapper) parseData(e []byte) []byte {
	//pos := 0
	pos := 2
	dataType := int(binary.BigEndian.Uint64(e[pos : pos+1]))
	pos++
	if dataType == BinRespTypes["CONNECTION_TYPE"] {
		return h.handleConnectionType(e, pos)
	} else if dataType == BinRespTypes["DATA_TYPE"] {
		return h.handleDataType(e, pos)
	} else if dataType == BinRespTypes["SUBSCRIBE_TYPE"] || dataType == BinRespTypes["UNSUBSCRIBE_TYPE"] {
		return h.handleSubscriptionType(e, pos, dataType)
	} else if dataType == BinRespTypes["SNAPSHOT"] {
		return h.handleSnapshot(e, pos)
	} else if dataType == BinRespTypes["CHPAUSE_TYPE"] || dataType == BinRespTypes["CHRESUME_TYPE"] {
		return h.handleChannelType(e, pos, dataType)
	} else if dataType == BinRespTypes["OPC_SUBSCRIBE"] {
		return h.handleOPCSubscribe(e, pos)
	}
	return nil
}

func (h *HSWrapper) handleConnectionType(e []byte, pos int) []byte {
	jsonRes := make(map[string]interface{})
	fCount := e[pos]
	pos++
	if fCount >= 2 {
		pos++
		valLen := buf2long(e[pos : pos+2])
		pos += 2
		status := string(e[pos : pos+valLen])
		pos += valLen
		pos++
		valLen = buf2long(e[pos : pos+2])
		pos += 2
		ackCount := buf2long(e[pos : pos+valLen])
		pos += valLen
		if status == BinRespStat["OK"] {
			jsonRes["stat"] = STAT["OK"]
			jsonRes["type"] = RespTypeValues["CONN"]
			jsonRes["msg"] = "successful"
			jsonRes["stCode"] = RespCodes["SUCCESS"]
		} else {
			jsonRes["stat"] = STAT["NOT_OK"]
			jsonRes["type"] = RespTypeValues["CONN"]
			jsonRes["msg"] = "failed"
			jsonRes["stCode"] = RespCodes["CONNECTION_FAILED"]
		}
		h.ackNum = ackCount
	} else if fCount == 1 {
		fid1 := e[pos]
		pos++
		valLen := buf2long(e[pos : pos+2])
		pos += 2
		status := string(e[pos : pos+valLen])
		pos += valLen
		if status == BinRespStat["OK"] {
			jsonRes["stat"] = STAT["OK"]
			jsonRes["type"] = RespTypeValues["CONN"]
			jsonRes["msg"] = "successful"
			jsonRes["stCode"] = RespCodes["SUCCESS"]
		} else {
			jsonRes["stat"] = STAT["NOT_OK"]
			jsonRes["type"] = RespTypeValues["CONN"]
			jsonRes["msg"] = "failed"
			jsonRes["stCode"] = RespCodes["CONNECTION_FAILED"]
		}
	} else {
		jsonRes["stat"] = STAT["NOT_OK"]
		jsonRes["type"] = RespTypeValues["CONN"]
		jsonRes["msg"] = "invalid field count"
		jsonRes["stCode"] = RespCodes["CONNECTION_INVALID"]
	}
	return sendJSONArrResp(jsonRes)
}

func (h *HSWrapper) handleDataType(e []byte, pos int) []byte {
	if h.ackNum > 0 {
		h.counter++
		msgNum := buf2long(e[pos : pos+4])
		pos += 4
		if h.counter == h.ackNum {
			req := getAcknowledgementReq(msgNum)
			if ws != nil {
				ws.send(req, 0x2)
				h.counter = 0
			}
		}
	}
	hList := make([]interface{}, 0)
	g := buf2long(e[pos : pos+2])
	pos += 2
	for n := 0; n < g; n++ {
		pos += 2
		c := buf2long(e[pos : pos+1])
		pos++
		if c == ResponseTypes["SNAP"] {
			f := buf2long(e[pos : pos+4])
			pos += 4
			nameLen := buf2long(e[pos : pos+1])
			pos++
			topicName := buf2string(e[pos : pos+nameLen])
			pos += nameLen
			d := h.getNewTopicData(topicName)
			if d != nil {
				topicList[f] = d
				fCount := buf2long(e[pos : pos+1])
				pos++
				for index := 0; index < fCount; index++ {
					fValue := buf2long(e[pos : pos+4])
					d.setLongValues(index, fValue)
					pos += 4
				}
				d.setMultiplierAndPrec()
				fCount = buf2long(e[pos : pos+1])
				pos++
				for index := 0; index < fCount; index++ {
					fid := buf2long(e[pos : pos+1])
					pos++
					dataLen := buf2long(e[pos : pos+1])
					pos++
					strVal := buf2string(e[pos : pos+dataLen])
					pos += dataLen
					d.setStringValues(fid, strVal)
				}
				hList = append(hList, d.prepareData("SNAP"))
			} else {
				fmt.Println("Invalid topic feed type !")
			}
		} else if c == ResponseTypes["UPDATE"] {
			f := buf2long(e[pos : pos+4])
			pos += 4
			d := topicList[f]
			if d == nil {
				fmt.Println("Topic Not Available in TopicList!")
			} else {
				fCount := buf2long(e[pos : pos+1])
				pos++
				for index := 0; index < fCount; index++ {
					fValue := buf2long(e[pos : pos+4])
					d.setLongValues(index, fValue)
					pos += 4
				}
			}
			hList = append(hList, d.prepareData("SUB"))
		} else {
			fmt.Println("Invalid ResponseType:", c)
		}
	}
	return sendJSONArrResp(hList)
}

func (h *HSWrapper) handleSubscriptionType(e []byte, pos int, dataType int) []byte {
	status := h.getStatus(e, pos)
	jsonRes := make(map[string]interface{})
	if status == BinRespStat["OK"] {
		jsonRes["stat"] = STAT["OK"]
		if dataType == BinRespTypes["SUBSCRIBE_TYPE"] {
			jsonRes["type"] = RespTypeValues["SUBS"]
		} else {
			jsonRes["type"] = RespTypeValues["UNSUBS"]
		}
		jsonRes["msg"] = "successful"
		jsonRes["stCode"] = RespCodes["SUCCESS"]
	} else {
		jsonRes["stat"] = STAT["NOT_OK"]
		if dataType == BinRespTypes["SUBSCRIBE_TYPE"] {
			jsonRes["type"] = RespTypeValues["SUBS"]
			jsonRes["msg"] = "subscription failed"
			jsonRes["stCode"] = RespCodes["SUBSCRIPTION_FAILED"]
		} else {
			jsonRes["type"] = RespTypeValues["UNSUBS"]
			jsonRes["msg"] = "unsubscription failed"
			jsonRes["stCode"] = RespCodes["UNSUBSCRIPTION_FAILED"]
		}
	}
	return sendJSONArrResp(jsonRes)
}

func (h *HSWrapper) handleSnapshot(e []byte, pos int) []byte {
	status := h.getStatus(e, pos)
	jsonRes := make(map[string]interface{})
	if status == BinRespStat["OK"] {
		jsonRes["stat"] = STAT["OK"]
		jsonRes["type"] = RespTypeValues["SNAP"]
		jsonRes["msg"] = "successful"
		jsonRes["stCode"] = RespCodes["SUCCESS"]
	} else {
		jsonRes["stat"] = STAT["NOT_OK"]
		jsonRes["type"] = RespTypeValues["SNAP"]
		jsonRes["msg"] = "failed"
		jsonRes["stCode"] = RespCodes["SNAPSHOT_FAILED"]
	}
	return sendJSONArrResp(jsonRes)
}

func (h *HSWrapper) handleChannelType(e []byte, pos int, dataType int) []byte {
	status := h.getStatus(e, pos)
	jsonRes := make(map[string]interface{})
	if status == BinRespStat["OK"] {
		jsonRes["stat"] = STAT["OK"]
		if dataType == BinRespTypes["CHPAUSE_TYPE"] {
			jsonRes["type"] = RespTypeValues["CHANNELP"]
		} else {
			jsonRes["type"] = RespTypeValues["CHANNELR"]
		}
		jsonRes["msg"] = "successful"
		jsonRes["stCode"] = RespCodes["SUCCESS"]
	} else {
		jsonRes["stat"] = STAT["NOT_OK"]
		if dataType == BinRespTypes["CHPAUSE_TYPE"] {
			jsonRes["type"] = RespTypeValues["CHANNELP"]
		} else {
			jsonRes["type"] = RespTypeValues["CHANNELR"]
		}
		jsonRes["msg"] = "failed"
		if dataType == BinRespTypes["CHPAUSE_TYPE"] {
			jsonRes["stCode"] = RespCodes["CHANNELP_FAILED"]
		} else {
			jsonRes["stCode"] = RespCodes["CHANNELR_FAILED"]
		}
	}
	return sendJSONArrResp(jsonRes)
}

func (h *HSWrapper) handleOPCSubscribe(e []byte, pos int) []byte {
	status := h.getStatus(e, pos)
	pos += 5
	jsonRes := make(map[string]interface{})
	if status == BinRespStat["OK"] {
		jsonRes["stat"] = STAT["OK"]
		jsonRes["type"] = RespTypeValues["OPC"]
		jsonRes["msg"] = "successful"
		jsonRes["stCode"] = RespCodes["SUCCESS"]
		fld := buf2long(e[pos : pos+1])
		pos++
		fieldLength := buf2long(e[pos : pos+2])
		pos += 2
		opcKey := buf2string(e[pos : pos+fieldLength])
		pos += fieldLength
		jsonRes["key"] = opcKey
		fld = buf2long(e[pos : pos+1])
		pos++
		fieldLength = buf2long(e[pos : pos+2])
		pos += 2
		data := buf2string(e[pos : pos+fieldLength])
		pos += fieldLength
		var scrips map[string]interface{}
		json.Unmarshal([]byte(data), &scrips)
		jsonRes["scrips"] = scrips["data"]
	} else {
		jsonRes["stat"] = STAT["NOT_OK"]
		jsonRes["type"] = RespTypeValues["OPC"]
		jsonRes["msg"] = "failed"
		jsonRes["stCode"] = 11040
	}
	return sendJSONArrResp(jsonRes)
}
