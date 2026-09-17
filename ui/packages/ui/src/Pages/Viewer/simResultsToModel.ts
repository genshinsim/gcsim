import type {
	Enemy,
	model,
	Shields,
	SimResults,
	Statistics,
	TargetBucketStats,
} from "@gcsim/types";

export function simResultsToModel(data: SimResults): model.SimulationResult {
	return {
		schema_version: data.schema_version,
		sim_version: data.sim_version,
		modified: data.modified,
		build_date: data.build_date,
		sample_seed: data.sample_seed,
		config_file: data.config_file,
		simulator_settings: data.simulator_settings as
			| model.SimulatorSettings
			| undefined,
		energy_settings: data.energy_settings as model.EnergySettings | undefined,
		initial_character: data.initial_character,
		character_details: data.character_details,
		target_details: data.target_details?.map(enemyToModel),
		player_position: data.player_position,
		incomplete_characters: data.incomplete_characters,
		statistics:
			data.statistics != null ? statisticsToModel(data.statistics) : undefined,
		mode: data.mode as model.SimMode | undefined,
		key_type: data.key_type,
	};
}

function enemyToModel(enemy: Enemy): model.Enemy {
	return {
		level: enemy.level,
		hp: enemy.hp,
		resist: enemy.resist,
		position: enemy.position,
		particle_drop_threshold: enemy.particle_drop_threshold,
		particle_drop_count: enemy.particle_drop_count,
		particle_element:
			enemy.particle_element != null
				? String(enemy.particle_element)
				: undefined,
		name: enemy.name,
		modified: enemy.modified,
	};
}

function statisticsToModel(stats: Statistics): model.SimulationStatistics {
	return {
		min_seed: stats.min_seed,
		max_seed: stats.max_seed,
		p25_seed: stats.p25_seed,
		p50_seed: stats.p50_seed,
		p75_seed: stats.p75_seed,
		iterations: stats.iterations,
		duration: stats.duration,
		dps: stats.dps,
		rps: stats.rps,
		eps: stats.eps,
		hps: stats.hps,
		shp: stats.shp,
		warnings: stats.warnings,
		failed_actions: stats.failed_actions,
		element_dps: stats.element_dps,
		target_dps: numberKeyed(stats.target_dps),
		character_dps: stats.character_dps,
		dps_by_element: stats.dps_by_element,
		dps_by_target: stats.dps_by_target?.map((byTarget) => ({
			targets: numberKeyed(byTarget.targets),
		})),
		source_dps: stats.source_dps,
		source_damage_instances: stats.source_damage_instances,
		damage_buckets: stats.damage_buckets,
		cumu_damage_contrib: stats.cumu_damage_contrib,
		cumu_damage:
			stats.cumu_damage != null
				? targetBucketStatsToModel(stats.cumu_damage)
				: undefined,
		shields: shieldsToModel(stats.shields),
		field_time: stats.field_time,
		total_source_energy: stats.total_source_energy,
		source_reactions: stats.source_reactions,
		character_actions: stats.character_actions,
		target_aura_uptime: stats.target_aura_uptime,
		end_stats: stats.end_stats,
	};
}

function targetBucketStatsToModel(
	bucketStats: TargetBucketStats,
): model.TargetBucketStats {
	return {
		bucket_size: bucketStats.bucket_size,
		targets: numberKeyed(bucketStats.targets),
	};
}

function shieldsToModel(
	shields: Shields | undefined,
): { [key: string]: model.ShieldInfo } | undefined {
	if (shields == null) {
		return undefined;
	}
	const out: { [key: string]: model.ShieldInfo } = {};
	for (const [key, info] of Object.entries(shields)) {
		out[key] = {
			hp: info.hp as unknown as
				| { [key: string]: model.DescriptiveStats }
				| undefined,
			uptime: info.uptime,
		};
	}
	return out;
}

function numberKeyed<T>(
	m: { [key: string]: T } | undefined,
): { [key: number]: T } | undefined {
	if (m == null) {
		return undefined;
	}
	const out: { [key: number]: T } = {};
	for (const [key, value] of Object.entries(m)) {
		out[Number(key)] = value;
	}
	return out;
}
