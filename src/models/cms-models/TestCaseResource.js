import { DataTypes } from "sequelize";
import AbstractModel from "../AbstractResource.js";
import issueResource from "./IssueResource.js";

class TestCaseResource extends AbstractModel {
  constructor(
    tableName = 'testcases',
    modelName = 'testcases',
    schema = {
      id: {
        type: DataTypes.INTEGER,
        allowNull: false,
        primaryKey: true,
        autoIncrement: true,
      },
      issueId: {
        type: DataTypes.INTEGER,
      },
      outputPath: {
        type: DataTypes.STRING,
      }
    },
  ) {
    super(tableName, modelName, schema);
  }

  afterInitModel(){
    const testCaseModel = this.getDataModel();
    const issueModel = issueResource.getDataModel();
    testCaseModel.belongsTo(
      issueModel,
      {
        foreignKey: {
          field: 'issueId',
        },
        targetKey: 'id',
        onDelete: 'CASCADE',
        onUpdate: 'CASCADE',
      },
    );
  }
}

const testCaseResource = new TestCaseResource();
export default testCaseResource;