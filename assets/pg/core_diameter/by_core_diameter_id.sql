SELECT core_diameter_id AS id,
	core_diameter,
	name,
	unit
FROM core_diameter
WHERE core_diameter_id = $1
