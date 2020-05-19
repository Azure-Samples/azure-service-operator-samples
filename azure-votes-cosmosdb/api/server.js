// configure application insights
// - setup early so it can track other packages
// - uses environment variable APPINSIGHTS_INSTRUMENTATIONKEY
const appInsights = require("applicationinsights");
appInsights.setup()
  .setSendLiveMetrics(true)
  .start();
const startTime = Date.now();

const path = require("path");
const express = require("express");
const bodyParser = require("body-parser");
const morgan = require("morgan");

const { create: createDatabase } = require("./database");
const { create: createValidator } = require("./validator");
const { create: createVotesRouter } = require("./votesRouter");

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
  const aiClient = appInsights.defaultClient;

  const app = express();
  app.set('view engine', 'pug');

  app.use(morgan('combined'));
  app.use(bodyParser.json());
  app.use(express.static(path.join(__dirname, 'public')));

  app.use('/votes', createVotesRouter({
    database,
    validate,
    aiClient,
  }));

  app.get('/', function (req, res) {
    res.render('index', { title: 'Express' });
  });

  // application-wide exception logging
  app.use((err, req, res, next) => {
    aiClient.trackException({
      exception: err,
    });
    next(err)
  });

  app.listen(config.server.port, () => {
    const startupDuration = Date.now() - startTime;
    aiClient.trackMetric({ name: "application-startup-duration", value: startupDuration });
    console.log("App now running on port", config.server.port);
  });
};

const initializers = [
  createDatabase(config.database),
  createValidator(config.schema),
];
Promise.all(initializers).then(startApp);