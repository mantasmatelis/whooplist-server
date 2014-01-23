DROP SCHEMA IF EXISTS wl CASCADE;
CREATE SCHEMA wl AUTHORIZATION whooplist;

CREATE TABLE wl.user
(
  id serial NOT NULL PRIMARY KEY,
  email text NOT NULL,
  name text NOT NULL,
  fname text,
  lname text,
  phone text,
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
CREATE INDEX session_key_index ON wl.session(key);

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

CREATE TABLE wl.friend
(
  id serial NOT NULL PRIMARY KEY,
  from_id integer REFERENCES wl.user(id),
  to_id integer REFERENCES wl.user(id),
  timestamp timestamp without time zone DEFAULT LOCALTIMESTAMP NOT NULL
);
CREATE INDEX friend_from_id ON wl.friend(from_id);
CREATE INDEX friend_to_id ON wl.friend(to_id);

CREATE TABLE wl.whooplist_item
(
  id serial NOT NULL PRIMARY KEY,
  list_id integer REFERENCES wl.list(id),
  place_id integer REFERENCES wl.list(id),
  score integer
);
CREATE INDEX whooplist_item_list_id_index ON wl.whooplist_item(list_id);
CREATE INDEX whooplist_item_place_id_index ON wl.whooplist_item(place_id);

INSERT INTO wl.user
  (id, email, name, password_hash, role) VALUES
  (1, 'base@whooplist.com', 'Base Data', '', 'b');

INSERT INTO wl.user
  (id, email, name, fname, lname, school, password_hash, role) VALUES
  (11, 'mantas@whooplist.com', 'Mantas Matelis', 'Mantas', 'Matelis',
    'University of Waterloo', 'IV4r5xGvqeot7jWZiW4wUcnxUW/h4TEFzbT2COTpvv4=', 'a'),
  (12, 'dev@whooplist.com', 'Dev Chakraborty', 'Dev', 'Chakraborty',
    'Western University', 'qZggmcGnjNRlMUjwI4kWmCL56sJnQfvqA32JjtXPAuA=', 'a'),
  (13, 'jitesh@whooplist.com', 'Jitesh Vyas', 'Jitesh', 'Vyas',
    'Western University', '4JM/xLOrGrbFG0CyEcN8bjKAL0wpAfHzqcexi0d0jHQ=', 'a');

INSERT INTO wl.friend
  (from_id, to_id) VALUES
  (11, 12), (11, 13),
  (12, 11), (12, 13),
  (13, 11), (13, 12);

ALTER SEQUENCE wl.user_id_seq RESTART WITH 10000;

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

GRANT ALL ON wl.user, wl.session, wl.list, wl.place, wl.list_item,
  wl.feed_item, wl.friend, wl.whooplist_item TO whooplist;

GRANT ALL ON ALL SEQUENCES IN SCHEMA wl TO whooplist;

-- BEGIN BASE PLACE DATA

-- BEGIN TEST DATA

INSERT INTO wl.session (user_id, key) VALUES (11, 'mantas');
INSERT INTO wl.session (user_id, key) VALUES (12, 'dev');
INSERT INTO wl.session (user_id, key) VALUES (13, 'jitesh');

INSERT INTO wl.user
  (id, email, name, password_hash, role) VALUES
  (1001, 'test1@whooplist.com', 'Test 1', 'P', 't'),
  (1002, 'test2@whooplist.com', 'Test 2', 'P', 't'),
  (1003, 'test3@whooplist.com', 'Test 3', 'P', 't'),
  (1004, 'test4@whooplist.com', 'Test 4', 'P', 't'),
  (1005, 'test5@whooplist.com', 'Test 5', 'P', 't'),
  (1006, 'test6@whooplist.com', 'Test 6', 'P', 't'),
  (1007, 'test7@whooplist.com', 'Test 7', 'P', 't'),
  (1008, 'test8@whooplist.com', 'Test 8', 'P', 't'),
  (1009, 'test9@whooplist.com', 'Test 9', 'P', 't'),
  (1010, 'test10@whooplist.com', 'Test 10', 'P', 't');

INSERT INTO wl.session (user_id, key) VALUES 
  (1001, 'test1'), (1002, 'test2'), (1003, 'test3'),
  (1004, 'test4'), (1005, 'test5'), (1006, 'test6'),
  (1007, 'test7'), (1008, 'test8'), (1009, 'test9'),
  (1010, 'test10');

INSERT INTO wl.friend
  (from_id, to_id) VALUES
  (1001, 1002), (1001, 1003), (1001, 1004), (1001, 1005), (1001, 1006), (1001, 1007), (1001, 1008), (1001, 1009),
  (1002, 1001), (1002, 1003), (1002, 1004),
  (1003, 1001), (1003, 1008), (1003, 1009),
  (1004, 1005), (1004, 1007);
