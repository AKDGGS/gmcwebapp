CREATE INDEX borehole_prospect_id_idx ON borehole(prospect_id);

CREATE INDEX borehole_point_borehole_id_idx ON borehole_point(borehole_id);

CREATE INDEX borehole_point_point_id_idx ON borehole_point(point_id);

CREATE INDEX collection_name_idx ON collection(name);

CREATE INDEX container_barcode_idx ON container(barcode);

CREATE INDEX container_container_type_id_idx ON container(container_type_id);

CREATE INDEX container_name_idx ON container(name);

CREATE INDEX container_alt_barcode_idx ON container(alt_barcode);

CREATE INDEX container_active_idx ON container(active);

CREATE INDEX container_parent_container_id_idx
	ON container(parent_container_id);

CREATE INDEX container_material_name_idx ON container_material(name);

CREATE INDEX container_type_name_idx ON container_type(name);

CREATE INDEX core_diameter_name_idx ON core_diameter(name);

CREATE INDEX core_diameter_core_diameter_idx ON core_diameter(core_diameter);

CREATE INDEX inventory_active_idx ON inventory(active);

CREATE INDEX inventory_parent_id_idx ON inventory(parent_id);

CREATE INDEX inventory_dimension_id_idx ON inventory(dimension_id);

CREATE INDEX inventory_collection_id_idx ON inventory(collection_id);

CREATE INDEX inventory_project_id_idx ON inventory(project_id);

CREATE INDEX inventory_barcode_idx ON inventory(barcode);

CREATE INDEX inventory_alt_barcode_idx ON inventory(alt_barcode);

CREATE INDEX inventory_coalesce_barcode_idx
	ON inventory(coalesce(barcode, alt_barcode));

CREATE INDEX inventory_container_id_idx ON inventory(container_id);

CREATE INDEX inventory_container_material_id_idx
	ON inventory(container_material_id);

CREATE INDEX inventory_container_log_inventory_id_idx
	ON inventory_container_log(inventory_id);

CREATE INDEX inventory_container_log_destination_idx
	ON inventory_container_log(destination);

CREATE INDEX inventory_container_log_log_date_idx
	ON inventory_container_log(log_date DESC);

CREATE INDEX container_log_container_id_idx
	ON container_log(container_id);

CREATE INDEX container_log_destination_idx
	ON container_log(destination);

CREATE INDEX container_log_log_date_idx
	ON container_log(log_date DESC);

CREATE INDEX inventory_borehole_borehole_id_idx
	ON inventory_borehole(borehole_id);

CREATE INDEX inventory_borehole_inventory_id_idx
	ON inventory_borehole(inventory_id);

CREATE INDEX inventory_outcrop_inventory_id_idx
	ON inventory_outcrop(inventory_id);

CREATE INDEX inventory_outcrop_outcrop_id_idx ON inventory_outcrop(outcrop_id);

CREATE INDEX inventory_publication_publication_id_idx
	ON inventory_publication(publication_id);

CREATE INDEX inventory_publication_inventory_id_idx
	ON inventory_publication(inventory_id);

CREATE INDEX inventory_well_well_id_idx ON inventory_well(well_id);

CREATE INDEX inventory_well_inventory_id_idx ON inventory_well(inventory_id);

CREATE INDEX outcrop_point_outcrop_id_idx ON outcrop_point(outcrop_id);

CREATE INDEX outcrop_point_point_id_idx ON outcrop_point(point_id);

CREATE INDEX outcrop_place_place_id_idx ON outcrop_place(place_id);

CREATE INDEX outcrop_place_outcrop_id_idx ON outcrop_place(outcrop_id);

CREATE INDEX organization_name_idx ON organization(name);

CREATE INDEX organization_abbr_idx ON organization(abbr);

CREATE INDEX organization_organization_type_id_idx
	ON organization(organization_type_id);

CREATE INDEX publication_citation_id_idx ON publication(citation_id);

CREATE INDEX place_name_idx ON place(name);

CREATE INDEX place_type_idx ON place(type);

CREATE INDEX quadrangle_name_idx ON quadrangle(name);

CREATE INDEX mining_district_name_idx ON mining_district(name);

CREATE INDEX mining_district_lower_name_idx ON mining_district(LOWER(name));

CREATE INDEX quadrangle_scale_idx ON quadrangle(scale);

CREATE INDEX note_active_idx ON note(active);

CREATE INDEX note_note_date_idx ON note(note_date);

CREATE INDEX note_note_type_id_idx ON note(note_type_id);

CREATE INDEX note_type_name_idx ON note_type(name);

CREATE INDEX dimension_name_idx ON dimension(name);

CREATE INDEX inventory_quality_check_date_idx ON inventory_quality(check_date);

CREATE INDEX prospect_ardf_number_idx ON prospect(ardf_number);

CREATE INDEX prospect_lower_ardf_number_idx ON prospect(LOWER(ardf_number));

CREATE INDEX well_api_number_idx ON well(api_number);

CREATE INDEX well_name_idx ON well(name);

CREATE INDEX well_well_number_idx ON well(well_number);

CREATE INDEX well_api_number_int_idx ON well(CAST (api_number AS bigint));

CREATE INDEX well_point_point_id_idx ON well_point(point_id);

CREATE INDEX well_point_well_id_idx ON well_point(well_id);

CREATE INDEX well_place_place_id_idx ON well_place(place_id);

CREATE INDEX well_place_well_id_idx ON well_place(well_id);

CREATE INDEX shotpoint_shotline_id_idx ON shotpoint(shotline_id);

CREATE INDEX shotpoint_shotpoint_number_idx ON shotpoint(shotpoint_number);

CREATE INDEX inventory_shotpoint_inventory_id_idx
	ON inventory_shotpoint(inventory_id);

CREATE INDEX inventory_shotpoint_shotpoint_id_idx
	ON inventory_shotpoint(shotpoint_id);

CREATE INDEX audit_barcode_idx ON audit(barcode);

CREATE INDEX plss_geog_idx ON plss USING GIST(geog);

CREATE INDEX point_geog_idx ON point USING GIST(geog);

-- B-tree to optimize against IS NULL/IS NOT NULL
CREATE INDEX point_geog_btree_idx ON point(geog);

CREATE INDEX place_geog_idx ON place USING GIST(geog);

-- B-tree to optimize against IS NULL/IS NOT NULL
CREATE INDEX place_geog_btree_idx ON place(geog);

CREATE INDEX region_geog_idx ON region USING GIST(geog);

CREATE INDEX quadrangle_geog_idx ON quadrangle USING GIST(geog);

CREATE INDEX energy_district_geog_idx ON energy_district USING GIST(geog);

CREATE INDEX mining_district_geog_idx ON mining_district USING GIST(geog);

CREATE INDEX gmc_region_geog_idx  ON gmc_region USING GIST(geog);

CREATE INDEX utm_geog_idx ON utm USING GIST(geog);

CREATE INDEX ON inventory USING gin(keywords);

-- Materialized view indexes
CREATE INDEX inventory_geog_inventory_id_idx ON inventory_geog(inventory_id);

CREATE INDEX inventory_geog_geog_idx ON inventory_geog USING GIST(geog);
