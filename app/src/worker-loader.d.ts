// worker-loader.d.ts
declare module "worker-loader!*" {
    class SimWorker extends Worker {
        constructor();
    }

    // Uncomment this if you set the `esModule` option to `false`
    // export = WebpackWorker;
    export default SimWorker;
}