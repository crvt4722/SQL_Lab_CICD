/* eslint-disable no-undef */
import Sequelize from "sequelize";
import 'dotenv/config';

export const systemSequelize = new Sequelize(
  process.env.DB_SYSTEM_NAME,
  process.env.DB_SYSTEM_USERNAME,
  process.env.DB_SYSTEM_PASSWORD,
  {
    host: process.env.DB_SYSTEM_HOST,
    port: process.env.DB_SYSTEM_PORT,
    dialect: process.env.DB_SYSTEM_TYPE,
  },
);


async function authenticateSequelize() {
  await systemSequelize.authenticate();
  console.log("[SQL_LAB_SERVER] Connect to systemSequelize successfully!");
}

export default authenticateSequelize;