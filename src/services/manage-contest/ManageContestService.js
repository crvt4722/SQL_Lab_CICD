import { Op } from "sequelize";
import contestIssueResource from "../../models/cms-models/ContestIssueResource.js";
import submitionResource from "../../models/cms-models/SubmitionResource.js";
import issueResource from "../../models/cms-models/IssueResource.js";
import contestUserResource from "../../models/cms-models/ContestUserResource.js";
import contestResource from "../../models/cms-models/ContestResource.js";
import userResource from "../../models/cms-models/UserResource.js";

class ManageContestService {
  PAGE_SIZE = 50;
  
  async getContestDetail(params){
    try {
      const contestDetail = await contestResource.getDataModel().findOne(
        {
          where: {
            id: params.id,
          },
          include: [
            {
              model: issueResource.getDataModel(),
              through: contestIssueResource.getDataModel(),
            },
            {
              model: contestUserResource.getDataModel(),
              include: [
                {
                  model: userResource.getDataModel(),
                  attributes: ['id', 'username', 'firstName', 'lastName'],
                },
                {
                  model: submitionResource.getDataModel(),
                  required: false,
                  where: {
                    status: 'AC'
                  },
                }
              ]
            }
          ]
        }
      );  
      return contestDetail;
    } catch (error) {
      console.log(error.message);
      return null;
    }
  }
  
  async handleSyncIssuesToContest(data){
    if(!data.addIssues){
      return null;
    }
    const issues = await issueResource.getDataModel().findAll({
      where: {
        code: {
          [Op.in]: data.addIssues,
        }
      }
    });

    const issueIds = issues.map(issue => issue.id);
    const contestIssueData = issueIds.map(issueId => {
      return {
        contestId: data.contestId,
        issueId,
      };
    });
    const result = await contestIssueResource.bulkCreate(contestIssueData);
    return result;
  }
}

const manageContestService = new ManageContestService();
export default manageContestService;

