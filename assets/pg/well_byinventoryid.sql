SELECT w.well_id, w.name, w.alt_names,
  w.well_number, w.api_number,
  w.is_onshore, w.is_federal,
  w.spud_date, w.completion_date,
  w.measured_depth, w.vertical_depth,
  w.elevation, w.elevation_kb,
  w.permit_status, w.completion_status, w.permit_number,
  w.unit::text
  FROM well AS w
  JOIN inventory_well AS iw
    ON iw.well_id = w.well_id
  WHERE iw.inventory_id = $1
