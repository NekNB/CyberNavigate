package convert

func AnyToStringSlice[T any](input *[]T, converter func(T) string) *[]string {
	return AnyToAnySlice(input, converter)
}
func StringToAnySlice[T any](input *[]string, converter func(string) T) *[]T {
	return AnyToAnySlice(input, converter)
}

func AnyToAnySlice[I any, O any](input *[]I, converter func(I) O) *[]O {
	if input == nil {
		return nil
	}
	var output = make([]O, len(*input))
	for i, entity := range *input {
		output[i] = converter(entity)
	}

	return &output
}
