package ear

type NAETSSInfo struct {
	SessionID      *string `json:"sessionid"`
	Infrastructure *string `json:"infrastructure"`
	Identity       *string `json:"identity"`
}
