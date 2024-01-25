import { db } from "@gcsim/types";
import eula from "images/eula.png";
import InfiniteScroll from "react-infinite-scroll-component";
import { ActionBar } from "SharedComponents/ActionBar";
import { Warning } from "SharedComponents/Warning";
import { ListView } from "../../SharedComponents/ListView";
import { useTranslation } from "react-i18next";

type Props = {
  data: db.IEntry[];
  fetchData: () => void;
  hasMore: boolean;
};

export const DBView = (props: Props) => {
  const { t } = useTranslation();
  return (
    <div className="flex flex-col gap-4 m-8 my-4 items-center">
      <ActionBar simCount={props.data.length} />
      <Warning />
      {props.data.length === 0 ? (
        <div className="6 flex flex-col justify-center items-center h-screen">
          <img src={eula} className=" object-contain opacity-50 w-32 h-32" />
        </div>
      ) : (
        <InfiniteScroll
          dataLength={props.data.length} //This is important field to render the next data
          next={props.fetchData}
          hasMore={props.hasMore}
          loader={<h4>{t<string>("sim.loading")}</h4>}
          endMessage={
            <>
              <p className="text-center mt-4">
                <b>Yay! You have seen it all.</b>
              </p>
              <p className="text-center">
                {` Didn't find what you're looking for? Create and submit your own on our discord `}
              </p>
            </>
          }
          //TODO: enable pull down functionality for refreshing maybe??
        >
          <ListView data={props.data} />
        </InfiniteScroll>
      )}
    </div>
  );
};
