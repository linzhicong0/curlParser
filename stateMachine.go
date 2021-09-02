package curlParser

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
)

const (
	RequestStateName = "RequestState"
	HeadersStateName = "HeadersState"
)

// This interface contain the method for state
type IState interface {
	Handler(text string, request *CurlRequest) (error, IState)
	AddNextState(name string, state IState) error
}

// This struct describe the state
type state struct {
	name      string
	nextState map[string]IState
}

// This is for parsing the header
type headerState struct {
	state
}

// This is for parsing the request url and query params
type requestState struct {
	state
}

func (hs *headerState) Handler(text string, request *CurlRequest) (error, IState) {

	expr := `'(.*?')`
	re, _ := regexp.Compile(expr)
	result := string(re.Find([]byte(text)))
	result = result[1 : len(result)-1]

	if request.Headers == nil {
		request.Headers = make(map[string]interface{})
	}

	splitString := strings.Split(result, ":")
	headerName := strings.TrimSpace(splitString[0])
	headerValue := strings.TrimSpace(splitString[1])

	if _, ok := request.Headers[headerName]; !ok {
		request.Headers[headerName] = headerValue
	}
	return nil, hs.nextState[HeadersStateName]
}

func (rs *requestState) Handler(text string, request *CurlRequest) (error, IState) {
	expr := `'(.*?')`
	re, _ := regexp.Compile(expr)
	result := string(re.Find([]byte(text)))
	result = result[1 : len(result)-1]
	urlPath, err := url.Parse(result)
	if err != nil {
		return err, nil
	}

	request.Host = urlPath.Host
	request.Path = urlPath.Path
	request.Port = urlPath.Port()

	if request.QueryParams == nil {
		request.QueryParams = make(map[string]interface{})
	}

	for k, v := range urlPath.Query() {
		request.QueryParams[k] = v
	}
	return nil, rs.nextState[HeadersStateName]
}

func (rs *requestState) AddNextState(name string, state IState) error {
	if rs.nextState == nil {
		rs.nextState = make(map[string]IState)
	}

	if _, ok := rs.nextState[name]; ok {
		return fmt.Errorf("the state already existed: %s", name)
	}
	rs.nextState[name] = state
	return nil
}

func (hs *headerState) AddNextState(name string, state IState) error {
	if hs.nextState == nil {
		hs.nextState = make(map[string]IState)
	}

	if _, ok := hs.nextState[name]; ok {
		return fmt.Errorf("the state already existed: %s", name)
	}
	hs.nextState[name] = state
	return nil
}

type stateMachine struct {
	states map[string]IState
}

// Add the given state to the state machine
func (sm *stateMachine) AddState(name string, state IState, nextStateName string, nextState IState) error {
	if _, ok := sm.states[name]; ok {
		return sm.states[name].AddNextState(nextStateName, nextState)
	}
	state.AddNextState(nextStateName, nextState)
	sm.states[name] = state
	return nil
}

// Start the state machine
func (sm *stateMachine) Start(name string, curlText string) (error, *CurlRequest) {

	if _, ok := sm.states[name]; !ok {
		return fmt.Errorf("can not find state: %s", name), nil
	}

	currentState := sm.states[name]
	curlRequest := new(CurlRequest)

	reader := bufio.NewReader(strings.NewReader(curlText))

	for {

		// Read line from curl cmd
		buffer, _, err := reader.ReadLine()
		if err == io.EOF {
			return nil, curlRequest
		}
		if err != nil {
			return err, nil
		}
		line := string(buffer)

		// Handle the line
		err, nextState := currentState.Handler(line, curlRequest)
		if err != nil {
			return err, nil
		}
		if nextState != nil {
			currentState = nextState
		} else {
			return nil, curlRequest
		}
	}

}
