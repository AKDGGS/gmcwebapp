INSERT INTO file (filename, description, size, mimetype, content, content_md5)
VALUES ($1, $2, $3, $4, ''::bytea, $5) RETURNING file_id
