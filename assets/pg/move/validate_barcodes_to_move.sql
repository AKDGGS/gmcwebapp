SELECT ARRAY(
	SELECT COALESCE(barcode, alt_barcode) FROM container
	WHERE barcode = ANY($1::varchar[])
		OR alt_barcode = ANY($1::varchar[])

	UNION

	SELECT COALESCE(barcode, alt_barcode) FROM inventory
	WHERE barcode = ANY($1::varchar[])
		OR alt_barcode = ANY($1::varchar[])
) @> $1::varchar[];
