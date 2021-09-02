package curlParser

type CurlParser struct{}

func (cp *CurlParser) Parse(curlText string) (error, *CurlRequest) {

	sm := new(stateMachine)
	sm.states = make(map[string]IState)
	requestState := new(requestState)
	headersState := new(headerState)
	sm.AddState(RequestStateName, requestState, HeadersStateName, headersState)
	sm.AddState(HeadersStateName, headersState, HeadersStateName, headersState)

	return sm.Start(RequestStateName, curlText)
}
