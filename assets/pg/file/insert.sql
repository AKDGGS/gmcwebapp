INSERT INTO file (filename, description, size, mimetype)
VALUES ($1, $2, $3, $4) RETURNING file_id
