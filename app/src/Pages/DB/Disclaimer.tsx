import { Button, Classes, Dialog, Drawer, H3, H4 } from "@blueprintjs/core";

type Props = {
  isOpen: boolean;
  onClose: () => void;
  hideAlways: () => void;
};

export function Disclaimer(props: Props) {
  return (
    <Dialog
      className="w-screen"
      isOpen={props.isOpen}
      onClose={props.onClose}
      canEscapeKeyClose
      canOutsideClickClose
      style={{ width: "90%", maxWidth: "800px" }}
    >
      <div className={Classes.DIALOG_BODY}>
        <div className="flex flex-col place-items-start ">
          <H3>FAQs</H3>
          <p>
            Welcome to the gcsim rotation database. Here you'll find many{" "}
            <b>user submitted</b> rotations and their calculated DPS.
          </p>
          <H4>Purpose of this database</H4>
          <p>
            This database primarily serves two purposes:
            <li className="ml-6">
              It serves as a collection of gcsim configs that others can use as
              an example to write their own configs
            </li>
            <li className="ml-6">
              It is sorted by DPS in order to provoke a reaction for the purpose
              of driving code quality. We subscribe to same philsophy as used by
              World of Warcraft's SimulationCraft. See{" "}
              <a href="https://github.com/simulationcraft/simc/wiki/PremedititatedProvocation">
                here
              </a>
            </li>
          </p>
          <p>
            The purpose of this database <b>is not to provide a tier list</b>.
            There are too many qualitative considerations that cannot be
            captured in these calculations. In addition, all of these
            calculations makes certain substat assumptions that in general does
            not apply to all users.
          </p>
          <H4>Why is this database sorted by DPS?!</H4>
          <p>See point 2 in the purpose above</p>
          <H4>
            I cannot replicate these numbers in game! The calc must be wrong!
          </H4>
          <p>
            In the first place, gcsim cannot and will not ever be a good way to
            calculate your in game damage. There does not (current) exist any
            content in Genshin that even remotely resemble a dps dummy. In
            addition to that, most simulations generally assume perfect inputs
            and zero skill issues. In reality, probably most of the population
            cannot pull that off with any sort of consistency. End of the day,
            if you want a dps meter, it's probably best to petition Hoyoverse
            for that.
          </p>
          <p>
            In addition, Genshin combat relies upon reactions which are very
            timing specific and just missing a single one could potentially
            result in very different results. (Obviously this differs from team
            to team.)
          </p>
          <p>
            By design, the primary purpose of gcsim is to provide relative
            comparisons, helping answer questions such as "is weapon x better
            than weapon y for <b>my team</b>". It helps fill a niche that
            Genshin Optimizer struggles with (i.e. team dps calcs), but is
            getting notably better at (Thanks to Waverider) Seriously, we love
            Genshin Optimizer.
          </p>
          <p>
            Finally, all of these calculations are user submitted. Meaning some
            of the rotations are far more optimized than others due to their
            popularity. Some rotations are extremely tight and very difficult to
            play while others are far more forgiving. This is not something that
            can be captured by numbers alone.
          </p>
          <H4>
            Why don't you simulate skill issue or make these simulations more
            "realistic"
          </H4>
          <p>
            It's not that we don't want to see more realistic simulations. It's
            that "realistic" is way too subjective. What's realistic for one
            player may not be realistic for another. So where do we draw the
            line?
          </p>
          <p>
            In the first place, gcsim is just a fancy calculator. Calculators
            don't have opinions on its own. The whole idea is that you can
            always take an existing rotation someone else has put together and
            modify it to something you can pull off yourself.
          </p>
          <p>
            In addition to that, highly optimized rotation has a lot of
            educational value. You can take a look at the rotation and see all
            the little tips and tricks you can use to increase your dps ever so
            slightly.
          </p>
          <H4>
            Why does character x uses artifact or weapon y?! We all know z is
            better!
          </H4>
          <p>
            As mentioned previously, these are all user submitted calculations.
            Currently, due to the design limitations of the website, only one
            calculation is kept per team. To keep things simple, the maintainer
            of this database (a volunteer) has decided to shown the team that
            has the highest dps to keep things unopinionated. This does mean
            that unfortunately some very valid weapon/artifact combinations that
            are just as competitive are not shown.
          </p>
          <p>
            In an ideal world, we would have the interface to show all the
            different possible team/weapon/artifact combinations that users have
            submitted. But due to dev time constraints this is the best we got.
          </p>
          <p>
            Of course, that doesn't mean you cannot take an existing rotation on
            here and modify it to see for yourself what it would look like with
            weapon x or artifact y. In fact, that's the primary purpose of this
            database as outlined in the Purpose section.
          </p>
          <p>
            <b className="text-amber-600">
              Finally, if you are a web developer (Javascript/Typescript/React)
              and would like to help out, we're in desparate need of help.
              Please help us make this site better by adding in some of these
              features. Please reach out to us on{" "}
              <a href="https://discord.gg/gcsim">discord</a>.
            </b>
          </p>
          <H4>Why is this database single target only?!</H4>
          <p>
            For the same reason as the question above. We all want to see
            multiple targets. But we are severely limited by the available
            developer time.{" "}
            <b className="text-amber-600">
              So again, if you're a web developer and would like to help out,
              please please reach out to us on{" "}
              <a href="https://discord.gg/gcsim">discord</a>.
            </b>
          </p>
          <H4>Why does this database load so damn slow?!</H4>
          <p>
            Because some idiot dev{" "}
            <span className="line-through">blame srl</span> designed this super
            poorly and loads the entire (ever growing) database in one fetch.
            Said dev clearly did not think about pagination or any of the
            typical good web practices and simply slapped this page together as
            an after thought.
          </p>
          <p>
            Also, in case it wasn't made clear before...{" "}
            <b className="text-amber-600">
              PLEASE SEND <a href="https://discord.gg/gcsim"> HALP</a>
            </b>
          </p>
          <H4>
            I want to see more rotations in the future! I have a rotation that I
            want to share!
          </H4>
          <p>
            This is a great way to share your rotation with the community.
            Please do not hesitate to submit your rotation on our{" "}
            <a href="https://discord.gg/gcsim">discord</a>.
          </p>
          <H4>
            I have a suggestion for a new rotation or a correction to a
            rotation.
          </H4>
          <p>
            Please do not hesitate to submit your suggestion on our{" "}
            <a href="https://discord.gg/gcsim">discord</a>.
          </p>
          <H4>I have a question that is not answered here.</H4>
          <p>
            Please do not hesitate to ask your question on the{" "}
            <a href="https://discord.gg/gcsim">discord</a>.
          </p>
          <H4>I have a problem with the database.</H4>
          <p>
            <span className=" line-through">Blame srl.</span> Please do not
            hesitate to contact us on{" "}
            <a href="https://discord.gg/gcsim">discord</a>.
          </p>
          <H4>
            In case it wasn't obvious, come talk to us on{" "}
            <a href="https://discord.gg/gcsim">discord</a> about the meaning of
            sim, life, 42, anything.
          </H4>
        </div>
      </div>
      <div className={Classes.DIALOG_FOOTER}>
        <div className={Classes.DIALOG_FOOTER_ACTIONS}>
          <Button intent="primary" onClick={props.onClose}>
            Close
          </Button>
          <Button intent="danger" onClick={props.hideAlways}>
            Don't show again
          </Button>
        </div>
      </div>
    </Dialog>
  );
}
