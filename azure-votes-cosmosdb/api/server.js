const express = require("express");
const bodyParser = require("body-parser");
const morgan = require("morgan");
const appInsights = require("applicationinsights");

const votesRouter = require("./votesRouter");
const validator = require("./validator");
const database = require("./database");

appInsights.setup(process.env.APPINSIGHTS_INSTKEY);
appInsights.start();

const config = {
  schema: "./vote.schema.json",
  cosmos: {
    endpoint: process.env.COSMOSDB_ACCOUNTURI,
    key: process.env.COSMOSDB_ACCOUNTKEY,
    databaseId: "Voting",
    containerId: "Votes",
    partitionKey: { kind: "Hash", paths: ["/createdAt"] },
  },
  server: {
    port: process.env.PORT || 8080,
  },
};

const startApp = (services) => {
  const database = services[0];
  const validate = services[1];

  const app = express();
  app.use(morgan('combined'));
  app.use(bodyParser.json());
  app.use('/', votesRouter.create({
    name: config.cosmos.containerId,
    database,
    validate,
  }));

  app.listen(config.server.port, () => {
    console.log("App now running on port", config.server.port);
  });
};

const initializers = [
  database.create(config.cosmos),
  validator.create(config.schema),
];
Promise.all(initializers).then(startApp);