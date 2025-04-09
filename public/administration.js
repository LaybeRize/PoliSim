let updateQuery;
window.onload = function () {
    const selectEl = document.body.querySelector("#parameter-query-selection");
    const textAreaEl = document.body.querySelector("#query");
    const textInfoEl = document.body.querySelector("#query-parameter-info");
    function updateOnQuery() {
        textInfoEl.innerHTML = "Info:<br>";
        switch (selectEl.value) {
            case "del-chatroom":
                textInfoEl.innerHTML += "- Parameter 1: Chatroom ID";
                textAreaEl.value = "DELETE FROM chat_messages WHERE room_id = $1; --PARAM=0:1\n"+
                    "DELETE FROM chat_rooms_to_account WHERE room_id = $1; --PARAM=0:1\n"+
                    "DELETE FROM chat_rooms WHERE room_id = $1; --PARAM=0:1\n"+
                    "SELECT ARRAY[COUNT(room_id)::TEXT] FROM chat_rooms WHERE room_id = $1; --PARAM=0:1";
                break;
            case "del-article":
                textInfoEl.innerHTML += "- Parameter 1: Article ID<br>- Parameter 2: Publication ID";
                //noinspection ALL
                textAreaEl.value = "DELETE FROM newspaper_article WHERE id = $1 AND publication_id = $2; --PARAM=0:2\n"+
                    "WITH news_amount AS ( SELECT COUNT(id) AS number FROM newspaper_article WHERE publication_id = $1 )"+
                    "DELETE FROM newspaper_publication WHERE 0 == (SELECT number FROM news_amount) AND id = $1 AND (special = true or published = true); --PARAM=1:2"+
                    "SELECT ARRAY[( SELECT COUNT(id) FROM newspaper_article WHERE id = $1 )::TEXT,\n"+
                    "           ( SELECT COUNT(id) FROM newspaper_publication WHERE id = $2 )::TEXT] FROM version_management; --PARAM=0:2";
                break;
            case "del-publication":
                textInfoEl.innerHTML += "- Parameter 1: Publication ID";
                textAreaEl.value = "DELETE FROM newspaper_article WHERE publication_id = $1; --PARAM=0:1\n"+
                    "DELETE FROM newspaper_publication WHERE id = $1; --PARAM=0:1\n"+
                    "SELECT ARRAY[COUNT(id)::TEXT] FROM newspaper_publication WHERE id = $1; --PARAM=0:1";
                break;
            case "del-newspaper":
                textInfoEl.innerHTML += "- Parameter 1: Newspaper Name";
                textAreaEl.value = "DELETE FROM newspaper_article WHERE publication_id =\n"+
                    "    ANY(ARRAY(SELECT id FROM newspaper_publication WHERE newspaper_name = $1));  --PARAM=0:1\n"+
                    "DELETE FROM newspaper_publication WHERE newspaper_name = $1; --PARAM=0:1\n"+
                    "DELETE FROM newspaper_to_account WHERE newspaper_name = $1; --PARAM=0:1\n"+
                    "DELETE FROM newspaper WHERE name = $1; --PARAM=0:1\n"+
                    "SELECT ARRAY[COUNT(name)::TEXT] FROM newspaper WHERE name = $1; --PARAM=0:1";
                break;
            case "del-document":
                textInfoEl.innerHTML += "- Parameter 1: Document ID";
                textAreaEl.value = "DELETE FROM has_voted WHERE vote_id = \n"+
                    "    ANY(ARRAY(SELECT id FROM document_to_vote WHERE document_id = $1)); --PARAM=0:1\n"+
                    "DELETE FROM comment_to_document WHERE document_id = $1; --PARAM=0:1\n"+
                    "DELETE FROM document_to_vote WHERE document_id = $1; --PARAM=0:1\n"+
                    "DELETE FROM document_to_account WHERE document_id = $1; --PARAM=0:1\n"+
                    "DELETE FROM document WHERE id = $1; --PARAM=0:1\n"+
                    "SELECT ARRAY[COUNT(id)::TEXT] FROM document WHERE id = $1; --PARAM=0:1";
                break;
            case "del-letter":
                textInfoEl.innerHTML += "- Parameter 1: Letter ID";
                textAreaEl.value = "DELETE FROM letter_to_account WHERE letter_id = $1; --PARAM=0:1\n"+
                    "DELETE FROM letter WHERE id = $1; --PARAM=0:1\n"+
                    "SELECT ARRAY[COUNT(id)::TEXT] FROM letter WHERE id = $1; --PARAM=0:1";
                break;
            case "del-note":
                textInfoEl.innerHTML += "- Parameter 1: Note ID";
                textAreaEl.value = "DELETE FROM blackboard_references WHERE base_note_id = $1 OR reference_id = $1; --PARAM=0:1\n"+
                    "DELETE FROM blackboard_note WHERE id = $1; --PARAM=0:1\n"+
                    "SELECT ARRAY[COUNT(id)::TEXT] FROM blackboard_note WHERE id = $1; --PARAM=0:1";
                break;
            case "del-doc-tag":
                textInfoEl.innerHTML += "Parameter 1: Tag ID";
                //noinspection ALL
                textAreaEl.value = "UPDATE document SET extra_info = subquery.element\n"+
                    "FROM  (SELECT document.id AS doc_id, extra_info #- ARRAY['tags', (position - 1)::TEXT] AS element\n"+
                    "       FROM document, jsonb_array_elements(extra_info->'tags') WITH ORDINALITY arr(elem, position)\n"+
                    "       WHERE elem->>'id' = $1) as subquery\n"+
                    "WHERE id = subquery.doc_id RETURNING ARRAY[id]; --PARAM=0:1";
                break;
            case "":
                textInfoEl.innerHTML += "";
                textAreaEl.value = "\n"+
                    "\n"+
                    "";
                break;
        }
    }
    updateQuery = updateOnQuery;
}