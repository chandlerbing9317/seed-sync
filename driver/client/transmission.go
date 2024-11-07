package client

// transmission交互层
// ref:
// https://github.com/transmission/transmission/blob/4.0.5/docs/rpc-spec.md
// https://github.com/hekmon/transmissionrpc

const csrfHeader = "X-Transmission-Session-Id"

type requestBody struct {
	Method    string      `json:"method"`
	Arguments interface{} `json:"arguments"`
	Tag       int         `json:"tag,omitempty"`
}

type responseBody struct {
	Arguments interface{} `json:"arguments"`
	Result    string      `json:"result"`
	Tag       *int        `json:"tag"`
}
