package websocket

import (
	"fmt"
	"math"
)

type IndexTopicData struct {
	*TopicData
	precision      int
	precisionValue int
	multiplier     int
}

func NewIndexTopicData() *IndexTopicData {
	t := &IndexTopicData{
		TopicData: NewTopicData(TopicTypes["INDEX"]),
	}
	t.updatedFieldsArray = [100]interface{}{}
	return t
}

func (t *IndexTopicData) setMultiplierAndPrec() {
	if t.updatedFieldsArray[INDEX_INDEX["PRECISION"]] != nil {
		t.precision = t.fieldDataArray[INDEX_INDEX["PRECISION"]].(int)
		t.precisionValue = int(math.Pow(10, float64(t.precision)))
	}
	if t.updatedFieldsArray[INDEX_INDEX["MULTIPLIER"]] != nil {
		t.multiplier = t.fieldDataArray[INDEX_INDEX["MULTIPLIER"]].(int)
	}
}

func (t *IndexTopicData) prepareData(reqType interface{}) map[string]interface{} {
	t.prepareCommonData()
	if t.updatedFieldsArray[INDEX_INDEX["LTP"]] != nil || t.updatedFieldsArray[INDEX_INDEX["CLOSE"]] != nil {
		ltp := t.fieldDataArray[INDEX_INDEX["LTP"]]
		close := t.fieldDataArray[INDEX_INDEX["CLOSE"]]
		if ltp != nil && close != nil {
			change := ltp.(int) - close.(int)
			t.fieldDataArray[INDEX_INDEX["CHANGE"]] = change
			t.updatedFieldsArray[INDEX_INDEX["CHANGE"]] = true
			perChange := round(float64(change)/float64(close.(int))*100, t.precision)
			t.fieldDataArray[INDEX_INDEX["PERCHANGE"]] = perChange
			t.updatedFieldsArray[INDEX_INDEX["PERCHANGE"]] = true
		}
	}

	jsonRes := make(map[string]interface{})
	for index := range IndexMapping {
		dataType := IndexMapping[index]
		val := t.fieldDataArray[index]
		if t.updatedFieldsArray[index] != nil && val != nil && dataType != nil {
			switch dataType.Type {
			case FieldTypes["FLOAT32"]:
				val = round(float64(val.(int))/(float64(t.multiplier)*float64(t.precisionValue)), t.precision)
			case FieldTypes["DATE"]:
				val = getFormatDate(int64(val.(int)))
			}
			jsonRes[dataType.Name] = fmt.Sprintf("%v", val)
		}
	}
	t.updatedFieldsArray = [100]interface{}{}
	if reqType != nil {
		jsonRes["request_type"] = reqType
	}
	return jsonRes
}
