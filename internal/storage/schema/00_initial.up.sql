

-- stores users
create table users(
  id INTEGER not null primary key,
  Name text
);

create table user_session(
  user_id INTEGER not null primary key,
  chat_id INTEGER not null,
  last_user_action DATETIME,
  data string
);


create table user_links(
  link_id INTEGER not null primary key,
  user_id INTEGER not null,
  link text,
  headline text,
  added DATETIME
);



create table settings(
  key string primary key,
  value string
);