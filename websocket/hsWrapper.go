package websocket

import (
	"encoding/json"
	"strings"
	"fmt"
	"encoding/binary"
	"context"
	"github.com/shikharvaish28/kotak-neo-api/api"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
	"strconv"
)

// BrokerEvent to be sent in channel
type BrokerEvent struct {
	Event     string
	Timestamp time.Time
}

type HSWrapper struct {
	counter       int
	ackNum        int
	ws            *websocket.Conn
	channel       chan BrokerEvent
	channelTokens map[int]interface{}
}

func NewHSWrapper() (*HSWrapper, chan BrokerEvent) {
	ctx := context.Background()
	conn, _, err := websocket.Dial(ctx, api.WebsocketUrl, nil)
	if err != nil {
		panic(fmt.Sprintf("websocket error - %s", err.Error()))
		return nil, nil
	}
	channel := make(chan BrokerEvent, 100) // bounded channel of 100 events.

	go func() {
		for {
			var msg string
			err := wsjson.Read(ctx, conn, &msg)
			if err != nil {
				log.Println("Failed to read message:", err)
				return
			}
			fmt.Println("Received:", msg)
			channel <- BrokerEvent{
				Event:     msg,
				Timestamp: time.Now(),
			}
		}
	}()

	return &HSWrapper{
		counter:       0,
		ackNum:        0,
		ws:            conn,
		channel:       channel,
		channelTokens: make(map[int]interface{}),
	}, channel
}

func (h *HSWrapper) GetLiveFeed(instrumentTokens []map[string]string, isIndex bool, isDepth bool) error {
	// TODO: perform total instrument handling.
	subscriptionType := ReqTypeValues["SCRIP_SUBS"]
	if isIndex {
		subscriptionType = ReqTypeValues["INDEX_SUBS"]
	}
	if isDepth {
		subscriptionType = ReqTypeValues["DEPTH_SUBS"]
	}
	tempTokenList := []map[string]interface{}{}

	// validate and push all instrumentation tokens.
	if validInstrumentationTokens(instrumentTokens) {
		for _, token := range instrumentTokens {
			key := token["instrument_token"]
			val := map[string]string{
				"instrument_token":  token["instrument_token"],
				"exchange_segment":  token["exchange_segment"],
				"subscription_type": subscriptionType,
			}
			tempTokenList = append(tempTokenList, map[string]interface{}{
				key: val,
			})
		}
		// is map[int][]map[string]interface{}
		channelTokens := h.channelSegregation(tempTokenList)
		h.subscribeScrips(channelTokens)
	}

	return nil
}

// channelSegregation internally returns a map[int][]map[string]interface{}
func (h *HSWrapper) channelSegregation(tmpTokenList []map[string]interface{}) map[int]interface{} {
	outChannelList := map[int]interface{}{}
	for channelNum := 2; channelNum < 17; channelNum++ {
		// Check if there is an existing channel array for this channel number
		if _, ok := h.channelTokens[channelNum]; !ok {
			h.channelTokens[channelNum] = []map[string]interface{}{}
		}
		if _, ok := outChannelList[channelNum]; !ok {
			outChannelList[channelNum] = []map[string]interface{}{}
		}

		// Note: I don't care about the length checks for now.
		if values, ok := h.channelTokens[channelNum].([]map[string]interface{}); ok {
			h.channelTokens[channelNum] = append(values, tmpTokenList...)
		}
		if values, ok := outChannelList[channelNum].([]map[string]interface{}); ok {
			outChannelList[channelNum] = append(values, tmpTokenList...)
		}
	}

	return outChannelList
}

func validInstrumentationTokens(tokens []map[string]string) bool {
	validParams := []string{"instrument_token", "exchange_segment"}
	for _, item := range tokens {
		for _, param := range validParams {
			if _, ok := item[param]; ok {
			} else {
				return false
			}
		}
	}
	return true
}

func (h *HSWrapper) subscribeScrips(tokens map[int]interface{}) {
	for _, v := range tokens {
		if values, ok := v.([]map[string]interface{}); ok {
			for channel, tokens := range values {
				for tokenMap := range tokens {
					values := []interface{}{}
					for _, val := range tokenMap {
						values = append(values, val)
					}
					instrumentMap := (values[0]).(map[string]string)
					scrips := h.formatTokenLive(instrumentMap)
					// TODO: Ensure that json being sent from here is mapped correctly in wsSend()
					requestParams := map[string]interface{}{
						"TYPE":        instrumentMap["subscription_type"],
						"SCRIPS":      scrips,
						"CHANNEL_NUM": channel,
					}
					h.wsSend(requestParams)
				}
			}
		}
	}
}

func (h *HSWrapper) formatTokenLive(instrumentTokens map[string]string) interface{} {
	scrips := ""

	if exchangeSegment, ok1 := instrumentTokens["exchange_segment"]; ok1 {
		if instrumentToken, ok2 := instrumentTokens["instrument_token"]; ok2 {
			scrips += exchangeSegment + "|" + instrumentToken
		}
	}
	return scrips
}

// TODO: continue completing this.
func (h *HSWrapper) wsSend(reqJson map[string]interface{}) {
	reqType := reqJson["TYPE"].(string)
	scrips := ""
	channelNum := 1
	if val, ok := reqJson["SCRIPS"]; ok {
		scrips = val.(string)
		channelNum = int(reqJson["CHANNEL_NUM"].(float64))
	}

	var req []byte

	switch reqType {
	case ReqTypeValues["CONNECTION"]:
		if userId, ok := reqJson["USER_ID"]; ok {
			req = prepareConnectionRequest(userId.(string))
		} else if sessionId, ok := reqJson["SESSION_ID"]; ok {
			req = prepareConnectionRequest(sessionId.(string))
		} else if auth, ok := reqJson["AUTHORIZATION"]; ok {
			if sid, ok := reqJson["SID"]; ok {
				req = prepareConnectionRequest2(auth.(string), sid.(string))
			} else {
				fmt.Println("Authorization mode is enabled: Authorization or Sid not found !")
			}
		} else {
			fmt.Println("Invalid conn mode !")
		}
	case ReqTypeValues["SCRIP_SUBS"]:
		req = prepareSubsUnSubsRequest(scrips, strconv.Itoa(BinRespTypes["SUBSCRIBE_TYPE"]), ScripPrefix, channelNum)
	case ReqTypeValues["SCRIP_UNSUBS"]:
		req = prepareSubsUnSubsRequest(scrips, strconv.Itoa(BinRespTypes["UNSUBSCRIBE_TYPE"]), ScripPrefix, channelNum)
	case ReqTypeValues["INDEX_SUBS"]:
		req = prepareSubsUnSubsRequest(scrips, strconv.Itoa(BinRespTypes["SUBSCRIBE_TYPE"]), IndexPrefix, channelNum)
	case ReqTypeValues["INDEX_UNSUBS"]:
		req = prepareSubsUnSubsRequest(scrips, strconv.Itoa(BinRespTypes["UNSUBSCRIBE_TYPE"]), IndexPrefix, channelNum)
	case ReqTypeValues["DEPTH_SUBS"]:
		req = prepareSubsUnSubsRequest(scrips, strconv.Itoa(BinRespTypes["SUBSCRIBE_TYPE"]), DepthPrefix, channelNum)
	case ReqTypeValues["DEPTH_UNSUBS"]:
		req = prepareSubsUnSubsRequest(scrips, strconv.Itoa(BinRespTypes["UNSUBSCRIBE_TYPE"]), DepthPrefix, channelNum)
	case ReqTypeValues["CHANNEL_PAUSE"]:
		req = prepareChannelRequest(BinRespTypes["CHPAUSE_TYPE"], channelNum)
	case ReqTypeValues["CHANNEL_RESUME"]:
		req = prepareChannelRequest(BinRespTypes["CHRESUME_TYPE"], channelNum)
	case ReqTypeValues["SNAP_MW"]:
		req = prepareSnapshotRequest(scrips, strconv.Itoa(BinRespTypes["SNAPSHOT"]), ScripPrefix)
	case ReqTypeValues["SNAP_DP"]:
		req = prepareSnapshotRequest(scrips, strconv.Itoa(BinRespTypes["SNAPSHOT"]), DepthPrefix)
	case ReqTypeValues["SNAP_IF"]:
		req = prepareSnapshotRequest(scrips, strconv.Itoa(BinRespTypes["SNAPSHOT"]), IndexPrefix)
	// TODO: complete opc chain subscription, throttling and logging requests later
	//case ReqTypeValues["OPC_SUBS"]:
	//	req = getOpcChainSubsRequest(reqJson["OPC_KEY"].(string), reqJson["STK_PRC"].(string),
	//		reqJson["HIGH_STK"].(string), reqJson["LOW_STK"].(string), channelNum)
	//case ReqTypeValues["THROTTLING_INTERVAL"]:
	//	req = prepareThrottlingIntervalRequest(scrips)
	//case ReqTypeValues["LOG"]:
	//	if enable, ok := reqJson["enable"].(bool); ok {
	//		enableLog(enable)
	//	}
	default:
		fmt.Println("Unknown request type!")
	}

	if h.ws != nil && len(req) > 0 {
		if err := h.ws.Write(context.Background(), websocket.MessageBinary, req); err != nil {
			fmt.Println("Unable to send request! Reason: Connection faulty or request not valid!")
		}
	} else {
		fmt.Println("Unable to send request! Reason: Connection faulty or request not valid!")
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
