package evm

import (
	"evm/util"
	"runtime"
	"sync"
)

// This file defines some parallel interfaces which allow the user parallel run operations

// CallParameter define parameter of Call to simplfiy ParallelCall
type CallParameter struct {
	Caller Address
	Callee Address
	Code   []byte
	Gas    uint64
	Input  []byte
	Value  uint64
}

// CallResult define result of Call to simplifiy ParallelCall
type CallResult struct {
	Output []byte
	Err    error
}

// ParallelCall provide an interface for user parallel call
func ParallelCall(bc Blockchain, db DB, ctx *Context, params []*CallParameter) []*CallResult {
	var result = make([]*CallResult, len(params))
	// first we will parallel run codes
	threadSize := runtime.NumCPU()
	if threadSize < 2 {
		threadSize = 2
	}
	var ch = make(chan bool, threadSize)
	var wg sync.WaitGroup
	var vms = make([]*EVM, len(params))
	for i := range params {
		wg.Add(1)
		ch <- true
		go func(i int) {
			defer func() {
				<-ch
				wg.Done()
			}()
			var param = params[i]
			var gas = param.Gas
			// TODO: fill context content
			var context = &Context{
				Input: param.Input,
				Value: param.Value,
				Gas:   &gas,
			}
			vms[i] = New(bc, db, context)
			vms[i].sync = false
			output, err := vms[i].Call(param.Caller, param.Callee, param.Code)
			result[i] = &CallResult{
				Output: output,
				Err:    err,
			}
		}(i)
	}
	wg.Wait()
	// then we find out conflicts
	var setMap = make(map[string]bool)
	// TODO: If we re run txs we need to make sure that rerun will get right answer
	for i := range vms {
		cache := vms[i].cache
		// Note: We regard account update will cause potential conflict, which is not accuracy.
		// TODO: We may do this better.
		if !cache.accountUpdate {
			var isConflict = false
			for key := range cache.reads {
				if util.Contain(setMap, key) {
					isConflict = true
					break
				}
			}
			if !isConflict {
				for i := range cache.sets {
					setMap[i] = true
				}
				continue
			}
		}
		// log.Infof("rerun tx %d because of conflict", i)
		// we need to rerun
		var gas = params[i].Gas
		var context = &Context{
			Input: params[i].Input,
			Value: params[i].Value,
			Gas:   &gas,
		}
		vms[i] = New(bc, db, context)
		vms[i].sync = false
		param := params[i]
		output, err := vms[i].Call(param.Caller, param.Callee, param.Code)
		result[i] = &CallResult{
			Output: output,
			Err:    err,
		}
	}
	// then we will sync changes to db
	for i := range vms {
		if result[i].Err == nil {
			vms[i].cache.Sync()
		}
	}
	return result
}
