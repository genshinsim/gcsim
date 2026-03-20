// @gcsim/types public API

// Generated protobuf types — backend (namespaced to avoid conflicts)
export * as db from "./generated/protos/backend/db_pb.js";
export * as preview from "./generated/protos/backend/preview_pb.js";
export * as share from "./generated/protos/backend/share_pb.js";

// Generated protobuf types — model (top-level, most commonly used)
export * from "./generated/protos/model/data_pb.js";
export * from "./generated/protos/model/db_pb.js";
export * from "./generated/protos/model/enums_pb.js";
export * from "./generated/protos/model/notification_pb.js";
export * from "./generated/protos/model/result_pb.js";
export * from "./generated/protos/model/sample_pb.js";
export * from "./generated/protos/model/sim_pb.js";

// Custom interfaces (JSON result shapes from the simulator)
// Namespaced to avoid collisions with proto-generated types of the same name
export type * as Sim from "./sim.js";
