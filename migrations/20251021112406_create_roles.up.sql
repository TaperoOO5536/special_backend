create role rl_user with password 'qwer1234';
alter role rl_user login;
grant connect on database 'specialappdb' to rl_user;

create role rl_admin with password 'qwer1234';
alter role rl_admin login;
grant connect on database 'specialappdb' to rl_admin;


grant select on items to rl_user;
grant select, update, insert on events to rl_user;
grant select on item_pictures to rl_user;
grant select, insert on event_pictures to rl_user;
grant select, insert, update on users to rl_user;
grant select, insert on orders to rl_user;
grant usage, select on sequence order_number to rl_user;
grant select, insert, update on order_items to rl_user;
grant select, insert, update, delete on user_events to rl_user;

grant select, insert, update, delete on items to rl_admin;
grant select, insert, update, delete on events to rl_admin;
grant select, insert, update, delete on item_pictures to rl_admin;
grant select, insert, update, delete on event_pictures to rl_admin;
grant select, insert, update on users to rl_admin;
grant select, update, delete on orders to rl_admin;
grant usage, select on sequence order_number to rl_admin;
grant select, insert, update, delete on order_items to rl_admin;
grant select, delete on user_events to rl_admin;
grant select, insert, update, delete on admins to rl_admin;
