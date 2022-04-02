package swap

import "errors"

type TradeType int

const (
	TradeTypeExactInput  TradeType = 0
	TradeTypeExactOutput TradeType = 1
)

type Trade struct {
	Route        Route
	TradeType    TradeType
	InputAmount  TokenAmount
	OutputAmount TokenAmount
}

func NewTrade(route Route, amount TokenAmount, tradeType TradeType) (*Trade, error) {
	amounts := make([]TokenAmount, len(route.Path))

	var inputAmount, outputAmount TokenAmount
	if tradeType == TradeTypeExactInput {
		amounts[0] = amount
		for i := 0; i < len(route.Path)-1; i++ {
			tokenAmount, err := route.Pairs[i].GetOutputAmount(amounts[i])
			if err != nil {
				return nil, err
			}

			amounts[i+1] = tokenAmount
		}

		inputAmount, outputAmount = amount, amounts[len(amounts)-1]
	} else {
		amounts[len(amounts)-1] = amount
		for i := len(route.Path) - 1; i > 0; i-- {
			tokenAmount, err := route.Pairs[i-1].GetInputAmount(amounts[i])
			if err != nil {
				return nil, err
			}

			amounts[i-1] = tokenAmount
		}

		outputAmount, inputAmount = amount, amounts[0]
	}

	if inputAmount.Amount <= 0 || outputAmount.Amount <= 0 {
		return nil, ErrInsufficientReserve
	}

	return &Trade{
		Route:        route,
		TradeType:    tradeType,
		InputAmount:  inputAmount,
		OutputAmount: outputAmount,
	}, nil
}

func NewAccurateTrade(route Route, amount TokenAmount, tradeType TradeType) (*Trade, error) {
	amounts := make([]TokenAmount, len(route.Path))

	var inputAmount, outputAmount TokenAmount
	if tradeType == TradeTypeExactInput {
		amounts[0] = amount
		for i := 0; i < len(route.Path)-1; i++ {
			tokenAmount, err := route.Pairs[i].GetOutputAmountInPip(amounts[i])
			if err != nil {
				return nil, err
			}

			amounts[i+1] = tokenAmount
		}

		inputAmount, outputAmount = amount, amounts[len(amounts)-1]
	} else {
		amounts[len(amounts)-1] = amount
		for i := len(route.Path) - 1; i > 0; i-- {
			tokenAmount, err := route.Pairs[i-1].GetInputAmountInPip(amounts[i])
			if err != nil {
				return nil, err
			}

			amounts[i-1] = tokenAmount
		}

		outputAmount, inputAmount = amount, amounts[0]
	}

	if outputAmount.GetAmountInPip().Sign() != 1 || inputAmount.GetAmountInPip().Sign() != 1 {
		return nil, nil
	}

	return &Trade{
		Route:        route,
		TradeType:    tradeType,
		InputAmount:  inputAmount,
		OutputAmount: outputAmount,
	}, nil
}

func inputOutputComparator(tradeA, tradeB *Trade) int {
	if tradeA.OutputAmount.GetAmount() == tradeB.OutputAmount.GetAmount() {
		if tradeA.InputAmount.GetAmount() == tradeB.InputAmount.GetAmount() {
			return 0
		}

		// trade A requires less input than trade B, so A should come first
		if tradeA.InputAmount.GetAmount() < tradeB.InputAmount.GetAmount() {
			return -1
		} else {
			return 1
		}
	} else {
		// tradeA has less output than trade B, so should come second
		if tradeA.OutputAmount.GetAmount() < tradeB.OutputAmount.GetAmount() {
			return 1
		} else {
			return -1
		}
	}
}

func tradeComparator(tradeA, tradeB *Trade) bool {
	ioComp := inputOutputComparator(tradeA, tradeB)
	if ioComp != 0 {
		return ioComp == 1
	}

	// finally consider the number of hops since each hop costs gas
	if len(tradeA.Route.Path) > len(tradeB.Route.Path) {
		return false
	}

	return true
}

func GetBestTradeExactIn(pairs []*PairTrade, currencyOut Token, currencyAmountIn TokenAmount, maxHops int) (*Trade, error) {
	trade, err := getBestTradeExactIn(pairs, currencyOut, currencyAmountIn, maxHops, make([]*PairTrade, 0), currencyAmountIn, nil)
	if err != nil {
		return nil, err
	}

	if trade == nil {
		return nil, errors.New("path not found")
	}

	return NewAccurateTrade(trade.Route, currencyAmountIn, TradeTypeExactInput)
}

func getBestTradeExactIn(
	pairs []*PairTrade,
	currencyOut Token,
	currencyAmountIn TokenAmount,
	maxHops int,
	currentPairs []*PairTrade,
	originalAmountIn TokenAmount,
	bestTrade *Trade,
) (*Trade, error) {
	if maxHops <= 0 {
		return bestTrade, nil
	}

	tokenOut, tokenAmountIn := currencyOut, currencyAmountIn

	for i, pair := range pairs {
		if !pair.Token0.Token.IsEqual(tokenAmountIn.Token) && !pair.Token1.Token.IsEqual(tokenAmountIn.Token) {
			continue
		}

		if maxHops == 1 && !pair.Token0.Token.IsEqual(tokenOut) && !pair.Token1.Token.IsEqual(tokenOut) {
			continue
		}

		if pair.getReserve0() <= 0 || pair.getReserve1() <= 0 {
			continue
		}

		amountOut, err := pair.GetOutputAmount(tokenAmountIn)
		if err != nil {
			if err == ErrInsufficientReserve {
				continue
			}
			return bestTrade, err
		}

		// we have arrived at the output token, so this is the final trade of one of the paths
		if amountOut.Token.IsEqual(tokenOut) {
			trade, err := NewTrade(
				NewRoute(append(currentPairs, pair), originalAmountIn.GetCurrency(), &currencyOut),
				originalAmountIn,
				TradeTypeExactInput,
			)

			if err != nil {
				continue
			}

			if bestTrade == nil || tradeComparator(bestTrade, trade) {
				bestTrade = trade
			}
		} else if maxHops > 1 && len(pairs) > 1 {
			// otherwise, consider all the other paths that lead from this token as long as we have not exceeded maxHops
			temp := make([]*PairTrade, len(pairs))
			copy(temp, pairs)
			temp[i] = temp[len(temp)-1]
			pairsExcludingThisPair := temp[:len(temp)-1]

			newCurrentPairs := append(currentPairs, pair)

			var err error
			bestTrade, err = getBestTradeExactIn(
				pairsExcludingThisPair,
				currencyOut,
				amountOut,
				maxHops-1,
				newCurrentPairs,
				originalAmountIn,
				bestTrade,
			)

			if err != nil {
				return nil, err
			}
		}
	}

	return bestTrade, nil
}

func GetBestTradeExactOut(pairs []*PairTrade, currencyIn Token, amountOut TokenAmount, maxhops int) (*Trade, error) {
	trade, err := getBestTradeExactOut(pairs, currencyIn, amountOut, maxhops, make([]*PairTrade, 0), amountOut, nil)
	if err != nil {
		return nil, err
	}

	if trade == nil {
		return nil, errors.New("path not found")
	}

	return NewAccurateTrade(trade.Route, amountOut, TradeTypeExactOutput)
}

func getBestTradeExactOut(
	pairs []*PairTrade,
	currencyIn Token,
	currencyAmountOut TokenAmount,
	maxHops int,
	currentPairs []*PairTrade,
	originalAmountOut TokenAmount,
	bestTrade *Trade,
) (*Trade, error) {
	if maxHops <= 0 {
		return bestTrade, nil
	}

	tokenIn, amountOut, currencyOut := currencyIn, currencyAmountOut, originalAmountOut.GetCurrency()

	for i, pair := range pairs {
		if !pair.Token0.Token.IsEqual(amountOut.Token) && !pair.Token1.Token.IsEqual(amountOut.Token) {
			continue
		}

		if maxHops == 1 && !pair.Token0.Token.IsEqual(tokenIn) && !pair.Token1.Token.IsEqual(tokenIn) {
			continue
		}

		if pair.getReserve0() <= 0 || pair.getReserve1() <= 0 {
			continue
		}

		amountIn, err := pair.GetInputAmount(amountOut)
		if err != nil {
			if err == ErrInsufficientReserve {
				continue
			}

			return nil, err
		}

		if amountIn.Token.IsEqual(tokenIn) {
			trade, err := NewTrade(
				NewRoute(append([]*PairTrade{pair}, currentPairs...), currencyIn, &currencyOut),
				originalAmountOut,
				TradeTypeExactOutput,
			)

			if err != nil {
				continue
			}

			if bestTrade == nil || tradeComparator(bestTrade, trade) {
				bestTrade = trade
			}
		} else if maxHops > 1 && len(pairs) > 1 {
			// otherwise, consider all the other paths that lead from this token as long as we have not exceeded maxHops
			temp := make([]*PairTrade, len(pairs))
			copy(temp, pairs)
			temp[i] = temp[len(temp)-1]
			pairsExcludingThisPair := temp[:len(temp)-1]

			newCurrentPairs := append([]*PairTrade{pair}, currentPairs...)

			var err error
			bestTrade, err = getBestTradeExactOut(
				pairsExcludingThisPair,
				currencyIn,
				amountIn,
				maxHops-1,
				newCurrentPairs,
				originalAmountOut,
				bestTrade,
			)

			if err != nil {
				return nil, err
			}
		}
	}

	return bestTrade, nil
}
