SELECT f.file_id AS "ID", f.filename AS "Name", f.mimetype as "Type",
	f.size::numeric AS "Size"
FROM prospect_file AS pf
JOIN file AS f
	ON f.file_id = pf.file_id
WHERE pf.prospect_id = $1
