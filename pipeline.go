package mediator

// buildPipeline iteratively applies a function to combine a seed value with elements from a behavior slice into a final value.
func buildPipeline[TBehavior any, TFunc any](behaviors []TBehavior, seed TFunc, f func(TFunc, TBehavior) TFunc) TFunc {
	result := seed
	for _, behavior := range behaviors {
		result = f(result, behavior)
	}
	return result
}
