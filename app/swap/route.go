package swap

type Route struct {
	Pairs  []*PairTrade
	Path   []Token
	Input  Token
	Output Token
}

func NewRoute(pairs []*PairTrade, input Token, output *Token) Route {
	path := make([]Token, len(pairs)+1)
	path[0] = input

	for i, pair := range pairs {
		currentInput, currentOutput := path[i], pair.Token0
		if currentInput.IsEqual(pair.Token0.Token) {
			currentOutput = pair.Token1
		}

		path[i+1] = currentOutput.Token
	}

	if output == nil {
		output = new(Token)
		*output = path[len(path)-1]
	}

	route := Route{
		Pairs:  pairs,
		Path:   path,
		Input:  input,
		Output: *output,
	}

	return route
}
