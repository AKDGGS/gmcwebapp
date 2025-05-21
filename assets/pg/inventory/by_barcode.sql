SELECT
	iv.inventory_id AS id,
	jsonb_build_object(
		'id', cl.collection_id,
		'name', cl.name,
		'description', cl.description
	) AS collection,
	jsonb_build_object(
		'id', co.container_id,
		'name', co.name,
		'path_cache', co.path_cache
	) AS container,
	jsonb_build_object(
		'id', cd.core_diameter_id,
		'name', cd.name,
		'core_diameter', cd.core_diameter,
		'unit', COALESCE(cd.unit::text, 'ft')
	) AS core_diameter,
	iv.dggs_sample_id AS sample_id,
	iv.sample_number,
	iv.sample_number_prefix,
	iv.alt_sample_number,
	iv.published_sample_number,
	iv.published_number_has_suffix,
	iv.barcode,
	iv.alt_barcode,
	iv.state_number,
	iv.box_number,
	iv.set_number,
	iv.split_number,
	iv.slide_number,
	iv.slip_number,
	iv.lab_number,
	iv.lab_report_id,
	iv.map_number,
	iv.description,
	iv.remark,
	iv.tray,
	iv.interval_top,
	iv.interval_bottom,
	iv.keywords::text[],
	COALESCE(iv.interval_unit::text, 'ft') AS interval_unit,
	iv.core_number,
	iv.weight,
	iv.weight_unit::text,
	iv.sample_frequency,
	iv.recovery,
	iv.can_publish,
	iv.radiation_msvh,
	iv.received_date,
	iv.entered_date,
	iv.modified_date,
	iv.modified_user,
	iv.active,
	(
		SELECT
		jsonb_agg(
			jsonb_strip_nulls(
				jsonb_build_object(
					'id', b.borehole_id,
					'name', b.name,
					'prospect',
						jsonb_build_object(
						'id', b.prospect_id,
						'name', p.name,
						'ardf', p.ardf_number
					)
				)
			)
		)
		FROM inventory_borehole AS ib
		JOIN borehole AS b ON b.borehole_id = ib.borehole_id
		LEFT JOIN prospect AS p ON p.prospect_id = b.prospect_id
		WHERE ib.inventory_id = iv.inventory_id
		) AS boreholes,
	(
		SELECT jsonb_agg(
			jsonb_build_object(
				'id', o.outcrop_id,
				'name', o.name,
				'number', o.outcrop_number,
				'year', o.year
			)
		)
		FROM inventory_outcrop AS io
		JOIN outcrop AS o ON o.outcrop_id = io.outcrop_id
		WHERE io.inventory_id = iv.inventory_id
	) AS outcrops,
	(
		SELECT jsonb_agg(
			jsonb_strip_nulls(
				jsonb_build_object(
				'id', sq.shotpoint_id,
				'number', sq.shotpoint_number,
				'shotline',
					jsonb_build_object(
					'id', sq.shotline_id,
					'name', sq.name,
					'alt_names', sq.alt_names,
					'year', sq.year,
					'remark', sq.remark
					)
				)
			) ORDER BY sq.shotline_id, sq.shotpoint_number
		) AS shotpoints
		FROM (
			SELECT isp.inventory_id, sp.shotpoint_id, sp.shotpoint_number,
				sl.shotline_id, sl.name, sl.alt_names, sl.year, sl.remark
			FROM inventory_shotpoint AS isp
			LEFT OUTER JOIN shotpoint AS sp ON sp.shotpoint_id = isp.shotpoint_id
			LEFT OUTER JOIN shotline AS sl ON sl.shotline_id = sp.shotline_id
			WHERE isp.inventory_id = iv.inventory_id
			ORDER BY sl.shotline_id, sp.shotpoint_number
		) AS sq
	) AS shotpoints,
	(
		SELECT jsonb_agg(
			jsonb_strip_nulls(
				jsonb_build_object(
					'id', wq.well_id,
					'name', wq.name,
					'alt_names', wq.alt_names,
					'number', wq.well_number,
					'api_number', wq.api_number,
					'onshore', wq.is_onshore,
					'federal', wq.is_federal,
					'spud_date', to_char(wq.spud_date, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
					'completion date', to_char(wq.completion_date, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
					'measured_depth', wq.measured_depth,
					'vertical_depth', wq.vertical_depth,
					'elevation', wq.elevation,
					'elevation_kb', wq.elevation_kb,
					'permit_status', wq.permit_status,
					'completion_status', wq.completion_status,
					'permit_number', wq.permit_number,
					'unit', COALESCE(wq.unit::text, 'ft'),
					'organizations', (
						SELECT json_agg(
							jsonb_build_object(
								'id', org.organization_id,
								'name', org.name,
								'type', json_build_object('name', org.ot_name),
								'remark', org.remark,
								'is_current', org.is_current
							)
						)
						FROM (
							SELECT o.organization_id, o.name, ot.name AS ot_name, o.remark,
							 	wo.is_current, wo.well_id
							FROM organization AS o
							JOIN organization_type AS ot ON o.organization_type_id = ot.organization_type_id
							JOIN well_operator AS wo ON o.organization_id = wo.organization_id
							WHERE wo.well_id = wq.well_id
							ORDER BY wo.is_current DESC, o.name
						) AS org
					)
				)
			)
		) AS wells
		FROM (
			SELECT iw.inventory_id, w.well_id, w.name, w.alt_names, w.well_number,
				w.api_number, w.is_onshore, w.is_federal, w.spud_date, w.completion_date,
				w.measured_depth, w.vertical_depth, w.elevation, w.elevation_kb,
				w.permit_status, w.completion_status, w.permit_number, w.unit
			FROM inventory_well AS iw
			JOIN well AS w ON w.well_id = iw.well_id
			WHERE iw.inventory_id = iv.inventory_id
		) AS wq
	) AS wells,
	(
		SELECT jsonb_agg(
			jsonb_strip_nulls(
				jsonb_build_object(
					'id', iq.inventory_quality_id,
					'date', to_char(iq.check_date, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
					'remark', iq.remark,
					'username', iq.username,
					'issues', iq.issues::text[]
				)
			)
		)
		FROM inventory_quality AS iq
		WHERE inventory_id = iv.inventory_id
		GROUP BY iq.inventory_quality_id, iq.check_date
		ORDER BY iq.check_date DESC
		LIMIT 1
	) AS qualities
	FROM inventory AS iv
	LEFT OUTER JOIN collection AS cl
		ON cl.collection_id = iv.collection_id
	LEFT OUTER JOIN container AS co
		ON co.container_id = iv.container_id
	LEFT OUTER JOIN core_diameter AS cd
		ON cd.core_diameter_id = iv.core_diameter_id
	WHERE iv.active AND
	(
		iv.barcode = $1
		OR iv.barcode = ('GMC-' || $1 )
		OR iv.alt_barcode = $1
		OR iv.container_id IN (
			WITH RECURSIVE r AS (
				SELECT container_id
				FROM container WHERE barcode = $1 OR alt_barcode = $1

				UNION ALL

				SELECT co.container_id
				FROM r
				JOIN container AS co
				ON r.container_id = co.parent_container_id
			) SELECT container_id FROM r
		)
	)
	GROUP BY iv.inventory_id, cl.collection_id, co.container_id, cd.core_diameter_id
	LIMIT 100
