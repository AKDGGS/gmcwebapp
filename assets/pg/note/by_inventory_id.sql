SELECT n.note_id, n.note, n.note_date, n.is_public, n.username,
	nt.note_type_id, nt.name, nt.description
FROM note AS n
JOIN note_type AS nt
	ON nt.note_type_id = n.note_type_id
JOIN inventory_note AS iv
	ON iv.note_id = n.note_id
WHERE iv.inventory_id = $1
	AND n.active
ORDER BY note_date DESC
