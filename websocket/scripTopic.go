package websocket

import (
	"fmt"
	"math"
	"strconv"
)

type ScripTopicData struct {
	*TopicData
	precision      int
	precisionValue int
	multiplier     int
}

func NewScripTopicData() *ScripTopicData {
	t := &ScripTopicData{
		TopicData: NewTopicData(TopicTypes["SCRIP"]),
	}
	t.updatedFieldsArray = [100]interface{}{}
	return t
}

func (t *ScripTopicData) setMultiplierAndPrecision() {
	if t.updatedFieldsArray[ScripIndex["PRECISION"]] != nil {
		t.precision = t.fieldDataArray[ScripIndex["PRECISION"]].(int)
		t.precisionValue = int(math.Pow(10, float64(t.precision)))
	}
	if t.updatedFieldsArray[ScripIndex["MULTIPLIER"]] != nil {
		t.multiplier = t.fieldDataArray[ScripIndex["MULTIPLIER"]].(int)
	}
}

func (t *ScripTopicData) prepareData(reqType string) map[string]interface{} {
	t.prepareCommonData()
	precisionFormat := "%." + strconv.Itoa(t.precision) + "f"

	if t.updatedFieldsArray[ScripIndex["LTP"]] != nil || t.updatedFieldsArray[ScripIndex["CLOSE"]] != nil {
		ltp := t.fieldDataArray[ScripIndex["LTP"]]
		close := t.fieldDataArray[ScripIndex["CLOSE"]]
		if ltp != nil && close != nil {
			change := ltp.(int) - close.(int)
			t.fieldDataArray[ScripIndex["CHANGE"]] = change
			t.updatedFieldsArray[ScripIndex["CHANGE"]] = true
			t.fieldDataArray[ScripIndex["PERCHANGE"]] = fmt.Sprintf(precisionFormat, float64(change)/float64(close.(int))*100)
			t.updatedFieldsArray[ScripIndex["PERCHANGE"]] = true
		}
	}

	if t.updatedFieldsArray[ScripIndex["VOLUME"]] != nil || t.updatedFieldsArray[ScripIndex["VWAP"]] != nil {
		volume := t.fieldDataArray[ScripIndex["VOLUME"]]
		vwap := t.fieldDataArray[ScripIndex["VWAP"]]
		if volume != nil && vwap != nil {
			t.fieldDataArray[ScripIndex["TURNOVER"]] = volume.(int) * vwap.(int)
			t.updatedFieldsArray[ScripIndex["TURNOVER"]] = true
		}
	}

	jsonRes := make(map[string]interface{})
	for index := range ScripMapping {
		dataType := ScripMapping[index]
		val := t.fieldDataArray[index]
		if t.updatedFieldsArray[index] != nil && val != nil && dataType != nil {
			switch dataType.Type {
			case FieldTypes["FLOAT32"]:
				val = fmt.Sprintf(precisionFormat, float64(val.(int))/float64(t.multiplier*t.precisionValue))
			case FieldTypes["DATE"]:
				val = getFormatDate(int64(val.(int)))
			}
			jsonRes[dataType.Name] = fmt.Sprintf("%v", val)
		}
	}
	t.updatedFieldsArray = [100]interface{}{}
	if reqType != "" {
		jsonRes["request_type"] = reqType
	}
	return jsonRes
}
