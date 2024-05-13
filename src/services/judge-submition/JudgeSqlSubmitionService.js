import JudgeSqlSubmitionConstant from '../../constants/JudgeSqlSubmitionConstant.js';
import ApiService from '../ApiService.js';
class JudgeSqlSubmitionServiceClass {
  async sendSqlSubmition(data){
    const compiler = data.compiler;
    if(compiler === JudgeSqlSubmitionConstant.COMPILER_MYSQL){
      const url = `${process.env.MYSQL_JUDGE_BASE_URL}/judge-submition/submit`;
      const body = {
        user: data.user,
        issue: data.issue,
        statement: data.statement,
        executeType: data.executeType,
      };
      const response = await ApiService.post(url, body);
      return response.data;
    }
    return null;
  }
}

const judgeSqlSubmitionService = new JudgeSqlSubmitionServiceClass();
export default judgeSqlSubmitionService;
