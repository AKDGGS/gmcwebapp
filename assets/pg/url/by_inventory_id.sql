SELECT u.url_id AS "ID", u.url, u.description, u.url_type::text AS "Type"
From url AS u
JOIN inventory_url AS iu
	ON iu.url_id = u.url_id
WHERE iu.inventory_id = $1
