SELECT ARRAY(
    SELECT barcode
    FROM container
    WHERE barcode = ANY($1::varchar[])
    UNION
    SELECT barcode
    FROM inventory
    WHERE barcode = ANY($1::varchar[])
) @> $1::varchar[]
