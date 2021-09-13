import { simConfig } from './simSlice'
import './workerHack'

onmessage = async (ev: { data: simConfig }) => {
    const t1 = performance.now()


    let t2 = performance.now()

    console.log("finished in: ", t2 - t1)

}