CREATE EXTENSION IF NOT EXISTS citext;
CREATE type tvote as ENUM(
    '-1',
    '1'
);

create table IF NOT EXISTS author(
    id serial primary key,
    fullname citext not null ,
    nickname citext collate "ucs_basic" not null unique,
    email citext not null unique,
    about citext default('')
);
create unique index nickname_index on author(nickname);

create table IF NOT EXISTS forum(
    slug citext primary key,
    title citext not null ,
    author citext references author (nickname) on delete cascade,
    posts integer default ('0'),
    threads integer default ('0')
);
create unique index slug_index on forum(slug);

create table IF NOT EXISTS thread(
    id bigserial primary key,
    author citext collate "ucs_basic" references author (nickname) on delete cascade,
    created timestamp with time zone not null default now(),
    forum citext references forum (slug) on delete cascade,
    title citext not null,
    message citext not null,
    slug citext ,
    votes bigint default('0')
);

create table IF NOT EXISTS post(
    id bigserial primary key,
    author citext collate "ucs_basic" references author (nickname) on delete cascade,
    created timestamp with time zone not null default now(),
    message citext not null,
    forum citext references forum (slug) on delete cascade,
    thread bigint references thread(id) on delete cascade,
    isedited boolean default ('false'),
    parent bigint default('0')
);

create table IF NOT EXISTS vote(
    id serial primary key,
    nickname citext collate "ucs_basic" references author(nickname) on delete cascade,
    voice tvote not null,
    thread bigint references thread(id) on delete cascade,

    unique(nickname,thread)
);

CREATE FUNCTION vote_insert() RETURNS trigger AS '
BEGIN
update thread set votes = votes + NEW.voice::citext::integer where id = NEW.thread;
return NEW;
END;
' 
LANGUAGE  plpgsql;

CREATE FUNCTION vote_update() RETURNS trigger AS '
BEGIN
update thread set votes = votes + NEW.voice::citext::integer - OLD.voice::citext::integer where id = NEW.thread;
return NEW;
END;
' 
LANGUAGE  plpgsql;

CREATE FUNCTION thread_insert() RETURNS trigger AS '
BEGIN
update forum set threads = threads + 1 where slug = NEW.forum;
return NEW;
END;
' 
LANGUAGE  plpgsql;

CREATE FUNCTION post_insert() RETURNS trigger AS '
BEGIN
update forum set posts= posts + 1 where slug = NEW.forum;
return NEW;
END;
' 
LANGUAGE  plpgsql;


CREATE TRIGGER vote_trig_insert
AFTER INSERT ON vote FOR EACH ROW
EXECUTE PROCEDURE  vote_insert();

CREATE TRIGGER vote_trig_update
AFTER update ON vote FOR EACH ROW
EXECUTE PROCEDURE vote_update();

CREATE TRIGGER thread_trig_insert
AFTER insert ON thread FOR EACH ROW
EXECUTE PROCEDURE thread_insert();

CREATE TRIGGER post_trig_insert
AFTER insert ON post FOR EACH ROW
EXECUTE PROCEDURE post_insert();