import tanuki from "images/tanuki.png";

export function Warning() {
  return (
    <div className="relative bg-slate-900 px-5 sm:px-5 py-0 border border-blue-800 md:min-w-wmd sm:min-w-wsm md:min-w-wmd sm:max-w-sm md:max-w-2xl lg:max-w-4xl xl:max-w-[900px]">
      <div className="inline-flex pt-4">
        <img src={tanuki} className="w-15 h-10 mx-0"/>
        <div className="font-semibold px-3 pt-2 text-xl w-50 text-gray-200">Please read!</div>
        <img src={tanuki} className="w-15 h-10 mx-0"/>
      </div>
      <div className="space-y-3 pb-3 text-s leading-5 text-gray-400">
        <p>Unless tagged otherwise, the sims in this database are here to provide examples of interesting, user-submitted configs.</p>
        <p>Be aware that these configs will vary with different standards, levels of optimization, and reasonable application to in-game performance. You are encouraged to write your own, or edit these and click rerun.</p>
        <p className="font-semibold leading-6 text-gray-200">gcsim does not exist solely to output standardized theoretical mean dps. It is entirely up to the user to input what they deem useful. So stop blindly comparing the dps of the configs here, you're doing it wrong.</p>
      </div>
    </div>
  );
}