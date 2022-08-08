package types

type Dict map[string]string
type List []string
type Json map[string]any

type Hook func(object any) (error, any)
type HooksDict map[string][]Hook
