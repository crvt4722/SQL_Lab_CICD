import contestService from "../services/contest/ContestService.js";
import issueService from "../services/issue/IssueService.js";
import judgeSqlSubmitionService from "../services/judge-submition/JudgeSqlSubmitionService.js";

class ContestController {
  async getContestList(req, res){
    try {
      const contestList = await contestService.getContestList();
      console.log(contestList);
      res.status(200).send(contestList);
    } catch (error) {
      res.status(400).send(error);
    }
  }

  async getContestDetail(req, res) {
    try {
      const params = req.params;
      const contestDetail = await contestService.getContestDetail(params);
      for(const contestIssue of contestDetail.contest.contest_issues){
        let submitStatus;
        const submitions = contestIssue.dataValues.submitions;
        if(submitions && submitions.length){
          const submitAC = submitions.find(item => item.status === 'AC');
          submitStatus = submitAC ? 1 : 2;
        } else {
          submitStatus = 0;
        }
        contestIssue.dataValues.submitStatus = submitStatus;
        // delete contestIssue.dataValues.submitions;
      }
      res.status(200).send(contestDetail);
    } catch (error) {
      console.log(error.message);
      res.status(400).send(error);
    }
  }

  async getIssue(req, res){
    try {
      const params = req.params;
      const issue = await issueService.getIssue(params);
      console.log(params);
      if(!issue){
        return res.status(400).send({
          message: "Bài tập không tồn tại"
        });
      }
      return res.status(200).send(issue);
    } catch (error) {
      console.log(">>>> ERROR");
      res.status(400).send(error);
    }
  }

  async submitIssue(req, res){
    try {
      const {statement, compiler} = req.body;
      const {contestId, issueCode} = req.params;
      const user = {
        id: 1,
        username: 'B20DCCN736',
      }; 

      const contestIssue = await issueService.getIssueHandleSubmit({contestId, issueCode, userId: user.id});
      const executeType = contestIssue.issue.executeType;
      // console.log("contestIssue.issue", contestIssue.issue.testcases);
  
      
      const data = {
        user,
        issue: contestIssue.issue,
        executeType,
        statement,
        compiler,
      };

      const result = await judgeSqlSubmitionService.sendSqlSubmition(data);
      res.status(200).send(result);
    } catch (error) {
      res.status(400).send({
        meesage: error
      });
    }
  }
}
  
const contestController = new ContestController();
export default contestController;