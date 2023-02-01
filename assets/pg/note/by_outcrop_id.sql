SELECT n.note_id AS "ID", n.note AS "Note", n.note_date AS "Date", n.is_public AS "IsPublic", n.username AS "Username",
	nt.note_type_id AS "Type", nt.name AS "Name", nt.description AS "Description"
FROM note AS n
JOIN note_type AS nt
	ON nt.note_type_id = n.note_type_id
JOIN outcrop_note AS ocn
	ON ocn.note_id = n.note_id
WHERE ocn.outcrop_id = $1
	AND n.active
ORDER BY note_date DESC
