SELECT
	u.url_id AS id,
	u.url,
	u.description,
	u.url_type::text AS type
From url AS u
JOIN well_url AS wu
	ON wu.url_id = u.url_id
WHERE wu.well_id = $1
