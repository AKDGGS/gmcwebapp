WITH ids AS (
	SELECT note_id
	FROM borehole_note

	UNION

	SELECT note_id
	FROM inventory_note

	UNION

	SELECT note_id
	FROM outcrop_note

	UNION

	SELECT note_id
	FROM publication_note

	UNION

	SELECT note_id
	FROM shotline_note

	UNION

	SELECT note_id
	FROM well_note

	UNION

	SELECT note_id
	FROM sample_note
)
SELECT n.note_id AS "Note ID",
	username AS "User",
	TO_CHAR(note_date, 'MM/dd/yyyy') AS "Date"
FROM note AS n
LEFT JOIN ids ON ids.note_id = n.note_id
WHERE ids.note_id IS NULL
