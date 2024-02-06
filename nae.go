package ear

import (
	"errors"
	"fmt"
)

type NAETTSInfo struct {
	SessionID      *string `json:"sessionid"`
	Infrastructure *string `json:"infrastructure"`
	Identity       *string `json:"identity"`
}

func ToNAETTSInfo(v interface{}) (*NAETTSInfo, error) {
	vMap, ok := v.(map[string]interface{})
	if !ok {
		return nil, errors.New(`unexpected format for "tee-info"`)
	}

	var info NAETTSInfo

	for key, val := range vMap {
		s := str(val)
		switch key {
		case "sessionid":
			info.SessionID = &s
		case "infrastructure":
			info.Infrastructure = &s
		case "identity":
			info.Identity = &s
		default:
			return nil, fmt.Errorf(`found unknown key %q in "naetts" object`, key)
		}
	}
	return &info, nil
}
