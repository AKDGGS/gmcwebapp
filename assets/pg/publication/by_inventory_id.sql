SELECT p.publication_id AS "ID", p.title, p.description, p.year,
	p.publication_type AS "Type", p.publication_number AS "PublicationNumber",
	p.publication_series AS "PublicationSeries", p.can_publish AS "CanPublish"
	FROM inventory_publication AS ip
	JOIN publication AS p
		ON p.publication_id = ip.publication_id
	WHERE ip.inventory_id = $1
