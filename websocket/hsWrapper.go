package websocket

import (
	"encoding/json"
	"strings"
	"fmt"
	"encoding/binary"
	"nhooyr.io/websocket"
	"context"
	"github.com/shikharvaish28/kotak-neo-api/api"
)

type HSWrapper struct {
	counter int
	ackNum  int
	ws      *websocket.Conn
}

func NewHSWrapper() *HSWrapper {
	ctx := context.Background()
	conn, _, err := websocket.Dial(ctx, api.WebsocketUrl, nil)
	if err != nil {
		panic(fmt.Sprintf("websocket error - %s", err.Error()))
		return nil
	}

	return &HSWrapper{
		counter: 0,
		ackNum:  0,
		ws:      conn,
	}
}

func (h *HSWrapper) getNewTopicData(c string) Topic {
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
		valLen := int(buf2long(e[pos : pos+2]))
		pos += 2
		status := string(e[pos : pos+valLen])
		pos += valLen
		pos++
		valLen = int(buf2long(e[pos : pos+2]))
		pos += 2
		ackCount := int(buf2long(e[pos : pos+valLen]))
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
		pos++
		valLen := int(buf2long(e[pos : pos+2]))
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
		msgNum := int32(buf2long(e[pos : pos+4]))
		pos += 4
		if h.counter == h.ackNum {
			req := getAcknowledgementReq(msgNum)
			_ = h.ws.Write(context.Background(), websocket.MessageBinary, req)
		}
	}
	hList := make([]interface{}, 0)
	g := int(buf2long(e[pos : pos+2]))
	pos += 2
	for n := 0; n < g; n++ {
		pos += 2
		c := int(buf2long(e[pos : pos+1]))
		pos++
		if c == ResponseTypes["SNAP"] {
			f := int(buf2long(e[pos : pos+4]))
			pos += 4
			nameLen := int(buf2long(e[pos : pos+1]))
			pos++
			topicName := buf2string(e[pos : pos+nameLen])
			pos += nameLen
			topicData := h.getNewTopicData(topicName)
			if topicData != nil {
				topicList[f] = topicData
				fCount := int(buf2long(e[pos : pos+1]))
				pos++
				for index := 0; index < fCount; index++ {
					fValue := buf2long(e[pos : pos+4])
					topicData.setLongValues(index, fValue)
					pos += 4
				}
				topicData.setMultiplierAndPrecision()
				fCount = int(buf2long(e[pos : pos+1]))
				pos++
				for index := 0; index < fCount; index++ {
					fid := int(buf2long(e[pos : pos+1]))
					pos++
					dataLen := int(buf2long(e[pos : pos+1]))
					pos++
					strVal := buf2string(e[pos : pos+dataLen])
					pos += dataLen
					topicData.setStringValues(fid, strVal)
				}
				hList = append(hList, topicData.prepareData("SNAP"))
			} else {
				fmt.Println("Invalid topic feed type !")
			}
		} else if c == ResponseTypes["UPDATE"] {
			f := int(buf2long(e[pos : pos+4]))
			pos += 4
			topicData := topicList[f]
			if topicData == nil {
				fmt.Println("Topic Not Available in TopicList!")
			} else {
				fCount := int(buf2long(e[pos : pos+1]))
				pos++
				for index := 0; index < fCount; index++ {
					fValue := buf2long(e[pos : pos+4])
					topicData.setLongValues(index, fValue)
					pos += 4
				}
			}
			hList = append(hList, topicData.prepareData("SUB"))
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
		pos++
		fieldLength := int(buf2long(e[pos : pos+2]))
		pos += 2
		opcKey := buf2string(e[pos : pos+fieldLength])
		pos += fieldLength
		jsonRes["key"] = opcKey
		pos++
		fieldLength = int(buf2long(e[pos : pos+2]))
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
