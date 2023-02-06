import * as $protobuf from "protobufjs";
import Long = require("long");
/** Namespace model. */
export namespace model {

    /** Properties of a DBEntry. */
    interface IDBEntry {

        /** DBEntry id */
        id?: (string|null);

        /** DBEntry key */
        key?: (string|null);

        /** DBEntry create_date */
        create_date?: (number|Long|null);

        /** DBEntry run_date */
        run_date?: (number|Long|null);

        /** DBEntry sim_duration */
        sim_duration?: (model.IDescriptiveStats|null);

        /** DBEntry config */
        config?: (string|null);

        /** DBEntry hash */
        hash?: (string|null);

        /** DBEntry mode */
        mode?: (model.SimMode|null);

        /** DBEntry total_damage */
        total_damage?: (model.IDescriptiveStats|null);

        /** DBEntry char_names */
        char_names?: (string[]|null);

        /** DBEntry target_count */
        target_count?: (number|null);

        /** DBEntry mean_dps_per_target */
        mean_dps_per_target?: (number|null);

        /** DBEntry team */
        team?: (model.ICharacter[]|null);

        /** DBEntry dps_by_target */
        dps_by_target?: ({ [k: string]: model.IDescriptiveStats }|null);

        /** DBEntry description */
        description?: (string|null);

        /** DBEntry accepted_tags */
        accepted_tags?: (string[]|null);

        /** DBEntry rejected_tags */
        rejected_tags?: (string[]|null);

        /** DBEntry is_db_valid */
        is_db_valid?: (boolean|null);
    }

    /** Represents a DBEntry. */
    class DBEntry implements IDBEntry {

        /**
         * Constructs a new DBEntry.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IDBEntry);

        /** DBEntry id. */
        public id: string;

        /** DBEntry key. */
        public key: string;

        /** DBEntry create_date. */
        public create_date: (number|Long);

        /** DBEntry run_date. */
        public run_date: (number|Long);

        /** DBEntry sim_duration. */
        public sim_duration?: (model.IDescriptiveStats|null);

        /** DBEntry config. */
        public config: string;

        /** DBEntry hash. */
        public hash: string;

        /** DBEntry mode. */
        public mode: model.SimMode;

        /** DBEntry total_damage. */
        public total_damage?: (model.IDescriptiveStats|null);

        /** DBEntry char_names. */
        public char_names: string[];

        /** DBEntry target_count. */
        public target_count: number;

        /** DBEntry mean_dps_per_target. */
        public mean_dps_per_target: number;

        /** DBEntry team. */
        public team: model.ICharacter[];

        /** DBEntry dps_by_target. */
        public dps_by_target: { [k: string]: model.IDescriptiveStats };

        /** DBEntry description. */
        public description: string;

        /** DBEntry accepted_tags. */
        public accepted_tags: string[];

        /** DBEntry rejected_tags. */
        public rejected_tags: string[];

        /** DBEntry is_db_valid. */
        public is_db_valid: boolean;

        /**
         * Gets the default type url for DBEntry
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a DBEntries. */
    interface IDBEntries {

        /** DBEntries data */
        data?: (model.IDBEntry[]|null);
    }

    /** Represents a DBEntries. */
    class DBEntries implements IDBEntries {

        /**
         * Constructs a new DBEntries.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IDBEntries);

        /** DBEntries data. */
        public data: model.IDBEntry[];

        /**
         * Gets the default type url for DBEntries
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a DBQueryOpt. */
    interface IDBQueryOpt {

        /** DBQueryOpt query */
        query?: (google.protobuf.IStruct|null);

        /** DBQueryOpt sort */
        sort?: (google.protobuf.IStruct|null);

        /** DBQueryOpt project */
        project?: (google.protobuf.IStruct|null);

        /** DBQueryOpt skip */
        skip?: (number|Long|null);

        /** DBQueryOpt limit */
        limit?: (number|Long|null);
    }

    /** Represents a DBQueryOpt. */
    class DBQueryOpt implements IDBQueryOpt {

        /**
         * Constructs a new DBQueryOpt.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IDBQueryOpt);

        /** DBQueryOpt query. */
        public query?: (google.protobuf.IStruct|null);

        /** DBQueryOpt sort. */
        public sort?: (google.protobuf.IStruct|null);

        /** DBQueryOpt project. */
        public project?: (google.protobuf.IStruct|null);

        /** DBQueryOpt skip. */
        public skip: (number|Long);

        /** DBQueryOpt limit. */
        public limit: (number|Long);

        /**
         * Gets the default type url for DBQueryOpt
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a Submission. */
    interface ISubmission {

        /** Submission id */
        id?: (string|null);

        /** Submission config */
        config?: (string|null);

        /** Submission submitter */
        submitter?: (string|null);

        /** Submission description */
        description?: (string|null);

        /** Submission preview */
        preview?: (string|null);

        /** Submission tags */
        tags?: (string[]|null);
    }

    /** Represents a Submission. */
    class Submission implements ISubmission {

        /**
         * Constructs a new Submission.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.ISubmission);

        /** Submission id. */
        public id: string;

        /** Submission config. */
        public config: string;

        /** Submission submitter. */
        public submitter: string;

        /** Submission description. */
        public description: string;

        /** Submission preview. */
        public preview: string;

        /** Submission tags. */
        public tags: string[];

        /**
         * Gets the default type url for Submission
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a ComputeWork. */
    interface IComputeWork {

        /** ComputeWork key */
        key?: (string|null);

        /** ComputeWork cfg */
        cfg?: (string|null);
    }

    /** Represents a ComputeWork. */
    class ComputeWork implements IComputeWork {

        /**
         * Constructs a new ComputeWork.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IComputeWork);

        /** ComputeWork key. */
        public key: string;

        /** ComputeWork cfg. */
        public cfg: string;

        /**
         * Gets the default type url for ComputeWork
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a Version. */
    interface IVersion {

        /** Version major */
        major?: (number|Long|null);

        /** Version minor */
        minor?: (number|Long|null);
    }

    /** Represents a Version. */
    class Version implements IVersion {

        /**
         * Constructs a new Version.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IVersion);

        /** Version major. */
        public major: (number|Long);

        /** Version minor. */
        public minor: (number|Long);

        /**
         * Gets the default type url for Version
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a SimulationResult. */
    interface ISimulationResult {

        /** SimulationResult schema_version */
        schema_version?: (model.IVersion|null);

        /** SimulationResult sim_version */
        sim_version?: (string|null);

        /** SimulationResult modified */
        modified?: (boolean|null);

        /** SimulationResult build_date */
        build_date?: (string|null);

        /** SimulationResult sample_seed */
        sample_seed?: (string|null);

        /** SimulationResult config */
        config?: (string|null);

        /** SimulationResult simulator_settings */
        simulator_settings?: (model.ISimulatorSettings|null);

        /** SimulationResult energy_settings */
        energy_settings?: (model.IEnergySettings|null);

        /** SimulationResult initial_character */
        initial_character?: (string|null);

        /** SimulationResult character_details */
        character_details?: (model.ICharacter[]|null);

        /** SimulationResult target_details */
        target_details?: (model.IEnemy[]|null);

        /** SimulationResult statistics */
        statistics?: (model.ISimulationStatistics|null);

        /** SimulationResult mode */
        mode?: (model.SimMode|null);
    }

    /** Represents a SimulationResult. */
    class SimulationResult implements ISimulationResult {

        /**
         * Constructs a new SimulationResult.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.ISimulationResult);

        /** SimulationResult schema_version. */
        public schema_version?: (model.IVersion|null);

        /** SimulationResult sim_version. */
        public sim_version: string;

        /** SimulationResult modified. */
        public modified: boolean;

        /** SimulationResult build_date. */
        public build_date: string;

        /** SimulationResult sample_seed. */
        public sample_seed: string;

        /** SimulationResult config. */
        public config: string;

        /** SimulationResult simulator_settings. */
        public simulator_settings?: (model.ISimulatorSettings|null);

        /** SimulationResult energy_settings. */
        public energy_settings?: (model.IEnergySettings|null);

        /** SimulationResult initial_character. */
        public initial_character: string;

        /** SimulationResult character_details. */
        public character_details: model.ICharacter[];

        /** SimulationResult target_details. */
        public target_details: model.IEnemy[];

        /** SimulationResult statistics. */
        public statistics?: (model.ISimulationStatistics|null);

        /** SimulationResult mode. */
        public mode: model.SimMode;

        /**
         * Gets the default type url for SimulationResult
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a SimulationStatistics. */
    interface ISimulationStatistics {

        /** SimulationStatistics min_seed */
        min_seed?: (string|null);

        /** SimulationStatistics max_seed */
        max_seed?: (string|null);

        /** SimulationStatistics p25_seed */
        p25_seed?: (string|null);

        /** SimulationStatistics p50_seed */
        p50_seed?: (string|null);

        /** SimulationStatistics p75_seed */
        p75_seed?: (string|null);

        /** SimulationStatistics runtime */
        runtime?: (number|null);

        /** SimulationStatistics iterations */
        iterations?: (number|Long|null);

        /** SimulationStatistics duration */
        duration?: (model.IOverviewStats|null);

        /** SimulationStatistics DPS */
        DPS?: (model.IOverviewStats|null);

        /** SimulationStatistics RPS */
        RPS?: (model.IOverviewStats|null);

        /** SimulationStatistics EPS */
        EPS?: (model.IOverviewStats|null);

        /** SimulationStatistics HPS */
        HPS?: (model.IOverviewStats|null);

        /** SimulationStatistics SHP */
        SHP?: (model.IOverviewStats|null);

        /** SimulationStatistics total_damage */
        total_damage?: (model.IOverviewStats|null);

        /** SimulationStatistics warnings */
        warnings?: (model.IWarnings|null);

        /** SimulationStatistics failed_actions */
        failed_actions?: (model.IFailedActions[]|null);

        /** SimulationStatistics element_dps */
        element_dps?: ({ [k: string]: model.IDescriptiveStats }|null);

        /** SimulationStatistics target_dps */
        target_dps?: ({ [k: string]: model.IDescriptiveStats }|null);

        /** SimulationStatistics character_dps */
        character_dps?: (model.IDescriptiveStats[]|null);

        /** SimulationStatistics breakdown_by_element_dps */
        breakdown_by_element_dps?: (model.IElementStats[]|null);

        /** SimulationStatistics breakdown_by_target_dps */
        breakdown_by_target_dps?: (model.ITargetStats[]|null);

        /** SimulationStatistics damage_buckets */
        damage_buckets?: (model.IBucketStats|null);

        /** SimulationStatistics cumulative_damage_contribution */
        cumulative_damage_contribution?: (model.ICharacterBucketStats|null);

        /** SimulationStatistics shields */
        shields?: ({ [k: string]: model.IShieldInfo }|null);
    }

    /** Represents a SimulationStatistics. */
    class SimulationStatistics implements ISimulationStatistics {

        /**
         * Constructs a new SimulationStatistics.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.ISimulationStatistics);

        /** SimulationStatistics min_seed. */
        public min_seed: string;

        /** SimulationStatistics max_seed. */
        public max_seed: string;

        /** SimulationStatistics p25_seed. */
        public p25_seed: string;

        /** SimulationStatistics p50_seed. */
        public p50_seed: string;

        /** SimulationStatistics p75_seed. */
        public p75_seed: string;

        /** SimulationStatistics runtime. */
        public runtime: number;

        /** SimulationStatistics iterations. */
        public iterations: (number|Long);

        /** SimulationStatistics duration. */
        public duration?: (model.IOverviewStats|null);

        /** SimulationStatistics DPS. */
        public DPS?: (model.IOverviewStats|null);

        /** SimulationStatistics RPS. */
        public RPS?: (model.IOverviewStats|null);

        /** SimulationStatistics EPS. */
        public EPS?: (model.IOverviewStats|null);

        /** SimulationStatistics HPS. */
        public HPS?: (model.IOverviewStats|null);

        /** SimulationStatistics SHP. */
        public SHP?: (model.IOverviewStats|null);

        /** SimulationStatistics total_damage. */
        public total_damage?: (model.IOverviewStats|null);

        /** SimulationStatistics warnings. */
        public warnings?: (model.IWarnings|null);

        /** SimulationStatistics failed_actions. */
        public failed_actions: model.IFailedActions[];

        /** SimulationStatistics element_dps. */
        public element_dps: { [k: string]: model.IDescriptiveStats };

        /** SimulationStatistics target_dps. */
        public target_dps: { [k: string]: model.IDescriptiveStats };

        /** SimulationStatistics character_dps. */
        public character_dps: model.IDescriptiveStats[];

        /** SimulationStatistics breakdown_by_element_dps. */
        public breakdown_by_element_dps: model.IElementStats[];

        /** SimulationStatistics breakdown_by_target_dps. */
        public breakdown_by_target_dps: model.ITargetStats[];

        /** SimulationStatistics damage_buckets. */
        public damage_buckets?: (model.IBucketStats|null);

        /** SimulationStatistics cumulative_damage_contribution. */
        public cumulative_damage_contribution?: (model.ICharacterBucketStats|null);

        /** SimulationStatistics shields. */
        public shields: { [k: string]: model.IShieldInfo };

        /**
         * Gets the default type url for SimulationStatistics
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a SignedSimulationStatistics. */
    interface ISignedSimulationStatistics {

        /** SignedSimulationStatistics stats */
        stats?: (model.ISimulationStatistics|null);

        /** SignedSimulationStatistics hash */
        hash?: (string|null);
    }

    /** Represents a SignedSimulationStatistics. */
    class SignedSimulationStatistics implements ISignedSimulationStatistics {

        /**
         * Constructs a new SignedSimulationStatistics.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.ISignedSimulationStatistics);

        /** SignedSimulationStatistics stats. */
        public stats?: (model.ISimulationStatistics|null);

        /** SignedSimulationStatistics hash. */
        public hash: string;

        /**
         * Gets the default type url for SignedSimulationStatistics
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an OverviewStats. */
    interface IOverviewStats {

        /** OverviewStats min */
        min?: (number|null);

        /** OverviewStats max */
        max?: (number|null);

        /** OverviewStats mean */
        mean?: (number|null);

        /** OverviewStats SD */
        SD?: (number|null);

        /** OverviewStats Q1 */
        Q1?: (number|null);

        /** OverviewStats Q2 */
        Q2?: (number|null);

        /** OverviewStats Q3 */
        Q3?: (number|null);

        /** OverviewStats hist */
        hist?: ((number|Long)[]|null);
    }

    /** Represents an OverviewStats. */
    class OverviewStats implements IOverviewStats {

        /**
         * Constructs a new OverviewStats.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IOverviewStats);

        /** OverviewStats min. */
        public min: number;

        /** OverviewStats max. */
        public max: number;

        /** OverviewStats mean. */
        public mean: number;

        /** OverviewStats SD. */
        public SD: number;

        /** OverviewStats Q1. */
        public Q1: number;

        /** OverviewStats Q2. */
        public Q2: number;

        /** OverviewStats Q3. */
        public Q3: number;

        /** OverviewStats hist. */
        public hist: (number|Long)[];

        /**
         * Gets the default type url for OverviewStats
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a DescriptiveStats. */
    interface IDescriptiveStats {

        /** DescriptiveStats min */
        min?: (number|null);

        /** DescriptiveStats max */
        max?: (number|null);

        /** DescriptiveStats mean */
        mean?: (number|null);

        /** DescriptiveStats SD */
        SD?: (number|null);
    }

    /** Represents a DescriptiveStats. */
    class DescriptiveStats implements IDescriptiveStats {

        /**
         * Constructs a new DescriptiveStats.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IDescriptiveStats);

        /** DescriptiveStats min. */
        public min: number;

        /** DescriptiveStats max. */
        public max: number;

        /** DescriptiveStats mean. */
        public mean: number;

        /** DescriptiveStats SD. */
        public SD: number;

        /**
         * Gets the default type url for DescriptiveStats
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an ElementStats. */
    interface IElementStats {

        /** ElementStats elements */
        elements?: ({ [k: string]: model.IDescriptiveStats }|null);
    }

    /** Represents an ElementStats. */
    class ElementStats implements IElementStats {

        /**
         * Constructs a new ElementStats.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IElementStats);

        /** ElementStats elements. */
        public elements: { [k: string]: model.IDescriptiveStats };

        /**
         * Gets the default type url for ElementStats
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a TargetStats. */
    interface ITargetStats {

        /** TargetStats targets */
        targets?: ({ [k: string]: model.IDescriptiveStats }|null);
    }

    /** Represents a TargetStats. */
    class TargetStats implements ITargetStats {

        /**
         * Constructs a new TargetStats.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.ITargetStats);

        /** TargetStats targets. */
        public targets: { [k: string]: model.IDescriptiveStats };

        /**
         * Gets the default type url for TargetStats
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a BucketStats. */
    interface IBucketStats {

        /** BucketStats bucket_size */
        bucket_size?: (number|Long|null);

        /** BucketStats buckets */
        buckets?: (model.IDescriptiveStats[]|null);
    }

    /** Represents a BucketStats. */
    class BucketStats implements IBucketStats {

        /**
         * Constructs a new BucketStats.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IBucketStats);

        /** BucketStats bucket_size. */
        public bucket_size: (number|Long);

        /** BucketStats buckets. */
        public buckets: model.IDescriptiveStats[];

        /**
         * Gets the default type url for BucketStats
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a CharacterBucketStats. */
    interface ICharacterBucketStats {

        /** CharacterBucketStats bucket_size */
        bucket_size?: (number|Long|null);

        /** CharacterBucketStats characters */
        characters?: (model.ICharacterBuckets[]|null);
    }

    /** Represents a CharacterBucketStats. */
    class CharacterBucketStats implements ICharacterBucketStats {

        /**
         * Constructs a new CharacterBucketStats.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.ICharacterBucketStats);

        /** CharacterBucketStats bucket_size. */
        public bucket_size: (number|Long);

        /** CharacterBucketStats characters. */
        public characters: model.ICharacterBuckets[];

        /**
         * Gets the default type url for CharacterBucketStats
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a CharacterBuckets. */
    interface ICharacterBuckets {

        /** CharacterBuckets buckets */
        buckets?: (model.IDescriptiveStats[]|null);
    }

    /** Represents a CharacterBuckets. */
    class CharacterBuckets implements ICharacterBuckets {

        /**
         * Constructs a new CharacterBuckets.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.ICharacterBuckets);

        /** CharacterBuckets buckets. */
        public buckets: model.IDescriptiveStats[];

        /**
         * Gets the default type url for CharacterBuckets
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a Warnings. */
    interface IWarnings {

        /** Warnings target_overlap */
        target_overlap?: (boolean|null);

        /** Warnings insufficient_energy */
        insufficient_energy?: (boolean|null);

        /** Warnings insufficient_stamina */
        insufficient_stamina?: (boolean|null);

        /** Warnings swap_cd */
        swap_cd?: (boolean|null);

        /** Warnings skill_cd */
        skill_cd?: (boolean|null);
    }

    /** Represents a Warnings. */
    class Warnings implements IWarnings {

        /**
         * Constructs a new Warnings.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IWarnings);

        /** Warnings target_overlap. */
        public target_overlap: boolean;

        /** Warnings insufficient_energy. */
        public insufficient_energy: boolean;

        /** Warnings insufficient_stamina. */
        public insufficient_stamina: boolean;

        /** Warnings swap_cd. */
        public swap_cd: boolean;

        /** Warnings skill_cd. */
        public skill_cd: boolean;

        /**
         * Gets the default type url for Warnings
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a FailedActions. */
    interface IFailedActions {

        /** FailedActions insufficient_energy */
        insufficient_energy?: (model.IDescriptiveStats|null);

        /** FailedActions insufficient_stamina */
        insufficient_stamina?: (model.IDescriptiveStats|null);

        /** FailedActions swap_cd */
        swap_cd?: (model.IDescriptiveStats|null);

        /** FailedActions skill_cd */
        skill_cd?: (model.IDescriptiveStats|null);
    }

    /** Represents a FailedActions. */
    class FailedActions implements IFailedActions {

        /**
         * Constructs a new FailedActions.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IFailedActions);

        /** FailedActions insufficient_energy. */
        public insufficient_energy?: (model.IDescriptiveStats|null);

        /** FailedActions insufficient_stamina. */
        public insufficient_stamina?: (model.IDescriptiveStats|null);

        /** FailedActions swap_cd. */
        public swap_cd?: (model.IDescriptiveStats|null);

        /** FailedActions skill_cd. */
        public skill_cd?: (model.IDescriptiveStats|null);

        /**
         * Gets the default type url for FailedActions
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a ShieldInfo. */
    interface IShieldInfo {

        /** ShieldInfo hp */
        hp?: ({ [k: string]: model.IDescriptiveStats }|null);

        /** ShieldInfo uptime */
        uptime?: (model.IDescriptiveStats|null);
    }

    /** Represents a ShieldInfo. */
    class ShieldInfo implements IShieldInfo {

        /**
         * Constructs a new ShieldInfo.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IShieldInfo);

        /** ShieldInfo hp. */
        public hp: { [k: string]: model.IDescriptiveStats };

        /** ShieldInfo uptime. */
        public uptime?: (model.IDescriptiveStats|null);

        /**
         * Gets the default type url for ShieldInfo
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a Sample. */
    interface ISample {

        /** Sample sim_version */
        sim_version?: (string|null);

        /** Sample build_date */
        build_date?: (string|null);

        /** Sample modified */
        modified?: (boolean|null);

        /** Sample config */
        config?: (string|null);

        /** Sample initial_character */
        initial_character?: (string|null);

        /** Sample character_details */
        character_details?: (model.ICharacter[]|null);

        /** Sample target_details */
        target_details?: (model.IEnemy[]|null);

        /** Sample seed */
        seed?: (string|null);

        /** Sample logs */
        logs?: (google.protobuf.IStruct[]|null);
    }

    /** Represents a Sample. */
    class Sample implements ISample {

        /**
         * Constructs a new Sample.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.ISample);

        /** Sample sim_version. */
        public sim_version: string;

        /** Sample build_date. */
        public build_date: string;

        /** Sample modified. */
        public modified: boolean;

        /** Sample config. */
        public config: string;

        /** Sample initial_character. */
        public initial_character: string;

        /** Sample character_details. */
        public character_details: model.ICharacter[];

        /** Sample target_details. */
        public target_details: model.IEnemy[];

        /** Sample seed. */
        public seed: string;

        /** Sample logs. */
        public logs: google.protobuf.IStruct[];

        /**
         * Gets the default type url for Sample
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a Character. */
    interface ICharacter {

        /** Character name */
        name?: (string|null);

        /** Character element */
        element?: (string|null);

        /** Character level */
        level?: (number|Long|null);

        /** Character max_level */
        max_level?: (number|Long|null);

        /** Character cons */
        cons?: (number|Long|null);

        /** Character weapon */
        weapon?: (model.IWeapon|null);

        /** Character talents */
        talents?: (model.ICharacterTalents|null);

        /** Character sets */
        sets?: ({ [k: string]: (number|Long) }|null);

        /** Character stats */
        stats?: (number[]|null);

        /** Character snapshot_stats */
        snapshot_stats?: (number[]|null);
    }

    /** Represents a Character. */
    class Character implements ICharacter {

        /**
         * Constructs a new Character.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.ICharacter);

        /** Character name. */
        public name: string;

        /** Character element. */
        public element: string;

        /** Character level. */
        public level: (number|Long);

        /** Character max_level. */
        public max_level: (number|Long);

        /** Character cons. */
        public cons: (number|Long);

        /** Character weapon. */
        public weapon?: (model.IWeapon|null);

        /** Character talents. */
        public talents?: (model.ICharacterTalents|null);

        /** Character sets. */
        public sets: { [k: string]: (number|Long) };

        /** Character stats. */
        public stats: number[];

        /** Character snapshot_stats. */
        public snapshot_stats: number[];

        /**
         * Gets the default type url for Character
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a CharacterTalents. */
    interface ICharacterTalents {

        /** CharacterTalents attack */
        attack?: (number|Long|null);

        /** CharacterTalents skill */
        skill?: (number|Long|null);

        /** CharacterTalents burst */
        burst?: (number|Long|null);
    }

    /** Represents a CharacterTalents. */
    class CharacterTalents implements ICharacterTalents {

        /**
         * Constructs a new CharacterTalents.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.ICharacterTalents);

        /** CharacterTalents attack. */
        public attack: (number|Long);

        /** CharacterTalents skill. */
        public skill: (number|Long);

        /** CharacterTalents burst. */
        public burst: (number|Long);

        /**
         * Gets the default type url for CharacterTalents
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a Weapon. */
    interface IWeapon {

        /** Weapon name */
        name?: (string|null);

        /** Weapon refine */
        refine?: (number|Long|null);

        /** Weapon level */
        level?: (number|Long|null);

        /** Weapon max_level */
        max_level?: (number|Long|null);
    }

    /** Represents a Weapon. */
    class Weapon implements IWeapon {

        /**
         * Constructs a new Weapon.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IWeapon);

        /** Weapon name. */
        public name: string;

        /** Weapon refine. */
        public refine: (number|Long);

        /** Weapon level. */
        public level: (number|Long);

        /** Weapon max_level. */
        public max_level: (number|Long);

        /**
         * Gets the default type url for Weapon
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an Enemy. */
    interface IEnemy {

        /** Enemy level */
        level?: (number|Long|null);

        /** Enemy HP */
        HP?: (number|null);

        /** Enemy resist */
        resist?: ({ [k: string]: number }|null);

        /** Enemy pos */
        pos?: (model.ICoord|null);

        /** Enemy particle_drop_threshold */
        particle_drop_threshold?: (number|null);

        /** Enemy particle_drop_count */
        particle_drop_count?: (number|null);

        /** Enemy particle_element */
        particle_element?: (string|null);
    }

    /** Represents an Enemy. */
    class Enemy implements IEnemy {

        /**
         * Constructs a new Enemy.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IEnemy);

        /** Enemy level. */
        public level: (number|Long);

        /** Enemy HP. */
        public HP: number;

        /** Enemy resist. */
        public resist: { [k: string]: number };

        /** Enemy pos. */
        public pos?: (model.ICoord|null);

        /** Enemy particle_drop_threshold. */
        public particle_drop_threshold: number;

        /** Enemy particle_drop_count. */
        public particle_drop_count: number;

        /** Enemy particle_element. */
        public particle_element: string;

        /**
         * Gets the default type url for Enemy
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a Coord. */
    interface ICoord {

        /** Coord x */
        x?: (number|null);

        /** Coord y */
        y?: (number|null);

        /** Coord r */
        r?: (number|null);
    }

    /** Represents a Coord. */
    class Coord implements ICoord {

        /**
         * Constructs a new Coord.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.ICoord);

        /** Coord x. */
        public x: number;

        /** Coord y. */
        public y: number;

        /** Coord r. */
        public r: number;

        /**
         * Gets the default type url for Coord
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a SimulatorSettings. */
    interface ISimulatorSettings {

        /** SimulatorSettings duration */
        duration?: (number|null);

        /** SimulatorSettings damage_mode */
        damage_mode?: (boolean|null);

        /** SimulatorSettings enable_hitlag */
        enable_hitlag?: (boolean|null);

        /** SimulatorSettings def_halt */
        def_halt?: (boolean|null);

        /** SimulatorSettings number_of_workers */
        number_of_workers?: (number|Long|null);

        /** SimulatorSettings iterations */
        iterations?: (number|Long|null);

        /** SimulatorSettings delays */
        delays?: (model.IDelays|null);
    }

    /** Represents a SimulatorSettings. */
    class SimulatorSettings implements ISimulatorSettings {

        /**
         * Constructs a new SimulatorSettings.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.ISimulatorSettings);

        /** SimulatorSettings duration. */
        public duration: number;

        /** SimulatorSettings damage_mode. */
        public damage_mode: boolean;

        /** SimulatorSettings enable_hitlag. */
        public enable_hitlag: boolean;

        /** SimulatorSettings def_halt. */
        public def_halt: boolean;

        /** SimulatorSettings number_of_workers. */
        public number_of_workers: (number|Long);

        /** SimulatorSettings iterations. */
        public iterations: (number|Long);

        /** SimulatorSettings delays. */
        public delays?: (model.IDelays|null);

        /**
         * Gets the default type url for SimulatorSettings
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a Delays. */
    interface IDelays {

        /** Delays skill */
        skill?: (number|Long|null);

        /** Delays burst */
        burst?: (number|Long|null);

        /** Delays attack */
        attack?: (number|Long|null);

        /** Delays charge */
        charge?: (number|Long|null);

        /** Delays aim */
        aim?: (number|Long|null);

        /** Delays dash */
        dash?: (number|Long|null);

        /** Delays jump */
        jump?: (number|Long|null);

        /** Delays swap */
        swap?: (number|Long|null);
    }

    /** Represents a Delays. */
    class Delays implements IDelays {

        /**
         * Constructs a new Delays.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IDelays);

        /** Delays skill. */
        public skill: (number|Long);

        /** Delays burst. */
        public burst: (number|Long);

        /** Delays attack. */
        public attack: (number|Long);

        /** Delays charge. */
        public charge: (number|Long);

        /** Delays aim. */
        public aim: (number|Long);

        /** Delays dash. */
        public dash: (number|Long);

        /** Delays jump. */
        public jump: (number|Long);

        /** Delays swap. */
        public swap: (number|Long);

        /**
         * Gets the default type url for Delays
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an EnergySettings. */
    interface IEnergySettings {

        /** EnergySettings active */
        active?: (boolean|null);

        /** EnergySettings once */
        once?: (boolean|null);

        /** EnergySettings start */
        start?: (number|Long|null);

        /** EnergySettings end */
        end?: (number|Long|null);

        /** EnergySettings amount */
        amount?: (number|Long|null);

        /** EnergySettings last_energy_drop */
        last_energy_drop?: (number|Long|null);
    }

    /** Represents an EnergySettings. */
    class EnergySettings implements IEnergySettings {

        /**
         * Constructs a new EnergySettings.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IEnergySettings);

        /** EnergySettings active. */
        public active: boolean;

        /** EnergySettings once. */
        public once: boolean;

        /** EnergySettings start. */
        public start: (number|Long);

        /** EnergySettings end. */
        public end: (number|Long);

        /** EnergySettings amount. */
        public amount: (number|Long);

        /** EnergySettings last_energy_drop. */
        public last_energy_drop: (number|Long);

        /**
         * Gets the default type url for EnergySettings
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an AvatarData. */
    interface IAvatarData {

        /** AvatarData rarity */
        rarity?: (number|Long|null);

        /** AvatarData body */
        body?: (model.BodyType|null);

        /** AvatarData region */
        region?: (model.ZoneType|null);

        /** AvatarData element */
        element?: (model.Element|null);

        /** AvatarData weapon_class */
        weapon_class?: (model.WeaponClass|null);

        /** AvatarData image_name */
        image_name?: (string|null);

        /** AvatarData base_stats */
        base_stats?: (model.IAvatarStatsData|null);
    }

    /** Represents an AvatarData. */
    class AvatarData implements IAvatarData {

        /**
         * Constructs a new AvatarData.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IAvatarData);

        /** AvatarData rarity. */
        public rarity: (number|Long);

        /** AvatarData body. */
        public body: model.BodyType;

        /** AvatarData region. */
        public region: model.ZoneType;

        /** AvatarData element. */
        public element: model.Element;

        /** AvatarData weapon_class. */
        public weapon_class: model.WeaponClass;

        /** AvatarData image_name. */
        public image_name: string;

        /** AvatarData base_stats. */
        public base_stats?: (model.IAvatarStatsData|null);

        /**
         * Gets the default type url for AvatarData
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an AvatarStatsData. */
    interface IAvatarStatsData {

        /** AvatarStatsData base_hp */
        base_hp?: (number|null);

        /** AvatarStatsData base_atk */
        base_atk?: (number|null);

        /** AvatarStatsData base_def */
        base_def?: (number|null);

        /** AvatarStatsData hp_curve */
        hp_curve?: (model.AvatarCurveType|null);

        /** AvatarStatsData atk_curve */
        atk_curve?: (model.AvatarCurveType|null);

        /** AvatarStatsData def_cruve */
        def_cruve?: (model.AvatarCurveType|null);

        /** AvatarStatsData specialized */
        specialized?: (model.StatType|null);

        /** AvatarStatsData promo_data */
        promo_data?: (model.IPromotionData[]|null);
    }

    /** Represents an AvatarStatsData. */
    class AvatarStatsData implements IAvatarStatsData {

        /**
         * Constructs a new AvatarStatsData.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IAvatarStatsData);

        /** AvatarStatsData base_hp. */
        public base_hp: number;

        /** AvatarStatsData base_atk. */
        public base_atk: number;

        /** AvatarStatsData base_def. */
        public base_def: number;

        /** AvatarStatsData hp_curve. */
        public hp_curve: model.AvatarCurveType;

        /** AvatarStatsData atk_curve. */
        public atk_curve: model.AvatarCurveType;

        /** AvatarStatsData def_cruve. */
        public def_cruve: model.AvatarCurveType;

        /** AvatarStatsData specialized. */
        public specialized: model.StatType;

        /** AvatarStatsData promo_data. */
        public promo_data: model.IPromotionData[];

        /**
         * Gets the default type url for AvatarStatsData
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a PromotionData. */
    interface IPromotionData {

        /** PromotionData max_level */
        max_level?: (number|Long|null);

        /** PromotionData hp */
        hp?: (number|null);

        /** PromotionData atk */
        atk?: (number|null);

        /** PromotionData def */
        def?: (number|null);

        /** PromotionData special */
        special?: (number|null);
    }

    /** Represents a PromotionData. */
    class PromotionData implements IPromotionData {

        /**
         * Constructs a new PromotionData.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IPromotionData);

        /** PromotionData max_level. */
        public max_level: (number|Long);

        /** PromotionData hp. */
        public hp: number;

        /** PromotionData atk. */
        public atk: number;

        /** PromotionData def. */
        public def: number;

        /** PromotionData special. */
        public special: number;

        /**
         * Gets the default type url for PromotionData
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** AvatarCurveType enum. */
    enum AvatarCurveType {
        GROW_CURVE_HP_S4 = 0,
        GROW_CURVE_ATTACK_S4 = 1,
        GROW_CURVE_HP_S5 = 2,
        GROW_CURVE_ATTACK_S5 = 3
    }

    /** WeaponCurveType enum. */
    enum WeaponCurveType {
        GROW_CURVE_ATTACK_101 = 0,
        GROW_CURVE_ATTACK_102 = 1,
        GROW_CURVE_ATTACK_103 = 2,
        GROW_CURVE_ATTACK_104 = 3,
        GROW_CURVE_ATTACK_105 = 4,
        GROW_CURVE_CRITICAL_101 = 5,
        GROW_CURVE_ATTACK_201 = 6,
        GROW_CURVE_ATTACK_202 = 7,
        GROW_CURVE_ATTACK_203 = 8,
        GROW_CURVE_ATTACK_204 = 9,
        GROW_CURVE_ATTACK_205 = 10,
        GROW_CURVE_CRITICAL_201 = 11,
        GROW_CURVE_ATTACK_301 = 12,
        GROW_CURVE_ATTACK_302 = 13,
        GROW_CURVE_ATTACK_303 = 14,
        GROW_CURVE_ATTACK_304 = 15,
        GROW_CURVE_ATTACK_305 = 16,
        GROW_CURVE_CRITICAL_301 = 17
    }

    /** WeaponClass enum. */
    enum WeaponClass {
        WEAPON_SWORD_ONE_HAND = 0,
        WEAPON_CLAYMORE = 1,
        WEAPON_POLE = 2,
        WEAPON_BOW = 3,
        WEAPON_CATALYST = 4
    }

    /** BodyType enum. */
    enum BodyType {
        BODY_UNKNOWN = 0,
        BODY_BOY = 1,
        BODY_GIRL = 2,
        BODY_MALE = 3,
        BODY_LADY = 4,
        BODY_LOLI = 5
    }

    /** ZoneType enum. */
    enum ZoneType {
        ASSOC_TYPE_UNKNOWN = 0,
        ASSOC_TYPE_MONDSTADT = 1,
        ASSOC_TYPE_LIYUE = 2,
        ASSOC_TYPE_INAZUMA = 3,
        ASSOC_TYPE_SUMERU = 4,
        ASSOC_TYPE_FATUI = 5
    }

    /** Element enum. */
    enum Element {
        Electric = 0,
        Fire = 1,
        Ice = 2,
        Water = 3,
        Grass = 4,
        ELEMENT_QUICKEN = 5,
        ELEMENT_FROZEN = 6,
        Wind = 7,
        Rock = 8,
        ELEMENT_NONE = 9,
        ELEMENT_PHYSICAL = 10,
        ELEMENT_UNKNOWN = 11
    }

    /** StatType enum. */
    enum StatType {
        UNSPECIFIED = 0,
        FIGHT_PROP_DEFENSE_PERCENT = 1,
        FIGHT_PROP_DEFENSE = 2,
        FIGHT_PROP_HP = 3,
        FIGHT_PROP_HP_PERCENT = 4,
        FIGHT_PROP_ATTACK = 5,
        FIGHT_PROP_ATTACK_PERCENT = 6,
        FIGHT_PROP_CHARGE_EFFICIENCY = 7,
        FIGHT_PROP_ELEMENT_MASTERY = 8,
        FIGHT_PROP_CRITICAL = 9,
        FIGHT_PROP_CRITICAL_HURT = 10,
        FIGHT_PROP_HEAL = 11,
        FIGHT_PROP_FIRE_ADD_HURT = 12,
        FIGHT_PROP_WATER_ADD_HURT = 13,
        FIGHT_PROP_GRASS_ADD_HURT = 14,
        FIGHT_PROP_ELEC_ADD_HURT = 15,
        FIGHT_PROP_WIND_ADD_HURT = 16,
        FIGHT_PROP_ICE_ADD_HURT = 17,
        FIGHT_PROP_ROCK_ADD_HURT = 18,
        FIGHT_PROP_PHYSICAL_ADD_HURT = 19,
        FIGHT_PROP_SHIELD_COST_MINUS_RATIO_ADD_HURT = 20,
        FIGHT_PROP_HEALED_ADD = 21,
        FIGHT_PROP_BASE_HP = 22,
        FIGHT_PROP_BASE_ATTACK = 23,
        FIGHT_PROP_BASE_DEFENSE = 24,
        FIGHT_PROP_MAX_HP = 25
    }

    /** SimMode enum. */
    enum SimMode {
        DURATION_MODE = 0,
        TTK_MODE = 1
    }
}

/** Namespace google. */
export namespace google {

    /** Namespace protobuf. */
    namespace protobuf {

        /** Properties of a Struct. */
        interface IStruct {

            /** Struct fields */
            fields?: ({ [k: string]: google.protobuf.IValue }|null);
        }

        /** Represents a Struct. */
        class Struct implements IStruct {

            /**
             * Constructs a new Struct.
             * @param [properties] Properties to set
             */
            constructor(properties?: google.protobuf.IStruct);

            /** Struct fields. */
            public fields: { [k: string]: google.protobuf.IValue };

            /**
             * Gets the default type url for Struct
             * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns The default type url
             */
            public static getTypeUrl(typeUrlPrefix?: string): string;
        }

        /** Properties of a Value. */
        interface IValue {

            /** Value nullValue */
            nullValue?: (google.protobuf.NullValue|null);

            /** Value numberValue */
            numberValue?: (number|null);

            /** Value stringValue */
            stringValue?: (string|null);

            /** Value boolValue */
            boolValue?: (boolean|null);

            /** Value structValue */
            structValue?: (google.protobuf.IStruct|null);

            /** Value listValue */
            listValue?: (google.protobuf.IListValue|null);
        }

        /** Represents a Value. */
        class Value implements IValue {

            /**
             * Constructs a new Value.
             * @param [properties] Properties to set
             */
            constructor(properties?: google.protobuf.IValue);

            /** Value nullValue. */
            public nullValue?: (google.protobuf.NullValue|null);

            /** Value numberValue. */
            public numberValue?: (number|null);

            /** Value stringValue. */
            public stringValue?: (string|null);

            /** Value boolValue. */
            public boolValue?: (boolean|null);

            /** Value structValue. */
            public structValue?: (google.protobuf.IStruct|null);

            /** Value listValue. */
            public listValue?: (google.protobuf.IListValue|null);

            /** Value kind. */
            public kind?: ("nullValue"|"numberValue"|"stringValue"|"boolValue"|"structValue"|"listValue");

            /**
             * Gets the default type url for Value
             * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns The default type url
             */
            public static getTypeUrl(typeUrlPrefix?: string): string;
        }

        /** NullValue enum. */
        enum NullValue {
            NULL_VALUE = 0
        }

        /** Properties of a ListValue. */
        interface IListValue {

            /** ListValue values */
            values?: (google.protobuf.IValue[]|null);
        }

        /** Represents a ListValue. */
        class ListValue implements IListValue {

            /**
             * Constructs a new ListValue.
             * @param [properties] Properties to set
             */
            constructor(properties?: google.protobuf.IListValue);

            /** ListValue values. */
            public values: google.protobuf.IValue[];

            /**
             * Gets the default type url for ListValue
             * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns The default type url
             */
            public static getTypeUrl(typeUrlPrefix?: string): string;
        }
    }
}
