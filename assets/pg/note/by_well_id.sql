SELECT n.note_id AS id, n.note, n.note_date AS date, n.is_public AS public,
n.username, nt.note_type_id as typeID, nt.name , nt.description
FROM note AS n
JOIN note_type AS nt
	ON nt.note_type_id = n.note_type_id
JOIN well_note AS wn
	ON wn.note_id = n.note_id
WHERE wn.well_id = $1
	AND n.active
ORDER BY note_date DESC
