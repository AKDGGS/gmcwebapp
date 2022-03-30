SELECT p.publication_id, p.title, p.description, p.year,
	p.publication_type, p.publication_number, p.publication_series,
	p.can_publish
	FROM inventory_publication AS ip
	JOIN publication AS p
		ON p.publication_id = ip.publication_id
	WHERE ip.inventory_id = $1
