import userResource from "../models/cms-models/UserResource.js";
import contestResource from "../models/cms-models/ContestResource.js";
import issueResource from "../models/cms-models/IssueResource.js";
import testCaseResource from "../models/cms-models/TestCaseResource.js";
import contestUserResource from "../models/cms-models/ContestUserResource.js";
import contestIssueResource from "../models/cms-models/ContestIssueResource.js";
import submitionResource from "../models/cms-models/SubmitionResource.js";
export class SchemaClass {
  getTableList() {
    return [
      userResource,
      contestResource,
      issueResource,
      testCaseResource,
      contestUserResource,
      contestIssueResource,
      submitionResource,
    ];
  }


  async syncAllTable() {
    const modelList = this.getTableList();
    for(const model of modelList){
      await model.sync();
    }
  }
}

const schemaObject = new SchemaClass();
export default schemaObject;