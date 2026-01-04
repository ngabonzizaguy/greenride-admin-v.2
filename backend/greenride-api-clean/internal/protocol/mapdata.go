package protocol

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

type MapData map[string]any

type MapArrays []MapData

func (t MapData) GetStringArray(field string) []string {
	if _v, ok := t[field]; ok {
		switch val := _v.(type) {
		case []string:
			return val
		}
	}
	return []string{}
}

func (t MapArrays) Map(handler func(r MapData) bool) (data MapArrays) {
	data = MapArrays{}
	if len(t) == 0 {
		return
	}
	for _, r := range t {
		if handler(r) {
			data = append(data, r)
		}
	}
	return
}

func (t MapArrays) ForEach(handler func(r MapData)) {
	for _, r := range t {
		handler(r)
	}
}

func (t MapArrays) ToJson() string {
	json, _ := json.Marshal(t)
	return string(json)
}

func NewMapData(src map[string]interface{}) (data MapData) {
	data = MapData{}
	data.Copy(src)
	return
}

func (t MapData) Set(field string, val interface{}) {
	t[field] = val
}

func (t MapData) CopyObject(obj interface{}) {
	_req, _ := json.Marshal(obj)
	t.Copy(cast.ToStringMap(string(_req)))
}
func (t MapData) Copy(data map[string]interface{}) {
	for k, v := range data {
		t.Set(k, v)
	}
}

func (t MapData) Combind(target MapData) {
	if len(target) == 0 {
		return
	}
	copy := NewMapData(target)
	copy.Copy(t)
	t.Copy(copy)
}
func (t MapData) CopyStringMap(data map[string]string) {
	for k, v := range data {
		t.Set(k, v)
	}
}

func (t MapData) Keys() (list []string) {
	list = []string{}
	for k := range t {
		list = append(list, k)
	}
	return list
}

func (t MapData) Reomve(key string) {
	delete(t, key)
}

func (t MapData) GetDatetime(field string) *time.Time {
	_v := t.Get(field)
	if _v != "" {
		_time, err := time.Parse(time.DateTime, _v)
		if err == nil {
			return &_time
		}
	}
	return nil
}
func (t MapData) GetQueryListHandler(field string) func(MapData) MapArrays {
	if v, ok := t[field]; ok {
		switch val := v.(type) {
		case func(MapData) MapArrays:
			return val
		}
	}
	return nil
}

func (t MapData) Get(field string) string {
	if v, ok := t[field]; ok {
		switch val := v.(type) {
		case []MapData:
			rs, _ := json.Marshal(val)
			return string(rs)
		case MapData:
			return val.ToJson()
		default:
			return strings.TrimSpace(cast.ToString(val))
		}
	}
	return ""
}

func (t MapData) GetUpper(field string) string {
	return strings.ToUpper(t.Get(field))
}

func (t MapData) GetLower(field string) string {
	return strings.ToLower(t.Get(field))
}
func (t MapData) GetDecimal(field string) *decimal.Decimal {
	if v, ok := t[field]; ok {
		_v, ok := v.(*decimal.Decimal)
		if ok {
			return _v
		}
		dv, err := decimal.NewFromString(cast.ToString(v))
		if err == nil {
			return &dv
		}
	}
	return nil
}

func (t MapData) GetBytes(field string) []byte {
	if v, ok := t[field]; ok {
		return []byte(cast.ToString(v))
	}
	return []byte{}
}

func (t MapData) GetInt64(field string) int64 {

	if v, ok := t[field]; ok {
		return cast.ToInt64(v)
	}
	return 0
}

func (t MapData) GetMillSecondtimeToTime(field string) (rs *time.Time) {

	if v, ok := t[field]; ok {
		tm := cast.ToInt64(v) / 1000
		if tm > 0 {
			_time := time.Unix(tm, 0)
			rs = &_time
			return
		}
		return
	}
	return
}
func (t MapData) GetUnixtimeToTime(field string) (rs *time.Time) {

	if v, ok := t[field]; ok {
		tm := cast.ToInt64(v)
		if tm > 0 {
			_time := time.Unix(tm, 0)
			rs = &_time
			return
		}
		return
	}
	return
}
func (t MapData) GetUnixtimeToDateTimeFormat(field string) string {

	if v, ok := t[field]; ok {
		tm := cast.ToInt64(v)
		if tm > 0 {
			_time := time.Unix(tm, 0).Format(time.DateTime)
			return _time
		}
		return ""
	}
	return ""
}

func (t MapData) GetInt(field string) int {
	if v, ok := t[field]; ok {
		return cast.ToInt(v)
	}
	return 0
}

func (t MapData) GetFloat64(field string) float64 {

	if v, ok := t[field]; ok {
		return cast.ToFloat64(v)
	}
	return 0
}
func (t MapData) GetBool(field string) bool {
	if v, ok := t[field]; ok {
		return cast.ToBool(v)
	}
	return false
}
func (t MapData) GetObject(field string) interface{} {

	if v, ok := t[field]; ok {
		return v
	}
	return nil
}

func (t MapData) GetHandler(field string) http.HandlerFunc {

	if v, ok := t[field]; ok {
		return v.(http.HandlerFunc)
	}
	return nil
}

func (t MapData) GetArrayFromString(field string) []string {
	v := t.Get(field)
	if v == "" {
		return []string{}
	}
	return strings.Split(v, ",")
}
func (t MapData) GetArrayMapDataFromJson(field string) MapArrays {
	if v, ok := t[field]; ok {
		rs := MapArrays{}
		json.Unmarshal([]byte(cast.ToString(v)), &rs)
		return rs
	}
	return []MapData{}
}

func (t MapData) GetMapDataFromJson(field string) MapData {
	if v, ok := t[field]; ok {
		rs := MapData{}
		json.Unmarshal([]byte(cast.ToString(v)), &rs)
		return rs
	}
	return MapData{}
}

func (t MapData) GetStringArrayFromJson(field string) []string {
	if v, ok := t[field]; ok {
		rs := []string{}
		json.Unmarshal([]byte(cast.ToString(v)), &rs)
		return rs
	}
	return []string{}
}

func (t MapData) GetMapArray(field string) MapArrays {
	rs := MapArrays{}
	if v, ok := t[field]; ok {
		switch val := v.(type) {
		case string:
			json.Unmarshal([]byte(val), &rs)
			return rs
		case []byte:
			json.Unmarshal(val, &rs)
		case []interface{}:
			rs = MapArrays{}
			for _, _v := range val {
				rs = append(rs, cast.ToStringMap(_v))
			}
		default:
			rs = v.(MapArrays)
		}
	}
	return rs
}

func (t MapData) ToQueryUrl() string {
	if len(t) == 0 {
		return ""
	}
	params := []string{}
	for k, v := range t {
		params = append(params, fmt.Sprintf("%v=%v", k, v))
	}
	return strings.Join(params, "&")
}
func (t MapData) Has(field string) bool {
	_, ok := t[field]
	return ok
}
func (t MapData) GetMapData(field string) MapData {
	if v, ok := t[field]; ok {
		switch _v := v.(type) {
		case MapData:
			return _v
		case interface{}:
			mv, ok := v.(*MapData)
			if ok {
				return *mv
			}
			__v, _ := cast.ToStringMapE(v)
			return __v
		}
	}
	return MapData{}
}
func (t MapData) ToString() string {
	return t.ToJson()
}
func (t MapData) ToStringMap() map[string]string {
	data := map[string]string{}
	for k := range t {
		data[k] = t.Get(k)
	}
	return data
}
func (t MapData) ToJson() string {
	_json, _ := json.Marshal(t)
	return string(_json)
}
func (t MapData) GetJsonObject(target_field string, target_obj interface{}) {
	v := t.GetBytes(target_field)
	if len(v) == 0 {
		return
	}
	json.Unmarshal(t.GetBytes(target_field), target_obj)
}

func (t MapData) ToObject(tg interface{}) {
	json_str := t.ToJson()
	err := json.Unmarshal([]byte(json_str), tg)
	if err != nil {
		return
	}
}
func (annotation *MapData) QueryList(val interface{}) error {
	if val == nil {
		return nil
	}
	switch val := val.(type) {
	case string:
		return json.Unmarshal([]byte(val), annotation)
	case []byte:
		return json.Unmarshal(val, annotation)
	default:
		return errors.New("not support")
	}
}

func (annotation MapData) Value() (driver.Value, error) {
	bytes, err := json.Marshal(annotation)
	return string(bytes), err
}

func (annotation *MapData) Scan(val interface{}) error {
	if val == nil {
		return nil
	}
	switch val := val.(type) {
	case string:
		return json.Unmarshal([]byte(val), annotation)
	case []byte:
		return json.Unmarshal(val, annotation)
	default:
		return errors.New("not support")
	}
}

func (t MapData) FirstVal(fields []string) (val string) {
	if len(fields) == 0 {
		return
	}
	for _, f := range fields {
		if _v := t.Get(f); _v != "" {
			val = _v
			return
		}
	}
	return
}
func (t MapData) SortAndJoin(simbo string) string {
	if len(t) <= 0 {
		return ""
	}
	var keys []string
	for k := range t {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	builder := strings.Builder{}
	for _, k := range keys {
		v := t.Get(k)
		if v == "" {
			continue
		}
		builder.WriteString(fmt.Sprintf("%v%s=%s", simbo, k, v))
	}
	return builder.String()[1:]
}

func (t *MapData) AppendToStringList(match_list_key, match_item string) {
	match_list := t.Get(match_list_key)
	if strings.Contains(match_list, match_item) {
		return
	}
	match_list = strings.TrimPrefix(fmt.Sprintf("%v,%v", match_list, match_item), ",")
	t.Set(match_list_key, match_list)
}
