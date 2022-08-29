import axios from 'axios';
import config from '../config';

export async function updateScores(teamID, action) {
  try {
    const result = await axios.post(`${config.apiBaseUrl}/goal?action=${action}&team=${teamID}`);

    return {
      status: result.status,
    };
  } catch (e) {
    return {
      error: e,
      status: e.response.status,
    };
  }
}

export async function updateScoresHandler(teamId, action) {
  await updateScores(teamId, action).then((data) => {
    if (data.error) alert(data.error);
  });
}
