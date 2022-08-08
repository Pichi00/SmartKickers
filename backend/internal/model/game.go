package model

import "errors"

const (
	teamWhite = 1
	teamBlue  = 2
)

type Game interface {
	AddGoal(int) error
	ResetScore()
	GetScore() gameScore
	GetChannel() chan bool
	SubGoal(int) error
}

type game struct {
	score        gameScore
	ScoreChannel chan bool
}

type gameScore struct {
	BlueScore  int `json:"blueScore"`
	WhiteScore int `json:"whiteScore"`
}

func NewGame() Game {
	return &game{ScoreChannel: make(chan bool)}
}

func (g *game) ResetScore() {
	g.score.BlueScore = 0
	g.score.WhiteScore = 0
}

func (g *game) AddGoal(teamID int) error {
	switch teamID {
	case teamWhite:
		g.score.WhiteScore++
	case teamBlue:
		g.score.BlueScore++
	default:
		return errors.New("bad team ID")
	}
	g.ScoreChannel <- true
	return nil
}

func (g *game) GetScore() gameScore {
	return g.score
}

func (g *game) GetChannel() chan bool {
	return g.ScoreChannel
}

func (g *game) SubGoal(teamID int) error {
	switch teamID {
	case teamWhite:
		if g.score.WhiteScore > 0 {
			g.score.WhiteScore--
		}
	case teamBlue:
		if g.score.BlueScore > 0 {
			g.score.BlueScore--
		}
	default:
		return errors.New("bad team ID")
	}
	return nil

}
