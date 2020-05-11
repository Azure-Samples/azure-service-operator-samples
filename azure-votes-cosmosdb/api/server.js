const path = require("path");
const express = require("express");
const bodyParser = require("body-parser");
const morgan = require("morgan");
const appInsights = require("applicationinsights");

const { create: createDatabase } = require("./database");
const { create: createValidator } = require("./validator");
const { create: createVotesRouter } = require("./votesRouter");

// configure application insights
appInsights.setup(process.env.APPINSIGHTS_INSTKEY);
appInsights.start();

const config = {
  schema: "./vote.schema.json",
  database: {
    kind: "cosmos",
    secretPath: process.env.COSMOSDB_SECRET_PATH,
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
  app.set('view engine', 'pug')

  app.use(morgan('combined'));
  app.use(bodyParser.json());
  app.use(express.static(path.join(__dirname, 'public')));

  app.use('/votes', createVotesRouter({
    database,
    validate,
  }));

  app.get('/', function (req, res) {
    res.render('index', { title: 'Express' });
  });

  app.listen(config.server.port, () => {
    console.log("App now running on port", config.server.port);
  });
};

const initializers = [
  createDatabase(config.database),
  createValidator(config.schema),
];
Promise.all(initializers).then(startApp);