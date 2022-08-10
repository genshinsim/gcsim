var PostgrestClient = require("@supabase/postgrest-js");
var client = new PostgrestClient.PostgrestClient(
  "https://handleronedev.gcsim.app"
);

async function test() {
  let res = await client.from("avatars").select();

  console.log(res);

  res = await client.rpc("get_or_insert_user", { key: "test1", name: "test1" });

  console.log(res);

  res = await client.from("users").select();

  console.log(res);
}

test();
