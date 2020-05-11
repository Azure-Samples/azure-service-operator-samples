const fs = require('fs');
const Ajv = require('ajv');

const { readJsonFile } = require('./util');

const create = async (schemaPath) => {
  const schema = await readJsonFile(schemaPath);
  const ajv = new Ajv({
    removeAdditional: true
  });
  return ajv.compile(schema);
};

module.exports = { create };