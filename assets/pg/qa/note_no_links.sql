SELECT note_id AS "Note ID",
	username AS "User",
	TO_CHAR(note_date, 'MM/dd/yyyy') AS "Date"
FROM note
WHERE note_id NOT IN (
	SELECT DISTINCT note_id
	FROM borehole_note
) AND note_id NOT IN (
	SELECT DISTINCT note_id
	FROM inventory_note
) AND note_id NOT IN (
	SELECT DISTINCT note_id
	FROM outcrop_note
) AND note_id NOT IN (
	SELECT DISTINCT note_id
	FROM publication_note
) AND note_id NOT IN (
	SELECT DISTINCT note_id
	FROM shotline_note
) AND note_id NOT IN (
	SELECT DISTINCT note_id
	FROM well_note
) AND note_id NOT IN (
	SELECT DISTINCT note_id
	FROM sample_note
)
