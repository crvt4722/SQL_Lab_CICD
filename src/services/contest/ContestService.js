import contestResource from "../../models/cms-models/ContestResource.js";

class ContestService {
  PAGE_SIZE = 50;
  async getContestList(filter = null, pageSize = this.PAGE_SIZE, currentPage = 1, sort = null){
    try {
      const contestList = await contestResource.getList(filter, pageSize, currentPage, sort);
      return contestList;
    } catch (error) {
      console.log(error.message);
      return null;
    }
  }

  async getContestDetail(params){
    try {
      const contestDetail = await contestResource.getContestDetail(params);
      return contestDetail;
    } catch (error) {
      console.log(">>> ERROR", error);
      return null;
    }
  }
  
}

const contestService = new ContestService();
export default contestService;

