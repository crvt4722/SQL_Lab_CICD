import { DataTypes } from "sequelize";
import AbstractResource from "../AbstractResource.js";
import testCaseResource from "./TestCaseResource.js";

class IssueResource extends AbstractResource {
  constructor(
    tableName = 'issues',
    modelName = 'issues',
    schema = {
      id: {
        type: DataTypes.INTEGER,
        allowNull: false,
        primaryKey: true,
        autoIncrement: true,
      },
      code: {
        type: DataTypes.STRING,
        allowNull: false,
        unique: true,
      },
      title: {
        type: DataTypes.STRING,
      },
      questionContent: {
        type: DataTypes.TEXT,
      },
      point: {
        type: DataTypes.INTEGER,
        default: 1,
      },
      limitedTime: {
        type: DataTypes.FLOAT,
        default: 2,
      },
      useTables: {
        type: DataTypes.STRING,
      },
      executeType: {
        type: DataTypes.STRING, // 1: readonly || 2: createtemptable || 3: sp
      }
    },
  ) {
    super(tableName, modelName, schema);
  }

  afterInitModel(){
    const issueModel = this.getDataModel();
    const testCaseModel = testCaseResource.getDataModel();
    issueModel.hasMany(testCaseModel);
  }
}

const issueResource = new IssueResource();
export default issueResource;