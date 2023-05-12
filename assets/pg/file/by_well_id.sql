SELECT f.file_id AS ID, f.filename AS Name, f.mimetype as "Type",
	f.size::numeric AS "Size"
FROM well_file AS wf
JOIN file AS f
	ON f.file_id = wf.file_id
WHERE wf.well_id = $1
