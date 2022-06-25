SELECT file_id AS "File ID",
	filename AS "Filename"
FROM file
WHERE file_id NOT IN (
	SELECT file_id
	FROM borehole_file
) AND file_id NOT IN (
	SELECT file_id
	FROM inventory_file
) AND file_id NOT IN (
	SELECT file_id
	FROM outcrop_file
) AND file_id NOT IN (
	SELECT file_id
	FROM container_file
) AND file_id NOT IN (
	SELECT file_id
	FROM prospect_file
) AND file_id NOT IN (
	SELECT file_id
	FROM sample_file
) AND file_id NOT IN (
	SELECT file_id
	FROM well_file
)
