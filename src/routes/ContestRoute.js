import express from 'express';
import ContestController from '../controllers/ContestController.js';
const contestRoute = express.Router();

contestRoute.get('/', ContestController.getContestList);
contestRoute.get('/:id', ContestController.getContestDetail);
contestRoute.get('/:contestId/issues/:issueCode', ContestController.getIssue);
contestRoute.post('/:contestId/issues/:issueCode/submit', ContestController.submitIssue);

export default contestRoute;

