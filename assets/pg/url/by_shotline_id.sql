SELECT u.url_id, u.url, u.description, u.url_type::text
From url AS u
JOIN shotline_url AS su
	ON su.url_id = u.url_id
WHERE su.shotline_id = $1
