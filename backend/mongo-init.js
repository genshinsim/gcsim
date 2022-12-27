db.auth("root", "root-password");

db = db.getSiblingDB("gcsim-database");

db.createUser({
  user: "gcsim",
  pwd: "gcsim-password",
  roles: [
    {
      role: "root",
      db: "gcsim-database",
    },
  ],
});

db.createCollection("data");
