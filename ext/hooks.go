package ext

type HooksDict map[string][]Hook
type Hook func(object any) (error, any)

var defaultHooksList = []string{"request", "client", "response"}

func DefaultHooks() HooksDict {
	hooks := HooksDict{}
	for _, v := range defaultHooksList {
		hooks[v] = make([]Hook, 0)
	}
	return hooks
}

func RegisterHook(hooksDict *HooksDict, key string, hook Hook) error {
	hooks := *hooksDict
	if hooks[key] != nil {
		hooks[key] = append(hooks[key], hook)
	} else {
		hooks[key] = []Hook{hook}
	}
	return nil
}

func DisPatchHook(key string, hooks HooksDict, data any) any {
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
