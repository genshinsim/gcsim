echo "add new user"

curl "http://localhost:3000/rpc/get_or_insert_user" \
  -X POST -H "Content-Type: application/json" \
  -d '{ "id": 156096226721398786, "name": "srl#2712" }'

# curl "http://localhost:3000/rpc/get_or_insert_user" \
#   -X POST -H "Content-Type: application/json" \
#   -d '{ "id": 320893091450191872, "name": "imring#3781" }'

echo "\n\n new sim"

curl "http://localhost:3000/rpc/share_sim" \
  -X POST -H "Content-Type: application/json" \
  -d '{ "metadata": "{}", "viewer_file": "fake", "user_id": 156096226721398786, "is_permanent": false, "is_public": false}'
echo ""
curl "http://localhost:3000/rpc/share_sim" \
  -X POST -H "Content-Type: application/json" \
  -d '{ "metadata": "{}", "viewer_file": "fake", "user_id": 156096226721398786, "is_permanent": true, "is_public": false}'

echo "\n\n list"
curl localhost:3000/active_user_simulations?user_id=eq.156096226721398786

echo "\n\n list perm sim count"
curl localhost:3000/user_simulation_count?user_id=eq.156096226721398786

# adding to db
gen_post_data()
{
cat <<EOF
{
    "simulation_key": "$1",
    "git_hash": "fakehash",
    "config_hash": "fakehash",
    "author": $2,
    "sim_description": "$3"
}
EOF
}

gen_replace_data()
{
cat <<EOF
{
    "old_key": "$1",
    "simulation_key": "$2",
    "git_hash": "fakehash",
    "config_hash": "fakehash",
    "author": $3,
    "sim_description": "$4"
}
EOF
}

echo "\n\n share then add to db"
shareKey=$(curl "http://localhost:3000/rpc/share_sim" \
  -X POST -H "Content-Type: application/json" \
  -d '{ "metadata": "{}", "viewer_file": "fake", "user_id": 156096226721398786, "is_permanent": false, "is_public": false}'\
)
shareKey=${shareKey//\"/}


echo "\n\n share key is: $shareKey"

dd=$(gen_post_data $shareKey 156096226721398786 boo)

echo "\n sending $dd"

id=$(curl "http://localhost:3000/rpc/add_db_sim" \
  -X POST -H "Content-Type: application/json" \
  -d "$dd" \
)
echo "\n\n id is $id"


echo "\n\n check db entry now permanent"

curl "http://localhost:3000/simulations?simulation_key=eq.$shareKey"

echo "\n\n check entry is in db"

curl "http://localhost:3000/db_simulations"


echo "\n\n update same entry"

dd=$(gen_replace_data $shareKey $shareKey 156096226721398786 updating)
echo "\n sending $dd"

curl "http://localhost:3000/rpc/replace_db_sim" \
  -X POST -H "Content-Type: application/json" \
  -d "$dd"

echo "\n\n check entry is updated"

curl "http://localhost:3000/db_simulations"


echo "\n\n add authors"

dd=$(gen_replace_data $shareKey $shareKey 320893091450191872 "more authors")
echo "\n sending $dd"

curl "http://localhost:3000/rpc/replace_db_sim" \
  -X POST -H "Content-Type: application/json" \
  -d "$dd"

echo "\n\n check author changed"

curl "http://localhost:3000/db_entry_authors"