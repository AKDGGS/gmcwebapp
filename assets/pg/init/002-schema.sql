CREATE TABLE token (
	token_id SERIAL PRIMARY KEY,
	token VARCHAR(128) NOT NULL UNIQUE,
	description VARCHAR(255)
);


CREATE TABLE note_type (
	note_type_id SERIAL PRIMARY KEY,
	name VARCHAR(50),
	description VARCHAR(255)
);


CREATE TABLE note (
	note_id SERIAL PRIMARY KEY,
	note_type_id INT REFERENCES note_type(note_type_id) NOT NULL,
	note TEXT NOT NULL,
	note_date DATE NOT NULL DEFAULT NOW(),
	is_public BOOLEAN NOT NULL DEFAULT true,
	username VARCHAR(25) NOT NULL,
	active BOOLEAN NOT NULL DEFAULT true
);


CREATE TABLE quadrangle (
	quadrangle_id SERIAL PRIMARY KEY,
	name VARCHAR(30) NOT NULL,
	alt_name VARCHAR(30) NULL,
	abbr VARCHAR(5) NULL,
	alt_abbr VARCHAR(5) NULL,
	scale INT NOT NULL,
	geog GEOGRAPHY(MultiPolygon) NULL
);


CREATE TABLE energy_district (
	energy_district_id SERIAL PRIMARY KEY,
	name VARCHAR(30) NOT NULL,
	geog GEOGRAPHY(MultiPolygon) NULL
);


CREATE TABLE place (
	place_id SERIAL PRIMARY KEY,
	name VARCHAR(150) NOT NULL,
	type VARCHAR(20) NOT NULL,
	geog GEOGRAPHY(Point) NULL
);


CREATE TABLE region (
	region_id SERIAL PRIMARY KEY,
	name VARCHAR(150) NOT NULL,
	type VARCHAR(20) NOT NULL,
	geog GEOGRAPHY(MultiPolygon) NULL
);


CREATE TABLE url (
	url_id SERIAL PRIMARY KEY,
	url_type urltype NOT NULL,
	description VARCHAR(255) NULL,
	url TEXT
);


CREATE TABLE file_type (
	file_type_id SERIAL PRIMARY KEY,
	name VARCHAR(50)
);


CREATE TABLE file (
	file_id SERIAL PRIMARY KEY,
	file_type_id INT REFERENCES file_type(file_type_id) NULL,
	description VARCHAR(255) NULL,
	mimetype VARCHAR(255) NOT NULL DEFAULT 'application/octet-stream',
	size INT NOT NULL,
	filename VARCHAR(255) NOT NULL,
	content BYTEA NOT NULL,
	content_md5 VARCHAR(32) UNIQUE NOT NULL
);


CREATE TABLE dimension (
	dimension_id SERIAL PRIMARY KEY,
	unit units NOT NULL,

	name VARCHAR(50) NOT NULL,
	height NUMERIC(10,2) NOT NULL,
	width NUMERIC(10,2) NOT NULL,
	depth NUMERIC(10,2) NOT NULL,
	remark TEXT NULL
);


CREATE TABLE utm_type (
	utm_type_id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL
);


CREATE TABLE utm (
	utm_id SERIAL PRIMARY KEY,
	utm_type_id INT REFERENCES utm_type(utm_type_id) NULL,
	unit units NULL,
	description VARCHAR(255) NULL,
	zone VARCHAR(3) NULL,
	easting INT NULL,
	northing INT NULL,
	srid INT NULL,
	geog GEOGRAPHY(MultiPolygon) NULL
);


CREATE TABLE mining_district (
	mining_district_id SERIAL PRIMARY KEY,
	name VARCHAR(35) NULL,
	geog GEOGRAPHY(MultiPolygon) NULL
);


CREATE TABLE gmc_region (
	gmc_region_id SERIAL PRIMARY KEY,
	name VARCHAR(35) NULL,
	geog GEOGRAPHY(MultiPolygon) NULL
);


CREATE TABLE plss (
	plss_id SERIAL PRIMARY KEY,
	meridian VARCHAR(4) NOT NULL,
	township INT NULL,
	township_dir VARCHAR(2) NULL,
	range INT NULL,
	range_dir VARCHAR(2) NULL,
	section INT NULL,
	quadrant VARCHAR(6) NULL,
	geog GEOGRAPHY(MultiPolygon) NULL
);


CREATE TABLE point_type (
	point_type_id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL
);


CREATE TABLE point (
	point_id SERIAL PRIMARY KEY,
	point_type_id INT REFERENCES point_type(point_type_id) NULL,
	description VARCHAR(255) NULL,
	geog GEOGRAPHY(Point) NOT NULL
);


CREATE TABLE organization_type (
	organization_type_id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL
);


CREATE TABLE organization (
	organization_id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	abbr VARCHAR(25) NULL,
	organization_type_id INT REFERENCES organization_type(organization_type_id) NOT NULL,
	remark TEXT NULL
);


CREATE TABLE person (
	person_id SERIAL PRIMARY KEY,
	first VARCHAR(100) NULL,
	middle VARCHAR(100) NULL,
	last VARCHAR(100) NOT NULL,
	suffix VARCHAR(50) NULL,

	phone VARCHAR(25) NULL,
	email VARCHAR(255) NULL,

	preferred_id INT REFERENCES person(person_id) NULL
);


CREATE TABLE person_organization (
	person_id INT REFERENCES person(person_id),
	organization_id INT REFERENCES organization(organization_id),
	log_date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	PRIMARY KEY(person_id, organization_id, log_date)
);


CREATE TABLE publication (
	publication_id SERIAL PRIMARY KEY,
	citation_id INT NULL,
	title VARCHAR(512) NOT NULL,
	description VARCHAR(1024) NULL,
	year INT NULL,
	publication_type VARCHAR(50) NULL,
	publication_number VARCHAR(50) NULL,
	publication_series VARCHAR(50) NULL,
	can_publish BOOLEAN NOT NULL DEFAULT false
);


CREATE TABLE publication_quadrangle (
	publication_id INT REFERENCES publication(publication_id) NOT NULL,
	quadrangle_id INT REFERENCES quadrangle(quadrangle_id) NOT NULL,
	PRIMARY KEY(publication_id, quadrangle_id)
);


CREATE TABLE publication_url (
	publication_id INT REFERENCES publication(publication_id) NOT NULL,
	url_id INT REFERENCES url(url_id) NOT NULL,
	PRIMARY KEY(publication_id, url_id)
);


CREATE TABLE publication_person (
	publication_id INT REFERENCES publication(publication_id),
	person_id INT REFERENCES person(person_id),
	PRIMARY KEY(publication_id, person_id)
);


CREATE TABLE publication_organization (
	publication_id INT REFERENCES publication(publication_id),
	organization_id INT REFERENCES organization(organization_id),
	PRIMARY KEY(publication_id, organization_id)
);


CREATE TABLE publication_note (
	publication_id INT REFERENCES publication(publication_id) NOT NULL,
	note_id INT REFERENCES note(note_id) NOT NULL,
	PRIMARY KEY(publication_id, note_id)
);


CREATE TABLE core_diameter (
	core_diameter_id SERIAL PRIMARY KEY,
	core_diameter NUMERIC(10,2) NOT NULL,
	name VARCHAR(100) NULL,
	unit units NULL
);


CREATE TABLE collection (
	collection_id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	description VARCHAR(255) NULL,
	organization_id INT REFERENCES organization(organization_id) NULL
);


CREATE TABLE project (
	project_id SERIAL PRIMARY KEY,
	dggs_project_id BIGINT NULL,
	organization_id INT REFERENCES organization(organization_id) NULL,
	name VARCHAR(100) NOT NULL,
	description TEXT NULL,
	remark TEXT NULL,
	start_date DATE NULL,
	end_date DATE NULL
);


CREATE TABLE stratigraphy_type (
	stratigraphy_type_id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL
);


CREATE TABLE stratigraphy (
	stratigraphy_id SERIAL PRIMARY KEY,
	parent_id INT REFERENCES stratigraphy(stratigraphy_id) NULL,
	stratigraphy_type_id INT REFERENCES stratigraphy_type(stratigraphy_type_id) NOT NULL,
	name VARCHAR(50) NOT NULL,
	alt_names VARCHAR(1024) NULL,
	remark TEXT NULL
);


CREATE TABLE well (
	well_id SERIAL PRIMARY KEY,

	name VARCHAR(255) NOT NULL,
	alt_names VARCHAR(1024) NULL,

	well_number VARCHAR(50) NULL,
	api_number VARCHAR(14) NULL,
	is_onshore BOOLEAN NOT NULL DEFAULT true,
	is_federal BOOLEAN NOT NULL DEFAULT false,

	spud_date DATE NULL,
	completion_date DATE NULL,

	measured_depth NUMERIC(10, 2) NULL,
	vertical_depth NUMERIC(10, 2) NULL,
	elevation NUMERIC(10, 2) NULL,
	elevation_kb NUMERIC(10, 2) NULL,
	unit units NULL,
	permit_number INT NULL,

	permit_status VARCHAR(6) NULL,
	completion_status VARCHAR(6) NULL,

	entered_date DATE NULL DEFAULT NOW(),
	modified_date TIMESTAMP WITHOUT TIME ZONE NULL,
	modified_user VARCHAR(64) NULL,
	stash JSONB NULL
);


CREATE TABLE well_place (
	well_id INT REFERENCES well(well_id) NOT NULL,
	place_id INT REFERENCES place(place_id) NOT NULL,
	PRIMARY KEY(well_id, place_id)
);


CREATE TABLE well_region (
	well_id INT REFERENCES well(well_id) NOT NULL,
	region_id INT REFERENCES region(region_id) NOT NULL,
	PRIMARY KEY(well_id, region_id)
);


CREATE TABLE well_point (
	well_id INT REFERENCES well(well_id) NOT NULL,
	point_id INT REFERENCES point(point_id) NOT NULL,
	PRIMARY KEY(well_id, point_id)
);


CREATE TABLE well_operator (
	well_id INT REFERENCES well(well_id) NOT NULL,
	organization_id INT REFERENCES organization(organization_id) NOT NULL,
	is_current BOOLEAN NOT NULL DEFAULT true,
	PRIMARY KEY(well_id, organization_id)
);


CREATE TABLE well_url (
	well_id INT REFERENCES well(well_id) NOT NULL,
	url_id INT REFERENCES url(url_id) NOT NULL,
	PRIMARY KEY(well_id, url_id)
);


CREATE TABLE well_note (
	well_id INT REFERENCES well(well_id) NOT NULL,
	note_id INT REFERENCES note(note_id) NOT NULL,
	PRIMARY KEY(well_id, note_id)
);


CREATE TABLE well_file (
	well_id INT REFERENCES well(well_id),
	file_id INT REFERENCES file(file_id),
	PRIMARY KEY(well_id, file_id)
);


CREATE TABLE well_stratigraphy (
	well_stratigraphy_id SERIAL PRIMARY KEY,
	stratigraphy_id INT REFERENCES stratigraphy(stratigraphy_id) NOT NULL,
	well_id INT REFERENCES well(well_id) NOT NULL,
	measured_depth_top NUMERIC(8, 2) NULL,
	measured_depth_bottom NUMERIC(8, 2) NULL,
	vertical_depth_top NUMERIC(8, 2) NULL,
	vertical_depth_bottom NUMERIC(8, 2) NULL,
	unit units NULL,
	published_date DATE NULL,
	remark TEXT NULL
);


CREATE TABLE well_stratigraphy_person (
	well_stratigraphy_id INT REFERENCES well_stratigraphy(well_stratigraphy_id),
	person_id INT REFERENCES person(person_id),
	PRIMARY KEY(well_stratigraphy_id, person_id)
);


CREATE TABLE well_stratigraphy_organization (
	well_stratigraphy_id INT REFERENCES well_stratigraphy(well_stratigraphy_id),
	organization_id INT REFERENCES organization(organization_id),
	PRIMARY KEY(well_stratigraphy_id, organization_id)
);


CREATE TABLE outcrop (
	outcrop_id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	outcrop_number VARCHAR(50) NULL, -- datatype?
	is_onshore BOOLEAN NOT NULL DEFAULT true,
	year SMALLINT NULL,
	entered_date DATE NULL DEFAULT NOW(),
	modified_date TIMESTAMP WITHOUT TIME ZONE NULL,
	modified_user VARCHAR(64) NULL,
	stash JSONB NULL
);


CREATE TABLE outcrop_quadrangle (
	outcrop_id INT REFERENCES outcrop(outcrop_id) NOT NULL,
	quadrangle_id INT REFERENCES quadrangle(quadrangle_id) NOT NULL,
	PRIMARY KEY(outcrop_id, quadrangle_id)
);


CREATE TABLE outcrop_point (
	outcrop_id INT REFERENCES outcrop(outcrop_id) NOT NULL,
	point_id INT REFERENCES point(point_id) NOT NULL,
	PRIMARY KEY(outcrop_id, point_id)
);


CREATE TABLE outcrop_stratigraphy (
	outcrop_id INT REFERENCES outcrop(outcrop_id) NOT NULL,
	stratigraphy_id INT REFERENCES stratigraphy(stratigraphy_id) NOT NULL,
	PRIMARY KEY(outcrop_id, stratigraphy_id)
);


CREATE TABLE outcrop_note (
	outcrop_id INT REFERENCES outcrop(outcrop_id) NOT NULL,
	note_id INT REFERENCES note(note_id) NOT NULL,
	PRIMARY KEY(outcrop_id, note_id)
);


CREATE TABLE outcrop_file (
	outcrop_id INT REFERENCES outcrop(outcrop_id),
	file_id INT REFERENCES file(file_id),
	PRIMARY KEY(outcrop_id, file_id)
);


CREATE TABLE outcrop_organization (
	outcrop_id INT REFERENCES outcrop(outcrop_id) NOT NULL,
	organization_id INT REFERENCES organization(organization_id) NOT NULL,
	PRIMARY KEY(outcrop_id, organization_id)
);


CREATE TABLE outcrop_place (
	outcrop_id INT REFERENCES outcrop(outcrop_id) NOT NULL,
	place_id INT REFERENCES place(place_id) NOT NULL,
	PRIMARY KEY(outcrop_id, place_id)
);


CREATE TABLE outcrop_region (
	outcrop_id INT REFERENCES outcrop(outcrop_id) NOT NULL,
	region_id INT REFERENCES region(region_id) NOT NULL,
	PRIMARY KEY(outcrop_id, region_id)
);


CREATE TABLE outcrop_plss (
	outcrop_id INT REFERENCES outcrop(outcrop_id) NOT NULL,
	plss_id INT REFERENCES plss(plss_id) NOT NULL,
	PRIMARY KEY(outcrop_id, plss_id)
);


CREATE TABLE prospect (
	prospect_id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	alt_names VARCHAR(1024) NULL,
	ardf_number VARCHAR(25) NULL
);


CREATE TABLE prospect_file (
	prospect_id INT REFERENCES prospect(prospect_id),
	file_id INT REFERENCES file(file_id),
	PRIMARY KEY(prospect_id, file_id)
);


CREATE TABLE borehole (
	borehole_id SERIAL PRIMARY KEY,
	prospect_id INT REFERENCES prospect(prospect_id) NULL,

	name VARCHAR(50) NOT NULL,
	alt_names VARCHAR(1024) NULL,

	is_onshore BOOLEAN NOT NULL DEFAULT true,
	completion_date DATE NULL,
	measured_depth NUMERIC(8, 2) NULL,
	measured_depth_unit units NULL,
	elevation NUMERIC(8, 2) NULL,
	elevation_unit units NULL,
	entered_date DATE NULL DEFAULT NOW(),
	modified_date TIMESTAMP WITHOUT TIME ZONE NULL,
	modified_user VARCHAR(64) NULL,
	stash JSONB NULL
);


CREATE TABLE borehole_point (
	borehole_id INT REFERENCES borehole(borehole_id) NOT NULL,
	point_id INT REFERENCES point(point_id) NOT NULL,
	PRIMARY KEY(borehole_id, point_id)
);


CREATE TABLE borehole_url (
	borehole_id INT REFERENCES borehole(borehole_id) NOT NULL,
	url_id INT REFERENCES url(url_id) NOT NULL,
	PRIMARY KEY(borehole_id, url_id)
);


CREATE TABLE borehole_note (
	borehole_id INT REFERENCES borehole(borehole_id) NOT NULL,
	note_id INT REFERENCES note(note_id) NOT NULL,
	PRIMARY KEY(borehole_id, note_id)
);


CREATE TABLE borehole_organization (
	borehole_id INT REFERENCES borehole(borehole_id) NOT NULL,
	organization_id INT REFERENCES organization(organization_id) NOT NULL,
	PRIMARY KEY(borehole_id, organization_id)
);


CREATE TABLE borehole_file (
	borehole_id INT REFERENCES borehole(borehole_id),
	file_id INT REFERENCES file(file_id),
	PRIMARY KEY(borehole_id, file_id)
);


CREATE TABLE shotline (
	shotline_id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	alt_names VARCHAR(1024) NULL,
	year SMALLINT NULL,
	remark TEXT NULL
);


CREATE TABLE shotline_url (
	shotline_id INT REFERENCES shotline(shotline_id) NOT NULL,
	url_id INT REFERENCES url(url_id) NOT NULL,
	PRIMARY KEY(shotline_id, url_id)
);


CREATE TABLE shotline_note (
	shotline_id INT REFERENCES shotline(shotline_id) NOT NULL,
	note_id INT REFERENCES note(note_id) NOT NULL,
	PRIMARY KEY(shotline_id, note_id)
);


CREATE TABLE shotpoint (
	shotpoint_id SERIAL PRIMARY KEY,
	shotline_id INT REFERENCES shotline(shotline_id) NOT NULL,
	shotpoint_number NUMERIC(8,2) NULL,
	stash JSONB NULL
);


CREATE TABLE shotpoint_point (
	shotpoint_id INT REFERENCES shotpoint(shotpoint_id) NOT NULL,
	point_id INT REFERENCES point(point_id) NOT NULL,
	PRIMARY KEY(shotpoint_id, point_id)
);


CREATE TABLE container_material (
	container_material_id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL
);


CREATE TABLE container_type (
	container_type_id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL
);


CREATE TABLE container (
	container_id SERIAL PRIMARY KEY,
	parent_container_id INT REFERENCES container(container_id) NULL,
	container_type_id INT REFERENCES container_type(container_type_id) NOT NULL,
	container_material_id INT REFERENCES container_material(container_material_id) NULL,

	name VARCHAR(50) NOT NULL,
	description VARCHAR(255) NULL,

	remark TEXT NULL,

	dimension_id INT REFERENCES dimension(dimension_id) NULL,

	barcode VARCHAR(25) NULL,
	alt_barcode VARCHAR(25) NULL,

	-- Trigger-induced cache of the container path
	path_cache VARCHAR(255) NULL,

	active BOOLEAN NOT NULL DEFAULT true
);


CREATE TABLE container_file (
	container_id INT REFERENCES container(container_id),
	file_id INT REFERENCES file(file_id),
	PRIMARY KEY(container_id, file_id)
);


CREATE TABLE container_log (
	container_log_id SERIAL PRIMARY KEY,
	container_id INT REFERENCES container(container_id) NOT NULL,
	destination TEXT NULL,
	log_date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);


CREATE TABLE inventory (
	inventory_id SERIAL PRIMARY KEY,
	parent_id INT REFERENCES inventory(inventory_id) NULL,
	collection_id INT REFERENCES collection(collection_id) NULL,
	project_id INT REFERENCES project(project_id) NULL,
	dimension_id INT REFERENCES dimension(dimension_id) NULL,
	container_id INT REFERENCES container(container_id) NULL,
	container_material_id INT REFERENCES container_material(container_material_id) NULL,

	keywords keyword[] NULL CHECK (cardinality(keywords) <> 0),

	dggs_sample_id BIGINT NULL,
	sample_number VARCHAR(50) NULL,
	sample_number_prefix VARCHAR(25) NULL,
	alt_sample_number VARCHAR(25) NULL,
	published_sample_number VARCHAR(50) NULL,
	published_number_has_suffix BOOLEAN NOT NULL DEFAULT false,
	barcode VARCHAR(25) NULL,
	alt_barcode VARCHAR(25) NULL,
	state_number VARCHAR(50) NULL,
	box_number VARCHAR(50) NULL,
	set_number VARCHAR(50) NULL,
	split_number VARCHAR(10) NULL,
	slide_number VARCHAR(10) NULL,
	slip_number INT NULL,
	lab_number VARCHAR(100),
	map_number VARCHAR(25), -- Stores BLM map number

	description TEXT NULL,
	remark TEXT NULL,

	tray SMALLINT NULL DEFAULT 1,

	lab_report_id VARCHAR(100) NULL,

	interval_top NUMERIC(8,2) NULL,
	interval_bottom NUMERIC(8,2) NULL,
	interval_unit units NULL,

	core_number VARCHAR(25) NULL,
	core_diameter_id INT REFERENCES core_diameter(core_diameter_id) NULL,

	weight NUMERIC(10, 2) NULL,
	weight_unit units NULL,

	-- This is a hack  This field should eventually be replaced
	-- with an integer field called "interval_frequency"
	sample_frequency VARCHAR(25) NULL,
	recovery VARCHAR(25) NULL,
	can_publish BOOLEAN NOT NULL DEFAULT false,
	-- Radiation level in milli-sievert per hour
	radiation_msvh NUMERIC(10, 4) NULL,
	received_date DATE NULL DEFAULT NULL,
	entered_date DATE NULL DEFAULT NOW(),
	modified_date TIMESTAMP WITHOUT TIME ZONE NULL,
	modified_user VARCHAR(64) NULL,

	stash JSONB NULL,
	active BOOLEAN NOT NULL DEFAULT true
);


CREATE TABLE inventory_container_log (
	inventory_container_log_id SERIAL PRIMARY KEY,
	inventory_id INT REFERENCES inventory(inventory_id) NOT NULL,
	destination TEXT NULL,
	log_date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);


CREATE TABLE inventory_collector (
	inventory_id INT REFERENCES inventory(inventory_id) NOT NULL,
	collector_id INT REFERENCES person(person_id) NOT NULL,
	PRIMARY KEY(inventory_id, collector_id)
);


CREATE TABLE inventory_url (
	inventory_id INT REFERENCES inventory(inventory_id) NOT NULL,
	url_id INT REFERENCES url(url_id) NOT NULL,
	PRIMARY KEY(inventory_id, url_id)
);


CREATE TABLE inventory_publication (
	inventory_id INT REFERENCES inventory(inventory_id),
	publication_id INT REFERENCES publication(publication_id),
	PRIMARY KEY(inventory_id, publication_id)
);


CREATE TABLE inventory_file (
	inventory_id INT REFERENCES inventory(inventory_id),
	file_id INT REFERENCES file(file_id),
	PRIMARY KEY(inventory_id, file_id)
);


CREATE TABLE inventory_note (
	inventory_id INT REFERENCES inventory(inventory_id) NOT NULL,
	note_id INT REFERENCES note(note_id) NOT NULL,
	PRIMARY KEY(inventory_id, note_id)
);


CREATE TABLE inventory_borehole (
	inventory_id INT REFERENCES inventory(inventory_id),
	borehole_id INT REFERENCES borehole(borehole_id),
	PRIMARY KEY(inventory_id, borehole_id)
);


CREATE TABLE inventory_well (
	inventory_id INT REFERENCES inventory(inventory_id),
	well_id INT REFERENCES well(well_id),
	PRIMARY KEY(inventory_id, well_id)
);


CREATE TABLE inventory_outcrop (
	inventory_id INT REFERENCES inventory(inventory_id),
	outcrop_id INT REFERENCES outcrop(outcrop_id),
	PRIMARY KEY(inventory_id, outcrop_id)
);


CREATE TABLE inventory_shotpoint (
	inventory_id INT REFERENCES inventory(inventory_id),
	shotpoint_id INT REFERENCES shotpoint(shotpoint_id),
	PRIMARY KEY(inventory_id, shotpoint_id)
);


CREATE TABLE inventory_quality (
	inventory_quality_id SERIAL PRIMARY KEY,
	inventory_id INT REFERENCES inventory(inventory_id) NOT NULL,
	check_date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	remark TEXT NULL,
	issues issue[] NULL CHECK (cardinality(issues) <> 0),
	username VARCHAR(25) NOT NULL
);


CREATE TABLE audit_group (
	audit_group_id SERIAL PRIMARY KEY,
	remark TEXT NULL,
	create_date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);


CREATE TABLE audit (
	audit_id SERIAL PRIMARY KEY,
	audit_group_id INT REFERENCES audit_group(audit_group_id) ON DELETE CASCADE NOT NULL,
	barcode VARCHAR(25) NOT NULL
);
