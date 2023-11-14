INSERT INTO inventory (
	barcode, remark, container_id
) VALUES ($1, $2, $3) RETURNING inventory_id
