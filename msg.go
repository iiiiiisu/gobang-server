package main

type RequestMsg struct {
	Cmd  string `json:"cmd"`
	Data string `json:"data"`
}

type ResponseMsg struct {
	Msg    string `json:"msg"`
	Result string `json:"result"`
}
