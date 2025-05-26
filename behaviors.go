package cqrs

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
)

type NextFunc func() (interface{}, error)

type IBehavior interface {
	Handle(ctx context.Context, request interface{}, next NextFunc) (interface{}, error)
}

var (
	commandBehaviors map[int]interface{}
	queryBehaviors   map[int]interface{}
	behaviorMu       sync.RWMutex
)

func init() {
	commandBehaviors = make(map[int]interface{})
	queryBehaviors = make(map[int]interface{})
}

func RegisterCommandBehavior(order int, behavior IBehavior) error {
	behaviorMu.Lock()
	defer behaviorMu.Unlock()

	_, found := commandBehaviors[order]
	if found {
		msg := fmt.Sprintf("position %d is taken by another command behavior.", order)
		return errors.New(msg)
	}

	commandBehaviors[order] = behavior
	return nil
}

func RegisterQueryBehavior(order int, behavior IBehavior) error {
	behaviorMu.Lock()
	defer behaviorMu.Unlock()

	_, found := queryBehaviors[order]
	if found {
		msg := fmt.Sprintf("position %d is taken by another query behavior.", order)
		return errors.New(msg)
	}

	queryBehaviors[order] = behavior
	return nil
}

func sortBehaviors(behaviors map[int]interface{}) []interface{} {
	keys := make([]int, 0)

	behaviorMu.RLock()
	defer behaviorMu.RUnlock()

	for key := range behaviors {
		keys = append(keys, key)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	sorted := make([]interface{}, 0)

	for _, key := range keys {
		sorted = append(sorted, behaviors[key])
	}

	return sorted
}
