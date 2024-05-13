import manageIssueService from "../services/manage-contest/ManageIssueService.js";
import issueResource from "../models/cms-models/IssueResource.js";
import testCaseResource from "../models/cms-models/TestCaseResource.js";
import manageContestService from "../services/manage-contest/ManageContestService.js";
class ManageContestController {
  async createIssue(req, res){
    try {
      const issue = req.body;
      console.log(issue);
      if(!issue.limitedTime || parseInt(issue.limitedTime) < 0){
        return res.status(400).send({
          message: "Trường giới hạn thời gian không hợp lệ"
        });
      }
      const createEnvironmentResult = await manageIssueService.createEnvironmentIssue({
        issue,
      });

      if(!createEnvironmentResult){
        return res.status(400).send({
          message: "Dữ liệu không hợp lệ"
        });
      }

      const createIssueResult = await issueResource.create({
        title: issue.title,
        code: 'SQL' + Date.now(),
        questionContent: issue.questionContent,
        point: issue.point,
        limitedTime: issue.limitedTime,
        useTables: issue.useTables.join(','),
        executeType: issue.executeType,
      });

      console.log(">>>> createIssueResult", createIssueResult);

      await testCaseResource.create({
        issueId: createIssueResult.id,
        outputPath: createEnvironmentResult.fileName
      });

      return res.status(200).send(createIssueResult);
    } catch (error) {
      console.log(">>> error", error.message);
      return res.status(400).send({
        message: error.message,
      });
    }
  }

  async getManageContestDetail(req, res) {
    try {
      const params = req.params;
      const contestDetail = await manageContestService.getContestDetail(params);
      for(const contestUser of contestDetail.contest_users){
        let submitIssueDict = {};
        let correctIssues = 0;
        let numberSubmitions = 0;
        if(contestUser.submitions){
          for(const submition of contestUser.submitions){
            if(!submitIssueDict[submition.contestIssueId]){
              submitIssueDict[submition.contestIssueId] = submition;
              correctIssues++;
            }
          }
          numberSubmitions = contestUser.submitions.length;
        }
        contestUser.dataValues.correctIssues = correctIssues;
        contestUser.dataValues.numberSubmitions = numberSubmitions;
      }
      res.status(200).send(contestDetail);
    } catch (error) {
      console.log(error.message);
      res.status(400).send(error);
    }
  }

  async handleAddIssuesToContest(req, res) {
    try {
      const data = req.body;
      const result = await manageContestService.handleSyncIssuesToContest(data);
      return res.status(200).send(result);
    } catch (error) {
      return res.status(400).send({
        message: error.message,
      });
    }
  }
}
  
const manageContestController = new ManageContestController();
export default manageContestController;