const express = require('express');

// Generic error handler used by all endpoints.
const handleError = (res, reason, message, code) => {
  console.log("ERROR: " + reason);
  res.status(code || 500).json({"error": message});
};

const create = (options) => {
  const router = express.Router();
  const {
    name,
    database,
    validate,
  } = options;

  router.get(`/${name}`, async (req, res, next) => {
    try {
      const items = await database.getAll();
      res.status(200).json(items);
    } catch (err) {
      next(err);
    }
  });

  router.get(`/${name}/:id`, async (req, res, next) => {
    try {
      const {id} = req.params;
      const item = await database.get(id);
      res.status(200).json(item);
    } catch (err) {
      next(err);
    }
  });

  router.post(`/${name}`, async (req, res, next) => {
    try {
      let item = req.body;
      if (validate(item)) {
        item = await database.insert(req.body);
        res.status(201).json(item);
      } else {
        res.status(422).json(validate.errors);
      }
    } catch (err) {
      next(err);
    }
  });

  router.put(`/${name}/:id`, async (req, res, next) => {
    try {
      const {id} = req.params;
      let item = req.body;
      if (validate(item)) {
        item = await database.update(id, req.body);
        res.status(200).json(item);
      } else {
        res.status(422).json(validate.errors);
      }
    } catch (err) {
      next(err);
    }
  });

  router.delete(`/${name}/:id`, async (req, res, next) => {
    try {
      const {id} = req.params;
      const item = await database.delete(id);
      res.status(204).send();
    } catch (err) {
      next(err);
    }
  });

  return router;
};

module.exports = { create };