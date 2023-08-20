package actions

import (
	"fmt"
	exinfra "src/exeiac/infra"
	extools "src/exeiac/tools"

	"github.com/buger/jsonparser"
)

type JsonValue struct {
	Value     string
	ValueType jsonparser.ValueType
}

type JsonFlatten map[string]JsonValue

type JsonValueCouple struct {
	V1 JsonValue
	V2 JsonValue
}

type JsonCompared map[string]JsonValueCouple

type ChangedOutputs map[string]JsonCompared

func (jv JsonValue) String() string {
	switch jv.ValueType {
	case jsonparser.NotExist, jsonparser.Null:
		return "Null"
	case jsonparser.String:
		return fmt.Sprintf("\"%s\"", jv.Value)
	case jsonparser.Boolean, jsonparser.Number:
		return string(jv.Value)
	default:
		return "TYPE_ERROR"
	}
}

func (jvc JsonValueCouple) String() string {
	return fmt.Sprintf("%s -> %s", jvc.V1, jvc.V2)
}

func FlattenJson(json []byte, parentPath string) (flattenedJson JsonFlatten, err error) {
	flattenedJson = make(map[string]JsonValue)

	if json == nil {
		flattenedJson[""] = JsonValue{Value: "", ValueType: jsonparser.Null}

		return
	}

	value, dataType, _, err := jsonparser.Get(json)
	if err != nil {
		return
	}

	switch dataType {
	case jsonparser.Object:
		jsonparser.ObjectEach(json, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			subKey := parentPath + "." + string(key)
			switch dataType {
			case jsonparser.NotExist, jsonparser.Null:
				flattenedJson[subKey] = JsonValue{Value: "", ValueType: dataType}
			case jsonparser.String, jsonparser.Boolean, jsonparser.Number:
				flattenedJson[subKey] = JsonValue{Value: string(value), ValueType: dataType}
			case jsonparser.Object, jsonparser.Array:
				subFlattenedJson, err := FlattenJson(value, subKey)
				if err != nil {
					return err
				}
				for k, v := range subFlattenedJson {
					flattenedJson[k] = v
				}
			default:
				return fmt.Errorf("unrecognized dataType %s for key %s", dataType.String(), subKey)
			}
			return nil
		})
	case jsonparser.Array:
		index := 0
		jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			subKey := fmt.Sprintf("%s[%d]", parentPath, index)
			switch dataType {
			case jsonparser.NotExist, jsonparser.Null:
				flattenedJson[subKey] = JsonValue{Value: "", ValueType: dataType}
			case jsonparser.String, jsonparser.Boolean, jsonparser.Number:
				flattenedJson[subKey] = JsonValue{Value: string(value), ValueType: dataType}
			case jsonparser.Object, jsonparser.Array:
				subFlattenedJson, err := FlattenJson(value, subKey)
				if err != nil {
					return
				}
				for k, v := range subFlattenedJson {
					flattenedJson[k] = v
				}
			default:
				return
			}
			index++
		})
	case jsonparser.NotExist, jsonparser.Null:
		flattenedJson[parentPath] = JsonValue{Value: "", ValueType: dataType}
	case jsonparser.String, jsonparser.Boolean, jsonparser.Number:
		flattenedJson[parentPath] = JsonValue{Value: string(value), ValueType: dataType}
	default:
		err = fmt.Errorf("unrecognized dataType %s for key %s", dataType.String(), parentPath)
		return
	}

	return
}

func CompareFlattenJsons(json1 *JsonFlatten, json2 *JsonFlatten) (comparedJson JsonCompared, areEqual bool) {
	areEqual = true
	comparedJson = make(JsonCompared)

	for k, v1 := range *json1 {
		v2, exist := (*json2)[k]
		if !exist {
			v2 = JsonValue{Value: "", ValueType: jsonparser.NotExist}
			areEqual = false
			(*json2)[k] = v2
			comparedJson["$"+k] = JsonValueCouple{V1: v1, V2: v2}
		} else if v1 != v2 {
			areEqual = false
			comparedJson["$"+k] = JsonValueCouple{V1: v1, V2: v2}
		}
	}

	for k, v2 := range *json2 {
		v1, exist := (*json1)[k]
		if !exist {
			v1 = JsonValue{Value: "", ValueType: jsonparser.NotExist}
			areEqual = false
			(*json1)[k] = v1
			comparedJson["$"+k] = JsonValueCouple{V1: v1, V2: v2}
		}
	}

	return
}

func CompareJsons(json1 []byte, json2 []byte) (comparedJson JsonCompared, areEqual bool, err error) {
	var flatJson1 JsonFlatten
	var flatJson2 JsonFlatten

	flatJson1, err = FlattenJson(json1, "")
	if err != nil {
		err = fmt.Errorf("pre-lay json output: %v", err)

		return
	}

	flatJson2, err = FlattenJson(json2, "")
	if err != nil {
		err = fmt.Errorf("after-lay json output: %v", err)

		return
	}

	comparedJson, areEqual = CompareFlattenJsons(&flatJson1, &flatJson2)

	return
}

func (changes ChangedOutputs) NeedToLayBrick(brick *exinfra.Brick) bool {
	for _, i := range brick.Inputs {
		if _, exist := changes[i.Dependency.From.Brick.Name]; exist {
			for jsonpath := range changes[i.Dependency.From.Brick.Name] {
				if extools.AreJsonPathsLinked(jsonpath, i.Dependency.From.JsonPath) {
					return true
				}
			}
		}
	}

	return false
}
