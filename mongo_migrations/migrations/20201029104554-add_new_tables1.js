module.exports = {
  async up(db, client) {
     await db.collection("foo_bar").createIndex({ foo: 1, bar: 1 }, { unique: true });
  },

  async down(db, client) {
     await db.collection("foo_bar").drop();
  }
};
