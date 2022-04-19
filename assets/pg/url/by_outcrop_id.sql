SELECT u.url_id, u.url, u.description, u.url_type::text
From url AS u
JOIN outcrop_url AS ou
	ON ou.url_id = u.url_id
WHERE ou.outcrop_id = $1
