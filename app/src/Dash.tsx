import Nav from "./Nav";
import { Icon } from "@blueprintjs/core";

export default function Dash() {
  return (
    <div>
      <Nav />
      <div className="w-full flex flex-col items-center ">
        <span className="font-bold text-md mt-4">
        <a href="https://github.com/genshinsim/gcsim" target="_blank" >gcsim</a> is a team dps simulator for Genshin Impact. Get started by choosing one of the following options!
        </span>
        <div className="flex flex-row flex-initial flex-wrap w-full lg:w-[60rem] mt-4">
          <div className="main-page-button-container">
            <div className="main-page-button">
              <span className="font-bold text-xl">
                <Icon icon="flame" className="mr-2" size={25} />
                Simple Mode
              </span>
            </div>
          </div>
          <div className="main-page-button-container">
            <div className="main-page-button">
              <span className="font-bold text-xl">
                <Icon icon="settings" className="mr-2" size={25} />
                Advanced Mode
              </span>
            </div>
          </div>

          <div className="main-page-button-container">
            <div className="main-page-button">
              <span className="font-bold text-xl">
                <Icon icon="chart" className="mr-2" size={25} />
                Viewer
              </span>
            </div>
          </div>
          <div className="main-page-button-container">
            <div className="main-page-button">
              <span className="font-bold text-xl">
              <Icon icon="database" className="mr-2" size={25} />
                Action Lists DB
              </span>
            </div>
          </div>
          <div className="main-page-button-container">
            <div className="main-page-button">
              <span className="font-bold text-xl">
              <Icon icon="database" className="mr-2" size={25} />
                Desktop Tool
              </span>
            </div>
          </div>
          <div className="main-page-button-container">
            <div className="main-page-button">
              <span className="font-bold text-xl">
                <Icon icon="document" className="mr-2" size={25} />
                Documentation
              </span>
            </div>
          </div>
          <div className="main-page-button-container">
            <div className="main-page-button">
              <span className="font-bold text-xl">
                <Icon icon="git-branch" className="mr-2" size={25} />
                Contribute
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
