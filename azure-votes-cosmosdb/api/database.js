const CosmosClient = require("@azure/cosmos").CosmosClient;

class DatabaseClient {
  constructor(endpoint, key, databaseId, containerId) {
    this.client = new CosmosClient({
      endpoint: endpoint,
      key: key
    });;
    this.database = this.client.database(databaseId);
    this.container = this.database.container(containerId);
  }

  async getAll() {
    //TODO: paging
    const resp = await this.container.items.readAll().fetchAll();
    return resp.resources;
  }

  async get(id) {
    const resp = await this.container.item(id).read();
    return resp.resource;
  }

  async insert(item) {
    const resp = await this.container.items.create(item);
    return resp.resource;
  }

  async update(id, item) {
    const oldItem = await this.get(id);
    const newItem = {
      ...oldItem,
      ...item
    };
    const resp = await this.container.item(id).replace(newItem);
    return resp.resource;
  }

  async delete(id) {
    const resp = await this.container.item(id).delete();
    return resp.resource;
  }
}

/*
 * This function ensures that the database is setup and populated correctly
 */
async function create(config) {
  const {
    endpoint,
    key,
    databaseId,
    containerId,
    partitionKey,
  } = config;
  const client = new CosmosClient({
    endpoint,
    key,
  });

  /**
   * Create the database if it does not exist
   */
  const { database } = await client.databases.createIfNotExists({
    id: databaseId
  });
  console.log(`Created database:\n${database.id}\n`);

  /**
   * Create the container if it does not exist
   */
  const { container } = await client
    .database(config.databaseId)
    .containers.createIfNotExists(
      { id: containerId, partitionKey },
      { offerThroughput: 400 }
    );
  console.log(`Created container:\n${container.id}\n`);

  return new DatabaseClient(endpoint, key, databaseId, containerId);
}

module.exports = { create };