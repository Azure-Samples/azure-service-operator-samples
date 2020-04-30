const fs = require('fs');
const Ajv = require('ajv');

const readJsonFile = (schemaPath) => {
  return new Promise((resolve, reject) => {
    fs.readFile(schemaPath, {encoding:"utf8"}, (err, data) => {
      if (err) reject(err);
      else resolve(JSON.parse(data));
    });
  });
}

const create = async (schemaPath) => {
  const schema = await readJsonFile(schemaPath);
  const ajv = new Ajv({
    removeAdditional: true
  });
  return ajv.compile(schema);
};

module.exports = { create };