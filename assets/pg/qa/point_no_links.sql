SELECT point_id AS "Point ID",
	description AS "Description"
FROM point
WHERE point_id NOT IN (
	SELECT point_id
	FROM borehole_point
) AND point_id NOT IN (
	SELECT point_id
	FROM outcrop_point
) AND point_id NOT IN (
	SELECT point_id
	FROM shotpoint_point
) AND point_id NOT IN (
	SELECT point_id
	FROM well_point
)
