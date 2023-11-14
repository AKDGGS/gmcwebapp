INSERT INTO container (
	barcode, name, remark, container_type_id
)
VALUES (
	$1, $2, $3, (
		SELECT container_type_id
		FROM container_type
		WHERE name = 'unknown'
		LIMIT 1
	)
) RETURNING container_type_id;
