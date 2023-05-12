SELECT f.file_id AS "ID", f.filename AS "Name", f.mimetype as "Type",
	f.size::numeric AS "Size"
FROM borehole_file AS bf
JOIN file AS f
	ON f.file_id = bf.file_id
WHERE bf.borehole_id = $1
