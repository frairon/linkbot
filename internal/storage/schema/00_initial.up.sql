

-- stores users
create table users(
  id INTEGER not null primary key,
  Name text not null
);

create table user_session(
  user_id INTEGER not null primary key,
  chat_id INTEGER not null,
  last_user_action DATETIME,
  data string
);


create table user_links(
  link_id text not null primary key,
  user_id INTEGER not null,
  category text not null default "default",
  hidden boolean not null default 0,
  link text not null,
  headline text not null,
  added DATETIME not null
);


create table settings(
  key string primary key,
  value string
);