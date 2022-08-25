package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResetScore(t *testing.T) {
	game := &game{score: GameScore{3, 1}, scoreChannel: make(chan GameScore, 32)}

	game.ResetScore()
	resultScore := <-game.scoreChannel

	if resultScore.BlueScore != 0 || resultScore.WhiteScore != 0 {
		t.Errorf("Score did not reset. Goals white: %v, Goals blue: %v", game.score.WhiteScore, game.score.BlueScore)
	}
}

func TestAddGoal(t *testing.T) {
	game := &game{score: GameScore{3, 1}, scoreChannel: make(chan GameScore, 32)}

	type args struct {
		name               string
		teamID             int
		expectedBlueScore  int
		expectedWhiteScore int
		expectedError      string
	}
	tests := []args{
		{name: "should increment team white score by one", teamID: TeamWhite, expectedBlueScore: 0, expectedWhiteScore: 1, expectedError: ""},
		{name: "should increment team blue score by one", teamID: TeamBlue, expectedBlueScore: 1, expectedWhiteScore: 0, expectedError: ""},
		{name: "should cause an error when invalid team ID", teamID: -1, expectedBlueScore: 0, expectedWhiteScore: 0, expectedError: "bad team ID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			game.score.WhiteScore = 0
			game.score.BlueScore = 0
			err := game.AddGoal(tt.teamID)
			if err == nil {
				resultScore := <-game.scoreChannel

				assert.Equal(t, resultScore.BlueScore, tt.expectedBlueScore, "blue team score changes incorrectly")
				assert.Equal(t, resultScore.WhiteScore, tt.expectedWhiteScore, "white team score changes incorrectly")
			}

			if tt.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestGameSubGoal(t *testing.T) {
	game := &game{score: GameScore{0, 2}, scoreChannel: make(chan GameScore, 32)}

	type args struct {
		name               string
		teamID             int
		expectedBlueScore  int
		expectedWhiteScore int
		expectedError      string
	}
	tests := []args{
		{name: "should decrement team white score by one", teamID: TeamWhite, expectedBlueScore: 2, expectedWhiteScore: 0, expectedError: ""},
		{name: "should decrement team blue score by one", teamID: TeamBlue, expectedBlueScore: 1, expectedWhiteScore: 1, expectedError: ""},
		{name: "should not decrement team blue score by one", teamID: TeamWhite, expectedBlueScore: 2, expectedWhiteScore: 0, expectedError: ""},
		{name: "should cause an error when invalid team ID", teamID: -1, expectedBlueScore: 2, expectedWhiteScore: 1, expectedError: "bad team ID"},
	}

	for id, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game.score.WhiteScore = 1
			game.score.BlueScore = 2
			if id == 2 {
				game.score.WhiteScore = 0
			}
			err := game.SubGoal(tt.teamID)
			if err == nil {
				resultScore := <-game.scoreChannel

				assert.Equal(t, resultScore.BlueScore, tt.expectedBlueScore, "blue team score changes incorrectly")
				assert.Equal(t, resultScore.WhiteScore, tt.expectedWhiteScore, "white team score changes incorrectly")
			}
			if tt.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func Test_game_UpdateManualGoals(t *testing.T) {
	game := &game{score: GameScore{0, 0}, manualGoals: ManualGoals{0, 0, 0, 0}}

	type args struct {
		name                string
		teamID              int
		action              string
		expectedManualGoals ManualGoals
		expectedError       string
	}
	tests := []args{
		// No error cases
		{name: "should increment blue manual goals by 1", teamID: TeamBlue, action: "add", expectedManualGoals: ManualGoals{1, 0, 0, 0}, expectedError: ""},
		{name: "should decrement blue manual goals by 1", teamID: TeamBlue, action: "sub", expectedManualGoals: ManualGoals{0, 1, 0, 0}, expectedError: ""},
		{name: "should increment white manual goals by 1", teamID: TeamWhite, action: "add", expectedManualGoals: ManualGoals{0, 0, 1, 0}, expectedError: ""},
		{name: "should decrement white manual goals by 1", teamID: TeamWhite, action: "sub", expectedManualGoals: ManualGoals{0, 0, 0, 1}, expectedError: ""},

		// Error cases
		{name: "should return bad team ID error", teamID: 0, action: "add", expectedManualGoals: ManualGoals{0, 0, 0, 0}, expectedError: "bad team ID"},
		{name: "should return bad team ID error", teamID: 0, action: "sub", expectedManualGoals: ManualGoals{0, 0, 0, 0}, expectedError: "bad team ID"},
		{name: "should return bad action error", teamID: TeamBlue, action: "addd", expectedManualGoals: ManualGoals{0, 0, 0, 0}, expectedError: "bad action type"},
		{name: "should return bad action error", teamID: TeamWhite, action: "addd", expectedManualGoals: ManualGoals{0, 0, 0, 0}, expectedError: "bad action type"},
		{name: "should return bad action error", teamID: -10, action: "addd", expectedManualGoals: ManualGoals{0, 0, 0, 0}, expectedError: "bad action type"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game.manualGoals = ManualGoals{0, 0, 0, 0}
			err := game.UpdateManualGoals(tt.teamID, tt.action)
			if tt.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
			assert.Equal(t, game.manualGoals.AddedBlue, tt.expectedManualGoals.AddedBlue, "Manual goals for blue team added incorrectly")
			assert.Equal(t, game.manualGoals.AddedWhite, tt.expectedManualGoals.AddedWhite, "Manual goals for white team added incorrectly")
			assert.Equal(t, game.manualGoals.SubtractedBlue, tt.expectedManualGoals.SubtractedBlue, "Manual goals for blue team subtracted incorrectly")
			assert.Equal(t, game.manualGoals.SubtractedWhite, tt.expectedManualGoals.SubtractedWhite, "Manual goals for white team subtracted incorrectly")
		})
	}
}
