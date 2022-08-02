SELECT ARRAY_AGG(unnest ORDER BY unnest)
FROM UNNEST(ENUM_RANGE(null::keyword)::TEXT[])
