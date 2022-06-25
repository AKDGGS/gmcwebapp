SELECT iv.inventory_id AS "Inventory ID",
	COALESCE(iv.barcode, iv.alt_barcode) AS "Barcode",
	co.name AS "Collection",
	ct.path_cache AS "Container"
FROM inventory AS iv
LEFT OUTER JOIN container AS ct
	ON ct.container_id = iv.inventory_id AND ct.active
LEFT OUTER JOIN collection AS co
	ON co.collection_id = iv.collection_id
WHERE iv.active
	AND iv.interval_top IS NOT NULL
	AND iv.interval_bottom IS NOT NULL
	AND iv.interval_bottom < iv.interval_top
