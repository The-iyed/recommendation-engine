package pkg

import "fmt"

func InterfaceToFloat64Slice(input any) ([]float64, error) {
	switch v := input.(type) {
	case []interface{}:
		result := make([]float64, len(v))
		for i, el := range v {
			if f, ok := el.(float64); ok {
				result[i] = f
			} else {
				return nil, fmt.Errorf("element at index %d is not a float64", i)
			}
		}
		return result, nil
	case []float64:
		result := make([]float64, len(v))
		for i, el := range v {
			result[i] = el
		}
		return result, nil
	default:
		return nil, fmt.Errorf("expected []interface{}, got %T", v)
	}
}
