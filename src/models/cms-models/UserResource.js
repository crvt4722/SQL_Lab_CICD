import { DataTypes } from "sequelize";
import AbstractResource from "../AbstractResource.js";


class UserResource extends AbstractResource {
  constructor(
    tableName = 'users',
    modelName = 'users',
    schema = {
      id: {
        type: DataTypes.INTEGER,
        allowNull: false,
        primaryKey: true,
        autoIncrement: true,
      },
      username: {
        type: DataTypes.STRING,
        allowNull: false,
        unique: true,
      },
      password: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      firstName: {
        type: DataTypes.STRING,
      },
      lastName: {
        type: DataTypes.STRING,
      },
    },
  ) {
    super(tableName, modelName, schema);
  }
}

const userResource = new UserResource();
export default userResource;