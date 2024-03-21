import { Button, ButtonGroup } from "@blueprintjs/core"
import { useLocation } from "wouter"
import tagData from "tags.json";

export const Home = () => {
    const [_, to] = useLocation()
    const sortedTagnames = Object.keys(tagData)
        .filter((key) => {
            return key !== "0" && key != "2";
        })
        .map((key) => {
            let name = tagData[key]["display_name"]
            if (key == "1") {
                name = "(Not Tagged)"
            }
            return <li key={key}>
                <span className="font-semibold text-rose-700">{name}</span>{`: ${tagData[key]["blurb"]}`}
            </li>
        });
    return (
        <div className="ml-2 mr-2 mt-2">
            <div className="text-center text-lg font-semibold text-indigo-600">Welcome to Simpact</div>
            <div className="mb-4">
                <p className="m-2">
                    Simpact is a database of gcsim simulations submitted and maintained by gcsim users.
                </p>
                <p className="m-2">
                    The database consists of various tags, each tag representing a collection of sims. These collections are maintained by volunteers and each collection follow their own rules.
                </p>
                <p className="m-2">Below are the current available tags:</p>
                <ul className="list-disc m-4 ml-8">
                    {sortedTagnames}
                </ul>
            </div>
            <ButtonGroup fill className="mb-4">
                <Button intent="primary" onClick={() => to("/database")}>Get started</Button>
            </ButtonGroup>

        </div>
    )

} 