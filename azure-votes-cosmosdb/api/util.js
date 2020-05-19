const fs = require('fs');

const readJsonFile = (schemaPath) => {
  return new Promise((resolve, reject) => {
    fs.readFile(schemaPath, {encoding:"utf8"}, (err, data) => {
      if (err) reject(err);
      else resolve(JSON.parse(data));
    });
  });
};

module.exports = { readJsonFile };