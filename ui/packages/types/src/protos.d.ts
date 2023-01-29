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

        /** DBEntry tags */
        tags?: (string[]|null);
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

        /** DBEntry tags. */
        public tags: string[];

        /**
         * Creates a new DBEntry instance using the specified properties.
         * @param [properties] Properties to set
         * @returns DBEntry instance
         */
        public static create(properties?: model.IDBEntry): model.DBEntry;

        /**
         * Encodes the specified DBEntry message. Does not implicitly {@link model.DBEntry.verify|verify} messages.
         * @param message DBEntry message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.IDBEntry, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified DBEntry message, length delimited. Does not implicitly {@link model.DBEntry.verify|verify} messages.
         * @param message DBEntry message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.IDBEntry, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a DBEntry message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns DBEntry
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.DBEntry;

        /**
         * Decodes a DBEntry message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns DBEntry
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.DBEntry;

        /**
         * Verifies a DBEntry message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a DBEntry message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns DBEntry
         */
        public static fromObject(object: { [k: string]: any }): model.DBEntry;

        /**
         * Creates a plain object from a DBEntry message. Also converts values to other types if specified.
         * @param message DBEntry
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.DBEntry, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this DBEntry to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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
         * Creates a new DBEntries instance using the specified properties.
         * @param [properties] Properties to set
         * @returns DBEntries instance
         */
        public static create(properties?: model.IDBEntries): model.DBEntries;

        /**
         * Encodes the specified DBEntries message. Does not implicitly {@link model.DBEntries.verify|verify} messages.
         * @param message DBEntries message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.IDBEntries, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified DBEntries message, length delimited. Does not implicitly {@link model.DBEntries.verify|verify} messages.
         * @param message DBEntries message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.IDBEntries, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a DBEntries message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns DBEntries
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.DBEntries;

        /**
         * Decodes a DBEntries message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns DBEntries
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.DBEntries;

        /**
         * Verifies a DBEntries message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a DBEntries message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns DBEntries
         */
        public static fromObject(object: { [k: string]: any }): model.DBEntries;

        /**
         * Creates a plain object from a DBEntries message. Also converts values to other types if specified.
         * @param message DBEntries
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.DBEntries, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this DBEntries to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for DBEntries
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** SimMode enum. */
    enum SimMode {
        DURATION_MODE = 0,
        TTK_MODE = 1
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
         * Creates a new Version instance using the specified properties.
         * @param [properties] Properties to set
         * @returns Version instance
         */
        public static create(properties?: model.IVersion): model.Version;

        /**
         * Encodes the specified Version message. Does not implicitly {@link model.Version.verify|verify} messages.
         * @param message Version message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.IVersion, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified Version message, length delimited. Does not implicitly {@link model.Version.verify|verify} messages.
         * @param message Version message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.IVersion, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a Version message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns Version
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.Version;

        /**
         * Decodes a Version message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns Version
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.Version;

        /**
         * Verifies a Version message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a Version message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns Version
         */
        public static fromObject(object: { [k: string]: any }): model.Version;

        /**
         * Creates a plain object from a Version message. Also converts values to other types if specified.
         * @param message Version
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.Version, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this Version to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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

        /** SimulationResult build_date */
        build_date?: (string|null);

        /** SimulationResult modified */
        modified?: (boolean|null);

        /** SimulationResult initial_character */
        initial_character?: (string|null);

        /** SimulationResult character_details */
        character_details?: (model.ICharacter[]|null);

        /** SimulationResult target_details */
        target_details?: (model.IEnemy[]|null);

        /** SimulationResult simulator_settings */
        simulator_settings?: (model.ISimulatorSettings|null);

        /** SimulationResult energy_settings */
        energy_settings?: (model.IEnergySettings|null);

        /** SimulationResult config */
        config?: (string|null);

        /** SimulationResult sample_seed */
        sample_seed?: (string|null);

        /** SimulationResult statistics */
        statistics?: (model.ISimulationStatistics|null);
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

        /** SimulationResult build_date. */
        public build_date: string;

        /** SimulationResult modified. */
        public modified: boolean;

        /** SimulationResult initial_character. */
        public initial_character: string;

        /** SimulationResult character_details. */
        public character_details: model.ICharacter[];

        /** SimulationResult target_details. */
        public target_details: model.IEnemy[];

        /** SimulationResult simulator_settings. */
        public simulator_settings?: (model.ISimulatorSettings|null);

        /** SimulationResult energy_settings. */
        public energy_settings?: (model.IEnergySettings|null);

        /** SimulationResult config. */
        public config: string;

        /** SimulationResult sample_seed. */
        public sample_seed: string;

        /** SimulationResult statistics. */
        public statistics?: (model.ISimulationStatistics|null);

        /**
         * Creates a new SimulationResult instance using the specified properties.
         * @param [properties] Properties to set
         * @returns SimulationResult instance
         */
        public static create(properties?: model.ISimulationResult): model.SimulationResult;

        /**
         * Encodes the specified SimulationResult message. Does not implicitly {@link model.SimulationResult.verify|verify} messages.
         * @param message SimulationResult message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.ISimulationResult, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified SimulationResult message, length delimited. Does not implicitly {@link model.SimulationResult.verify|verify} messages.
         * @param message SimulationResult message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.ISimulationResult, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a SimulationResult message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns SimulationResult
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.SimulationResult;

        /**
         * Decodes a SimulationResult message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns SimulationResult
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.SimulationResult;

        /**
         * Verifies a SimulationResult message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a SimulationResult message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns SimulationResult
         */
        public static fromObject(object: { [k: string]: any }): model.SimulationResult;

        /**
         * Creates a plain object from a SimulationResult message. Also converts values to other types if specified.
         * @param message SimulationResult
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.SimulationResult, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this SimulationResult to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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

        /** SimulationStatistics SPS */
        SPS?: (model.IOverviewStats|null);

        /** SimulationStatistics total_damage */
        total_damage?: (model.IOverviewStats|null);

        /** SimulationStatistics warnings */
        warnings?: (model.IWarnings|null);

        /** SimulationStatistics failed_actions */
        failed_actions?: (model.IFailedActions[]|null);
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

        /** SimulationStatistics SPS. */
        public SPS?: (model.IOverviewStats|null);

        /** SimulationStatistics total_damage. */
        public total_damage?: (model.IOverviewStats|null);

        /** SimulationStatistics warnings. */
        public warnings?: (model.IWarnings|null);

        /** SimulationStatistics failed_actions. */
        public failed_actions: model.IFailedActions[];

        /**
         * Creates a new SimulationStatistics instance using the specified properties.
         * @param [properties] Properties to set
         * @returns SimulationStatistics instance
         */
        public static create(properties?: model.ISimulationStatistics): model.SimulationStatistics;

        /**
         * Encodes the specified SimulationStatistics message. Does not implicitly {@link model.SimulationStatistics.verify|verify} messages.
         * @param message SimulationStatistics message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.ISimulationStatistics, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified SimulationStatistics message, length delimited. Does not implicitly {@link model.SimulationStatistics.verify|verify} messages.
         * @param message SimulationStatistics message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.ISimulationStatistics, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a SimulationStatistics message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns SimulationStatistics
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.SimulationStatistics;

        /**
         * Decodes a SimulationStatistics message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns SimulationStatistics
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.SimulationStatistics;

        /**
         * Verifies a SimulationStatistics message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a SimulationStatistics message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns SimulationStatistics
         */
        public static fromObject(object: { [k: string]: any }): model.SimulationStatistics;

        /**
         * Creates a plain object from a SimulationStatistics message. Also converts values to other types if specified.
         * @param message SimulationStatistics
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.SimulationStatistics, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this SimulationStatistics to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for SimulationStatistics
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

        /** OverviewStats q1 */
        q1?: (number|null);

        /** OverviewStats q2 */
        q2?: (number|null);

        /** OverviewStats q3 */
        q3?: (number|null);

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

        /** OverviewStats q1. */
        public q1: number;

        /** OverviewStats q2. */
        public q2: number;

        /** OverviewStats q3. */
        public q3: number;

        /** OverviewStats hist. */
        public hist: (number|Long)[];

        /**
         * Creates a new OverviewStats instance using the specified properties.
         * @param [properties] Properties to set
         * @returns OverviewStats instance
         */
        public static create(properties?: model.IOverviewStats): model.OverviewStats;

        /**
         * Encodes the specified OverviewStats message. Does not implicitly {@link model.OverviewStats.verify|verify} messages.
         * @param message OverviewStats message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.IOverviewStats, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified OverviewStats message, length delimited. Does not implicitly {@link model.OverviewStats.verify|verify} messages.
         * @param message OverviewStats message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.IOverviewStats, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes an OverviewStats message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns OverviewStats
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.OverviewStats;

        /**
         * Decodes an OverviewStats message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns OverviewStats
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.OverviewStats;

        /**
         * Verifies an OverviewStats message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates an OverviewStats message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns OverviewStats
         */
        public static fromObject(object: { [k: string]: any }): model.OverviewStats;

        /**
         * Creates a plain object from an OverviewStats message. Also converts values to other types if specified.
         * @param message OverviewStats
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.OverviewStats, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this OverviewStats to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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
         * Creates a new DescriptiveStats instance using the specified properties.
         * @param [properties] Properties to set
         * @returns DescriptiveStats instance
         */
        public static create(properties?: model.IDescriptiveStats): model.DescriptiveStats;

        /**
         * Encodes the specified DescriptiveStats message. Does not implicitly {@link model.DescriptiveStats.verify|verify} messages.
         * @param message DescriptiveStats message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.IDescriptiveStats, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified DescriptiveStats message, length delimited. Does not implicitly {@link model.DescriptiveStats.verify|verify} messages.
         * @param message DescriptiveStats message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.IDescriptiveStats, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a DescriptiveStats message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns DescriptiveStats
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.DescriptiveStats;

        /**
         * Decodes a DescriptiveStats message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns DescriptiveStats
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.DescriptiveStats;

        /**
         * Verifies a DescriptiveStats message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a DescriptiveStats message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns DescriptiveStats
         */
        public static fromObject(object: { [k: string]: any }): model.DescriptiveStats;

        /**
         * Creates a plain object from a DescriptiveStats message. Also converts values to other types if specified.
         * @param message DescriptiveStats
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.DescriptiveStats, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this DescriptiveStats to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for DescriptiveStats
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
         * Creates a new Warnings instance using the specified properties.
         * @param [properties] Properties to set
         * @returns Warnings instance
         */
        public static create(properties?: model.IWarnings): model.Warnings;

        /**
         * Encodes the specified Warnings message. Does not implicitly {@link model.Warnings.verify|verify} messages.
         * @param message Warnings message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.IWarnings, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified Warnings message, length delimited. Does not implicitly {@link model.Warnings.verify|verify} messages.
         * @param message Warnings message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.IWarnings, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a Warnings message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns Warnings
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.Warnings;

        /**
         * Decodes a Warnings message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns Warnings
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.Warnings;

        /**
         * Verifies a Warnings message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a Warnings message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns Warnings
         */
        public static fromObject(object: { [k: string]: any }): model.Warnings;

        /**
         * Creates a plain object from a Warnings message. Also converts values to other types if specified.
         * @param message Warnings
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.Warnings, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this Warnings to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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
         * Creates a new FailedActions instance using the specified properties.
         * @param [properties] Properties to set
         * @returns FailedActions instance
         */
        public static create(properties?: model.IFailedActions): model.FailedActions;

        /**
         * Encodes the specified FailedActions message. Does not implicitly {@link model.FailedActions.verify|verify} messages.
         * @param message FailedActions message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.IFailedActions, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified FailedActions message, length delimited. Does not implicitly {@link model.FailedActions.verify|verify} messages.
         * @param message FailedActions message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.IFailedActions, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a FailedActions message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns FailedActions
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.FailedActions;

        /**
         * Decodes a FailedActions message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns FailedActions
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.FailedActions;

        /**
         * Verifies a FailedActions message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a FailedActions message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns FailedActions
         */
        public static fromObject(object: { [k: string]: any }): model.FailedActions;

        /**
         * Creates a plain object from a FailedActions message. Also converts values to other types if specified.
         * @param message FailedActions
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.FailedActions, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this FailedActions to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for FailedActions
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a Character. */
    interface ICharacter {

        /** Character key */
        key?: (string|null);

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
        stats?: ({ [k: string]: number }|null);

        /** Character snapshot */
        snapshot?: ({ [k: string]: number }|null);
    }

    /** Represents a Character. */
    class Character implements ICharacter {

        /**
         * Constructs a new Character.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.ICharacter);

        /** Character key. */
        public key: string;

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
        public stats: { [k: string]: number };

        /** Character snapshot. */
        public snapshot: { [k: string]: number };

        /**
         * Creates a new Character instance using the specified properties.
         * @param [properties] Properties to set
         * @returns Character instance
         */
        public static create(properties?: model.ICharacter): model.Character;

        /**
         * Encodes the specified Character message. Does not implicitly {@link model.Character.verify|verify} messages.
         * @param message Character message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.ICharacter, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified Character message, length delimited. Does not implicitly {@link model.Character.verify|verify} messages.
         * @param message Character message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.ICharacter, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a Character message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns Character
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.Character;

        /**
         * Decodes a Character message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns Character
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.Character;

        /**
         * Verifies a Character message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a Character message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns Character
         */
        public static fromObject(object: { [k: string]: any }): model.Character;

        /**
         * Creates a plain object from a Character message. Also converts values to other types if specified.
         * @param message Character
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.Character, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this Character to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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
         * Creates a new CharacterTalents instance using the specified properties.
         * @param [properties] Properties to set
         * @returns CharacterTalents instance
         */
        public static create(properties?: model.ICharacterTalents): model.CharacterTalents;

        /**
         * Encodes the specified CharacterTalents message. Does not implicitly {@link model.CharacterTalents.verify|verify} messages.
         * @param message CharacterTalents message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.ICharacterTalents, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified CharacterTalents message, length delimited. Does not implicitly {@link model.CharacterTalents.verify|verify} messages.
         * @param message CharacterTalents message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.ICharacterTalents, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a CharacterTalents message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns CharacterTalents
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.CharacterTalents;

        /**
         * Decodes a CharacterTalents message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns CharacterTalents
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.CharacterTalents;

        /**
         * Verifies a CharacterTalents message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a CharacterTalents message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns CharacterTalents
         */
        public static fromObject(object: { [k: string]: any }): model.CharacterTalents;

        /**
         * Creates a plain object from a CharacterTalents message. Also converts values to other types if specified.
         * @param message CharacterTalents
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.CharacterTalents, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this CharacterTalents to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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
         * Creates a new Weapon instance using the specified properties.
         * @param [properties] Properties to set
         * @returns Weapon instance
         */
        public static create(properties?: model.IWeapon): model.Weapon;

        /**
         * Encodes the specified Weapon message. Does not implicitly {@link model.Weapon.verify|verify} messages.
         * @param message Weapon message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.IWeapon, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified Weapon message, length delimited. Does not implicitly {@link model.Weapon.verify|verify} messages.
         * @param message Weapon message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.IWeapon, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a Weapon message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns Weapon
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.Weapon;

        /**
         * Decodes a Weapon message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns Weapon
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.Weapon;

        /**
         * Verifies a Weapon message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a Weapon message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns Weapon
         */
        public static fromObject(object: { [k: string]: any }): model.Weapon;

        /**
         * Creates a plain object from a Weapon message. Also converts values to other types if specified.
         * @param message Weapon
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.Weapon, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this Weapon to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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
         * Creates a new Enemy instance using the specified properties.
         * @param [properties] Properties to set
         * @returns Enemy instance
         */
        public static create(properties?: model.IEnemy): model.Enemy;

        /**
         * Encodes the specified Enemy message. Does not implicitly {@link model.Enemy.verify|verify} messages.
         * @param message Enemy message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.IEnemy, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified Enemy message, length delimited. Does not implicitly {@link model.Enemy.verify|verify} messages.
         * @param message Enemy message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.IEnemy, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes an Enemy message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns Enemy
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.Enemy;

        /**
         * Decodes an Enemy message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns Enemy
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.Enemy;

        /**
         * Verifies an Enemy message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates an Enemy message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns Enemy
         */
        public static fromObject(object: { [k: string]: any }): model.Enemy;

        /**
         * Creates a plain object from an Enemy message. Also converts values to other types if specified.
         * @param message Enemy
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.Enemy, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this Enemy to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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
         * Creates a new Coord instance using the specified properties.
         * @param [properties] Properties to set
         * @returns Coord instance
         */
        public static create(properties?: model.ICoord): model.Coord;

        /**
         * Encodes the specified Coord message. Does not implicitly {@link model.Coord.verify|verify} messages.
         * @param message Coord message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.ICoord, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified Coord message, length delimited. Does not implicitly {@link model.Coord.verify|verify} messages.
         * @param message Coord message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.ICoord, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a Coord message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns Coord
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.Coord;

        /**
         * Decodes a Coord message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns Coord
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.Coord;

        /**
         * Verifies a Coord message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a Coord message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns Coord
         */
        public static fromObject(object: { [k: string]: any }): model.Coord;

        /**
         * Creates a plain object from a Coord message. Also converts values to other types if specified.
         * @param message Coord
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.Coord, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this Coord to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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

        /** SimulatorSettings number_of_iterations */
        number_of_iterations?: (number|Long|null);

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

        /** SimulatorSettings number_of_iterations. */
        public number_of_iterations: (number|Long);

        /** SimulatorSettings delays. */
        public delays?: (model.IDelays|null);

        /**
         * Creates a new SimulatorSettings instance using the specified properties.
         * @param [properties] Properties to set
         * @returns SimulatorSettings instance
         */
        public static create(properties?: model.ISimulatorSettings): model.SimulatorSettings;

        /**
         * Encodes the specified SimulatorSettings message. Does not implicitly {@link model.SimulatorSettings.verify|verify} messages.
         * @param message SimulatorSettings message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.ISimulatorSettings, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified SimulatorSettings message, length delimited. Does not implicitly {@link model.SimulatorSettings.verify|verify} messages.
         * @param message SimulatorSettings message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.ISimulatorSettings, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a SimulatorSettings message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns SimulatorSettings
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.SimulatorSettings;

        /**
         * Decodes a SimulatorSettings message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns SimulatorSettings
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.SimulatorSettings;

        /**
         * Verifies a SimulatorSettings message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a SimulatorSettings message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns SimulatorSettings
         */
        public static fromObject(object: { [k: string]: any }): model.SimulatorSettings;

        /**
         * Creates a plain object from a SimulatorSettings message. Also converts values to other types if specified.
         * @param message SimulatorSettings
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.SimulatorSettings, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this SimulatorSettings to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for SimulatorSettings
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a Delays. */
    interface IDelays {

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

        /** Delays swap. */
        public swap: (number|Long);

        /**
         * Creates a new Delays instance using the specified properties.
         * @param [properties] Properties to set
         * @returns Delays instance
         */
        public static create(properties?: model.IDelays): model.Delays;

        /**
         * Encodes the specified Delays message. Does not implicitly {@link model.Delays.verify|verify} messages.
         * @param message Delays message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.IDelays, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified Delays message, length delimited. Does not implicitly {@link model.Delays.verify|verify} messages.
         * @param message Delays message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.IDelays, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a Delays message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns Delays
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.Delays;

        /**
         * Decodes a Delays message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns Delays
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.Delays;

        /**
         * Verifies a Delays message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a Delays message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns Delays
         */
        public static fromObject(object: { [k: string]: any }): model.Delays;

        /**
         * Creates a plain object from a Delays message. Also converts values to other types if specified.
         * @param message Delays
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.Delays, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this Delays to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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
         * Creates a new EnergySettings instance using the specified properties.
         * @param [properties] Properties to set
         * @returns EnergySettings instance
         */
        public static create(properties?: model.IEnergySettings): model.EnergySettings;

        /**
         * Encodes the specified EnergySettings message. Does not implicitly {@link model.EnergySettings.verify|verify} messages.
         * @param message EnergySettings message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.IEnergySettings, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified EnergySettings message, length delimited. Does not implicitly {@link model.EnergySettings.verify|verify} messages.
         * @param message EnergySettings message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.IEnergySettings, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes an EnergySettings message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns EnergySettings
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.EnergySettings;

        /**
         * Decodes an EnergySettings message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns EnergySettings
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.EnergySettings;

        /**
         * Verifies an EnergySettings message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates an EnergySettings message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns EnergySettings
         */
        public static fromObject(object: { [k: string]: any }): model.EnergySettings;

        /**
         * Creates a plain object from an EnergySettings message. Also converts values to other types if specified.
         * @param message EnergySettings
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.EnergySettings, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this EnergySettings to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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
         * Creates a new AvatarData instance using the specified properties.
         * @param [properties] Properties to set
         * @returns AvatarData instance
         */
        public static create(properties?: model.IAvatarData): model.AvatarData;

        /**
         * Encodes the specified AvatarData message. Does not implicitly {@link model.AvatarData.verify|verify} messages.
         * @param message AvatarData message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.IAvatarData, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified AvatarData message, length delimited. Does not implicitly {@link model.AvatarData.verify|verify} messages.
         * @param message AvatarData message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.IAvatarData, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes an AvatarData message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns AvatarData
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.AvatarData;

        /**
         * Decodes an AvatarData message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns AvatarData
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.AvatarData;

        /**
         * Verifies an AvatarData message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates an AvatarData message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns AvatarData
         */
        public static fromObject(object: { [k: string]: any }): model.AvatarData;

        /**
         * Creates a plain object from an AvatarData message. Also converts values to other types if specified.
         * @param message AvatarData
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.AvatarData, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this AvatarData to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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
         * Creates a new AvatarStatsData instance using the specified properties.
         * @param [properties] Properties to set
         * @returns AvatarStatsData instance
         */
        public static create(properties?: model.IAvatarStatsData): model.AvatarStatsData;

        /**
         * Encodes the specified AvatarStatsData message. Does not implicitly {@link model.AvatarStatsData.verify|verify} messages.
         * @param message AvatarStatsData message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.IAvatarStatsData, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified AvatarStatsData message, length delimited. Does not implicitly {@link model.AvatarStatsData.verify|verify} messages.
         * @param message AvatarStatsData message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.IAvatarStatsData, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes an AvatarStatsData message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns AvatarStatsData
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.AvatarStatsData;

        /**
         * Decodes an AvatarStatsData message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns AvatarStatsData
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.AvatarStatsData;

        /**
         * Verifies an AvatarStatsData message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates an AvatarStatsData message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns AvatarStatsData
         */
        public static fromObject(object: { [k: string]: any }): model.AvatarStatsData;

        /**
         * Creates a plain object from an AvatarStatsData message. Also converts values to other types if specified.
         * @param message AvatarStatsData
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.AvatarStatsData, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this AvatarStatsData to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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
         * Creates a new PromotionData instance using the specified properties.
         * @param [properties] Properties to set
         * @returns PromotionData instance
         */
        public static create(properties?: model.IPromotionData): model.PromotionData;

        /**
         * Encodes the specified PromotionData message. Does not implicitly {@link model.PromotionData.verify|verify} messages.
         * @param message PromotionData message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: model.IPromotionData, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified PromotionData message, length delimited. Does not implicitly {@link model.PromotionData.verify|verify} messages.
         * @param message PromotionData message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: model.IPromotionData, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a PromotionData message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns PromotionData
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): model.PromotionData;

        /**
         * Decodes a PromotionData message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns PromotionData
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): model.PromotionData;

        /**
         * Verifies a PromotionData message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a PromotionData message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns PromotionData
         */
        public static fromObject(object: { [k: string]: any }): model.PromotionData;

        /**
         * Creates a plain object from a PromotionData message. Also converts values to other types if specified.
         * @param message PromotionData
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: model.PromotionData, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this PromotionData to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

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
}

/** Namespace db. */
export namespace db {

    /** Represents a DBStore */
    class DBStore extends $protobuf.rpc.Service {

        /**
         * Constructs a new DBStore service.
         * @param rpcImpl RPC implementation
         * @param [requestDelimited=false] Whether requests are length-delimited
         * @param [responseDelimited=false] Whether responses are length-delimited
         */
        constructor(rpcImpl: $protobuf.RPCImpl, requestDelimited?: boolean, responseDelimited?: boolean);

        /**
         * Creates new DBStore service using the specified rpc implementation.
         * @param rpcImpl RPC implementation
         * @param [requestDelimited=false] Whether requests are length-delimited
         * @param [responseDelimited=false] Whether responses are length-delimited
         * @returns RPC service. Useful where requests and/or responses are streamed.
         */
        public static create(rpcImpl: $protobuf.RPCImpl, requestDelimited?: boolean, responseDelimited?: boolean): DBStore;

        /**
         * Calls Create.
         * @param request CreateRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and CreateResponse
         */
        public create(request: db.ICreateRequest, callback: db.DBStore.CreateCallback): void;

        /**
         * Calls Create.
         * @param request CreateRequest message or plain object
         * @returns Promise
         */
        public create(request: db.ICreateRequest): Promise<db.CreateResponse>;

        /**
         * Calls Get.
         * @param request GetRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and GetResponse
         */
        public get(request: db.IGetRequest, callback: db.DBStore.GetCallback): void;

        /**
         * Calls Get.
         * @param request GetRequest message or plain object
         * @returns Promise
         */
        public get(request: db.IGetRequest): Promise<db.GetResponse>;
    }

    namespace DBStore {

        /**
         * Callback as used by {@link db.DBStore#create}.
         * @param error Error, if any
         * @param [response] CreateResponse
         */
        type CreateCallback = (error: (Error|null), response?: db.CreateResponse) => void;

        /**
         * Callback as used by {@link db.DBStore#get}.
         * @param error Error, if any
         * @param [response] GetResponse
         */
        type GetCallback = (error: (Error|null), response?: db.GetResponse) => void;
    }

    /** Properties of a CreateRequest. */
    interface ICreateRequest {

        /** CreateRequest data */
        data?: (model.IDBEntry|null);
    }

    /** Represents a CreateRequest. */
    class CreateRequest implements ICreateRequest {

        /**
         * Constructs a new CreateRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.ICreateRequest);

        /** CreateRequest data. */
        public data?: (model.IDBEntry|null);

        /**
         * Creates a new CreateRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns CreateRequest instance
         */
        public static create(properties?: db.ICreateRequest): db.CreateRequest;

        /**
         * Encodes the specified CreateRequest message. Does not implicitly {@link db.CreateRequest.verify|verify} messages.
         * @param message CreateRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: db.ICreateRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified CreateRequest message, length delimited. Does not implicitly {@link db.CreateRequest.verify|verify} messages.
         * @param message CreateRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: db.ICreateRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a CreateRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns CreateRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): db.CreateRequest;

        /**
         * Decodes a CreateRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns CreateRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): db.CreateRequest;

        /**
         * Verifies a CreateRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a CreateRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns CreateRequest
         */
        public static fromObject(object: { [k: string]: any }): db.CreateRequest;

        /**
         * Creates a plain object from a CreateRequest message. Also converts values to other types if specified.
         * @param message CreateRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: db.CreateRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this CreateRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for CreateRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a CreateResponse. */
    interface ICreateResponse {

        /** CreateResponse key */
        key?: (string|null);
    }

    /** Represents a CreateResponse. */
    class CreateResponse implements ICreateResponse {

        /**
         * Constructs a new CreateResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.ICreateResponse);

        /** CreateResponse key. */
        public key: string;

        /**
         * Creates a new CreateResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns CreateResponse instance
         */
        public static create(properties?: db.ICreateResponse): db.CreateResponse;

        /**
         * Encodes the specified CreateResponse message. Does not implicitly {@link db.CreateResponse.verify|verify} messages.
         * @param message CreateResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: db.ICreateResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified CreateResponse message, length delimited. Does not implicitly {@link db.CreateResponse.verify|verify} messages.
         * @param message CreateResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: db.ICreateResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a CreateResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns CreateResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): db.CreateResponse;

        /**
         * Decodes a CreateResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns CreateResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): db.CreateResponse;

        /**
         * Verifies a CreateResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a CreateResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns CreateResponse
         */
        public static fromObject(object: { [k: string]: any }): db.CreateResponse;

        /**
         * Creates a plain object from a CreateResponse message. Also converts values to other types if specified.
         * @param message CreateResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: db.CreateResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this CreateResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for CreateResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetRequest. */
    interface IGetRequest {

        /** GetRequest query */
        query?: (google.protobuf.IStruct|null);

        /** GetRequest limit */
        limit?: (number|Long|null);

        /** GetRequest page */
        page?: (number|Long|null);
    }

    /** Represents a GetRequest. */
    class GetRequest implements IGetRequest {

        /**
         * Constructs a new GetRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IGetRequest);

        /** GetRequest query. */
        public query?: (google.protobuf.IStruct|null);

        /** GetRequest limit. */
        public limit: (number|Long);

        /** GetRequest page. */
        public page: (number|Long);

        /**
         * Creates a new GetRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns GetRequest instance
         */
        public static create(properties?: db.IGetRequest): db.GetRequest;

        /**
         * Encodes the specified GetRequest message. Does not implicitly {@link db.GetRequest.verify|verify} messages.
         * @param message GetRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: db.IGetRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified GetRequest message, length delimited. Does not implicitly {@link db.GetRequest.verify|verify} messages.
         * @param message GetRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: db.IGetRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a GetRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns GetRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): db.GetRequest;

        /**
         * Decodes a GetRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns GetRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): db.GetRequest;

        /**
         * Verifies a GetRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a GetRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns GetRequest
         */
        public static fromObject(object: { [k: string]: any }): db.GetRequest;

        /**
         * Creates a plain object from a GetRequest message. Also converts values to other types if specified.
         * @param message GetRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: db.GetRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this GetRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for GetRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetResponse. */
    interface IGetResponse {

        /** GetResponse data */
        data?: (model.IDBEntries|null);
    }

    /** Represents a GetResponse. */
    class GetResponse implements IGetResponse {

        /**
         * Constructs a new GetResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IGetResponse);

        /** GetResponse data. */
        public data?: (model.IDBEntries|null);

        /**
         * Creates a new GetResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns GetResponse instance
         */
        public static create(properties?: db.IGetResponse): db.GetResponse;

        /**
         * Encodes the specified GetResponse message. Does not implicitly {@link db.GetResponse.verify|verify} messages.
         * @param message GetResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: db.IGetResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified GetResponse message, length delimited. Does not implicitly {@link db.GetResponse.verify|verify} messages.
         * @param message GetResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: db.IGetResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a GetResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns GetResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): db.GetResponse;

        /**
         * Decodes a GetResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns GetResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): db.GetResponse;

        /**
         * Verifies a GetResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a GetResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns GetResponse
         */
        public static fromObject(object: { [k: string]: any }): db.GetResponse;

        /**
         * Creates a plain object from a GetResponse message. Also converts values to other types if specified.
         * @param message GetResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: db.GetResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this GetResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for GetResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
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
             * Creates a new Struct instance using the specified properties.
             * @param [properties] Properties to set
             * @returns Struct instance
             */
            public static create(properties?: google.protobuf.IStruct): google.protobuf.Struct;

            /**
             * Encodes the specified Struct message. Does not implicitly {@link google.protobuf.Struct.verify|verify} messages.
             * @param message Struct message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encode(message: google.protobuf.IStruct, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Encodes the specified Struct message, length delimited. Does not implicitly {@link google.protobuf.Struct.verify|verify} messages.
             * @param message Struct message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encodeDelimited(message: google.protobuf.IStruct, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Decodes a Struct message from the specified reader or buffer.
             * @param reader Reader or buffer to decode from
             * @param [length] Message length if known beforehand
             * @returns Struct
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): google.protobuf.Struct;

            /**
             * Decodes a Struct message from the specified reader or buffer, length delimited.
             * @param reader Reader or buffer to decode from
             * @returns Struct
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): google.protobuf.Struct;

            /**
             * Verifies a Struct message.
             * @param message Plain object to verify
             * @returns `null` if valid, otherwise the reason why it is not
             */
            public static verify(message: { [k: string]: any }): (string|null);

            /**
             * Creates a Struct message from a plain object. Also converts values to their respective internal types.
             * @param object Plain object
             * @returns Struct
             */
            public static fromObject(object: { [k: string]: any }): google.protobuf.Struct;

            /**
             * Creates a plain object from a Struct message. Also converts values to other types if specified.
             * @param message Struct
             * @param [options] Conversion options
             * @returns Plain object
             */
            public static toObject(message: google.protobuf.Struct, options?: $protobuf.IConversionOptions): { [k: string]: any };

            /**
             * Converts this Struct to JSON.
             * @returns JSON object
             */
            public toJSON(): { [k: string]: any };

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
             * Creates a new Value instance using the specified properties.
             * @param [properties] Properties to set
             * @returns Value instance
             */
            public static create(properties?: google.protobuf.IValue): google.protobuf.Value;

            /**
             * Encodes the specified Value message. Does not implicitly {@link google.protobuf.Value.verify|verify} messages.
             * @param message Value message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encode(message: google.protobuf.IValue, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Encodes the specified Value message, length delimited. Does not implicitly {@link google.protobuf.Value.verify|verify} messages.
             * @param message Value message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encodeDelimited(message: google.protobuf.IValue, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Decodes a Value message from the specified reader or buffer.
             * @param reader Reader or buffer to decode from
             * @param [length] Message length if known beforehand
             * @returns Value
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): google.protobuf.Value;

            /**
             * Decodes a Value message from the specified reader or buffer, length delimited.
             * @param reader Reader or buffer to decode from
             * @returns Value
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): google.protobuf.Value;

            /**
             * Verifies a Value message.
             * @param message Plain object to verify
             * @returns `null` if valid, otherwise the reason why it is not
             */
            public static verify(message: { [k: string]: any }): (string|null);

            /**
             * Creates a Value message from a plain object. Also converts values to their respective internal types.
             * @param object Plain object
             * @returns Value
             */
            public static fromObject(object: { [k: string]: any }): google.protobuf.Value;

            /**
             * Creates a plain object from a Value message. Also converts values to other types if specified.
             * @param message Value
             * @param [options] Conversion options
             * @returns Plain object
             */
            public static toObject(message: google.protobuf.Value, options?: $protobuf.IConversionOptions): { [k: string]: any };

            /**
             * Converts this Value to JSON.
             * @returns JSON object
             */
            public toJSON(): { [k: string]: any };

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
             * Creates a new ListValue instance using the specified properties.
             * @param [properties] Properties to set
             * @returns ListValue instance
             */
            public static create(properties?: google.protobuf.IListValue): google.protobuf.ListValue;

            /**
             * Encodes the specified ListValue message. Does not implicitly {@link google.protobuf.ListValue.verify|verify} messages.
             * @param message ListValue message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encode(message: google.protobuf.IListValue, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Encodes the specified ListValue message, length delimited. Does not implicitly {@link google.protobuf.ListValue.verify|verify} messages.
             * @param message ListValue message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encodeDelimited(message: google.protobuf.IListValue, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Decodes a ListValue message from the specified reader or buffer.
             * @param reader Reader or buffer to decode from
             * @param [length] Message length if known beforehand
             * @returns ListValue
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): google.protobuf.ListValue;

            /**
             * Decodes a ListValue message from the specified reader or buffer, length delimited.
             * @param reader Reader or buffer to decode from
             * @returns ListValue
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): google.protobuf.ListValue;

            /**
             * Verifies a ListValue message.
             * @param message Plain object to verify
             * @returns `null` if valid, otherwise the reason why it is not
             */
            public static verify(message: { [k: string]: any }): (string|null);

            /**
             * Creates a ListValue message from a plain object. Also converts values to their respective internal types.
             * @param object Plain object
             * @returns ListValue
             */
            public static fromObject(object: { [k: string]: any }): google.protobuf.ListValue;

            /**
             * Creates a plain object from a ListValue message. Also converts values to other types if specified.
             * @param message ListValue
             * @param [options] Conversion options
             * @returns Plain object
             */
            public static toObject(message: google.protobuf.ListValue, options?: $protobuf.IConversionOptions): { [k: string]: any };

            /**
             * Converts this ListValue to JSON.
             * @returns JSON object
             */
            public toJSON(): { [k: string]: any };

            /**
             * Gets the default type url for ListValue
             * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns The default type url
             */
            public static getTypeUrl(typeUrlPrefix?: string): string;
        }
    }
}

/** Namespace result. */
export namespace result {

    /** Represents a SubmissionStore */
    class SubmissionStore extends $protobuf.rpc.Service {

        /**
         * Constructs a new SubmissionStore service.
         * @param rpcImpl RPC implementation
         * @param [requestDelimited=false] Whether requests are length-delimited
         * @param [responseDelimited=false] Whether responses are length-delimited
         */
        constructor(rpcImpl: $protobuf.RPCImpl, requestDelimited?: boolean, responseDelimited?: boolean);

        /**
         * Creates new SubmissionStore service using the specified rpc implementation.
         * @param rpcImpl RPC implementation
         * @param [requestDelimited=false] Whether requests are length-delimited
         * @param [responseDelimited=false] Whether responses are length-delimited
         * @returns RPC service. Useful where requests and/or responses are streamed.
         */
        public static create(rpcImpl: $protobuf.RPCImpl, requestDelimited?: boolean, responseDelimited?: boolean): SubmissionStore;

        /**
         * Calls Submit.
         * @param request SubmitRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and SubmitResponse
         */
        public submit(request: result.ISubmitRequest, callback: result.SubmissionStore.SubmitCallback): void;

        /**
         * Calls Submit.
         * @param request SubmitRequest message or plain object
         * @returns Promise
         */
        public submit(request: result.ISubmitRequest): Promise<result.SubmitResponse>;

        /**
         * Calls Update.
         * @param request UpdateRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and UpdateResponse
         */
        public update(request: result.IUpdateRequest, callback: result.SubmissionStore.UpdateCallback): void;

        /**
         * Calls Update.
         * @param request UpdateRequest message or plain object
         * @returns Promise
         */
        public update(request: result.IUpdateRequest): Promise<result.UpdateResponse>;

        /**
         * Calls Remove.
         * @param request RemoveRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and RemoveResponse
         */
        public remove(request: result.IRemoveRequest, callback: result.SubmissionStore.RemoveCallback): void;

        /**
         * Calls Remove.
         * @param request RemoveRequest message or plain object
         * @returns Promise
         */
        public remove(request: result.IRemoveRequest): Promise<result.RemoveResponse>;

        /**
         * Calls List.
         * @param request ListRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and ListResponse
         */
        public list(request: result.IListRequest, callback: result.SubmissionStore.ListCallback): void;

        /**
         * Calls List.
         * @param request ListRequest message or plain object
         * @returns Promise
         */
        public list(request: result.IListRequest): Promise<result.ListResponse>;

        /**
         * Calls Approve.
         * @param request ApproveRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and ApproveResponse
         */
        public approve(request: result.IApproveRequest, callback: result.SubmissionStore.ApproveCallback): void;

        /**
         * Calls Approve.
         * @param request ApproveRequest message or plain object
         * @returns Promise
         */
        public approve(request: result.IApproveRequest): Promise<result.ApproveResponse>;

        /**
         * Calls Replace.
         * @param request ReplaceRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and ReplaceResponse
         */
        public replace(request: result.IReplaceRequest, callback: result.SubmissionStore.ReplaceCallback): void;

        /**
         * Calls Replace.
         * @param request ReplaceRequest message or plain object
         * @returns Promise
         */
        public replace(request: result.IReplaceRequest): Promise<result.ReplaceResponse>;

        /**
         * Calls Reject.
         * @param request RejectRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and RejectResponse
         */
        public reject(request: result.IRejectRequest, callback: result.SubmissionStore.RejectCallback): void;

        /**
         * Calls Reject.
         * @param request RejectRequest message or plain object
         * @returns Promise
         */
        public reject(request: result.IRejectRequest): Promise<result.RejectResponse>;
    }

    namespace SubmissionStore {

        /**
         * Callback as used by {@link result.SubmissionStore#submit}.
         * @param error Error, if any
         * @param [response] SubmitResponse
         */
        type SubmitCallback = (error: (Error|null), response?: result.SubmitResponse) => void;

        /**
         * Callback as used by {@link result.SubmissionStore#update}.
         * @param error Error, if any
         * @param [response] UpdateResponse
         */
        type UpdateCallback = (error: (Error|null), response?: result.UpdateResponse) => void;

        /**
         * Callback as used by {@link result.SubmissionStore#remove}.
         * @param error Error, if any
         * @param [response] RemoveResponse
         */
        type RemoveCallback = (error: (Error|null), response?: result.RemoveResponse) => void;

        /**
         * Callback as used by {@link result.SubmissionStore#list}.
         * @param error Error, if any
         * @param [response] ListResponse
         */
        type ListCallback = (error: (Error|null), response?: result.ListResponse) => void;

        /**
         * Callback as used by {@link result.SubmissionStore#approve}.
         * @param error Error, if any
         * @param [response] ApproveResponse
         */
        type ApproveCallback = (error: (Error|null), response?: result.ApproveResponse) => void;

        /**
         * Callback as used by {@link result.SubmissionStore#replace}.
         * @param error Error, if any
         * @param [response] ReplaceResponse
         */
        type ReplaceCallback = (error: (Error|null), response?: result.ReplaceResponse) => void;

        /**
         * Callback as used by {@link result.SubmissionStore#reject}.
         * @param error Error, if any
         * @param [response] RejectResponse
         */
        type RejectCallback = (error: (Error|null), response?: result.RejectResponse) => void;
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
    }

    /** Represents a Submission. */
    class Submission implements ISubmission {

        /**
         * Constructs a new Submission.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.ISubmission);

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

        /**
         * Creates a new Submission instance using the specified properties.
         * @param [properties] Properties to set
         * @returns Submission instance
         */
        public static create(properties?: result.ISubmission): result.Submission;

        /**
         * Encodes the specified Submission message. Does not implicitly {@link result.Submission.verify|verify} messages.
         * @param message Submission message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.ISubmission, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified Submission message, length delimited. Does not implicitly {@link result.Submission.verify|verify} messages.
         * @param message Submission message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.ISubmission, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a Submission message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns Submission
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.Submission;

        /**
         * Decodes a Submission message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns Submission
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.Submission;

        /**
         * Verifies a Submission message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a Submission message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns Submission
         */
        public static fromObject(object: { [k: string]: any }): result.Submission;

        /**
         * Creates a plain object from a Submission message. Also converts values to other types if specified.
         * @param message Submission
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.Submission, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this Submission to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for Submission
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a SubmitRequest. */
    interface ISubmitRequest {

        /** SubmitRequest config */
        config?: (string|null);

        /** SubmitRequest submitter */
        submitter?: (string|null);

        /** SubmitRequest description */
        description?: (string|null);
    }

    /** Represents a SubmitRequest. */
    class SubmitRequest implements ISubmitRequest {

        /**
         * Constructs a new SubmitRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.ISubmitRequest);

        /** SubmitRequest config. */
        public config: string;

        /** SubmitRequest submitter. */
        public submitter: string;

        /** SubmitRequest description. */
        public description: string;

        /**
         * Creates a new SubmitRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns SubmitRequest instance
         */
        public static create(properties?: result.ISubmitRequest): result.SubmitRequest;

        /**
         * Encodes the specified SubmitRequest message. Does not implicitly {@link result.SubmitRequest.verify|verify} messages.
         * @param message SubmitRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.ISubmitRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified SubmitRequest message, length delimited. Does not implicitly {@link result.SubmitRequest.verify|verify} messages.
         * @param message SubmitRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.ISubmitRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a SubmitRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns SubmitRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.SubmitRequest;

        /**
         * Decodes a SubmitRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns SubmitRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.SubmitRequest;

        /**
         * Verifies a SubmitRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a SubmitRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns SubmitRequest
         */
        public static fromObject(object: { [k: string]: any }): result.SubmitRequest;

        /**
         * Creates a plain object from a SubmitRequest message. Also converts values to other types if specified.
         * @param message SubmitRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.SubmitRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this SubmitRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for SubmitRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a SubmitResponse. */
    interface ISubmitResponse {

        /** SubmitResponse id */
        id?: (string|null);
    }

    /** Represents a SubmitResponse. */
    class SubmitResponse implements ISubmitResponse {

        /**
         * Constructs a new SubmitResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.ISubmitResponse);

        /** SubmitResponse id. */
        public id: string;

        /**
         * Creates a new SubmitResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns SubmitResponse instance
         */
        public static create(properties?: result.ISubmitResponse): result.SubmitResponse;

        /**
         * Encodes the specified SubmitResponse message. Does not implicitly {@link result.SubmitResponse.verify|verify} messages.
         * @param message SubmitResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.ISubmitResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified SubmitResponse message, length delimited. Does not implicitly {@link result.SubmitResponse.verify|verify} messages.
         * @param message SubmitResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.ISubmitResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a SubmitResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns SubmitResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.SubmitResponse;

        /**
         * Decodes a SubmitResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns SubmitResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.SubmitResponse;

        /**
         * Verifies a SubmitResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a SubmitResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns SubmitResponse
         */
        public static fromObject(object: { [k: string]: any }): result.SubmitResponse;

        /**
         * Creates a plain object from a SubmitResponse message. Also converts values to other types if specified.
         * @param message SubmitResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.SubmitResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this SubmitResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for SubmitResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a ListRequest. */
    interface IListRequest {

        /** ListRequest user_filter */
        user_filter?: (string|null);
    }

    /** Represents a ListRequest. */
    class ListRequest implements IListRequest {

        /**
         * Constructs a new ListRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.IListRequest);

        /** ListRequest user_filter. */
        public user_filter: string;

        /**
         * Creates a new ListRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns ListRequest instance
         */
        public static create(properties?: result.IListRequest): result.ListRequest;

        /**
         * Encodes the specified ListRequest message. Does not implicitly {@link result.ListRequest.verify|verify} messages.
         * @param message ListRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.IListRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified ListRequest message, length delimited. Does not implicitly {@link result.ListRequest.verify|verify} messages.
         * @param message ListRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.IListRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a ListRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns ListRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.ListRequest;

        /**
         * Decodes a ListRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns ListRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.ListRequest;

        /**
         * Verifies a ListRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a ListRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns ListRequest
         */
        public static fromObject(object: { [k: string]: any }): result.ListRequest;

        /**
         * Creates a plain object from a ListRequest message. Also converts values to other types if specified.
         * @param message ListRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.ListRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this ListRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for ListRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a ListResponse. */
    interface IListResponse {

        /** ListResponse data */
        data?: (result.ISubmission[]|null);
    }

    /** Represents a ListResponse. */
    class ListResponse implements IListResponse {

        /**
         * Constructs a new ListResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.IListResponse);

        /** ListResponse data. */
        public data: result.ISubmission[];

        /**
         * Creates a new ListResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns ListResponse instance
         */
        public static create(properties?: result.IListResponse): result.ListResponse;

        /**
         * Encodes the specified ListResponse message. Does not implicitly {@link result.ListResponse.verify|verify} messages.
         * @param message ListResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.IListResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified ListResponse message, length delimited. Does not implicitly {@link result.ListResponse.verify|verify} messages.
         * @param message ListResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.IListResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a ListResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns ListResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.ListResponse;

        /**
         * Decodes a ListResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns ListResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.ListResponse;

        /**
         * Verifies a ListResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a ListResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns ListResponse
         */
        public static fromObject(object: { [k: string]: any }): result.ListResponse;

        /**
         * Creates a plain object from a ListResponse message. Also converts values to other types if specified.
         * @param message ListResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.ListResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this ListResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for ListResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an UpdateRequest. */
    interface IUpdateRequest {

        /** UpdateRequest id */
        id?: (string|null);

        /** UpdateRequest config */
        config?: (string|null);

        /** UpdateRequest submitter */
        submitter?: (string|null);

        /** UpdateRequest description */
        description?: (string|null);
    }

    /** Represents an UpdateRequest. */
    class UpdateRequest implements IUpdateRequest {

        /**
         * Constructs a new UpdateRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.IUpdateRequest);

        /** UpdateRequest id. */
        public id: string;

        /** UpdateRequest config. */
        public config: string;

        /** UpdateRequest submitter. */
        public submitter: string;

        /** UpdateRequest description. */
        public description: string;

        /**
         * Creates a new UpdateRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns UpdateRequest instance
         */
        public static create(properties?: result.IUpdateRequest): result.UpdateRequest;

        /**
         * Encodes the specified UpdateRequest message. Does not implicitly {@link result.UpdateRequest.verify|verify} messages.
         * @param message UpdateRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.IUpdateRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified UpdateRequest message, length delimited. Does not implicitly {@link result.UpdateRequest.verify|verify} messages.
         * @param message UpdateRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.IUpdateRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes an UpdateRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns UpdateRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.UpdateRequest;

        /**
         * Decodes an UpdateRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns UpdateRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.UpdateRequest;

        /**
         * Verifies an UpdateRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates an UpdateRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns UpdateRequest
         */
        public static fromObject(object: { [k: string]: any }): result.UpdateRequest;

        /**
         * Creates a plain object from an UpdateRequest message. Also converts values to other types if specified.
         * @param message UpdateRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.UpdateRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this UpdateRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for UpdateRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an UpdateResponse. */
    interface IUpdateResponse {

        /** UpdateResponse id */
        id?: (string|null);
    }

    /** Represents an UpdateResponse. */
    class UpdateResponse implements IUpdateResponse {

        /**
         * Constructs a new UpdateResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.IUpdateResponse);

        /** UpdateResponse id. */
        public id: string;

        /**
         * Creates a new UpdateResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns UpdateResponse instance
         */
        public static create(properties?: result.IUpdateResponse): result.UpdateResponse;

        /**
         * Encodes the specified UpdateResponse message. Does not implicitly {@link result.UpdateResponse.verify|verify} messages.
         * @param message UpdateResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.IUpdateResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified UpdateResponse message, length delimited. Does not implicitly {@link result.UpdateResponse.verify|verify} messages.
         * @param message UpdateResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.IUpdateResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes an UpdateResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns UpdateResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.UpdateResponse;

        /**
         * Decodes an UpdateResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns UpdateResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.UpdateResponse;

        /**
         * Verifies an UpdateResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates an UpdateResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns UpdateResponse
         */
        public static fromObject(object: { [k: string]: any }): result.UpdateResponse;

        /**
         * Creates a plain object from an UpdateResponse message. Also converts values to other types if specified.
         * @param message UpdateResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.UpdateResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this UpdateResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for UpdateResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a RemoveRequest. */
    interface IRemoveRequest {

        /** RemoveRequest id */
        id?: (string|null);

        /** RemoveRequest submitter */
        submitter?: (string|null);
    }

    /** Represents a RemoveRequest. */
    class RemoveRequest implements IRemoveRequest {

        /**
         * Constructs a new RemoveRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.IRemoveRequest);

        /** RemoveRequest id. */
        public id: string;

        /** RemoveRequest submitter. */
        public submitter: string;

        /**
         * Creates a new RemoveRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns RemoveRequest instance
         */
        public static create(properties?: result.IRemoveRequest): result.RemoveRequest;

        /**
         * Encodes the specified RemoveRequest message. Does not implicitly {@link result.RemoveRequest.verify|verify} messages.
         * @param message RemoveRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.IRemoveRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified RemoveRequest message, length delimited. Does not implicitly {@link result.RemoveRequest.verify|verify} messages.
         * @param message RemoveRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.IRemoveRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a RemoveRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns RemoveRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.RemoveRequest;

        /**
         * Decodes a RemoveRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns RemoveRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.RemoveRequest;

        /**
         * Verifies a RemoveRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a RemoveRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns RemoveRequest
         */
        public static fromObject(object: { [k: string]: any }): result.RemoveRequest;

        /**
         * Creates a plain object from a RemoveRequest message. Also converts values to other types if specified.
         * @param message RemoveRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.RemoveRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this RemoveRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for RemoveRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a RemoveResponse. */
    interface IRemoveResponse {

        /** RemoveResponse id */
        id?: (string|null);
    }

    /** Represents a RemoveResponse. */
    class RemoveResponse implements IRemoveResponse {

        /**
         * Constructs a new RemoveResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.IRemoveResponse);

        /** RemoveResponse id. */
        public id: string;

        /**
         * Creates a new RemoveResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns RemoveResponse instance
         */
        public static create(properties?: result.IRemoveResponse): result.RemoveResponse;

        /**
         * Encodes the specified RemoveResponse message. Does not implicitly {@link result.RemoveResponse.verify|verify} messages.
         * @param message RemoveResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.IRemoveResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified RemoveResponse message, length delimited. Does not implicitly {@link result.RemoveResponse.verify|verify} messages.
         * @param message RemoveResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.IRemoveResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a RemoveResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns RemoveResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.RemoveResponse;

        /**
         * Decodes a RemoveResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns RemoveResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.RemoveResponse;

        /**
         * Verifies a RemoveResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a RemoveResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns RemoveResponse
         */
        public static fromObject(object: { [k: string]: any }): result.RemoveResponse;

        /**
         * Creates a plain object from a RemoveResponse message. Also converts values to other types if specified.
         * @param message RemoveResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.RemoveResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this RemoveResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for RemoveResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a RejectRequest. */
    interface IRejectRequest {

        /** RejectRequest id */
        id?: (string|null);

        /** RejectRequest reason */
        reason?: (string|null);
    }

    /** Represents a RejectRequest. */
    class RejectRequest implements IRejectRequest {

        /**
         * Constructs a new RejectRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.IRejectRequest);

        /** RejectRequest id. */
        public id: string;

        /** RejectRequest reason. */
        public reason: string;

        /**
         * Creates a new RejectRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns RejectRequest instance
         */
        public static create(properties?: result.IRejectRequest): result.RejectRequest;

        /**
         * Encodes the specified RejectRequest message. Does not implicitly {@link result.RejectRequest.verify|verify} messages.
         * @param message RejectRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.IRejectRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified RejectRequest message, length delimited. Does not implicitly {@link result.RejectRequest.verify|verify} messages.
         * @param message RejectRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.IRejectRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a RejectRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns RejectRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.RejectRequest;

        /**
         * Decodes a RejectRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns RejectRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.RejectRequest;

        /**
         * Verifies a RejectRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a RejectRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns RejectRequest
         */
        public static fromObject(object: { [k: string]: any }): result.RejectRequest;

        /**
         * Creates a plain object from a RejectRequest message. Also converts values to other types if specified.
         * @param message RejectRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.RejectRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this RejectRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for RejectRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a RejectResponse. */
    interface IRejectResponse {

        /** RejectResponse id */
        id?: (string|null);
    }

    /** Represents a RejectResponse. */
    class RejectResponse implements IRejectResponse {

        /**
         * Constructs a new RejectResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.IRejectResponse);

        /** RejectResponse id. */
        public id: string;

        /**
         * Creates a new RejectResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns RejectResponse instance
         */
        public static create(properties?: result.IRejectResponse): result.RejectResponse;

        /**
         * Encodes the specified RejectResponse message. Does not implicitly {@link result.RejectResponse.verify|verify} messages.
         * @param message RejectResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.IRejectResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified RejectResponse message, length delimited. Does not implicitly {@link result.RejectResponse.verify|verify} messages.
         * @param message RejectResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.IRejectResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a RejectResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns RejectResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.RejectResponse;

        /**
         * Decodes a RejectResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns RejectResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.RejectResponse;

        /**
         * Verifies a RejectResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a RejectResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns RejectResponse
         */
        public static fromObject(object: { [k: string]: any }): result.RejectResponse;

        /**
         * Creates a plain object from a RejectResponse message. Also converts values to other types if specified.
         * @param message RejectResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.RejectResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this RejectResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for RejectResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an ApproveRequest. */
    interface IApproveRequest {

        /** ApproveRequest id */
        id?: (string|null);
    }

    /** Represents an ApproveRequest. */
    class ApproveRequest implements IApproveRequest {

        /**
         * Constructs a new ApproveRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.IApproveRequest);

        /** ApproveRequest id. */
        public id: string;

        /**
         * Creates a new ApproveRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns ApproveRequest instance
         */
        public static create(properties?: result.IApproveRequest): result.ApproveRequest;

        /**
         * Encodes the specified ApproveRequest message. Does not implicitly {@link result.ApproveRequest.verify|verify} messages.
         * @param message ApproveRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.IApproveRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified ApproveRequest message, length delimited. Does not implicitly {@link result.ApproveRequest.verify|verify} messages.
         * @param message ApproveRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.IApproveRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes an ApproveRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns ApproveRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.ApproveRequest;

        /**
         * Decodes an ApproveRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns ApproveRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.ApproveRequest;

        /**
         * Verifies an ApproveRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates an ApproveRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns ApproveRequest
         */
        public static fromObject(object: { [k: string]: any }): result.ApproveRequest;

        /**
         * Creates a plain object from an ApproveRequest message. Also converts values to other types if specified.
         * @param message ApproveRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.ApproveRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this ApproveRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for ApproveRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an ApproveResponse. */
    interface IApproveResponse {

        /** ApproveResponse id */
        id?: (string|null);

        /** ApproveResponse db_key */
        db_key?: (string|null);
    }

    /** Represents an ApproveResponse. */
    class ApproveResponse implements IApproveResponse {

        /**
         * Constructs a new ApproveResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.IApproveResponse);

        /** ApproveResponse id. */
        public id: string;

        /** ApproveResponse db_key. */
        public db_key: string;

        /**
         * Creates a new ApproveResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns ApproveResponse instance
         */
        public static create(properties?: result.IApproveResponse): result.ApproveResponse;

        /**
         * Encodes the specified ApproveResponse message. Does not implicitly {@link result.ApproveResponse.verify|verify} messages.
         * @param message ApproveResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.IApproveResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified ApproveResponse message, length delimited. Does not implicitly {@link result.ApproveResponse.verify|verify} messages.
         * @param message ApproveResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.IApproveResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes an ApproveResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns ApproveResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.ApproveResponse;

        /**
         * Decodes an ApproveResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns ApproveResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.ApproveResponse;

        /**
         * Verifies an ApproveResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates an ApproveResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns ApproveResponse
         */
        public static fromObject(object: { [k: string]: any }): result.ApproveResponse;

        /**
         * Creates a plain object from an ApproveResponse message. Also converts values to other types if specified.
         * @param message ApproveResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.ApproveResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this ApproveResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for ApproveResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a ReplaceRequest. */
    interface IReplaceRequest {

        /** ReplaceRequest id */
        id?: (string|null);

        /** ReplaceRequest db_key */
        db_key?: (string|null);
    }

    /** Represents a ReplaceRequest. */
    class ReplaceRequest implements IReplaceRequest {

        /**
         * Constructs a new ReplaceRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.IReplaceRequest);

        /** ReplaceRequest id. */
        public id: string;

        /** ReplaceRequest db_key. */
        public db_key: string;

        /**
         * Creates a new ReplaceRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns ReplaceRequest instance
         */
        public static create(properties?: result.IReplaceRequest): result.ReplaceRequest;

        /**
         * Encodes the specified ReplaceRequest message. Does not implicitly {@link result.ReplaceRequest.verify|verify} messages.
         * @param message ReplaceRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.IReplaceRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified ReplaceRequest message, length delimited. Does not implicitly {@link result.ReplaceRequest.verify|verify} messages.
         * @param message ReplaceRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.IReplaceRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a ReplaceRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns ReplaceRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.ReplaceRequest;

        /**
         * Decodes a ReplaceRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns ReplaceRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.ReplaceRequest;

        /**
         * Verifies a ReplaceRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a ReplaceRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns ReplaceRequest
         */
        public static fromObject(object: { [k: string]: any }): result.ReplaceRequest;

        /**
         * Creates a plain object from a ReplaceRequest message. Also converts values to other types if specified.
         * @param message ReplaceRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.ReplaceRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this ReplaceRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for ReplaceRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a ReplaceResponse. */
    interface IReplaceResponse {

        /** ReplaceResponse id */
        id?: (string|null);

        /** ReplaceResponse db_key */
        db_key?: (string|null);
    }

    /** Represents a ReplaceResponse. */
    class ReplaceResponse implements IReplaceResponse {

        /**
         * Constructs a new ReplaceResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.IReplaceResponse);

        /** ReplaceResponse id. */
        public id: string;

        /** ReplaceResponse db_key. */
        public db_key: string;

        /**
         * Creates a new ReplaceResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns ReplaceResponse instance
         */
        public static create(properties?: result.IReplaceResponse): result.ReplaceResponse;

        /**
         * Encodes the specified ReplaceResponse message. Does not implicitly {@link result.ReplaceResponse.verify|verify} messages.
         * @param message ReplaceResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.IReplaceResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified ReplaceResponse message, length delimited. Does not implicitly {@link result.ReplaceResponse.verify|verify} messages.
         * @param message ReplaceResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.IReplaceResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a ReplaceResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns ReplaceResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.ReplaceResponse;

        /**
         * Decodes a ReplaceResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns ReplaceResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.ReplaceResponse;

        /**
         * Verifies a ReplaceResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a ReplaceResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns ReplaceResponse
         */
        public static fromObject(object: { [k: string]: any }): result.ReplaceResponse;

        /**
         * Creates a plain object from a ReplaceResponse message. Also converts values to other types if specified.
         * @param message ReplaceResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.ReplaceResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this ReplaceResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for ReplaceResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Represents a Compute */
    class Compute extends $protobuf.rpc.Service {

        /**
         * Constructs a new Compute service.
         * @param rpcImpl RPC implementation
         * @param [requestDelimited=false] Whether requests are length-delimited
         * @param [responseDelimited=false] Whether responses are length-delimited
         */
        constructor(rpcImpl: $protobuf.RPCImpl, requestDelimited?: boolean, responseDelimited?: boolean);

        /**
         * Creates new Compute service using the specified rpc implementation.
         * @param rpcImpl RPC implementation
         * @param [requestDelimited=false] Whether requests are length-delimited
         * @param [responseDelimited=false] Whether responses are length-delimited
         * @returns RPC service. Useful where requests and/or responses are streamed.
         */
        public static create(rpcImpl: $protobuf.RPCImpl, requestDelimited?: boolean, responseDelimited?: boolean): Compute;

        /**
         * Calls Run.
         * @param request RunRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and RunResponse
         */
        public run(request: result.IRunRequest, callback: result.Compute.RunCallback): void;

        /**
         * Calls Run.
         * @param request RunRequest message or plain object
         * @returns Promise
         */
        public run(request: result.IRunRequest): Promise<result.RunResponse>;
    }

    namespace Compute {

        /**
         * Callback as used by {@link result.Compute#run}.
         * @param error Error, if any
         * @param [response] RunResponse
         */
        type RunCallback = (error: (Error|null), response?: result.RunResponse) => void;
    }

    /** Properties of a RunRequest. */
    interface IRunRequest {

        /** RunRequest key */
        key?: (string|null);

        /** RunRequest config */
        config?: (string|null);

        /** RunRequest api_key */
        api_key?: (string|null);

        /** RunRequest callback_url */
        callback_url?: (string|null);
    }

    /** Represents a RunRequest. */
    class RunRequest implements IRunRequest {

        /**
         * Constructs a new RunRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.IRunRequest);

        /** RunRequest key. */
        public key: string;

        /** RunRequest config. */
        public config: string;

        /** RunRequest api_key. */
        public api_key: string;

        /** RunRequest callback_url. */
        public callback_url: string;

        /**
         * Creates a new RunRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns RunRequest instance
         */
        public static create(properties?: result.IRunRequest): result.RunRequest;

        /**
         * Encodes the specified RunRequest message. Does not implicitly {@link result.RunRequest.verify|verify} messages.
         * @param message RunRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.IRunRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified RunRequest message, length delimited. Does not implicitly {@link result.RunRequest.verify|verify} messages.
         * @param message RunRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.IRunRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a RunRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns RunRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.RunRequest;

        /**
         * Decodes a RunRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns RunRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.RunRequest;

        /**
         * Verifies a RunRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a RunRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns RunRequest
         */
        public static fromObject(object: { [k: string]: any }): result.RunRequest;

        /**
         * Creates a plain object from a RunRequest message. Also converts values to other types if specified.
         * @param message RunRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.RunRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this RunRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for RunRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a RunResponse. */
    interface IRunResponse {
    }

    /** Represents a RunResponse. */
    class RunResponse implements IRunResponse {

        /**
         * Constructs a new RunResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: result.IRunResponse);

        /**
         * Creates a new RunResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns RunResponse instance
         */
        public static create(properties?: result.IRunResponse): result.RunResponse;

        /**
         * Encodes the specified RunResponse message. Does not implicitly {@link result.RunResponse.verify|verify} messages.
         * @param message RunResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: result.IRunResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified RunResponse message, length delimited. Does not implicitly {@link result.RunResponse.verify|verify} messages.
         * @param message RunResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: result.IRunResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a RunResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns RunResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): result.RunResponse;

        /**
         * Decodes a RunResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns RunResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): result.RunResponse;

        /**
         * Verifies a RunResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a RunResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns RunResponse
         */
        public static fromObject(object: { [k: string]: any }): result.RunResponse;

        /**
         * Creates a plain object from a RunResponse message. Also converts values to other types if specified.
         * @param message RunResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: result.RunResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this RunResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for RunResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }
}
