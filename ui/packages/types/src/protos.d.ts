import * as $protobuf from "protobufjs";
import Long = require("long");
/** Namespace db. */
export namespace db {

    /** Properties of an Entry. */
    interface IEntry {

        /** Entry id */
        id?: (string|null);

        /** Entry create_date */
        create_date?: (number|Long|null);

        /** Entry config */
        config?: (string|null);

        /** Entry description */
        description?: (string|null);

        /** Entry submitter */
        submitter?: (string|null);

        /** Entry accepted_tags */
        accepted_tags?: (model.DBTag[]|null);

        /** Entry rejected_tags */
        rejected_tags?: (model.DBTag[]|null);

        /** Entry is_db_valid */
        is_db_valid?: (boolean|null);

        /** Entry share_key */
        share_key?: (string|null);

        /** Entry last_update */
        last_update?: (number|Long|null);

        /** Entry hash */
        hash?: (string|null);

        /** Entry summary */
        summary?: (db.IEntrySummary|null);
    }

    /** Represents an Entry. */
    class Entry implements IEntry {

        /**
         * Constructs a new Entry.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IEntry);

        /** Entry id. */
        public id: string;

        /** Entry create_date. */
        public create_date: (number|Long);

        /** Entry config. */
        public config: string;

        /** Entry description. */
        public description: string;

        /** Entry submitter. */
        public submitter: string;

        /** Entry accepted_tags. */
        public accepted_tags: model.DBTag[];

        /** Entry rejected_tags. */
        public rejected_tags: model.DBTag[];

        /** Entry is_db_valid. */
        public is_db_valid: boolean;

        /** Entry share_key. */
        public share_key: string;

        /** Entry last_update. */
        public last_update: (number|Long);

        /** Entry hash. */
        public hash: string;

        /** Entry summary. */
        public summary?: (db.IEntrySummary|null);

        /**
         * Gets the default type url for Entry
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an EntrySummary. */
    interface IEntrySummary {

        /** EntrySummary sim_duration */
        sim_duration?: (model.IDescriptiveStats|null);

        /** EntrySummary mode */
        mode?: (model.SimMode|null);

        /** EntrySummary total_damage */
        total_damage?: (model.IDescriptiveStats|null);

        /** EntrySummary char_names */
        char_names?: (string[]|null);

        /** EntrySummary target_count */
        target_count?: (number|null);

        /** EntrySummary mean_dps_per_target */
        mean_dps_per_target?: (number|null);

        /** EntrySummary team */
        team?: (model.ICharacter[]|null);

        /** EntrySummary dps_by_target */
        dps_by_target?: ({ [k: string]: model.IDescriptiveStats }|null);
    }

    /** Represents an EntrySummary. */
    class EntrySummary implements IEntrySummary {

        /**
         * Constructs a new EntrySummary.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IEntrySummary);

        /** EntrySummary sim_duration. */
        public sim_duration?: (model.IDescriptiveStats|null);

        /** EntrySummary mode. */
        public mode: model.SimMode;

        /** EntrySummary total_damage. */
        public total_damage?: (model.IDescriptiveStats|null);

        /** EntrySummary char_names. */
        public char_names: string[];

        /** EntrySummary target_count. */
        public target_count: number;

        /** EntrySummary mean_dps_per_target. */
        public mean_dps_per_target: number;

        /** EntrySummary team. */
        public team: model.ICharacter[];

        /** EntrySummary dps_by_target. */
        public dps_by_target: { [k: string]: model.IDescriptiveStats };

        /**
         * Gets the default type url for EntrySummary
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an Entries. */
    interface IEntries {

        /** Entries data */
        data?: (db.IEntry[]|null);
    }

    /** Represents an Entries. */
    class Entries implements IEntries {

        /**
         * Constructs a new Entries.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IEntries);

        /** Entries data. */
        public data: db.IEntry[];

        /**
         * Gets the default type url for Entries
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a QueryOpt. */
    interface IQueryOpt {

        /** QueryOpt query */
        query?: (google.protobuf.IStruct|null);

        /** QueryOpt sort */
        sort?: (google.protobuf.IStruct|null);

        /** QueryOpt project */
        project?: (google.protobuf.IStruct|null);

        /** QueryOpt skip */
        skip?: (number|Long|null);

        /** QueryOpt limit */
        limit?: (number|Long|null);
    }

    /** Represents a QueryOpt. */
    class QueryOpt implements IQueryOpt {

        /**
         * Constructs a new QueryOpt.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IQueryOpt);

        /** QueryOpt query. */
        public query?: (google.protobuf.IStruct|null);

        /** QueryOpt sort. */
        public sort?: (google.protobuf.IStruct|null);

        /** QueryOpt project. */
        public project?: (google.protobuf.IStruct|null);

        /** QueryOpt skip. */
        public skip: (number|Long);

        /** QueryOpt limit. */
        public limit: (number|Long);

        /**
         * Gets the default type url for QueryOpt
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a ComputeWork. */
    interface IComputeWork {

        /** ComputeWork id */
        id?: (string|null);

        /** ComputeWork config */
        config?: (string|null);

        /** ComputeWork iterations */
        iterations?: (number|null);
    }

    /** Represents a ComputeWork. */
    class ComputeWork implements IComputeWork {

        /**
         * Constructs a new ComputeWork.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IComputeWork);

        /** ComputeWork id. */
        public id: string;

        /** ComputeWork config. */
        public config: string;

        /** ComputeWork iterations. */
        public iterations: number;

        /**
         * Gets the default type url for ComputeWork
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetRequest. */
    interface IGetRequest {

        /** GetRequest query */
        query?: (db.IQueryOpt|null);
    }

    /** Represents a GetRequest. */
    class GetRequest implements IGetRequest {

        /**
         * Constructs a new GetRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IGetRequest);

        /** GetRequest query. */
        public query?: (db.IQueryOpt|null);

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
        data?: (db.IEntries|null);
    }

    /** Represents a GetResponse. */
    class GetResponse implements IGetResponse {

        /**
         * Constructs a new GetResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IGetResponse);

        /** GetResponse data. */
        public data?: (db.IEntries|null);

        /**
         * Gets the default type url for GetResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetOneRequest. */
    interface IGetOneRequest {

        /** GetOneRequest id */
        id?: (string|null);
    }

    /** Represents a GetOneRequest. */
    class GetOneRequest implements IGetOneRequest {

        /**
         * Constructs a new GetOneRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IGetOneRequest);

        /** GetOneRequest id. */
        public id: string;

        /**
         * Gets the default type url for GetOneRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetOneResponse. */
    interface IGetOneResponse {

        /** GetOneResponse data */
        data?: (db.IEntry|null);
    }

    /** Represents a GetOneResponse. */
    class GetOneResponse implements IGetOneResponse {

        /**
         * Constructs a new GetOneResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IGetOneResponse);

        /** GetOneResponse data. */
        public data?: (db.IEntry|null);

        /**
         * Gets the default type url for GetOneResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetPendingRequest. */
    interface IGetPendingRequest {

        /** GetPendingRequest tag */
        tag?: (model.DBTag|null);

        /** GetPendingRequest query */
        query?: (db.IQueryOpt|null);
    }

    /** Represents a GetPendingRequest. */
    class GetPendingRequest implements IGetPendingRequest {

        /**
         * Constructs a new GetPendingRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IGetPendingRequest);

        /** GetPendingRequest tag. */
        public tag: model.DBTag;

        /** GetPendingRequest query. */
        public query?: (db.IQueryOpt|null);

        /**
         * Gets the default type url for GetPendingRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetPendingResponse. */
    interface IGetPendingResponse {

        /** GetPendingResponse data */
        data?: (db.IEntries|null);
    }

    /** Represents a GetPendingResponse. */
    class GetPendingResponse implements IGetPendingResponse {

        /**
         * Constructs a new GetPendingResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IGetPendingResponse);

        /** GetPendingResponse data. */
        public data?: (db.IEntries|null);

        /**
         * Gets the default type url for GetPendingResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetBySubmitterRequest. */
    interface IGetBySubmitterRequest {

        /** GetBySubmitterRequest submitter */
        submitter?: (string|null);

        /** GetBySubmitterRequest query */
        query?: (db.IQueryOpt|null);
    }

    /** Represents a GetBySubmitterRequest. */
    class GetBySubmitterRequest implements IGetBySubmitterRequest {

        /**
         * Constructs a new GetBySubmitterRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IGetBySubmitterRequest);

        /** GetBySubmitterRequest submitter. */
        public submitter: string;

        /** GetBySubmitterRequest query. */
        public query?: (db.IQueryOpt|null);

        /**
         * Gets the default type url for GetBySubmitterRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetBySubmitterResponse. */
    interface IGetBySubmitterResponse {

        /** GetBySubmitterResponse data */
        data?: (db.IEntries|null);
    }

    /** Represents a GetBySubmitterResponse. */
    class GetBySubmitterResponse implements IGetBySubmitterResponse {

        /**
         * Constructs a new GetBySubmitterResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IGetBySubmitterResponse);

        /** GetBySubmitterResponse data. */
        public data?: (db.IEntries|null);

        /**
         * Gets the default type url for GetBySubmitterResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetAllRequest. */
    interface IGetAllRequest {

        /** GetAllRequest query */
        query?: (db.IQueryOpt|null);
    }

    /** Represents a GetAllRequest. */
    class GetAllRequest implements IGetAllRequest {

        /**
         * Constructs a new GetAllRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IGetAllRequest);

        /** GetAllRequest query. */
        public query?: (db.IQueryOpt|null);

        /**
         * Gets the default type url for GetAllRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetAllResponse. */
    interface IGetAllResponse {

        /** GetAllResponse data */
        data?: (db.IEntries|null);
    }

    /** Represents a GetAllResponse. */
    class GetAllResponse implements IGetAllResponse {

        /**
         * Constructs a new GetAllResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IGetAllResponse);

        /** GetAllResponse data. */
        public data?: (db.IEntries|null);

        /**
         * Gets the default type url for GetAllResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an ApproveTagRequest. */
    interface IApproveTagRequest {

        /** ApproveTagRequest id */
        id?: (string|null);

        /** ApproveTagRequest tag */
        tag?: (model.DBTag|null);
    }

    /** Represents an ApproveTagRequest. */
    class ApproveTagRequest implements IApproveTagRequest {

        /**
         * Constructs a new ApproveTagRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IApproveTagRequest);

        /** ApproveTagRequest id. */
        public id: string;

        /** ApproveTagRequest tag. */
        public tag: model.DBTag;

        /**
         * Gets the default type url for ApproveTagRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an ApproveTagResponse. */
    interface IApproveTagResponse {

        /** ApproveTagResponse id */
        id?: (string|null);
    }

    /** Represents an ApproveTagResponse. */
    class ApproveTagResponse implements IApproveTagResponse {

        /**
         * Constructs a new ApproveTagResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IApproveTagResponse);

        /** ApproveTagResponse id. */
        public id: string;

        /**
         * Gets the default type url for ApproveTagResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a RejectTagRequest. */
    interface IRejectTagRequest {

        /** RejectTagRequest id */
        id?: (string|null);

        /** RejectTagRequest tag */
        tag?: (model.DBTag|null);
    }

    /** Represents a RejectTagRequest. */
    class RejectTagRequest implements IRejectTagRequest {

        /**
         * Constructs a new RejectTagRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IRejectTagRequest);

        /** RejectTagRequest id. */
        public id: string;

        /** RejectTagRequest tag. */
        public tag: model.DBTag;

        /**
         * Gets the default type url for RejectTagRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a RejectTagResponse. */
    interface IRejectTagResponse {

        /** RejectTagResponse id */
        id?: (string|null);
    }

    /** Represents a RejectTagResponse. */
    class RejectTagResponse implements IRejectTagResponse {

        /**
         * Constructs a new RejectTagResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IRejectTagResponse);

        /** RejectTagResponse id. */
        public id: string;

        /**
         * Gets the default type url for RejectTagResponse
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
        constructor(properties?: db.ISubmitRequest);

        /** SubmitRequest config. */
        public config: string;

        /** SubmitRequest submitter. */
        public submitter: string;

        /** SubmitRequest description. */
        public description: string;

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
        constructor(properties?: db.ISubmitResponse);

        /** SubmitResponse id. */
        public id: string;

        /**
         * Gets the default type url for SubmitResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a DeletePendingRequest. */
    interface IDeletePendingRequest {

        /** DeletePendingRequest id */
        id?: (string|null);

        /** DeletePendingRequest sender */
        sender?: (string|null);
    }

    /** Represents a DeletePendingRequest. */
    class DeletePendingRequest implements IDeletePendingRequest {

        /**
         * Constructs a new DeletePendingRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IDeletePendingRequest);

        /** DeletePendingRequest id. */
        public id: string;

        /** DeletePendingRequest sender. */
        public sender: string;

        /**
         * Gets the default type url for DeletePendingRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a DeletePendingResponse. */
    interface IDeletePendingResponse {

        /** DeletePendingResponse id */
        id?: (string|null);
    }

    /** Represents a DeletePendingResponse. */
    class DeletePendingResponse implements IDeletePendingResponse {

        /**
         * Constructs a new DeletePendingResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IDeletePendingResponse);

        /** DeletePendingResponse id. */
        public id: string;

        /**
         * Gets the default type url for DeletePendingResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetWorkRequest. */
    interface IGetWorkRequest {
    }

    /** Represents a GetWorkRequest. */
    class GetWorkRequest implements IGetWorkRequest {

        /**
         * Constructs a new GetWorkRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IGetWorkRequest);

        /**
         * Gets the default type url for GetWorkRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetWorkResponse. */
    interface IGetWorkResponse {

        /** GetWorkResponse data */
        data?: (db.IComputeWork[]|null);
    }

    /** Represents a GetWorkResponse. */
    class GetWorkResponse implements IGetWorkResponse {

        /**
         * Constructs a new GetWorkResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IGetWorkResponse);

        /** GetWorkResponse data. */
        public data: db.IComputeWork[];

        /**
         * Gets the default type url for GetWorkResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a RejectWorkRequest. */
    interface IRejectWorkRequest {

        /** RejectWorkRequest id */
        id?: (string|null);

        /** RejectWorkRequest reason */
        reason?: (string|null);

        /** RejectWorkRequest hash */
        hash?: (string|null);
    }

    /** Represents a RejectWorkRequest. */
    class RejectWorkRequest implements IRejectWorkRequest {

        /**
         * Constructs a new RejectWorkRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IRejectWorkRequest);

        /** RejectWorkRequest id. */
        public id: string;

        /** RejectWorkRequest reason. */
        public reason: string;

        /** RejectWorkRequest hash. */
        public hash: string;

        /**
         * Gets the default type url for RejectWorkRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a RejectWorkResponse. */
    interface IRejectWorkResponse {
    }

    /** Represents a RejectWorkResponse. */
    class RejectWorkResponse implements IRejectWorkResponse {

        /**
         * Constructs a new RejectWorkResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.IRejectWorkResponse);

        /**
         * Gets the default type url for RejectWorkResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a CompleteWorkRequest. */
    interface ICompleteWorkRequest {

        /** CompleteWorkRequest id */
        id?: (string|null);

        /** CompleteWorkRequest result */
        result?: (model.ISimulationResult|null);
    }

    /** Represents a CompleteWorkRequest. */
    class CompleteWorkRequest implements ICompleteWorkRequest {

        /**
         * Constructs a new CompleteWorkRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.ICompleteWorkRequest);

        /** CompleteWorkRequest id. */
        public id: string;

        /** CompleteWorkRequest result. */
        public result?: (model.ISimulationResult|null);

        /**
         * Gets the default type url for CompleteWorkRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a CompleteWorkResponse. */
    interface ICompleteWorkResponse {

        /** CompleteWorkResponse id */
        id?: (string|null);
    }

    /** Represents a CompleteWorkResponse. */
    class CompleteWorkResponse implements ICompleteWorkResponse {

        /**
         * Constructs a new CompleteWorkResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: db.ICompleteWorkResponse);

        /** CompleteWorkResponse id. */
        public id: string;

        /**
         * Gets the default type url for CompleteWorkResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }
}

/** Namespace preview. */
export namespace preview {

    /** Properties of a GetRequest. */
    interface IGetRequest {

        /** GetRequest id */
        id?: (string|null);

        /** GetRequest data */
        data?: (model.ISimulationResult|null);
    }

    /** Represents a GetRequest. */
    class GetRequest implements IGetRequest {

        /**
         * Constructs a new GetRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: preview.IGetRequest);

        /** GetRequest id. */
        public id: string;

        /** GetRequest data. */
        public data?: (model.ISimulationResult|null);

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
        data?: (Uint8Array|null);
    }

    /** Represents a GetResponse. */
    class GetResponse implements IGetResponse {

        /**
         * Constructs a new GetResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: preview.IGetResponse);

        /** GetResponse data. */
        public data: Uint8Array;

        /**
         * Gets the default type url for GetResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }
}

/** Namespace share. */
export namespace share {

    /** Properties of a ShareEntry. */
    interface IShareEntry {

        /** ShareEntry id */
        id?: (string|null);

        /** ShareEntry result */
        result?: (model.ISimulationResult|null);

        /** ShareEntry expires_at */
        expires_at?: (number|Long|null);

        /** ShareEntry submitter */
        submitter?: (string|null);
    }

    /** Represents a ShareEntry. */
    class ShareEntry implements IShareEntry {

        /**
         * Constructs a new ShareEntry.
         * @param [properties] Properties to set
         */
        constructor(properties?: share.IShareEntry);

        /** ShareEntry id. */
        public id: string;

        /** ShareEntry result. */
        public result?: (model.ISimulationResult|null);

        /** ShareEntry expires_at. */
        public expires_at: (number|Long);

        /** ShareEntry submitter. */
        public submitter: string;

        /**
         * Gets the default type url for ShareEntry
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a CreateRequest. */
    interface ICreateRequest {

        /** CreateRequest result */
        result?: (model.ISimulationResult|null);

        /** CreateRequest expires_at */
        expires_at?: (number|Long|null);

        /** CreateRequest submitter */
        submitter?: (string|null);
    }

    /** Represents a CreateRequest. */
    class CreateRequest implements ICreateRequest {

        /**
         * Constructs a new CreateRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: share.ICreateRequest);

        /** CreateRequest result. */
        public result?: (model.ISimulationResult|null);

        /** CreateRequest expires_at. */
        public expires_at: (number|Long);

        /** CreateRequest submitter. */
        public submitter: string;

        /**
         * Gets the default type url for CreateRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a CreateResponse. */
    interface ICreateResponse {

        /** CreateResponse id */
        id?: (string|null);
    }

    /** Represents a CreateResponse. */
    class CreateResponse implements ICreateResponse {

        /**
         * Constructs a new CreateResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: share.ICreateResponse);

        /** CreateResponse id. */
        public id: string;

        /**
         * Gets the default type url for CreateResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a ReadRequest. */
    interface IReadRequest {

        /** ReadRequest id */
        id?: (string|null);
    }

    /** Represents a ReadRequest. */
    class ReadRequest implements IReadRequest {

        /**
         * Constructs a new ReadRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: share.IReadRequest);

        /** ReadRequest id. */
        public id: string;

        /**
         * Gets the default type url for ReadRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a ReadResponse. */
    interface IReadResponse {

        /** ReadResponse id */
        id?: (string|null);

        /** ReadResponse result */
        result?: (model.ISimulationResult|null);

        /** ReadResponse expires_at */
        expires_at?: (number|Long|null);
    }

    /** Represents a ReadResponse. */
    class ReadResponse implements IReadResponse {

        /**
         * Constructs a new ReadResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: share.IReadResponse);

        /** ReadResponse id. */
        public id: string;

        /** ReadResponse result. */
        public result?: (model.ISimulationResult|null);

        /** ReadResponse expires_at. */
        public expires_at: (number|Long);

        /**
         * Gets the default type url for ReadResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an UpdateRequest. */
    interface IUpdateRequest {

        /** UpdateRequest id */
        id?: (string|null);

        /** UpdateRequest result */
        result?: (model.ISimulationResult|null);

        /** UpdateRequest expires_at */
        expires_at?: (number|Long|null);

        /** UpdateRequest submitter */
        submitter?: (string|null);
    }

    /** Represents an UpdateRequest. */
    class UpdateRequest implements IUpdateRequest {

        /**
         * Constructs a new UpdateRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: share.IUpdateRequest);

        /** UpdateRequest id. */
        public id: string;

        /** UpdateRequest result. */
        public result?: (model.ISimulationResult|null);

        /** UpdateRequest expires_at. */
        public expires_at: (number|Long);

        /** UpdateRequest submitter. */
        public submitter: string;

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
        constructor(properties?: share.IUpdateResponse);

        /** UpdateResponse id. */
        public id: string;

        /**
         * Gets the default type url for UpdateResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a SetTTLRequest. */
    interface ISetTTLRequest {

        /** SetTTLRequest id */
        id?: (string|null);

        /** SetTTLRequest expires_at */
        expires_at?: (number|Long|null);
    }

    /** Represents a SetTTLRequest. */
    class SetTTLRequest implements ISetTTLRequest {

        /**
         * Constructs a new SetTTLRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: share.ISetTTLRequest);

        /** SetTTLRequest id. */
        public id: string;

        /** SetTTLRequest expires_at. */
        public expires_at: (number|Long);

        /**
         * Gets the default type url for SetTTLRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a SetTTLResponse. */
    interface ISetTTLResponse {

        /** SetTTLResponse id */
        id?: (string|null);
    }

    /** Represents a SetTTLResponse. */
    class SetTTLResponse implements ISetTTLResponse {

        /**
         * Constructs a new SetTTLResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: share.ISetTTLResponse);

        /** SetTTLResponse id. */
        public id: string;

        /**
         * Gets the default type url for SetTTLResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a DeleteRequest. */
    interface IDeleteRequest {

        /** DeleteRequest id */
        id?: (string|null);
    }

    /** Represents a DeleteRequest. */
    class DeleteRequest implements IDeleteRequest {

        /**
         * Constructs a new DeleteRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: share.IDeleteRequest);

        /** DeleteRequest id. */
        public id: string;

        /**
         * Gets the default type url for DeleteRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a DeleteResponse. */
    interface IDeleteResponse {

        /** DeleteResponse id */
        id?: (string|null);
    }

    /** Represents a DeleteResponse. */
    class DeleteResponse implements IDeleteResponse {

        /**
         * Constructs a new DeleteResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: share.IDeleteResponse);

        /** DeleteResponse id. */
        public id: string;

        /**
         * Gets the default type url for DeleteResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a RandomRequest. */
    interface IRandomRequest {
    }

    /** Represents a RandomRequest. */
    class RandomRequest implements IRandomRequest {

        /**
         * Constructs a new RandomRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: share.IRandomRequest);

        /**
         * Gets the default type url for RandomRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a RandomResponse. */
    interface IRandomResponse {

        /** RandomResponse id */
        id?: (string|null);
    }

    /** Represents a RandomResponse. */
    class RandomResponse implements IRandomResponse {

        /**
         * Constructs a new RandomResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: share.IRandomResponse);

        /** RandomResponse id. */
        public id: string;

        /**
         * Gets the default type url for RandomResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }
}

/** Namespace model. */
export namespace model {

    /** Properties of an AvatarDataMap. */
    interface IAvatarDataMap {

        /** AvatarDataMap data */
        data?: ({ [k: string]: model.IAvatarData }|null);
    }

    /** Represents an AvatarDataMap. */
    class AvatarDataMap implements IAvatarDataMap {

        /**
         * Constructs a new AvatarDataMap.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IAvatarDataMap);

        /** AvatarDataMap data. */
        public data: { [k: string]: model.IAvatarData };

        /**
         * Gets the default type url for AvatarDataMap
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an AvatarData. */
    interface IAvatarData {

        /** AvatarData id */
        id?: (number|null);

        /** AvatarData sub_id */
        sub_id?: (number|null);

        /** AvatarData key */
        key?: (string|null);

        /** AvatarData rarity */
        rarity?: (model.QualityType|null);

        /** AvatarData body */
        body?: (model.BodyType|null);

        /** AvatarData region */
        region?: (model.ZoneType|null);

        /** AvatarData element */
        element?: (model.Element|null);

        /** AvatarData weapon_class */
        weapon_class?: (model.WeaponClass|null);

        /** AvatarData icon_name */
        icon_name?: (string|null);

        /** AvatarData stats */
        stats?: (model.IAvatarStatsData|null);

        /** AvatarData skill_details */
        skill_details?: (model.IAvatarSkillsData|null);
    }

    /** Represents an AvatarData. */
    class AvatarData implements IAvatarData {

        /**
         * Constructs a new AvatarData.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IAvatarData);

        /** AvatarData id. */
        public id: number;

        /** AvatarData sub_id. */
        public sub_id: number;

        /** AvatarData key. */
        public key: string;

        /** AvatarData rarity. */
        public rarity: model.QualityType;

        /** AvatarData body. */
        public body: model.BodyType;

        /** AvatarData region. */
        public region: model.ZoneType;

        /** AvatarData element. */
        public element: model.Element;

        /** AvatarData weapon_class. */
        public weapon_class: model.WeaponClass;

        /** AvatarData icon_name. */
        public icon_name: string;

        /** AvatarData stats. */
        public stats?: (model.IAvatarStatsData|null);

        /** AvatarData skill_details. */
        public skill_details?: (model.IAvatarSkillsData|null);

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

        /** AvatarStatsData promo_data. */
        public promo_data: model.IPromotionData[];

        /**
         * Gets the default type url for AvatarStatsData
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an AvatarSkillsData. */
    interface IAvatarSkillsData {

        /** AvatarSkillsData skill */
        skill?: (number|null);

        /** AvatarSkillsData burst */
        burst?: (number|null);

        /** AvatarSkillsData attack */
        attack?: (number|null);

        /** AvatarSkillsData burst_energy_cost */
        burst_energy_cost?: (number|null);
    }

    /** Represents an AvatarSkillsData. */
    class AvatarSkillsData implements IAvatarSkillsData {

        /**
         * Constructs a new AvatarSkillsData.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IAvatarSkillsData);

        /** AvatarSkillsData skill. */
        public skill: number;

        /** AvatarSkillsData burst. */
        public burst: number;

        /** AvatarSkillsData attack. */
        public attack: number;

        /** AvatarSkillsData burst_energy_cost. */
        public burst_energy_cost: number;

        /**
         * Gets the default type url for AvatarSkillsData
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a WeaponDataMap. */
    interface IWeaponDataMap {

        /** WeaponDataMap data */
        data?: ({ [k: string]: model.IWeaponData }|null);
    }

    /** Represents a WeaponDataMap. */
    class WeaponDataMap implements IWeaponDataMap {

        /**
         * Constructs a new WeaponDataMap.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IWeaponDataMap);

        /** WeaponDataMap data. */
        public data: { [k: string]: model.IWeaponData };

        /**
         * Gets the default type url for WeaponDataMap
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a WeaponData. */
    interface IWeaponData {

        /** WeaponData id */
        id?: (number|null);

        /** WeaponData key */
        key?: (string|null);

        /** WeaponData rarity */
        rarity?: (number|null);

        /** WeaponData weapon_class */
        weapon_class?: (model.WeaponClass|null);

        /** WeaponData image_name */
        image_name?: (string|null);

        /** WeaponData base_stats */
        base_stats?: (model.IWeaponStatsData|null);
    }

    /** Represents a WeaponData. */
    class WeaponData implements IWeaponData {

        /**
         * Constructs a new WeaponData.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IWeaponData);

        /** WeaponData id. */
        public id: number;

        /** WeaponData key. */
        public key: string;

        /** WeaponData rarity. */
        public rarity: number;

        /** WeaponData weapon_class. */
        public weapon_class: model.WeaponClass;

        /** WeaponData image_name. */
        public image_name: string;

        /** WeaponData base_stats. */
        public base_stats?: (model.IWeaponStatsData|null);

        /**
         * Gets the default type url for WeaponData
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a WeaponStatsData. */
    interface IWeaponStatsData {

        /** WeaponStatsData base_props */
        base_props?: (model.IWeaponProp[]|null);

        /** WeaponStatsData promo_data */
        promo_data?: (model.IPromotionData[]|null);
    }

    /** Represents a WeaponStatsData. */
    class WeaponStatsData implements IWeaponStatsData {

        /**
         * Constructs a new WeaponStatsData.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IWeaponStatsData);

        /** WeaponStatsData base_props. */
        public base_props: model.IWeaponProp[];

        /** WeaponStatsData promo_data. */
        public promo_data: model.IPromotionData[];

        /**
         * Gets the default type url for WeaponStatsData
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a WeaponProp. */
    interface IWeaponProp {

        /** WeaponProp prop_type */
        prop_type?: (model.StatType|null);

        /** WeaponProp initial_value */
        initial_value?: (number|null);

        /** WeaponProp curve */
        curve?: (model.WeaponCurveType|null);
    }

    /** Represents a WeaponProp. */
    class WeaponProp implements IWeaponProp {

        /**
         * Constructs a new WeaponProp.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IWeaponProp);

        /** WeaponProp prop_type. */
        public prop_type: model.StatType;

        /** WeaponProp initial_value. */
        public initial_value: number;

        /** WeaponProp curve. */
        public curve: model.WeaponCurveType;

        /**
         * Gets the default type url for WeaponProp
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an ArtifactData. */
    interface IArtifactData {

        /** ArtifactData id */
        id?: (number|null);
    }

    /** Represents an ArtifactData. */
    class ArtifactData implements IArtifactData {

        /**
         * Constructs a new ArtifactData.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IArtifactData);

        /** ArtifactData id. */
        public id: number;

        /**
         * Gets the default type url for ArtifactData
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a PromotionData. */
    interface IPromotionData {

        /** PromotionData max_level */
        max_level?: (number|null);

        /** PromotionData add_props */
        add_props?: (model.IPromotionAddProp[]|null);
    }

    /** Represents a PromotionData. */
    class PromotionData implements IPromotionData {

        /**
         * Constructs a new PromotionData.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IPromotionData);

        /** PromotionData max_level. */
        public max_level: number;

        /** PromotionData add_props. */
        public add_props: model.IPromotionAddProp[];

        /**
         * Gets the default type url for PromotionData
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a PromotionAddProp. */
    interface IPromotionAddProp {

        /** PromotionAddProp prop_type */
        prop_type?: (model.StatType|null);

        /** PromotionAddProp value */
        value?: (number|null);
    }

    /** Represents a PromotionAddProp. */
    class PromotionAddProp implements IPromotionAddProp {

        /**
         * Constructs a new PromotionAddProp.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IPromotionAddProp);

        /** PromotionAddProp prop_type. */
        public prop_type: model.StatType;

        /** PromotionAddProp value. */
        public value: number;

        /**
         * Gets the default type url for PromotionAddProp
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** AvatarCurveType enum. */
    enum AvatarCurveType {
        INVALID_AVATAR_CURVE = 0,
        GROW_CURVE_HP_S4 = 1,
        GROW_CURVE_ATTACK_S4 = 2,
        GROW_CURVE_HP_S5 = 3,
        GROW_CURVE_ATTACK_S5 = 4
    }

    /** QualityType enum. */
    enum QualityType {
        INVALID_QUALITY_TYPE = 0,
        QUALITY_ORANGE_SP = 6,
        QUALITY_ORANGE = 5,
        QUALITY_PURPLE = 4
    }

    /** WeaponCurveType enum. */
    enum WeaponCurveType {
        INVALID_WEAPON_CURVE = 0,
        GROW_CURVE_ATTACK_101 = 1,
        GROW_CURVE_ATTACK_102 = 2,
        GROW_CURVE_ATTACK_103 = 3,
        GROW_CURVE_ATTACK_104 = 4,
        GROW_CURVE_ATTACK_105 = 5,
        GROW_CURVE_CRITICAL_101 = 6,
        GROW_CURVE_ATTACK_201 = 7,
        GROW_CURVE_ATTACK_202 = 8,
        GROW_CURVE_ATTACK_203 = 9,
        GROW_CURVE_ATTACK_204 = 10,
        GROW_CURVE_ATTACK_205 = 11,
        GROW_CURVE_CRITICAL_201 = 12,
        GROW_CURVE_ATTACK_301 = 13,
        GROW_CURVE_ATTACK_302 = 14,
        GROW_CURVE_ATTACK_303 = 15,
        GROW_CURVE_ATTACK_304 = 16,
        GROW_CURVE_ATTACK_305 = 17,
        GROW_CURVE_CRITICAL_301 = 18
    }

    /** WeaponClass enum. */
    enum WeaponClass {
        INVALID_WEAPON_CLASS = 0,
        WEAPON_SWORD_ONE_HAND = 1,
        WEAPON_CLAYMORE = 2,
        WEAPON_POLE = 3,
        WEAPON_BOW = 4,
        WEAPON_CATALYST = 5
    }

    /** BodyType enum. */
    enum BodyType {
        INVALID_BODY_TYPE = 0,
        BODY_UNKNOWN = 1,
        BODY_BOY = 2,
        BODY_GIRL = 3,
        BODY_MALE = 4,
        BODY_LADY = 5,
        BODY_LOLI = 6
    }

    /** ZoneType enum. */
    enum ZoneType {
        INVALID_ZONE_TYPE = 0,
        ASSOC_TYPE_UNKNOWN = 1,
        ASSOC_TYPE_MONDSTADT = 2,
        ASSOC_TYPE_LIYUE = 3,
        ASSOC_TYPE_INAZUMA = 4,
        ASSOC_TYPE_SUMERU = 5,
        ASSOC_TYPE_FATUI = 6,
        ASSOC_TYPE_RANGER = 7,
        ASSOC_TYPE_MAINACTOR = 8
    }

    /** Element enum. */
    enum Element {
        INVALID_ELEMENT = 0,
        Electric = 1,
        Fire = 2,
        Ice = 3,
        Water = 4,
        Grass = 5,
        ELEMENT_QUICKEN = 6,
        ELEMENT_FROZEN = 7,
        Wind = 8,
        Rock = 9,
        ELEMENT_NONE = 10,
        ELEMENT_PHYSICAL = 11,
        ELEMENT_UNKNOWN = 12
    }

    /** StatType enum. */
    enum StatType {
        INVALID_STAT_TYPE = 0,
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
        INVALID_SIM_MODE = 0,
        DURATION_MODE = 1,
        TTK_MODE = 2
    }

    /** ComputeWorkSource enum. */
    enum ComputeWorkSource {
        InvalidWork = 0,
        DBWork = 1,
        SubmissionWork = 2
    }

    /** DBTag enum. */
    enum DBTag {
        DB_TAG_INVALID = 0,
        DB_TAG_GCSIM = 1,
        DB_TAG_TESTING = 2
    }

    /** Properties of a Version. */
    interface IVersion {

        /** Version major */
        major?: (string|null);

        /** Version minor */
        minor?: (string|null);
    }

    /** Represents a Version. */
    class Version implements IVersion {

        /**
         * Constructs a new Version.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IVersion);

        /** Version major. */
        public major: string;

        /** Version minor. */
        public minor: string;

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

        /** SimulationResult key_type */
        key_type?: (string|null);

        /** SimulationResult created_date */
        created_date?: (number|Long|null);
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
        public sim_version?: (string|null);

        /** SimulationResult modified. */
        public modified?: (boolean|null);

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

        /** SimulationResult key_type. */
        public key_type: string;

        /** SimulationResult created_date. */
        public created_date: (number|Long);

        /** SimulationResult _sim_version. */
        public _sim_version?: "sim_version";

        /** SimulationResult _modified. */
        public _modified?: "modified";

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

        /** SimulationStatistics iterations */
        iterations?: (number|null);

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
        total_damage?: (model.IDescriptiveStats|null);

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

        /** SimulationStatistics iterations. */
        public iterations: number;

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
        public total_damage?: (model.IDescriptiveStats|null);

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
        hist?: (number[]|null);
    }

    /** Represents an OverviewStats. */
    class OverviewStats implements IOverviewStats {

        /**
         * Constructs a new OverviewStats.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IOverviewStats);

        /** OverviewStats min. */
        public min?: (number|null);

        /** OverviewStats max. */
        public max?: (number|null);

        /** OverviewStats mean. */
        public mean?: (number|null);

        /** OverviewStats SD. */
        public SD?: (number|null);

        /** OverviewStats Q1. */
        public Q1?: (number|null);

        /** OverviewStats Q2. */
        public Q2?: (number|null);

        /** OverviewStats Q3. */
        public Q3?: (number|null);

        /** OverviewStats hist. */
        public hist: number[];

        /** OverviewStats _min. */
        public _min?: "min";

        /** OverviewStats _max. */
        public _max?: "max";

        /** OverviewStats _mean. */
        public _mean?: "mean";

        /** OverviewStats _SD. */
        public _SD?: "SD";

        /** OverviewStats _Q1. */
        public _Q1?: "Q1";

        /** OverviewStats _Q2. */
        public _Q2?: "Q2";

        /** OverviewStats _Q3. */
        public _Q3?: "Q3";

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
        public min?: (number|null);

        /** DescriptiveStats max. */
        public max?: (number|null);

        /** DescriptiveStats mean. */
        public mean?: (number|null);

        /** DescriptiveStats SD. */
        public SD?: (number|null);

        /** DescriptiveStats _min. */
        public _min?: "min";

        /** DescriptiveStats _max. */
        public _max?: "max";

        /** DescriptiveStats _mean. */
        public _mean?: "mean";

        /** DescriptiveStats _SD. */
        public _SD?: "SD";

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
        bucket_size?: (number|null);

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
        public bucket_size: number;

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
        bucket_size?: (number|null);

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
        public bucket_size: number;

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

        /** Sample build_date */
        build_date?: (string|null);

        /** Sample sim_version */
        sim_version?: (string|null);

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

        /** Sample build_date. */
        public build_date: string;

        /** Sample sim_version. */
        public sim_version?: (string|null);

        /** Sample modified. */
        public modified?: (boolean|null);

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

        /** Sample _sim_version. */
        public _sim_version?: "sim_version";

        /** Sample _modified. */
        public _modified?: "modified";

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
        level?: (number|null);

        /** Character max_level */
        max_level?: (number|null);

        /** Character cons */
        cons?: (number|null);

        /** Character weapon */
        weapon?: (model.IWeapon|null);

        /** Character talents */
        talents?: (model.ICharacterTalents|null);

        /** Character sets */
        sets?: ({ [k: string]: number }|null);

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
        public level: number;

        /** Character max_level. */
        public max_level: number;

        /** Character cons. */
        public cons: number;

        /** Character weapon. */
        public weapon?: (model.IWeapon|null);

        /** Character talents. */
        public talents?: (model.ICharacterTalents|null);

        /** Character sets. */
        public sets: { [k: string]: number };

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
        attack?: (number|null);

        /** CharacterTalents skill */
        skill?: (number|null);

        /** CharacterTalents burst */
        burst?: (number|null);
    }

    /** Represents a CharacterTalents. */
    class CharacterTalents implements ICharacterTalents {

        /**
         * Constructs a new CharacterTalents.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.ICharacterTalents);

        /** CharacterTalents attack. */
        public attack: number;

        /** CharacterTalents skill. */
        public skill: number;

        /** CharacterTalents burst. */
        public burst: number;

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
        refine?: (number|null);

        /** Weapon level */
        level?: (number|null);

        /** Weapon max_level */
        max_level?: (number|null);
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
        public refine: number;

        /** Weapon level. */
        public level: number;

        /** Weapon max_level. */
        public max_level: number;

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
        level?: (number|null);

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
        public level: number;

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
        number_of_workers?: (number|null);

        /** SimulatorSettings iterations */
        iterations?: (number|null);

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
        public number_of_workers: number;

        /** SimulatorSettings iterations. */
        public iterations: number;

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
        skill?: (number|null);

        /** Delays burst */
        burst?: (number|null);

        /** Delays attack */
        attack?: (number|null);

        /** Delays charge */
        charge?: (number|null);

        /** Delays aim */
        aim?: (number|null);

        /** Delays dash */
        dash?: (number|null);

        /** Delays jump */
        jump?: (number|null);

        /** Delays swap */
        swap?: (number|null);
    }

    /** Represents a Delays. */
    class Delays implements IDelays {

        /**
         * Constructs a new Delays.
         * @param [properties] Properties to set
         */
        constructor(properties?: model.IDelays);

        /** Delays skill. */
        public skill: number;

        /** Delays burst. */
        public burst: number;

        /** Delays attack. */
        public attack: number;

        /** Delays charge. */
        public charge: number;

        /** Delays aim. */
        public aim: number;

        /** Delays dash. */
        public dash: number;

        /** Delays jump. */
        public jump: number;

        /** Delays swap. */
        public swap: number;

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
        start?: (number|null);

        /** EnergySettings end */
        end?: (number|null);

        /** EnergySettings amount */
        amount?: (number|null);

        /** EnergySettings last_energy_drop */
        last_energy_drop?: (number|null);
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
        public start: number;

        /** EnergySettings end. */
        public end: number;

        /** EnergySettings amount. */
        public amount: number;

        /** EnergySettings last_energy_drop. */
        public last_energy_drop: number;

        /**
         * Gets the default type url for EnergySettings
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
