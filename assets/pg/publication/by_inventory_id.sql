SELECT p.publication_id AS id,
	p.title,
	p.description,
	p.year,
	p.publication_type AS type,
	p.publication_number AS publicationNumber,
	p.publication_series AS publicationSeries,
	p.can_publish AS canPublish
	FROM inventory_publication AS ip
	JOIN publication AS p
		ON p.publication_id = ip.publication_id
	WHERE ip.inventory_id = $1
