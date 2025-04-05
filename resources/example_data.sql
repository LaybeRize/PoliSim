-- DROP SCHEMA public CASCADE;
-- CREATE SCHEMA public;

-- Password for all Accounts is 1234567890
INSERT INTO account (name, username, password, role, blocked, font_size, time_zone) VALUES
                    ('Clarissa Stephens','Clarissa Stephens','$2a$10$99SOqMF93J320lJYWdfjjeGfdPSGnjG2QDiZcUFyNSbo4Fk7LwcMq',4, true,100,'UTC');
INSERT INTO account (name, username, password, role, blocked, font_size, time_zone) VALUES
                    ('Vania Walton','Cool Username','$2a$10$99SOqMF93J320lJYWdfjjeGfdPSGnjG2QDiZcUFyNSbo4Fk7LwcMq',4, false,100,'Europe/Berlin');
INSERT INTO account (name, username, password, role, blocked, font_size, time_zone) VALUES
                    ('Timothy Nunez','XxNavigatorXx','$2a$10$99SOqMF93J320lJYWdfjjeGfdPSGnjG2QDiZcUFyNSbo4Fk7LwcMq',4, false,100,'Asia/Shanghai');
INSERT INTO account (name, username, password, role, blocked, font_size, time_zone) VALUES
                    ('Lonnie Hampton','Lonnie Hampton','$2a$10$99SOqMF93J320lJYWdfjjeGfdPSGnjG2QDiZcUFyNSbo4Fk7LwcMq',4, false,100,'America/Juneau');
INSERT INTO account (name, username, password, role, blocked, font_size, time_zone) VALUES
                    ('Lane Langstaff','Lane Langstaff','$2a$10$99SOqMF93J320lJYWdfjjeGfdPSGnjG2QDiZcUFyNSbo4Fk7LwcMq',4, false,100,'Etc/GMT-5');
INSERT INTO account (name, username, password, role, blocked, font_size, time_zone) VALUES ('Gavin George','Gavin George','',5, false,100,'UTC');
INSERT INTO account (name, username, password, role, blocked, font_size, time_zone) VALUES ('Eden Wilcher','Eden Wilcher','',5, false,100,'UTC');
-- Ownership
INSERT INTO ownership (account_name, owner_name) VALUES ('Clarissa Stephens', 'Clarissa Stephens');
INSERT INTO ownership (account_name, owner_name) VALUES ('Eden Wilcher', 'Lonnie Hampton');
INSERT INTO ownership (account_name, owner_name) VALUES ('Gavin George', 'Timothy Nunez');
INSERT INTO ownership (account_name, owner_name) VALUES ('Lane Langstaff', 'Lane Langstaff');
INSERT INTO ownership (account_name, owner_name) VALUES ('Lonnie Hampton', 'Lonnie Hampton');
INSERT INTO ownership (account_name, owner_name) VALUES ('Timothy Nunez', 'Timothy Nunez');
INSERT INTO ownership (account_name, owner_name) VALUES ('Vania Walton', 'Vania Walton');
-- Organisations
INSERT INTO organisation (name, main_group, sub_group, visibility, flair, users, admins) VALUES
                         ('The Villa', 'Land based Objects', 'Houses', 1, '', ARRAY['Gavin George', 'Lane Langstaff'] , ARRAY['Timothy Nunez']);
INSERT INTO organisation_to_account (organisation_name, account_name, is_admin) VALUES ('The Villa', 'Lane Langstaff', false),
                                                                                       ('The Villa', 'Gavin George', false),
                                                                                       ('The Villa', 'Timothy Nunez', true);
INSERT INTO organisation (name, main_group, sub_group, visibility, flair, users, admins) VALUES
                         ('Toilet-House', 'Land based Objects', 'Houses', 0, 'Supt.', ARRAY[] , ARRAY['Eden Wilcher']);
INSERT INTO organisation_to_account (organisation_name, account_name, is_admin) VALUES ('Toilet-House', 'Eden Wilcher', true);
INSERT INTO organisation (name, main_group, sub_group, visibility, flair, users, admins) VALUES
                         ('Super-Bunker', 'Land based Objects', 'Underground', 2, '', ARRAY['Eden Wilcher', 'Timothy Nunez'] , ARRAY['Gavin George', 'Vania Walton']);
INSERT INTO organisation_to_account (organisation_name, account_name, is_admin) VALUES ('Super-Bunker', 'Vania Walton', true),
                                                                                       ('Super-Bunker', 'Gavin George', true),
                                                                                       ('Super-Bunker', 'Timothy Nunez', false),
                                                                                       ('Super-Bunker', 'Eden Wilcher', false);
INSERT INTO organisation (name, main_group, sub_group, visibility, flair, users, admins) VALUES
                         ('Freighter', 'Water based Objects', 'Ships', 1, 'Captain', ARRAY['Eden Wilcher'] , ARRAY['Gavin George', 'Timothy Nunez', 'Vania Walton']);
INSERT INTO organisation_to_account (organisation_name, account_name, is_admin) VALUES ('Freighter', 'Vania Walton', true),
                                                                                       ('Freighter', 'Gavin George', true),
                                                                                       ('Freighter', 'Timothy Nunez', true),
                                                                                       ('Freighter', 'Eden Wilcher', false);
-- Titles
INSERT INTO title (name, main_group, sub_group, flair) VALUES
                  ('Head Chief of Heating', 'Housing', 'Household Management', 'HCoH');
INSERT INTO title_to_account (title_name, account_name) VALUES ('Head Chief of Heating', 'Vania Walton'),
                                                               ('Head Chief of Heating', 'Gavin George'),
                                                               ('Head Chief of Heating', 'Timothy Nunez'),
                                                               ('Head Chief of Heating', 'Eden Wilcher');
INSERT INTO title (name, main_group, sub_group, flair) VALUES
                  ('Simple Worker', 'Housing', 'Construction', '');
INSERT INTO title (name, main_group, sub_group, flair) VALUES
                  ('Construction Overseer', 'Ships', 'Construction', 'Constr. Overseer');
INSERT INTO title (name, main_group, sub_group, flair) VALUES
                  ('Sailor', 'Ships', 'Usage', '');
INSERT INTO title_to_account (title_name, account_name) VALUES ('Sailor', 'Vania Walton'),
                                                               ('Sailor', 'Eden Wilcher');
-- Newspapers
INSERT INTO newspaper (name) VALUES ('Falling Times'), ('Quacker''s Manual'), ('TimTom Daily'), ('The Sunshine'), ('Nyan Cat News');
INSERT INTO newspaper_to_account (newspaper_name, account_name) VALUES ('Falling Times', 'Gavin George'),
                                                                       ('Falling Times', 'Vania Walton'),
                                                                       ('Quacker''s Manual', 'Timothy Nunez'),
                                                                       ('TimTom Daily', 'Vania Walton'),
                                                                       ('TimTom Daily', 'Gavin George'),
                                                                       ('The Sunshine', 'Eden Wilcher'),
                                                                       ('The Sunshine', 'Gavin George'),
                                                                       ('The Sunshine', 'Timothy Nunez'),
                                                                       ('Nyan Cat News', 'Eden Wilcher');

-- Publications
INSERT INTO newspaper_publication (id, newspaper_name, special, published, publish_date) VALUES ('ID-PUB-ABC123-DEF436', 'Falling Times', false, false, '2025-03-25 21:15:52.000000'),
                                                                                                ('ID-PUB-AET123-DEF636', 'Quacker''s Manual', false, false, '2025-03-22 22:42:51.000000'),
                                                                                                ('ID-PUB-ABR323-DEF436', 'TimTom Daily', false, false, '2025-03-24 12:45:50.000000'),
                                                                                                ('ID-PUB-QBC453-DZR936', 'The Sunshine', false, false, '2025-03-29 05:28:49.000000'),
                                                                                                ('ID-PUB-QBC453-ERT52A', 'The Sunshine', true, false, '2025-04-02 12:48:42.000000'),
                                                                                                ('ID-PUB-AVC963-ASQ176', 'Nyan Cat News', false, false, '2025-03-30 20:39:42.000000');
