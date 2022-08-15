truncate simulations, avatars, users cascade;

-- some basic data
insert into users(user_id, user_name, user_role) values (-1, 'admin', 999);
insert into avatars(avatar_id, avatar_name) values (10000002, 'ayaka');
insert into avatars(avatar_id, avatar_name) values (10000003, 'jean');
insert into avatars(avatar_id, avatar_name) values (10000006, 'lisa');