package convert

func AnyToStringSlice[T any](input *[]T, converter func(T) string) *[]string {
	var output = make([]string, len(*input))

	for i, entity := range *input {
		output[i] = converter(entity)
	}

	return &output
}
func StringToAnySlice[T any](input *[]string, converter func(string) T) *[]T {
	var output = make([]T, len(*input))

	for i, entity := range *input {
		output[i] = converter(entity)
	}

	return &output
}

func AnyToAnySlice[I any, O any](input *[]I, converter func(I) O) *[]O {
	var output = make([]O, len(*input))
	for i, entity := range *input {
		output[i] = converter(entity)
	}

	return &output
}
