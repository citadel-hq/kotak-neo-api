package websocket

import (
	"fmt"
	"time"
)

// Constants and mappings
const (
	MaxScrips    = 100
	TrashVal     = -2147483648
	isEncryptOut = false
	isEncryptIn  = true
)

var (
	topicList = map[int]Topic{}
	counter   = 0
)

var FieldTypes = map[string]int{
	"FLOAT32": 1,
	"LONG":    2,
	"DATE":    3,
	"STRING":  4,
}

var StringIndex = map[string]int{
	"NAME":    51,
	"SYMBOL":  52,
	"EXCHG":   53,
	"TSYMBOL": 54,
}

var DEPTH_INDEX = map[string]int{
	"MULTIPLIER": 32,
	"PRECISION":  33,
}

var BinRespTypes = map[string]int{
	"CONNECTION_TYPE":  1,
	"THROTTLING_TYPE":  2,
	"ACK_TYPE":         3,
	"SUBSCRIBE_TYPE":   4,
	"UNSUBSCRIBE_TYPE": 5,
	"DATA_TYPE":        6,
	"CHPAUSE_TYPE":     7,
	"CHRESUME_TYPE":    8,
	"SNAPSHOT":         9,
	"OPC_SUBSCRIBE":    10,
}

var BinRespStat = map[string]string{
	"OK":     "K",
	"NOT_OK": "N",
}

var ResponseTypes = map[string]int{
	"SNAP":   83,
	"UPDATE": 85,
}

var STAT = map[string]string{
	"OK":     "Ok",
	"NOT_OK": "NotOk",
}

var RespTypeValues = map[string]string{
	"CONN":     "cn",
	"SUBS":     "sub",
	"UNSUBS":   "unsub",
	"SNAP":     "snap",
	"CHANNELR": "cr",
	"CHANNELP": "cp",
	"OPC":      "opc",
}

var RespCodes = map[string]int{
	"SUCCESS":               200,
	"CONNECTION_FAILED":     11001,
	"CONNECTION_INVALID":    11002,
	"SUBSCRIPTION_FAILED":   11011,
	"UNSUBSCRIPTION_FAILED": 11012,
	"SNAPSHOT_FAILED":       11013,
	"CHANNELP_FAILED":       11031,
	"CHANNELR_FAILED":       11032,
}

var TopicTypes = map[string]string{
	"SCRIP": "sf",
	"INDEX": "if",
	"DEPTH": "dp",
}

var INDEX_INDEX = map[string]int{
	"LTP":        2,
	"CLOSE":      3,
	"CHANGE":     10,
	"PERCHANGE":  11,
	"MULTIPLIER": 8,
	"PRECISION":  9,
}

var ScripIndex = map[string]int{
	"VOLUME":     4,
	"LTP":        5,
	"CLOSE":      21,
	"VWAP":       13,
	"MULTIPLIER": 23,
	"PRECISION":  24,
	"CHANGE":     25,
	"PERCHANGE":  26,
	"TURNOVER":   27,
}

var Keys = map[string]string{
	"TYPE":           "type",
	"USER_ID":        "user",
	"SESSION_ID":     "sessionid",
	"SCRIPS":         "scrips",
	"CHANNEL_NUM":    "channelnum",
	"CHANNEL_NUMS":   "channelnums",
	"JWT":            "jwt",
	"REDIS_KEY":      "redis",
	"STK_PRC":        "stkprc",
	"HIGH_STK":       "highstk",
	"LOW_STK":        "lowstk",
	"OPC_KEY":        "key",
	"AUTHORIZATION":  "Authorization",
	"SID":            "Sid",
	"X_ACCESS_TOKEN": "x-access-token",
	"SOURCE":         "source",
}

var ReqTypeValues = map[string]string{
	"CONNECTION":          "cn",
	"SCRIP_SUBS":          "mws",
	"SCRIP_UNSUBS":        "mwu",
	"INDEX_SUBS":          "ifs",
	"INDEX_UNSUBS":        "ifu",
	"DEPTH_SUBS":          "dps",
	"DEPTH_UNSUBS":        "dpu",
	"CHANNEL_RESUME":      "cr",
	"CHANNEL_PAUSE":       "cp",
	"SNAP_MW":             "mwsp",
	"SNAP_DP":             "dpsp",
	"SNAP_IF":             "ifsp",
	"OPC_SUBS":            "opc",
	"THROTTLING_INTERVAL": "ti",
	"STR":                 "str",
	"FORCE_CONNECTION":    "fcn",
	"LOG":                 "log",
}

type DataType struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

var IndexMapping = make([]*DataType, 55)
var ScripMapping = make([]*DataType, 100)
var DepthMapping = make([]*DataType, 55)

func init() {
	IndexMapping[0] = &DataType{"ftm0", FieldTypes["DATE"]}
	IndexMapping[1] = &DataType{"dtm1", FieldTypes["DATE"]}
	IndexMapping[INDEX_INDEX["LTP"]] = &DataType{"iv", FieldTypes["FLOAT32"]}
	IndexMapping[INDEX_INDEX["CLOSE"]] = &DataType{"ic", FieldTypes["FLOAT32"]}
	IndexMapping[4] = &DataType{"tvalue", FieldTypes["DATE"]}
	IndexMapping[5] = &DataType{"highPrice", FieldTypes["FLOAT32"]}
	IndexMapping[6] = &DataType{"lowPrice", FieldTypes["FLOAT32"]}
	IndexMapping[7] = &DataType{"openingPrice", FieldTypes["FLOAT32"]}
	IndexMapping[8] = &DataType{"mul", FieldTypes["LONG"]}
	IndexMapping[INDEX_INDEX["PRECISION"]] = &DataType{"prec", FieldTypes["LONG"]}
	IndexMapping[INDEX_INDEX["CHANGE"]] = &DataType{"cng", FieldTypes["FLOAT32"]}
	IndexMapping[INDEX_INDEX["PERCHANGE"]] = &DataType{"nc", FieldTypes["STRING"]}
	IndexMapping[StringIndex["NAME"]] = &DataType{"name", FieldTypes["STRING"]}
	IndexMapping[StringIndex["SYMBOL"]] = &DataType{"tk", FieldTypes["STRING"]}
	IndexMapping[StringIndex["EXCHG"]] = &DataType{"e", FieldTypes["STRING"]}
	IndexMapping[StringIndex["TSYMBOL"]] = &DataType{"ts", FieldTypes["STRING"]}

	ScripMapping[0] = &DataType{"ftm0", FieldTypes["DATE"]}
	ScripMapping[1] = &DataType{"dtm1", FieldTypes["DATE"]}
	ScripMapping[2] = &DataType{"fdtm", FieldTypes["DATE"]}
	ScripMapping[3] = &DataType{"ltt", FieldTypes["DATE"]}
	ScripMapping[ScripIndex["VOLUME"]] = &DataType{"v", FieldTypes["LONG"]}
	ScripMapping[ScripIndex["LTP"]] = &DataType{"ltp", FieldTypes["FLOAT32"]}
	ScripMapping[6] = &DataType{"ltq", FieldTypes["LONG"]}
	ScripMapping[7] = &DataType{"tbq", FieldTypes["LONG"]}
	ScripMapping[8] = &DataType{"tsq", FieldTypes["LONG"]}
	ScripMapping[9] = &DataType{"bp", FieldTypes["FLOAT32"]}
	ScripMapping[10] = &DataType{"sp", FieldTypes["FLOAT32"]}
	ScripMapping[11] = &DataType{"bq", FieldTypes["LONG"]}
	ScripMapping[12] = &DataType{"bs", FieldTypes["LONG"]}
	ScripMapping[ScripIndex["VWAP"]] = &DataType{"ap", FieldTypes["FLOAT32"]}
	ScripMapping[14] = &DataType{"lo", FieldTypes["FLOAT32"]}
	ScripMapping[15] = &DataType{"h", FieldTypes["FLOAT32"]}
	ScripMapping[16] = &DataType{"lcl", FieldTypes["FLOAT32"]}
	ScripMapping[17] = &DataType{"ucl", FieldTypes["FLOAT32"]}
	ScripMapping[18] = &DataType{"yh", FieldTypes["FLOAT32"]}
	ScripMapping[19] = &DataType{"yl", FieldTypes["FLOAT32"]}
	ScripMapping[20] = &DataType{"op", FieldTypes["FLOAT32"]}
	ScripMapping[ScripIndex["CLOSE"]] = &DataType{"c", FieldTypes["FLOAT32"]}
	ScripMapping[22] = &DataType{"oi", FieldTypes["LONG"]}
	ScripMapping[ScripIndex["MULTIPLIER"]] = &DataType{"mul", FieldTypes["LONG"]}
	ScripMapping[ScripIndex["PRECISION"]] = &DataType{"prec", FieldTypes["LONG"]}
	ScripMapping[ScripIndex["CHANGE"]] = &DataType{"cng", FieldTypes["FLOAT32"]}
	ScripMapping[ScripIndex["PERCHANGE"]] = &DataType{"nc", FieldTypes["STRING"]}
	ScripMapping[ScripIndex["TURNOVER"]] = &DataType{"to", FieldTypes["FLOAT32"]}
	ScripMapping[StringIndex["NAME"]] = &DataType{"name", FieldTypes["STRING"]}
	ScripMapping[StringIndex["SYMBOL"]] = &DataType{"tk", FieldTypes["STRING"]}
	ScripMapping[StringIndex["EXCHG"]] = &DataType{"e", FieldTypes["STRING"]}
	ScripMapping[StringIndex["TSYMBOL"]] = &DataType{"ts", FieldTypes["STRING"]}

	DepthMapping[0] = &DataType{"ftm0", FieldTypes["DATE"]}
	DepthMapping[1] = &DataType{"dtm1", FieldTypes["DATE"]}
	DepthMapping[2] = &DataType{"bp", FieldTypes["FLOAT32"]}
	DepthMapping[3] = &DataType{"bp1", FieldTypes["FLOAT32"]}
	DepthMapping[4] = &DataType{"bp2", FieldTypes["FLOAT32"]}
	DepthMapping[5] = &DataType{"bp3", FieldTypes["FLOAT32"]}
	DepthMapping[6] = &DataType{"bp4", FieldTypes["FLOAT32"]}
	DepthMapping[7] = &DataType{"sp", FieldTypes["FLOAT32"]}
	DepthMapping[8] = &DataType{"sp1", FieldTypes["FLOAT32"]}
	DepthMapping[9] = &DataType{"sp2", FieldTypes["FLOAT32"]}
	DepthMapping[10] = &DataType{"sp3", FieldTypes["FLOAT32"]}
	DepthMapping[11] = &DataType{"sp4", FieldTypes["FLOAT32"]}
	DepthMapping[12] = &DataType{"bq", FieldTypes["LONG"]}
	DepthMapping[13] = &DataType{"bq1", FieldTypes["LONG"]}
	DepthMapping[14] = &DataType{"bq2", FieldTypes["LONG"]}
	DepthMapping[15] = &DataType{"bq3", FieldTypes["LONG"]}
	DepthMapping[16] = &DataType{"bq4", FieldTypes["LONG"]}
	DepthMapping[17] = &DataType{"bs", FieldTypes["LONG"]}
	DepthMapping[18] = &DataType{"bs1", FieldTypes["LONG"]}
	DepthMapping[19] = &DataType{"bs2", FieldTypes["LONG"]}
	DepthMapping[20] = &DataType{"bs3", FieldTypes["LONG"]}
	DepthMapping[21] = &DataType{"bs4", FieldTypes["LONG"]}
	DepthMapping[22] = &DataType{"bno1", FieldTypes["LONG"]}
	DepthMapping[23] = &DataType{"bno2", FieldTypes["LONG"]}
	DepthMapping[24] = &DataType{"bno3", FieldTypes["LONG"]}
	DepthMapping[25] = &DataType{"bno4", FieldTypes["LONG"]}
	DepthMapping[26] = &DataType{"bno5", FieldTypes["LONG"]}
	DepthMapping[27] = &DataType{"sno1", FieldTypes["LONG"]}
	DepthMapping[28] = &DataType{"sno2", FieldTypes["LONG"]}
	DepthMapping[29] = &DataType{"sno3", FieldTypes["LONG"]}
	DepthMapping[30] = &DataType{"sno4", FieldTypes["LONG"]}
	DepthMapping[31] = &DataType{"sno5", FieldTypes["LONG"]}
	DepthMapping[DEPTH_INDEX["MULTIPLIER"]] = &DataType{"mul", FieldTypes["LONG"]}
	DepthMapping[DEPTH_INDEX["PRECISION"]] = &DataType{"prec", FieldTypes["LONG"]}
	DepthMapping[StringIndex["NAME"]] = &DataType{"name", FieldTypes["STRING"]}
	DepthMapping[StringIndex["SYMBOL"]] = &DataType{"tk", FieldTypes["STRING"]}
	DepthMapping[StringIndex["EXCHG"]] = &DataType{"e", FieldTypes["STRING"]}
	DepthMapping[StringIndex["TSYMBOL"]] = &DataType{"ts", FieldTypes["STRING"]}
}

func leadingZero(a int) string {
	if a < 10 {
		return fmt.Sprintf("0%d", a)
	}
	return fmt.Sprintf("%d", a)
}

func getFormatDate(a int64) string {
	date := time.Unix(a, 0)
	formatDate := fmt.Sprintf(
		"%02d/%02d/%d %02d:%02d:%02d",
		date.Day(),
		date.Month(),
		date.Year(),
		date.Hour(),
		date.Minute(),
		date.Second(),
	)
	return formatDate
}
