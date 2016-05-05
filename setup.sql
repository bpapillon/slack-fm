CREATE DATABASE umbelfm;
CREATE USER umbelfm WITH PASSWORD 'umbelfm';
GRANT ALL PRIVILEGES ON DATABASE umbelfm TO umbelfm;

CREATE TABLE Recommendations (
	id serial,
	user_id int,
	post_time decimal,
	service_name char(256),
	title char(256),
	url char(256),
	thumb_url char(256),
	thumb_width int,
	thumb_height int,
	audio_html char(256),
	audio_height int,
	audio_width int,
	PRIMARY KEY (id)
);
GRANT ALL PRIVILEGES ON TABLE Recommendations TO umbelfm;
GRANT ALL PRIVILEGES ON TABLE recommendations_id_seq TO umbelfm;

CREATE TABLE Users (
	id serial,
	slack_id char(256),
	photo_url char(256),
	name char(256),
	slug char(256),
	PRIMARY KEY (id)
);
GRANT ALL PRIVILEGES ON TABLE Users TO umbelfm;
GRANT ALL PRIVILEGES ON TABLE users_id_seq TO umbelfm;

CREATE TABLE Tags (
	id serial,
	tag char(256),
	PRIMARY KEY (id)
);
GRANT ALL PRIVILEGES ON TABLE Tags TO umbelfm;
GRANT ALL PRIVILEGES ON TABLE tags_id_seq TO umbelfm;

CREATE TABLE Recommendation_tags (
	id serial,
	tag_id int,
	recommendation_id int,
	user_id int,
	PRIMARY KEY (id)
);
GRANT ALL PRIVILEGES ON TABLE Recommendation_tags TO umbelfm;
GRANT ALL PRIVILEGES ON TABLE recommendation_tags_id_seq TO umbelfm;
