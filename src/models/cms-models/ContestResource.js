import { DataTypes, Op, Sequelize } from "sequelize";
import AbstractResource from "../AbstractResource.js";
import issueResource from "./IssueResource.js";
import contestIssueResource from "./ContestIssueResource.js";
import submitionResource from "./SubmitionResource.js";
import contestUserResource from "./ContestUserResource.js";
class ContestResource extends AbstractResource {
  constructor(
    tableName = 'contests',
    modelName = 'contests',
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
      openTime: {
        type: DataTypes.DATE,
      },
      closeTime: {
        type: DataTypes.DATE,
      }
    },
  ) {
    super(tableName, modelName, schema);
  }

  afterInitModel(){
    const contestModel = this.getDataModel();
    const issueModel = issueResource.getDataModel();
    const contestIssueModel = contestIssueResource.getDataModel();
    contestModel.belongsToMany(issueModel, { through: {model: contestIssueModel, unique: false}, foreignKey: 'contestId' });
    issueModel.belongsToMany(contestModel, { through: {model: contestIssueModel, unique: false}, foreignKey: 'issueId' });

    contestModel.hasMany(contestIssueModel);
    contestIssueModel.belongsTo(contestModel);

    issueModel.hasMany(contestIssueModel);
    contestIssueModel.belongsTo(issueModel);
  }

  async getContestDetail(params){
    const contest = await this.getDataModel().findOne(
      {
        where: {
          id: params.id,
        },
        include: [
          {
            model: contestIssueResource.getDataModel(),
            required: false,
            attributes: ['id'],
            include: [
              {
                model: submitionResource.getDataModel(),
                required: false,
                attributes: ['status'],
                where: {
                  status: {
                    [Op.ne]: null
                  }
                },
                include: [
                  {
                    model: contestUserResource.getDataModel(),
                    attributes: [],
                    where: {
                      userId: 1, // TODO
                    }
                  }
                ]
              },
              {
                model: issueResource.getDataModel(),
                attributes: ['id', 'code', 'title'],
              },
            ],
          },
          // {
          //   model: contestUserResource.getDataModel(),
          //   attributes: ['id', 'startTime', 'endTime', 'userId', 'participantRole'],
          //   where: {
          //     userId: 1, // TO DO
          //   }
          // },
        ]
      }
    );  
    console.log(">>> contest", contest);
    return {
      contest,
    };
  }
  
}

const contestResource = new ContestResource();
export default contestResource;