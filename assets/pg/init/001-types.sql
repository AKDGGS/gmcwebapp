CREATE TYPE urltype AS ENUM (
	'blm report', 'bom report', 'dggs report',
	'usgs report', 'well history', 'well log',
	'geologic data', 'internal report', 'location',
	'gmc data report', 'commercial analysis'
);

CREATE TYPE units AS ENUM (
	'µm', 'g', 'kg', 'lbs',
	'cm', 'm', 'ft', 'in'
);

CREATE TYPE keyword AS ENUM (
	'autogenerated', 'energy', 'engineering', 'geophysical', 'mineral',
	'outreach', 'supplies', 'survey', 'media', 'processed', 'raw',
	'residues', 'artifact', 'core', 'cuttings', 'fluid', 'grab',
	'outcrop', 'print', 'specimen', 'undetermined type',
	'apatite fission track', 'archive', 'argon argon', 'blank',
	'chemostratigraphy', 'clay', 'conodont', 'detrital zircon',
	'diatom', 'foraminifera', 'geochemistry', 'geochronology',
	'grain and pebble', 'grain-size', 'heavy mineral',
	'high pressure methane adsorption', 'kerogen', 'lead lead',
	'lithology', 'major oxide', 'megafossil',
	'mercury injection capillary pressure', 'microfossil',
	'nannoplankton', 'organic geochemistry', 'ostracod',
	'paleomagnetic', 'palynology', 'pan concentrate', 'petrology',
	'porosity and permeability', 'potassium argon', 'powder',
	'radiocarbon', 'radiolarian', 'scanning electron microscope',
	'siliceous fossil', 'spiral concentrate', 'spiral tail',
	'strain analysis', 'tephrochronology', 'thin section',
	'total organic carbon', 'trace element', 'undetermined analysis',
	'vitrinite', 'whole rock', 'xray diffraction', 'xray fluorescence',
	'zircon', 'auger', 'billet', 'bulk pieces', 'bulk reject', 'butt',
	'canned', 'center cut', 'center split', 'chips', 'coal', 'crude oil',
	'crushed', 'debris flow', 'dredge', 'gravel', 'hand sample', 'lava',
	'lava flow', 'map', 'organic material', 'periodical', 'photo', 'plug',
	'pulp', 'pyroclastic flow', 'quality control', 'quarter round', 'report',
	'rocker', 'rubble', 'sand', 'sediment', 'seep oil', 'shothole',
	'sidewall', 'silt', 'slab', 'slide', 'soil', 'stream sediment',
	'tailing', 'tephra fall', 'trench', 'trimming', 'undetermined form',
	'unwashed', 'vial', 'video', 'volcanic', 'washed', 'water',
	'well log', 'whole', 'biomarker', 'clast', 'coal absorption',
	'detrital mica', 'display', 'fluid inclusions', 'holotype',
	'igneous', 'metamorphic', 'mineral species', 'pillbox',
	'polished', 'recent', 'reference', 'sedimentary',
	'skeleton', 'stained', 'teaching', 'thermochronology',
	'uranium lead', 'waxed', 'index', 'dart', 'preserved',
	'sieve fraction', 'foamed'
);

CREATE TYPE issue AS ENUM (
	'needs_inventory', 'unsorted', 'radiation_risk',
	'material_damaged', 'box_damaged', 'missing',
	'needs_metadata', 'barcode_missing',
	'label_obscured', 'insufficient_material'
);

CREATE TYPE container_materials AS ENUM (
	'cardboard', 'cardboard/plastic', 'cloth fiber', 'glass',
	'metal', 'metal/wood',  'paper', 'plastic', 'wood'
);
