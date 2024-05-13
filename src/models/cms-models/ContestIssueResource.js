import { DataTypes } from "sequelize";
import AbstractResource from "../AbstractResource.js";
class ContestIssueResource extends AbstractResource {
  constructor(
    tableName = 'contest_issues',
    modelName = 'contest_issues',
    schema = {
      id: {
        type: DataTypes.INTEGER,
        allowNull: false,
        primaryKey: true,
        autoIncrement: true,
      },
      contestId: {
        type: DataTypes.INTEGER,
        allowNull: false,
      },
      issueId: {
        type: DataTypes.INTEGER,
        allowNull: false,
      },
    },
  ) {
    super(tableName, modelName, schema);
  }

  
}

const contestIssueResource = new ContestIssueResource();
export default contestIssueResource;