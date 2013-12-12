DROP SCHEMA IF EXISTS wl CASCADE;
CREATE SCHEMA wl AUTHORIZATION whooplist;

CREATE TABLE wl.user
(
  id serial NOT NULL PRIMARY KEY,
  email text NOT NULL,
  name text NOT NULL,
  fname text,
  lname text,
  school text,
  birthday date,
  gender integer,
  picture text,
  password_hash text NOT NULL,
  role character NOT NULL
);
CREATE INDEX user_lower_email_index ON wl.user(lower(email));

CREATE TABLE wl.session
(
  id serial NOT NULL PRIMARY KEY,
  user_id integer NOT NULL REFERENCES wl.user(id),
  key text NOT NULL,
  last_auth timestamp without time zone DEFAULT LOCALTIMESTAMP NOT NULL,
  last_use timestamp without time zone DEFAULT LOCALTIMESTAMP NOT NULL
);
CREATE INDEX session_key ON wl.session(key);

CREATE TABLE wl.list
(
  id serial NOT NULL PRIMARY KEY,
  name text NOT NULL,
  icon text NOT NULL,
  children text
);

CREATE TABLE wl.place
(
  id serial NOT NULL PRIMARY KEY,
  latitude double precision NOT NULL,
  longitude double precision NOT NULL,
  factual_id text NOT NULL,
  name text NOT NULL,
  address text,
  locality text,
  region text,
  postcode text,
  country text,
  telephone text,
  website text,
  email text
);
CREATE INDEX place_factual_id_index ON wl.place(factual_id);
CREATE INDEX place_latitude_index ON wl.place(latitude);
CREATE INDEX place_longitude_index ON wl.place(longitude);

CREATE TABLE wl.list_item
(
  id serial NOT NULL PRIMARY KEY,
  user_id integer NOT NULL REFERENCES wl.user(id),
  list_id integer NOT NULL REFERENCES wl.list(id),
  place_id integer NOT NULL REFERENCES wl.place(id),
  rank integer NOT NULL
);
CREATE INDEX list_item_user_id_index ON wl.list_item(user_id);

CREATE TABLE wl.feed_item
(
  id serial NOT NULL PRIMARY KEY,
  timestamp timestamp without time zone DEFAULT LOCALTIMESTAMP NOT NULL,
  user_id integer REFERENCES wl.user(id),
  latitude double precision,
  longitude double precision,
  place_id integer REFERENCES wl.place(id),
  list_id integer REFERENCES wl.list(id),
  picture text,
  type integer NOT NULL,
  aux_string text,
  aux_int integer
);


CREATE TABLE wl.whooplist_item
(
  id serial NOT NULL PRIMARY KEY,
  list_id integer REFERENCES wl.list(id),
  place_id integer REFERENCES wl.list(id),
  score integer
);
CREATE INDEX whooplist_item_list_id_index ON wl.whooplist_item(list_id);
CREATE INDEX whooplist_item_place_id_index ON wl.whooplist_item(place_id);

INSERT INTO wl.list (id, name, children, icon) VALUES
	(1, 'Eat', '6,7,8,9,10', '$list_icons/eat.svg$'),
	(2, 'Drink', '11,12,13,14', '$list_icons/drinks.svg$'),
	(3, 'Date', '15,11,7,8,10,16,17', '$list_icons/date.svg$'),
	(4, 'Shop', '18,19', '$list_icons/shop.svg$'),
	(5, 'Discover', '14, 16, 20, 21, 22', '$list_icons/explore.svg$'), 
	
	(6, 'Breakfast', '', '$list_icons/eat.svg$'),
	(7, 'Lunch', '', '$list_icons/eat.svg$'),
	(8, 'Dinner', '', '$list_icons/eat.svg$'),
	(9, '24 Hour Food', '', '$list_icons/timedeat.svg$'), 
	(10, 'Dessert', '', '$list_icons/dessert.svg$'),
	(11, 'Coffee', '', '$list_icons/coffee.svg$'),
	(12, 'Specialty', '', '$list_icons/bubbletea.svg$'), 
	(13, 'Bars', '', '$list_icons/bars.svg$'),
	(14, 'Nightlife', '', '$list_icons/nightlife.svg$'),
	(15, 'Attractions', '', '$list_icons/attractions.svg$'), 
	(16, 'Hang Out', '', '$list_icons/hangout.svg$'), 
	(17, 'Explore', '', '$list_icons/explore.svg$'),
	(18, 'Grocery', '', '$list_icons/groceries.svg$'),
	(19, 'Retail', '', '$list_icons/shop.svg$'),
	(20, 'Recreation', '', '$list_icons/recreation.svg$'),
	(21, 'Hike', '', '$list_icons/hike.svg$'),
	(22, 'Meet', '', '$list_icons/meet.evg$');
