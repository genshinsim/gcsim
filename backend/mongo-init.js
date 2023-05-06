db.auth("root", "root-password");

db = db.getSiblingDB("gcsim-database");

db.createCollection("data");

db.createCollection("shares");

db.createView(
    "gcsimvaliddb",
    "data",
    [ { $match: { is_db_valid: true } }]
)

db.createView(
    "gcsimsubs",
    "data",
    [ { $match: { summary: {$exists: false} } }]
)