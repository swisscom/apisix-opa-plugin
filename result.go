package main

/*
	Result is a type representing the struct returned as a result of OPA's
	policy evaluation.

	The main parts of it are the status code and optionally a reason.
	If the status code is 0, the status code will not be set (original will be kept) and
	the request will be forwarded.

	In any other case, the request is blocked and the "reason" is returned (if provided)
*/
type Result struct {
	StatusCode int16  `json:"status_code"`
	Reason     string `json:"reason,omitempty"`
}

type Metrics struct {
	TimerRegoInputParseNs   int `json:"timer_rego_input_parse_ns"`
	TimerRegoQueryParseNs   int `json:"timer_rego_query_parse_ns"`
	TimerRegoQueryCompileNs int `json:"timer_rego_query_compile_ns"`
	TimerRegoQueryEvalNs    int `json:"timer_rego_query_eval_ns"`
	TimerRegoModuleParseNs  int `json:"timer_rego_module_parse_ns"`
	TimerServerHandlerNs    int `json:"timer_server_handler_ns"`
}

type OpaResponse struct {
	Result     Result  `json:"result"`
	Metrics    Metrics `json:"metrics"`
	DecisionId string  `json:"decision_id,omitempty"`
}
