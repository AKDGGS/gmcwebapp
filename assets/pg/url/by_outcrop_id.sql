SELECT u.url_id AS "ID", u.url, u.description, u.url_type::text AS "Type"
From url AS u
JOIN outcrop_url AS ou
	ON ou.url_id = u.url_id
WHERE ou.outcrop_id = $1
