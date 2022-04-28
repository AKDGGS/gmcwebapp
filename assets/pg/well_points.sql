SELECT
	w.name, w.well_id, ST_AsText(ST_SnapToGrid(geog::geometry, 0.00001)) AS geog
FROM
	well AS w
JOIN (
	SELECT DISTINCT ON (well_id)
		well_id, geog
	FROM((
		SELECT well_id, geog
		FROM well_point AS wp
		JOIN point AS p
			ON p.point_id = wp.point_id
			WHERE geog IS NOT NULL
		) UNION ALL (
		SELECT well_id, geog
		FROM well_place AS wp
		JOIN place AS p
			ON p.place_id = wp.place_id
		WHERE p.geog IS NOT NULL
	)) AS q
)AS wg ON wg.well_id = w.well_id
JOIN (
	SELECT DISTINCT ON (well_id)
		well_id
	FROM inventory_well AS iw
	JOIN inventory AS i ON i.inventory_id = iw.inventory_id
	WHERE i.active = true AND i.can_publish = true
) AS wi ON wi.well_id = w.well_id
