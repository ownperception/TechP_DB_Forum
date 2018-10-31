drop table if exists author, post, thread, forum, vote cascade;  

drop type tvote;

drop function if exists vote_insert() cascade;
drop function if exists vote_update() cascade;
drop function if exists thread_insert() cascade;
drop function if exists post_insert() cascade;

drop trigger if exists vote_trig_update on vote cascade;
drop trigger if exists vote_trig_insert on vote cascade;
drop trigger if exists thread_trig_insert on vote cascade;
drop trigger if exists post_trig_insert on vote cascade;