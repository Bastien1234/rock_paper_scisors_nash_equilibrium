package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Trainer struct {
	NumActions      int
	PossibleActions []int
	ActionUtility   [][]int

	RegretSum   []int
	StrategySum []float64

	OppRegretSum   []int
	OppStrategySum []float64
}

func NewTrainer() Trainer {
	t := Trainer{}
	t.NumActions = 3
	t.PossibleActions = []int{0, 1, 2}
	t.ActionUtility = [][]int{
		{0, -1, 1},
		{1, 0, -1},
		{-1, 1, 0},
	}
	t.RegretSum = []int{0, 0, 0}
	t.StrategySum = []float64{0, 0, 0}
	t.OppRegretSum = []int{0, 0, 0}
	t.OppStrategySum = []float64{0, 0, 0}

	return t
}

func GetStrategy(regretSum *[]int) []float64 {

	/*
		Put all elements over 0
		Get biggest number
		Normalize by divising by total number
	*/

	regretFloat := make([]float64, len(*regretSum))
	normalizedSum := 0
	biggestRegretInList := 0

	for index, el := range *regretSum {
		if el < 0 {
			regretFloat[index] = float64(0)
		} else {
			regretFloat[index] = float64(el)

			if el > biggestRegretInList {
				biggestRegretInList = el
			}
		}

		normalizedSum += el
	}

	if normalizedSum > 0 {
		for index, _ := range regretFloat {
			regretFloat[index] /= float64(normalizedSum)
		}
	} else {
		for index, _ := range regretFloat {
			regretFloat[index] += 1.0 / float64(len(*regretSum))
		}
	}

	return regretFloat

}

func GetAction(strategy []float64) int {
	// Get random number
	rand.Seed(time.Now().UnixNano())
	n := float64(rand.Intn(100))

	/* Change that */
	Hstrat := make([]int, len(strategy))
	for i := 0; i < len(strategy); i++ {
		Hstrat[i] = int(strategy[i] * 100)
	}

	index := 0
	cumul := 0

	for _, number := range Hstrat {

		if n == 0 {
			return 0
		}

		cumul += number
		if n <= float64(cumul) {
			return index
		}

		index++
	}

	maxLen := len(strategy) - 1

	if index > maxLen {
		index = maxLen
	}
	return index
}

func (t *Trainer) GetReward(heroAction, vilainAction int) int {
	return t.ActionUtility[heroAction][vilainAction]
}

func (t *Trainer) Train(iterations int) {

	for i := 0; i < iterations; i++ {
		strategy := GetStrategy(&t.RegretSum)
		oppStrategy := GetStrategy(&t.OppRegretSum)
		for i := 0; i < len(strategy); i++ {
			t.StrategySum[i] += strategy[i]
			t.OppStrategySum[i] += oppStrategy[i]
		}

		opponentAction := GetAction(oppStrategy)
		heroAction := GetAction(strategy)

		oppReward := t.GetReward(opponentAction, heroAction)
		heroReward := t.GetReward(heroAction, opponentAction)

		for i := 0; i < t.NumActions; i++ {
			// Regrets adding
			heroRegret := t.GetReward(i, opponentAction) - heroReward
			vilainRegret := t.GetReward(i, heroAction) - oppReward
			// CFR + here
			if heroRegret > 0 {
				t.RegretSum[i] += heroRegret

			}
			if vilainRegret > 0 {
				t.OppRegretSum[i] += vilainRegret

			}
		}
	}
}

func (t *Trainer) PrintAverageStrategy(strategySum []float64) {
	avgStrat := []float64{0.0, 0.0, 0.0}
	var normalizingSum float64 = 0.0
	for _, el := range strategySum {
		normalizingSum += el
	}
	for i := 0; i < len(strategySum); i++ {
		if normalizingSum > 0 {
			avgStrat[i] = strategySum[i] / normalizingSum
		} else {
			avgStrat[i] = float64(1.0 / len(strategySum))
		}
	}
	fmt.Println("Strategy :")
	for _, el := range avgStrat {
		fmt.Printf("  %f  ", el)
	}
}

func main() {
	start := time.Now()
	trainer := NewTrainer()
	trainer.Train(100000)
	elapsed := time.Since(start)
	trainer.PrintAverageStrategy(trainer.StrategySum)
	fmt.Printf("Took : %s", elapsed)
}
