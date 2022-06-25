SELECT iw.inventory_id AS id,
	iv.interval_bottom::TEXT AS "Interval Bottom",
	GREATEST(we.measured_depth, we.vertical_depth)::TEXT AS "Well Depth"
FROM inventory_well AS iw
JOIN inventory AS iv
	ON iv.inventory_id = iw.inventory_id
JOIN well AS we
	ON we.well_id = iw.well_id
WHERE iv.active
	AND iv.interval_bottom IS NOT NULL
	AND COALESCE(we.measured_depth, we.vertical_depth) IS NOT NULL
	AND GREATEST(we.measured_depth, we.vertical_depth) < iv.interval_bottom 
ORDER BY (
	iv.interval_bottom - GREATEST(we.measured_depth, we.vertical_depth)
) DESC
