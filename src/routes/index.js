// import adminRouter from "./AdminRoute.js";
import contestRouter from "./ContestRoute.js";
import manageContestRouter from "./ManageContestRoute.js";
const initWebRoute = (app) => {
  // app.use('/api/v1/admin', adminRouter);
  app.use('/api/v1/contests', contestRouter);
  app.use('/api/v1/admin/manage-contests', manageContestRouter);
  return app.use('/', async function (req, res) {
    res.send("hello from server");
  });
};

export default initWebRoute;