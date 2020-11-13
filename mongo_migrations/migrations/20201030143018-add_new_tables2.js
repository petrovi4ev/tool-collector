module.exports = {
  async up(db, client) {
     await db.collection("table2").createIndex({ foo: 1 }, { unique: true });
  },

  async down(db, client) {
     await db.collection("table2").drop();
  }
};
