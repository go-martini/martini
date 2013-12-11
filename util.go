package martini

type zeroInterface interface {
	IsZero() bool
}

func isZero(arg interface{}) bool {
	var is bool
	switch arg := arg.(type) {
	case string:
		is = (len(arg) == 0)
	case bool:
		is = (arg == false)
	case int:
		is = (arg == 0)
	case int8:
		is = (arg == 0)
	case int16:
		is = (arg == 0)
	case int64:
		is = (arg == 0)
	case uint:
		is = (arg == 0)
	case uint8:
		is = (arg == 0)
	case uint16:
		is = (arg == 0)
	case uint64:
		is = (arg == 0)
	case zeroInterface:
		is = arg.IsZero()
	default:
		is = (arg == nil)
	}
	return is
}
