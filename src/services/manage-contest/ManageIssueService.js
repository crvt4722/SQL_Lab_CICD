import JudgeSqlSubmitionConstant from '../../constants/JudgeSqlSubmitionConstant.js';
import ApiService from '../ApiService.js';
class ManageIssueService {
  async createEnvironmentIssue(data){
    // if(data.compiler === JudgeSqlSubmitionConstant.COMPILER_MYSQL){
    const url = `${process.env.MYSQL_JUDGE_BASE_URL}/manage-issue/create-issue`;
    const response = await ApiService.post(url, data);
    return response.data;
    // }
  }
}

const manageIssueService = new ManageIssueService();
export default manageIssueService;

