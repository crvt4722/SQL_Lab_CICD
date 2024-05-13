import express from 'express';
import ManageContestController from '../controllers/ManageContestController.js';
const manageContestRoute = express.Router();

manageContestRoute.post('/create-issue', ManageContestController.createIssue);
manageContestRoute.post('/add-issues-to-contest', ManageContestController.handleAddIssuesToContest);
manageContestRoute.get('/:id', ManageContestController.getManageContestDetail);

export default manageContestRoute;

