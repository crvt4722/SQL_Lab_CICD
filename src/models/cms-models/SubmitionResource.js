import { DataTypes } from "sequelize";
import AbstractResource from "../AbstractResource.js";
import contestUserResource from "./ContestUserResource.js";
import contestIssueResource from "./ContestIssueResource.js";

class SubmitionResource extends AbstractResource {
  constructor(
    tableName = 'submitions',
    modelName = 'submitions',
    schema = {
      id: {
        type: DataTypes.INTEGER,
        allowNull: false,
        primaryKey: true,
        autoIncrement: true,
      },
      contestUserId: {
        type: DataTypes.INTEGER,
        allowNull: false,
      },
      contestIssueId: {
        type: DataTypes.INTEGER,
        allowNull: false,
      },
      status: {
        type: DataTypes.STRING,
      },
      compiler: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      srcCode: {
        type: DataTypes.TEXT,
      },
      submitTime: {
        type: DataTypes.DATE,
      },
      executionTime: {
        type: DataTypes.FLOAT,
      }
    },
  ) {
    super(tableName, modelName, schema);
  }

  afterInitModel(){
    const contestUserModel = contestUserResource.getDataModel();
    const contestIssueModel = contestIssueResource.getDataModel();
    const submitionModel = this.getDataModel();
    contestUserModel.belongsToMany(contestIssueModel, { through: { model: submitionModel, unique: false}, foreignKey: 'contestUserId' });
    contestIssueModel.belongsToMany(contestUserModel, { through: { model: submitionModel, unique: false}, foreignKey: 'contestIssueId' });

    contestIssueModel.hasMany(submitionModel);
    submitionModel.belongsTo(contestIssueModel);

    contestUserModel.hasMany(submitionModel);
    submitionModel.belongsTo(contestUserModel);
  }
}

const submitionResource = new SubmitionResource();
export default submitionResource;