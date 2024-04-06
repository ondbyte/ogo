package ogo

import "strconv"

func (v *RequestValidator[resData, respData]) setInt64(val string, i *int64) error {
	vv, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return err
	}
	*i = vv
	return nil
}

func (v *RequestValidator[resData, respData]) setBool(val string, i *bool) error {
	vv, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	*i = vv
	return nil
}

func (v *RequestValidator[resData, respData]) setFloat64(val string, i *float64) error {
	vv, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}
	*i = vv
	return nil
}

func (v *RequestValidator[resData, respData]) setUint64(val string, i *uint64) error {
	vv, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return err
	}
	*i = vv
	return nil
}
