export const TeamID = {
  Team_white: 0,
  Team_blue: 1,
};

export const ScoreChange = {
  Add_goal: 'add',
  Sub_goal: 'sub',
};

export class Goal {
  constructor(teamID, timestamp) {
    this.teamID = teamID;
    this.timestamp = timestamp;
  }
}
