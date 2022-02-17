SELECT well_id, name, alt_names, well_number, api_number,
  is_onshore, is_federal, permit_status, completion_status,
  spud_date,
  completion_date,
  measured_depth::double precision,
  vertical_depth::double precision,
  elevation::double precision,
  elevation_kb::double precision,
  COALESCE(unit::text, 'ft') AS unit
FROM well
WHERE well_id = $1
