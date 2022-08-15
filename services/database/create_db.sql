create table avatars (
    avatar_id serial primary key
    , avatar_name text not NULL unique
);
create table users (
    user_id bigint not null unique primary key-- discord id (twitter snowflake uint64)
    , user_name text not null default '' --discord tag
    , user_role integer not null default 1 -- 999 for admin, registered user should be min 1
);
-- https://stackoverflow.com/questions/26046816/is-there-a-way-to-set-an-expiry-time-after-which-a-data-entry-is-automaticall
-- https://stackoverflow.com/questions/61506272/does-postgresql-have-a-built-in-feature-to-delete-old-records-every-minute-or-ho
create table simulations (
    simulation_key uuid not null default gen_random_uuid() unique primary key
    , metadata JSONB --unstructured metadata; 
    , viewer_file text
    , is_permanent boolean default false
    , create_time timestamp with time zone DEFAULT current_timestamp NOT NULL
 );

-- for user shared simulations
create table user_simulations (
    user_simulation_id serial primary key
    , simulation_key uuid references simulations(simulation_key)
    , user_id bigint references users(user_id)
    , is_public boolean default false -- cannot be true if perm is not true
);

-- for database simulations
create table db_simulations (
    db_id serial primary key
    , simulation_key uuid not null unique references simulations(simulation_key)
    , git_hash text
    , config_hash text
    , sim_description text
);
create table tags (
    tag_id serial primary key
    , tag_name text
);
-- join tables
create table db_entry_authors (
    db_id int not null references db_simulations(db_id)
    , user_id bigint not null references users(user_id)
    , primary key (db_id, user_id)
);
create table simulationtags (
    tag_id integer references tags(tag_id)
    , simulation_key uuid references simulations(simulation_key) 
    , primary key(tag_id, simulation_key)
);
create table avatarsimulations (
    avatar_id integer references avatars(avatar_id)
    , simulation_key uuid references simulations(simulation_key) 
    , primary key(avatar_id, simulation_key)
);

-- view
create or replace view active_sim as 
    select
        *
    from simulations as s
    where create_time > current_timestamp - interval '30 days' or s.is_permanent
;
 
create or replace view active_user_simulations as
    select 
        s.simulation_key
        , s.metadata
        , s.is_permanent
        , u.user_id
        , case when u.user_name is null then 'anon' else u.user_name end as user_name 
    from simulations s
    left outer join user_simulations us
        on us.simulation_key = s.simulation_key
    left outer join users u
        on us.user_id = u.user_id
    where create_time > current_timestamp - interval '30 days' or s.is_permanent
;

create or replace view user_simulation_count as
    select
        us.user_id
        , count(*) 
    from user_simulations us
    left outer join simulations s
        on us.simulation_key = s.simulation_key
    where
        s.is_permanent
    group by us.user_id
;

create or replace view active_user_sims_by_avatar as
    select 
        s.simulation_key
        , s.metadata
        , s.is_permanent
        , s.create_time
        , a.avatar_id
        , a.avatar_name
    from avatarsimulations x
    join user_simulations u
        on x.simulation_key = u.simulation_key
    left outer join avatars a
        on a.avatar_id = x.avatar_id
    left outer join simulations s
        on x.simulation_key = s.simulation_key
    where create_time > current_timestamp - interval '30 days' or s.is_permanent
;

create or replace view db_sims_by_avatar as
    select 
        s.simulation_key
        , s.metadata
        , s.is_permanent
        , s.create_time
        , d.git_hash
        , d.sim_description
        , a.avatar_id
        , a.avatar_name
    from avatarsimulations x
    join db_simulations d
        on x.simulation_key = d.simulation_key
    left outer join simulations s
        on d.simulation_key = s.simulation_key
    left outer join avatars a
        on a.avatar_id = x.avatar_id
;

-- functions for commonly used queries
create or replace function public.share_sim(
    metadata JSONB
    , viewer_file text
    , user_id bigint
    , is_permanent boolean
    , is_public boolean
)
returns uuid
language plpgsql
as
$$
declare
    key uuid;
begin
    insert into simulations as s (metadata, viewer_file, is_permanent)
    values (
        metadata
        , viewer_file
        , is_permanent
    ) returning s.simulation_key into key;

    insert into user_simulations as u (simulation_key, user_id, is_public)
    values (
        key
        , user_id
        , is_public
    );

    return key;
end;
$$;

create or replace function public.link_avatar_to_sim(
    avatar text
    , key uuid
)
returns int
language plpgsql
as
$$
declare
    id int;
begin
    if exists(select 1 from avatars where avatar_name = avatar) then
        select avatar_id into id from avatars where avatar_name = avatar;
        insert into avatarsimulations (avatar_id, simulation_key) values (id, key);
        return id;
    end if;
    return -1;
end;
$$;

create or replace function public.add_db_sim(
    simulation_key uuid
    , git_hash text
    , config_hash text
    , author bigint
    , sim_description text
)
returns int 
language plpgsql
as
$$
declare
    key int;
begin
    insert into db_simulations as d (simulation_key, git_hash, config_hash, sim_description)
    values (
        add_db_sim.simulation_key
        , add_db_sim.git_hash
        , add_db_sim.config_hash
        , add_db_sim.sim_description
    ) returning d.db_id into key;

    -- add to authors list

    if exists(select 1 from users as u where u.user_id = add_db_sim.author) then
        insert into db_entry_authors as d (db_id, user_id)
        values (
            key
            , add_db_sim.author
        ) on conflict do nothing;
    end if;

    -- delete from user sims

    delete from user_simulations as u where u.simulation_key = add_db_sim.simulation_key;

    update simulations as s set is_permanent = 'true' where s.simulation_key = add_db_sim.simulation_key;

    return key;
end;
$$;

create or replace function public.replace_db_sim(
    old_key uuid
    , simulation_key uuid
    , git_hash text
    , config_hash text
    , author bigint
    , sim_description text
)
returns int 
language plpgsql
as
$$
declare
    key int;
begin
    update db_simulations as d set
        simulation_key = replace_db_sim.simulation_key
        , git_hash = replace_db_sim.git_hash
        , config_hash = replace_db_sim.config_hash
        , sim_description = replace_db_sim.sim_description
    where d.simulation_key = replace_db_sim.old_key
    returning d.db_id into key;

    -- add to authors list

    if exists(select 1 from users as u where u.user_id = replace_db_sim.author) then
        insert into db_entry_authors as d (db_id, user_id)
        values (
            key
            , replace_db_sim.author
        ) on conflict do nothing;
    end if;

    -- delete from user sims

    delete from user_simulations as u where u.simulation_key = replace_db_sim.simulation_key;

    update simulations as s set is_permanent = 'true' where s.simulation_key = replace_db_sim.simulation_key;

    update simulations as s set is_permanent = 'false' where s.simulation_key = replace_db_sim.old_key;

    return key;
end;
$$;

create or replace function public.get_or_insert_user(id bigint, name text)
returns table (
	user_id bigint,
	user_name text,
	user_role INTEGER
)
language plpgsql
as
$$
begin
	return query
		with s as (
			select users.user_id, users.user_name, users.user_role
			from users
			where users.user_id = id
		), i as (
			insert into users (user_id, user_name)
			select id, name
			where not exists (select 1 from s)
			returning users.user_id, users.user_name, users.user_role
		)
		select i.user_id, i.user_name, i.user_role
		from i
		union all
		select s.user_id, s.user_name, s.user_role
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

NOTIFY pgrst, 'reload schema';