module.exports = {
  async up(db, client) {
     await db.collection("existing_table1").createIndex({ blockchain_id: 1 }, { unique: true });
     await db.collection("existing_table2").createIndex({ address: 1 }, { unique: true });
     await db.createCollection("existing_table3",{capped: true, size:999999, max:1});
     await db.createCollection("existing_table4",{capped: true, size:999999, max:1});

  },

  async down(db, client) {
     await db.collection("existing_table1").drop();
     await db.collection("existing_table2").drop();
     await db.collection("existing_table3").drop();
     await db.collection("existing_table4").drop();

  }
};
