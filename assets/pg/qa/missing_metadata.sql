SELECT iv.inventory_id,
	COALESCE(iv.barcode, iv.alt_barcode, '') || ' ' || COALESCE(co.path_cache, '') AS description
FROM inventory AS iv
LEFT OUTER JOIN container AS co
	ON co.container_id = iv.inventory_id
LEFT OUTER JOIN inventory_well AS ivw
	ON ivw.inventory_id = iv.inventory_id
LEFT OUTER JOIN inventory_borehole AS ivb
	ON ivb.inventory_id = iv.inventory_id
LEFT OUTER JOIN inventory_outcrop AS ivo
	ON ivo.inventory_id = iv.inventory_id
LEFT OUTER JOIN inventory_shotpoint AS ivs
	ON ivs.inventory_id = iv.inventory_id
LEFT OUTER JOIN inventory_publication AS ivp
	ON ivp.inventory_id = iv.inventory_id
WHERE iv.active
	AND ivw.inventory_id IS NULL
	AND ivb.inventory_id IS NULL
	AND ivo.inventory_id IS NULL
	AND ivs.inventory_id IS NULL
	AND ivp.inventory_id IS NULL
