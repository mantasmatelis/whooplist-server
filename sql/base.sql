CREATE TABLE user
(
  id serial NOT NULL PRIMARY KEY,
  email text NOT NULL,
  name text NOT NULL,
  fname text,
  lname text,
  birthday date,
  gender integer,
  password_hash text NOT NULL,
  role character NOT NULL,
);

CREATE TABLE session
(
  id serial NOT NULL PRIMARY KEY,
  user_id integer NOT NULL REFERENCES user(id),
  key text NOT NULL,
  last_auth timestamp without time zone NOT NULL,
  last_use timestamp without time zone NOT NULL,
);

CREATE TABLE list
(
  id serial NOT NULL PRIMARY KEY,
  name text NOT NULL,
  icon text NOT NULL,
  children text,
);


CREATE TABLE place
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
  email text,
);

CREATE TABLE list_item
(
  id serial NOT NULL PRIMARY KEY,
  user_id integer NOT NULL REFERENCES user(id),
  list_id integer NOT NULL REFERENCES list(id),
  place_id integer NOT NULL REFERENCES place(id),
  rank integer NOT NULL,
);

CREATE TABLE feed_item
(
  id serial NOT NULL PRIMARY KEY,
  timestamp timestamp without time zone NOT NULL,
  user_id integer REFERENCES user(id),
  place_id integer REFERENCES place(id),
  list_id integer REFERENCES list(id),
  picture text,
  type integer NOT NULL,
  aux_string text,
  aux_int integer,
);


INSERT INTO list (id, name, children, icon) VALUES
	(1, "Eat", "6,7,8,9,10", "$assets/"), (2, "Drink", "11,12,13,14", ""),
	(3, "Date", "15,11,7,8,10,16,17", ""), (4, "Shop", "18,19", ""),
	(5, "Discover", "14, 16, 20, 21, 22", ""), 
	
	(6, "Breakfast", "", ""), (7, "Lunch", "", ""), (8, "Dinner", "", ""),
	(9, "24 Hour Food", "", ""), (10, "Dessert", "", ""), (11, "Coffee", "", ""),
	(12, "Specialty", "", ""), (13, "Bars", "", ""), (14, "Nightlife", "", ""),
	(15, "Attractions", "", ""), (16, "Hang Out", "", ""), (17, "Explore", "", ""),
	(18, "Grocery", "", ""), (19, "Retail", "", ""), (20, "Recreation", "", ""),
	(21, "Hike", "", ""), (22, "Meet", "", "");
