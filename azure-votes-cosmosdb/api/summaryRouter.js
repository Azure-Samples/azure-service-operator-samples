const express = require('express');

const create = (options) => {
  const router = express.Router();
  const {database} = options;

  router.get(`/summary`, async (req, res, next) => {
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

  return router;
};

module.exports = { create };