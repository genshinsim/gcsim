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
          <H3>FAQs (and Disclaimers). Seriously. READ THIS FIRST.</H3>
          <p>Last updated July 7th, 2022.</p>
          <p>
            Welcome to the gcsim rotation database. Here you'll find many{" "}
            <b>user submitted</b> rotations and their calculated DPS.
          </p>
          <p className=" text-red-700 font-semibold">
            You will want to read these FAQs in its entirety before you draw any
            conclusions from this database. Seriously. Read it. Don't say I
            didn't warn you.
          </p>
          <H4>Purpose of this database</H4>
          <p>
            This database primarily serves two purposes:
            <li className="ml-6">
              It serves as a collection of gcsim configs that others can use as
              an example to write their own configs.
            </li>
            <li className="ml-6">
              It is sorted by DPS in order to provoke a reaction for the purpose
              of driving code quality. We subscribe to same philosophy as used
              by World of Warcraft's SimulationCraft. See{" "}
              <a href="https://github.com/simulationcraft/simc/wiki/PremedititatedProvocation">
                here
              </a>
            </li>
          </p>
          <p>
            The purpose of this database <b>is not to provide a tier list</b>.
            There are too many qualitative considerations that cannot be
            captured in these calculations. In addition, all of these
            calculations make certain substat assumptions that in general do not
            apply to all users.
          </p>
          <H4>Why is this database sorted by DPS?!</H4>
          <p>See point 2 in the purpose above</p>
          <H4>
            I cannot replicate these numbers in game! The calc must be wrong!
          </H4>
          <p>
            Addendum (July 7, 2022): It appears that some people are
            misinterpreting this explanation. No gcsim cannot calculate your in
            game damage, but that does not make the results useless or
            inaccurate. What it means is that what <b>YOU</b> do in game can
            differ from simulated results due to skill issues, ping, enemy
            movement, dodging, other game mechanics etc... It also means that
            when looking at teams you must also consider qualitative factors
            such as ease to play, healing, etc... Numbers by itself does not
            tell the whole story.
          </p>
          <p>
            In the first place, gcsim cannot and will not ever be a good way to
            calculate your in game damage. There does not (currently) exist any
            content in Genshin that even remotely resembles a dps dummy. In
            addition to that, most simulations generally assume perfect inputs
            and zero skill issues. In reality, probably most of the population
            cannot pull that off with any sort of consistency. End of the day,
            if you want a dps meter, it's probably best to petition Hoyoverse
            for that.
          </p>
          <p>
            In addition, Genshin combat relies upon reactions which are timing
            specific and just missing a single one could alter your results.
            (Obviously this differs from team to team.)
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
            Finally, all of these calculations are user submitted, and anyone is
            capable of submitting one no matter how unoptimized the rotation may
            appear. Meaning some teams and their rotations are more optimized
            than others due to their popularity. Some rotations are more
            difficult to execute than others, this is not something than can be
            captured by numbers alone.
          </p>
          <H4>
            Why don't you simulate skill issue or make these simulations more
            "realistic"?
          </H4>
          <p>
            It's not that we don't want to see more realistic simulations. It's
            that "realistic" is way too subjective. What's realistic for one
            player may not be realistic for another. So where do we draw the
            line?
          </p>
          <p>
            In the first place, gcsim is just a fancy calculator. Calculators
            don't have opinions on their own. The whole idea is that you can
            always take an existing rotation someone else has put together and
            modify it to something you can pull off yourself.
          </p>
          <p>
            In addition to that, highly optimized rotations can have educational
            value. You can analyze the rotation to see tips and tricks that
            increase your dps ever so slightly.
          </p>
          <H4>What about hitlag?!</H4>
          <p>
            Yes hitlag is not implemented in gcsim. It is something we are
            working very hard on to add in. However, due to the complexity (it
            literally affects anything that relies on time, but not
            consistently), this is taking a while and is forcing us to rewrite
            our entire back end. This is also why you have not seen any updates
            to the UI. All my time is spent on dealing with this.
          </p>
          <p>
            This does have some implications on the results:
            <li className="ml-6">
              The exact rotation will not translate to the same amount of
              execution time in game because all frames currently used in gcsim
              are based on hitlag free frames. You can roughly estimate the
              result of hitlag on rotation length by adding roughly 3 frames to
              every melee hit. Give or take.
            </li>
            <li className="ml-6">
              It will cause the dps to be potentially inflated compared to a
              world with hitlag due to longer actual rotation time in game.
              However, the expected inflation is not as large as some rumours
              would have it. My current napkin math places it somewhere around
              10% at most. For teams that primarily rely upon either off-field
              or fully ranged sources of damage, the effects of hitlag should be
              even smaller.
            </li>
            <li className="ml-6">
              However, since gcsim is meant for relative comparisons, the
              current lack of hitlag should not affect the calculations
              relatively speaking (excluding certain teams that are purely
              ranged).
            </li>
          </p>
          <p>
            An additional note on the rotation time. Hitlag is not the only
            reason why rotation times can be longer than in-game. As I mentioned
            previously, gcsim assumes frame perfect. However, most people cannot
            play frame perfect with any sort of consistency. Every tool has its
            intended purpose and limitations. That is no different for gcsim.
            Always keep that in mind.
          </p>
          <H4>
            Why does character x uses artifact or weapon y?! We all know z is
            better!
          </H4>
          <p>
            As mentioned previously, these are all user submitted calculations.
            Currently, due to the design limitations of the website, only one
            calculation is kept per team. To keep things simple, the maintainer
            of this database (a volunteer) has decided to show the team that has
            the highest dps to keep things unopinionated. This does mean that
            unfortunately some very valid weapon/artifact combinations that are
            just as competitive are not shown.
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
              and would like to help out, we're in desperate need of help.
              Please help us make this site better by adding in some of these
              features. Please reach out to us on{" "}
              <a href="https://discord.gg/m7jvjdxx7q">discord</a>.
            </b>
          </p>
          <H4>So you're saying I should farm artifact x right</H4>
          <p>
            I sincerely hope that is not the conclusion you will draw from this
            database.{" "}
            <span className="line-through">
              And if it is please consider re-reading this FAQ again.
            </span>
          </p>
          <p>
            Seriously. Due to limitations mentioned above we only show one
            possible weapon/artifact configuration per team. It does not make
            said configuration optimal to farm given resin constraints. It also
            does not mean said configuration is the absolute best (because we're
            not fighting dps dummies here in Genshin).
          </p>
          <H4>Why is this database single target only?!</H4>
          <p>
            For the same reason as the question above. We all want to see
            multiple targets. But we are severely limited by the available
            developer time.{" "}
            <b className="text-amber-600">
              So again, if you're a web developer and would like to help out,
              please please reach out to us on{" "}
              <a href="https://discord.gg/m7jvjdxx7q">discord</a>.
            </b>
          </p>
          <p>
            Edited Jun 23, 2022: It has been brought to my attention that
            somehow everyone thinks this means that gcsim can only handle single
            target calculations. That is <b>NOT</b> what I'm saying here. gcsim
            itself can handle multiple targets just fine (in fact as many as
            your computer can handle; just copy and paste the target line in the
            config). However, we simply don't have the UI to show both single
            target and multi target simulations on the database nicely. So it
            was decided by the maintainer to only show single target for now.
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
              PLEASE SEND <a href="https://discord.gg/m7jvjdxx7q"> HALP</a>
            </b>
          </p>
          <H4>
            I want to see more rotations in the future! I have a rotation that I
            want to share!
          </H4>
          <p>
            This is a great way to share your rotation with the community.
            Please do not hesitate to submit your rotation on our{" "}
            <a href="https://discord.gg/m7jvjdxx7q">discord</a>.
          </p>
          <H4>
            I have a suggestion for a new rotation or a correction to a
            rotation.
          </H4>
          <p>
            Please do not hesitate to submit your suggestion on our{" "}
            <a href="https://discord.gg/m7jvjdxx7q">discord</a>.
          </p>
          <H4>I have a question that is not answered here.</H4>
          <p>
            Please do not hesitate to ask your question on the{" "}
            <a href="https://discord.gg/m7jvjdxx7q">discord</a>.
          </p>
          <H4>I have a problem with the database.</H4>
          <p>
            <span className=" line-through">Blame srl.</span> Please do not
            hesitate to contact us on{" "}
            <a href="https://discord.gg/m7jvjdxx7q">discord</a>.
          </p>
          <H4>
            In case it wasn't obvious, come talk to us on{" "}
            <a href="https://discord.gg/m7jvjdxx7q">discord</a> about the
            meaning of sim, life, 42, anything.
          </H4>
        </div>
      </div>
      <div className={Classes.DIALOG_FOOTER}>
        Written by srl.
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
