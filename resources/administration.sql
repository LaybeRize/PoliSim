-- noinspection ALL

-- Delete Chatroom
/*
Parameter 1: Chatroom ID
*/
DELETE FROM chat_messages WHERE room_id = $1; --PARAM=0:1
DELETE FROM chat_rooms_to_account WHERE room_id = $1; --PARAM=0:1
DELETE FROM chat_rooms WHERE room_id = $1; --PARAM=0:1
SELECT ARRAY[COUNT(room_id)::TEXT] FROM chat_rooms WHERE room_id = $1; --PARAM=0:1

-- Delete Article
/*
Parameter 1: Article ID<br>
Parameter 2: Publication ID
*/
DELETE FROM newspaper_article WHERE id = $1 AND publication_id = $2; --PARAM=0:2
WITH news_amount AS ( SELECT COUNT(id) AS number FROM newspaper_article WHERE publication_id = $2 )
DELETE FROM newspaper_publication WHERE 0 = (SELECT number FROM news_amount) AND id = $2 AND (special = true or published = true); --PARAM=1:2
SELECT ARRAY[( SELECT COUNT(id) FROM newspaper_article WHERE id = $1 )::TEXT,
           ( SELECT COUNT(id) FROM newspaper_publication WHERE id = $2 )::TEXT] FROM version_management; --PARAM=0:2

-- Delete Publication
/*
Parameter 1: Publication ID
*/
DELETE FROM newspaper_article WHERE publication_id = $1; --PARAM=0:1
DELETE FROM newspaper_publication WHERE id = $1; --PARAM=0:1
SELECT ARRAY[COUNT(id)::TEXT] FROM newspaper_publication WHERE id = $1; --PARAM=0:1

-- Delete Newspaper
/*
Parameter 1: Newspaper Name
*/
DELETE FROM newspaper_article WHERE publication_id =
    ANY(ARRAY(SELECT id FROM newspaper_publication WHERE newspaper_name = $1));  --PARAM=0:1
DELETE FROM newspaper_publication WHERE newspaper_name = $1; --PARAM=0:1
DELETE FROM newspaper_to_account WHERE newspaper_name = $1; --PARAM=0:1
DELETE FROM newspaper WHERE name = $1; --PARAM=0:1
SELECT ARRAY[COUNT(name)::TEXT] FROM newspaper WHERE name = $1; --PARAM=0:1

-- Delete Letter
/*
Parameter 1: Letter ID
*/
DELETE FROM letter_to_account WHERE letter_id = $1; --PARAM=0:1
DELETE FROM letter WHERE id = $1; --PARAM=0:1
SELECT ARRAY[COUNT(id)::TEXT] FROM letter WHERE id = $1; --PARAM=0:1

-- Delete Note
/*
Parameter 1: Note ID
*/
DELETE FROM blackboard_references WHERE base_note_id = $1 OR reference_id = $1; --PARAM=0:1
DELETE FROM blackboard_note WHERE id = $1; --PARAM=0:1
SELECT ARRAY[COUNT(id)::TEXT] FROM blackboard_note WHERE id = $1; --PARAM=0:1

-- Delete Document
/*
Parameter 1: Note ID
*/
DELETE FROM has_voted WHERE vote_id =
    ANY(ARRAY(SELECT id FROM document_to_vote WHERE document_id = $1)); --PARAM=0:1
DELETE FROM comment_to_document WHERE document_id = $1; --PARAM=0:1
DELETE FROM document_to_vote WHERE document_id = $1; --PARAM=0:1
DELETE FROM document_to_account WHERE document_id = $1; --PARAM=0:1
DELETE FROM document WHERE id = $1; --PARAM=0:1
SELECT ARRAY[COUNT(id)::TEXT] FROM document WHERE id = $1; --PARAM=0:1

-- Remove Document Tag
/*
Parameter 1: Tag ID
*/
UPDATE document SET extra_info = subquery.element
FROM  (SELECT document.id AS doc_id, extra_info #- ARRAY['tags', (position - 1)::TEXT] AS element
       FROM document, jsonb_array_elements(extra_info->'tags') WITH ORDINALITY arr(elem, position)
       WHERE elem->>'id' = $1) as subquery
WHERE id = subquery.doc_id RETURNING ARRAY[id]; --PARAM=0:1