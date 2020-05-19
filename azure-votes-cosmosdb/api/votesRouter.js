const express = require('express');

const create = (options) => {
  const router = express.Router();
  const {
    database,
    validate,
  } = options;

  router.get("/", async (req, res, next) => {
    try {
      const items = await database.getAll();
      res.status(200).json(items);
    } catch (err) {
      next(err);
    }
  });

  router.get("/summary", async (req, res, next) => {
    try {
      const query =  `SELECT COUNT(1) AS candidateCount, v.candidate
                      FROM Votes v
                      GROUP BY v.candidate`;
      const result = await database.query(query);
      res.status(200).json(result);
    } catch (err) {
      next(err);
    }
  });

  router.get("/:id", async (req, res, next) => {
    try {
      const {id} = req.params;
      const item = await database.get(id);
      res.status(200).json(item);
    } catch (err) {
      next(err);
    }
  });

  router.post("/", async (req, res, next) => {
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

  router.put("/:id", async (req, res, next) => {
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

  router.delete("/:id", async (req, res, next) => {
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