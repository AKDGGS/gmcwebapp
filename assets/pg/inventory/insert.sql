INSERT INTO inventory (
	barcode, remark, container_id, keywords
) VALUES ($1, $2, $3, $4::keyword[]) RETURNING inventory_id
