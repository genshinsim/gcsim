create table avatars (
    avatar_id serial primary key
    , avatar_name text not NULL
);
create table users (
    user_id serial primary key
    , user_key text not null unique -- this should correspond to something?
    , user_name text not null default '' --discord tag
    , user_role integer not null default 0 -- 999 for admin, registered user should be min 1
);
create table simulations (
    simulation_id serial primary key
    , simulation_key uuid not null unique
    , metadata JSONB --unstructured metadata; 
    , viewer_file text
    , fk_user_id integer references users(user_id)
    , is_permanent boolean default false
    , is_public boolean default false
    , is_shared boolean default false -- if this is true, forms part of db
    , review_state integer default 0 -- applies to shared only; must be reviewed to make it live
    , create_time timestamp with time zone DEFAULT current_timestamp NOT NULL
 );
create table tags (
    tag_id serial primary key
    , tag_name text
);
-- join tables
create table simulationtags (
    simulation_tag_id serial primary key
    , fk_tag_id integer references tags(tag_id)
    , fk_simulation integer references simulations(simulation_id)
);
create table avatarsimulations (
    avatar_simulation_id serial primary key
    , fk_avatar integer references avatars(avatar_id)
    , fk_simulation integer references simulations(simulation_id)
);

-- view
create view active_simulations as
    select * from simulations
    where create_time > current_timestamp - interval '14 days'
;

create view user_simulation_count as
    select
        u.user_id,
        u.user_key,
        u.user_name,
        u.user_role,
        count(*) 
    from users u join
        simulations s
        on u.user_id = s.fk_user_id
    group by u.user_id, u.user_key, u.user_name, u.user_role
;

-- functions for commonly used queries

create or replace function public.get_or_insert_user(key text, name text)
returns table (
	user_id integer,
	user_key text,
	user_name text,
	user_role INTEGER
)
language plpgsql
as
$$
begin
	return query
		with s as (
			select users.user_id, users.user_key, users.user_name, users.user_role
			from users
			where users.user_key = key
		), i as (
			insert into users (user_key, user_name)
			select key, name
			where not exists (select 1 from s)
			returning users.user_id, users.user_key, users.user_name, users.user_role
		)
		select i.user_id, i.user_key, i.user_name, i.user_role
		from i
		union all
		select s.user_id, s.user_key, s.user_name, s.user_role
		from s;
end;
$$;

-- auto notify psotgrest on schema change

-- Create an event trigger function
CREATE OR REPLACE FUNCTION public.pgrst_watch() RETURNS event_trigger
  LANGUAGE plpgsql
  AS $$
BEGIN
  NOTIFY pgrst, 'reload schema';
END;
$$;

-- This event trigger will fire after every ddl_command_end event
CREATE EVENT TRIGGER pgrst_watch
  ON ddl_command_end
  EXECUTE PROCEDURE public.pgrst_watch();

-- some basic data
insert into users(user_key, user_role) values ('admin', 999);
insert into avatars(avatar_id, avatar_name) values (10000002, 'ayaka');
insert into avatars(avatar_id, avatar_name) values (10000003, 'jean');
insert into avatars(avatar_id, avatar_name) values (10000006, 'lisa');


insert into users(user_key, user_name, user_role) values ('user1', 'one', 999);
insert into users(user_key, user_name, user_role) values ('user2', 'two', 999);

insert into simulations(fk_user_id, simulation_key) values (1, '8a00ce82-608c-4de9-8ef3-8613d7700b7b');
insert into simulations(fk_user_id, simulation_key) values (1, 'c6000c89-e216-407b-9306-3804c9726c91');
insert into simulations(fk_user_id, simulation_key) values (1, '3ea33d6b-bbb0-4e66-9e02-ca119277f5bf');
insert into simulations(fk_user_id, simulation_key) values (2, 'afd3d702-c99b-47a2-acb6-f09bd7c83d3f');
insert into simulations(fk_user_id, simulation_key) values (2, 'ab4679ef-36e5-4393-97a7-3644af748178');
insert into simulations(fk_user_id, simulation_key) values (3, '6a0be3de-947a-45f2-bb81-c456cdec9672');
insert into simulations(fk_user_id, simulation_key) values (3, 'c8156ab8-bb73-43a0-975c-5d1c656c5259');
insert into simulations(fk_user_id, simulation_key) values (3, 'd8761162-8365-46f0-b695-786250025690');
insert into simulations(fk_user_id, simulation_key) values (3, 'ced532d0-7ad9-4a4f-8f88-e8550f689c69');
insert into simulations(fk_user_id, simulation_key) values (3, '1f503708-4581-4aab-a638-20926ec72101');
insert into simulations(fk_user_id, simulation_key) values (3, 'af57c74f-e06e-44f8-abfa-35eb45b134ca');