/** @type {import('dependency-cruiser').IConfiguration} */
module.exports = {
  forbidden: [
    {
      name: "no-deep-package-imports",
      comment:
        "Import from package index only (@gcsim/<pkg>), never from internal paths like @gcsim/<pkg>/src/...",
      severity: "error",
      from: {},
      to: {
        pathNot: "node_modules",
        path: "packages/[^/]+/src/.+",
        via: {
          pathNot: "^packages/[^/]+/src/",
        },
      },
    },
    {
      name: "no-circular",
      comment: "No circular dependencies allowed between packages",
      severity: "error",
      from: {},
      to: {
        circular: true,
      },
    },
    {
      name: "no-app-to-app",
      comment: "Apps must not import from other apps",
      severity: "error",
      from: {
        path: "^apps/[^/]+/",
      },
      to: {
        path: "^apps/[^/]+/",
        pathNot: "^apps/({FROM_APP})/",
      },
    },
    {
      name: "no-package-to-app",
      comment: "Packages must not import from apps",
      severity: "error",
      from: {
        path: "^packages/",
      },
      to: {
        path: "^apps/",
      },
    },
  ],
  options: {
    doNotFollow: {
      path: "node_modules",
    },
    tsPreCompilationDeps: true,
    enhancedResolveOptions: {
      exportsFields: ["exports"],
      conditionNames: ["import", "require", "node", "default"],
    },
    reporterOptions: {
      text: {
        highlightFocused: true,
      },
    },
  },
};
