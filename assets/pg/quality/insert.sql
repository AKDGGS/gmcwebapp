INSERT INTO inventory_quality (
	inventory_id, remark, username, issues
) VALUES (
	$1, $2, $3, $4::issue[]
)
