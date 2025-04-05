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
                                                                       ('The Sunshine', 'Lane Langstaff'),
                                                                       ('Nyan Cat News', 'Eden Wilcher');

-- Publications
INSERT INTO newspaper_publication (id, newspaper_name, special, published, publish_date) VALUES ('ID-PUB-ABC123-DEF436', 'Falling Times', false, false, '2025-03-25 21:15:52.000000'),
                                                                                                ('ID-PUB-AET123-DEF636', 'Quacker''s Manual', false, false, '2025-03-22 22:42:51.000000'),
                                                                                                ('ID-PUB-ABR323-DEF436', 'TimTom Daily', false, false, '2025-03-24 12:45:50.000000'),
                                                                                                ('ID-PUB-QBC453-DZR936', 'The Sunshine', false, false, '2025-03-29 05:28:49.000000'),
                                                                                                ('ID-PUB-QBC453-ERT52A', 'The Sunshine', true, false, '2025-04-02 12:48:42.000000'),
                                                                                                ('ID-PUB-AVC963-ASQ176', 'Nyan Cat News', false, false, '2025-03-30 20:39:42.000000');

-- Articles
INSERT INTO newspaper_article (id, title, subtitle, author, flair, written, publication_id, html_body, raw_body) VALUES
                              ('ID-ARTICLE-ABC123-ABC125', 'Example Titel #1', '', 'Lane Langstaff', '',
                               '2025-04-02 12:48:31.000000', 'ID-PUB-QBC453-ERT52A',
                               '<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.<br>Aenean volutpat dignissim metus, non eleifend tortor cursus eget.</p>',
                               'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'||E'\n\n'||'Aenean volutpat dignissim metus, non eleifend tortor cursus eget.'),
                              ('ID-ARTICLE-DBC123-ABC126', 'Example Titel #2', '', 'Eden Wilcher', 'Test, abc, LoL',
                               '2025-04-02 13:12:07.000000', 'ID-PUB-QBC453-ERT52A',
                               '<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.<br>Aenean volutpat dignissim metus, non eleifend tortor cursus eget.</p>',
                               'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'||E'\n\n'||'Aenean volutpat dignissim metus, non eleifend tortor cursus eget.'),
                              ('ID-ARTICLE-EBC123-ABC127', 'Example Titel #3', '', 'Timothy Nunez', '',
                               '2025-03-25 11:21:01.000000', 'ID-PUB-AET123-DEF636',
                               '<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.<br>Aenean volutpat dignissim metus, non eleifend tortor cursus eget.</p>',
                               'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'||E'\n\n'||'Aenean volutpat dignissim metus, non eleifend tortor cursus eget.'),
                              ('ID-ARTICLE-FBC123-ABC128', 'Example Titel #4', '', 'Lane Langstaff', 'Another Flair',
                               '2025-04-01 15:48:49.000000', 'ID-PUB-QBC453-DZR936',
                               '<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.<br>Aenean volutpat dignissim metus, non eleifend tortor cursus eget.</p>',
                               'Lorem ipsum dolor sit amet, consectetur adipiscing elit.'||E'\n\n'||'Aenean volutpat dignissim metus, non eleifend tortor cursus eget.');

-- Todo: Documents

-- Todo: Discussions

-- Todo: Votes

-- Notes
INSERT INTO blackboard_note (id, title, author, flair, posted, body, blocked) VALUES
                            ('ID-NOTE-ABC123-ABC923', 'Example Note #1', 'Eden Wilcher', '', '2025-03-01 12:00:00.000000',
                             '<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.<br>Aenean volutpat dignissim metus, non eleifend tortor cursus eget.</p>', false),
                            ('ID-NOTE-ABC223-ABC823', 'Example Note #2', 'Gavin George', '', '2025-03-02 12:00:00.000000',
                             '<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.<br>Aenean volutpat dignissim metus, non eleifend tortor cursus eget.</p>', false),
                            ('ID-NOTE-ABC323-ABC723', 'Example Note #3', 'Lane Langstaff', 'Test-Flair', '2025-03-03 12:00:00.000000',
                             '<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.<br>Aenean volutpat dignissim metus, non eleifend tortor cursus eget.</p>', false),
                            ('ID-NOTE-ABC423-ABC623', 'Example Note #4', 'Lonnie Hampton', '', '2025-03-04 12:00:00.000000',
                             '<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.<br>Aenean volutpat dignissim metus, non eleifend tortor cursus eget.</p>', false),
                            ('ID-NOTE-ABC523-ABC523', 'Example Note #5', 'Timothy Nunez', '', '2025-03-05 12:00:00.000000',
                             '<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.<br>Aenean volutpat dignissim metus, non eleifend tortor cursus eget.</p>', true),
                            ('ID-NOTE-ABC623-ABC423', 'Example Note #6', 'Vania Walton', 'Ltn., Army Commander', '2025-03-06 12:00:00.000000',
                             '<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.<br>Aenean volutpat dignissim metus, non eleifend tortor cursus eget.</p>', false),
                            ('ID-NOTE-ABC723-ABC323', 'Example Note #7', 'Eden Wilcher', '', '2025-03-07 12:00:00.000000',
                             '<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.<br>Aenean volutpat dignissim metus, non eleifend tortor cursus eget.</p>', false),
                            ('ID-NOTE-ABC823-ABC223', 'Example Note #8', 'Vania Walton', 'Cptn.', '2025-03-08 12:00:00.000000',
                             '<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.<br>Aenean volutpat dignissim metus, non eleifend tortor cursus eget.</p>', false),
                            ('ID-NOTE-ABC923-ABC123', 'Example Note #9', 'Lonnie Hampton', '', '2025-03-09 12:00:00.000000',
                             '<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.<br>Aenean volutpat dignissim metus, non eleifend tortor cursus eget.</p>', false);
INSERT INTO blackboard_references (base_note_id, reference_id)  VALUES ('ID-NOTE-ABC823-ABC223', 'ID-NOTE-ABC323-ABC723'),
                                                                       ('ID-NOTE-ABC823-ABC223', 'ID-NOTE-ABC623-ABC423'),
                                                                       ('ID-NOTE-ABC823-ABC223', 'ID-NOTE-ABC423-ABC623'),
                                                                       ('ID-NOTE-ABC623-ABC423', 'ID-NOTE-ABC123-ABC923'),
                                                                       ('ID-NOTE-ABC723-ABC323', 'ID-NOTE-ABC323-ABC723'),
                                                                       ('ID-NOTE-ABC723-ABC323', 'ID-NOTE-ABC423-ABC623'),
                                                                       ('ID-NOTE-ABC523-ABC523', 'ID-NOTE-ABC123-ABC923'),
                                                                       ('ID-NOTE-ABC923-ABC123', 'ID-NOTE-ABC623-ABC423'),
                                                                       ('ID-NOTE-ABC923-ABC123', 'ID-NOTE-ABC423-ABC623'),
                                                                       ('ID-NOTE-ABC923-ABC123', 'ID-NOTE-ABC123-ABC923'),
                                                                       ('ID-NOTE-ABC923-ABC123', 'ID-NOTE-ABC523-ABC523');

-- Chats
INSERT INTO chat_rooms (room_id, name, created, member) VALUES ('ID-CHAT-ATD341-TWAS11', 'Plan for next Election', '2025-03-30 20:04:13.000000', ARRAY['Vania Walton', 'Clarissa Stephens']),
                                                               ('ID-CHAT-ARD541-TWAS12', 'Privat Conversation', '2025-02-12 13:42:22.000000', ARRAY['Gavin George', 'Lonnie Hampton']),
                                                               ('ID-CHAT-ASD241-TWAS13', 'Big Chat', '2025-04-01 08:13:31.000000', ARRAY['Gavin George', 'Clarissa Stephens', 'Lonnie Hampton', 'Eden Wilcher']);
INSERT INTO chat_rooms_to_account (room_id, account_name, new_message) VALUES ('ID-CHAT-ATD341-TWAS11', 'Vania Walton', false),
                                                                              ('ID-CHAT-ATD341-TWAS11', 'Lane Langstaff', false),
                                                                              ('ID-CHAT-ARD541-TWAS12', 'Gavin George', true),
                                                                              ('ID-CHAT-ARD541-TWAS12', 'Lonnie Hampton', false),
                                                                              ('ID-CHAT-ASD241-TWAS13', 'Gavin George', false),
                                                                              ('ID-CHAT-ASD241-TWAS13', 'Lane Langstaff', true),
                                                                              ('ID-CHAT-ASD241-TWAS13', 'Lonnie Hampton', true),
                                                                              ('ID-CHAT-ASD241-TWAS13', 'Eden Wilcher', false);
INSERT INTO chat_messages (room_id, sender, message, send_time) VALUES ('ID-CHAT-ARD541-TWAS12', 'Lonnie Hampton', 'Hello Gavin,<br>I haven''t heard from you in a long time.<br>Do you want to catch up soon?', '2025-04-02 13:42:22.000000'),
                                                                       ('ID-CHAT-ARD541-TWAS12', 'Lonnie Hampton', 'Just hit me up if you are ready to talk.', '2025-04-02 13:44:22.000000'),
                                                                       ('ID-CHAT-ASD241-TWAS13', 'Gavin George', 'Welcome, welcome guys.<br>I hope we can do great things together!', '2025-04-02 13:41:12.000000'),
                                                                       ('ID-CHAT-ASD241-TWAS13', 'Eden Wilcher', 'Not sure why I am here ...', '2025-04-02 13:42:34.000000'),
                                                                       ('ID-CHAT-ASD241-TWAS13', 'Gavin George', 'Oh come on, Eden, you will have fun too, right?<br>I really hope we can all become good friends.', '2025-04-02 13:46:42.000000'),
                                                                       ('ID-CHAT-ASD241-TWAS13', 'Eden Wilcher', 'Not happening, bro.', '2025-04-02 13:52:55.000000');