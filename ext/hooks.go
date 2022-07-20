package ext

import "github.com/sari3l/requests/types"

var defaultHooksList = []string{"request", "client", "response"}

func DefaultHooks() types.HooksDict {
	hooks := types.HooksDict{}
	for _, v := range defaultHooksList {
		hooks[v] = make([]types.Hook, 0)
	}
	return hooks
}

func RegisterHook(hooksDict *types.HooksDict, key string, hook types.Hook) error {
	hooks := *hooksDict
	if hooks[key] != nil {
		hooks[key] = append(hooks[key], hook)
	} else {
		hooks[key] = []types.Hook{hook}
	}
	return nil
}

func DisPatchHook(key string, hooks types.HooksDict, data any) any {
	if hooks[key] != nil {
		for _, hook := range hooks[key] {
			err, _hookData := hook(data)
			if err == nil && _hookData != nil {
				data = _hookData
			}
		}
	}
	return data
}
