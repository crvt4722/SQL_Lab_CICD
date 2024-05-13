import { DataTypes } from "sequelize";
import AbstractResource from "../AbstractResource.js";
import contestResource from "./ContestResource.js";
import userResource from "./UserResource.js";
class ContestUserResource extends AbstractResource {
  constructor(
    tableName = 'contest_users',
    modelName = 'contest_users',
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
      userId: {
        type: DataTypes.INTEGER,
        allowNull: false,
      },
      participantRole: {
        type: DataTypes.TINYINT,
        allowNull: false,
        default: 1, // participant: 1, host: 2
      },
      startTime: {
        type: DataTypes.DATE,
      },
      endTime: {
        type: DataTypes.DATE,
      }
    },
  ) {
    super(tableName, modelName, schema);
  }

  afterInitModel(){
    const contestModel = contestResource.getDataModel();
    const userModel = userResource.getDataModel();
    const contestUserModel = this.getDataModel();
    contestModel.belongsToMany(userModel, { through: {model: contestUserModel, unique: false}, foreignKey: 'contestId', unique: false, });
    userModel.belongsToMany(contestModel, { through: {model: contestUserModel, unique: false}, foreignKey: 'userId', unique: false, });

    contestModel.hasMany(contestUserModel);
    contestUserModel.belongsTo(contestModel);

    userModel.hasMany(contestUserModel);
    contestUserModel.belongsTo(userModel);
  }
}

const contestUserResource = new ContestUserResource();
export default contestUserResource;