import contestIssueResource from "../../models/cms-models/ContestIssueResource.js";
import contestUserResource from "../../models/cms-models/ContestUserResource.js";
import issueResource from "../../models/cms-models/IssueResource.js";
import submitionResource from "../../models/cms-models/SubmitionResource.js";
import testCaseResource from "../../models/cms-models/TestCaseResource.js";

class IssueService {
  PAGE_SIZE = 50;

  async getIssue(params){
    const issue = await contestIssueResource.getDataModel().findOne({
      where: {
        contestId: params.contestId,
      },
      include: [
        {
          model: issueResource.getDataModel(),
          where: {
            code: params.issueCode,
          }
        },
        {
          model: submitionResource.getDataModel(),
          required: false,
          include: [
            {
              model: contestUserResource.getDataModel(),
              required: false,
              attributes: [],
              where: {
                userId: params.userId || 1, // TODO
              }
            }
          ]
        }
      ]
    });
    return issue;
  }
  
  async getIssueHandleSubmit(params){
    const issue = await contestIssueResource.getDataModel().findOne({
      where: {
        contestId: params.contestId,
      },
      include: [
        {
          model: issueResource.getDataModel(),
          where: {
            code: params.issueCode,
          },
          include: [
            {
              model: testCaseResource.getDataModel()  
            }
          ]
        },
        {
          model: submitionResource.getDataModel(),
          required: false,
          include: [
            {
              model: contestUserResource.getDataModel(),
              required: false,
              attributes: [],
              where: {
                userId: params.userId || 1, // TODO
              }
            }
          ]
        }
      ]
    });
    return issue;
  }
}

const issueService = new IssueService();
export default issueService;

