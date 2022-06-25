SELECT outcrop_id AS "Outcrop ID",
	name AS "Outcrop Name"
FROM outcrop
WHERE outcrop_id NOT IN (
	SELECT outcrop_id
	FROM outcrop_point
) AND outcrop_id NOT IN (
	SELECT op.outcrop_id
	FROM outcrop_place AS op
	JOIN place AS pl
		ON pl.place_id = op.place_id
	WHERE pl.geog IS NOT NULL
) AND outcrop_id NOT IN (
	SELECT op.outcrop_id
	FROM outcrop_plss AS op
	JOIN plss AS pl
		ON pl.plss_id = op.plss_id
	WHERE pl.geog IS NOT NULL
) AND outcrop_id NOT IN (
	SELECT oq.outcrop_id
	FROM outcrop_quadrangle AS oq
	JOIN quadrangle AS qu
		ON qu.quadrangle_id = oq.quadrangle_id
	WHERE qu.geog IS NOT NULL
) AND outcrop_id NOT IN (
	SELECT oe.outcrop_id
	FROM outcrop_region AS oe
	JOIN region AS re
		ON re.region_id = oe.region_id
	WHERE re.geog IS NOT NULL
)
ORDER BY outcrop_id ASC
