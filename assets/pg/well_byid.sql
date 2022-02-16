SELECT well_id, name, alt_names, well_number, api_number,
  is_onshore, is_federal, spud_date, completion_date,
  measured_depth, vertical_depth, elevation, elevation_kb,
  permit_status, completion_status
FROM well
WHERE well_id = $1
