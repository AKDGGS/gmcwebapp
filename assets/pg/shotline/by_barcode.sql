SELECT sl.shotline_id AS id,
	sl.name,
	sl.alt_names,
	sl.year,
	sl.remark,
	array_to_string(ARRAY_AGG(COALESCE(
		sp.shotpoint_number::TEXT, 'Unknown'
	) ORDER BY shotpoint_number), ', ') AS numbers
FROM shotpoint AS sp
LEFT OUTER JOIN shotline AS sl
	ON sl.shotline_id = sp.shotline_id
WHERE sp.shotline_id IN (
  SELECT DISTINCT sp.shotline_id
  FROM inventory_shotpoint AS isp
  JOIN shotpoint as sp
    ON sp.shotpoint_id = isp.shotpoint_id
  JOIN inventory AS i
	ON i.inventory_id = isp.inventory_id
  WHERE barcode =ANY($1)
)
GROUP BY sl.shotline_id;
