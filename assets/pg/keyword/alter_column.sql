ALTER TABLE inventory ALTER COLUMN keywords TYPE keyword[] USING (
	(keywords::text[])::keyword[]
)
