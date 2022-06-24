package rsql

func GetValue(value Value) interface{} {
	var v interface{}
	switch value.(type) {
	case StringValue:
		sv, _ := value.(StringValue)
		v = sv.Value
	case IntegerValue:
		sv, _ := value.(IntegerValue)
		v = sv.Value
	case DateValue:
		sv, _ := value.(DateValue)
		v = sv.Value
	case DoubleValue:
		sv, _ := value.(DoubleValue)
		v = sv.Value
	case DateTimeValue:
		sv, _ := value.(DateTimeValue)
		v = sv.Value
	case BooleanValue:
		sv, _ := value.(BooleanValue)
		v = sv.Value
	case ListValue:
		sv, _ := value.(ListValue)
		v = GetValueList(sv)
	default:
		v = value
	}
	return v
}

func GetValueList(listValue ListValue) []interface{} {
	list := make([]interface{}, 0)
	for _, v := range listValue.Value {
		list = append(list, GetValue(v))
	}
	return list
}
