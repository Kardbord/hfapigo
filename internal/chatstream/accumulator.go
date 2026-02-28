// Package chatstream provides helpers for working with streamed chat responses.
package chatstream

// ToolCallAccumulator tracks the last known metadata for a streamed tool call.
// It fills in missing id/type/function-name fields across incremental deltas so
// callers always observe complete tool-call identifiers.
type ToolCallAccumulator struct {
	cache map[key]*toolCallState
}

// key identifies a tool call within a streaming response choice.
type key struct {
	choice int
	call   int
}

// toolCallState holds the sticky metadata captured for a tool call index.
type toolCallState struct {
	ID           string
	Type         string
	FunctionName string
}

// Merge records any provided metadata for a tool call and returns the
// accumulated values for that (choice, call) pair.
func (a *ToolCallAccumulator) Merge(
	choiceIndex int,
	callIndex int,
	toolID string,
	toolType string,
	functionName string,
) (finalID, finalType, finalName string) {
	state := a.state(choiceIndex, callIndex)
	if toolID != "" {
		state.ID = toolID
	}
	if toolType != "" {
		state.Type = toolType
	}
	if functionName != "" {
		state.FunctionName = functionName
	}

	return state.ID, state.Type, state.FunctionName
}

// state returns the cached metadata for the provided tool call, creating a new
// entry when this (choice, call) pair is observed for the first time.
func (a *ToolCallAccumulator) state(choiceIndex, callIndex int) *toolCallState {
	if a.cache == nil {
		a.cache = make(map[key]*toolCallState)
	}
	cacheKey := key{choice: choiceIndex, call: callIndex}
	if state, ok := a.cache[cacheKey]; ok {
		return state
	}
	state := &toolCallState{
		ID:           "",
		Type:         "",
		FunctionName: "",
	}
	a.cache[cacheKey] = state

	return state
}
