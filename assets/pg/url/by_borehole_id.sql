SELECT
	u.url_id AS id,
	u.url,
	u.description,
	u.url_type::text AS type
From url AS u
JOIN borehole_url AS bu
	ON bu.url_id = u.url_id
WHERE bu.borehole_id = $1
