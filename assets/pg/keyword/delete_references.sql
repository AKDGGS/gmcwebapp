UPDATE inventory SET
	keywords = array_remove(keywords, $1::keyword)
WHERE keywords @> ARRAY[$1::keyword]
