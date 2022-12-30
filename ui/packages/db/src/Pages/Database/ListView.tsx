import React from "react";

type ListViewProps = {
  query: any;
  sort: any;
  skip: any;
  limit: any;
};

export function ListView(props: ListViewProps) {
  const [data, setData] = React.useState<any[]>([]);
  React.useEffect(() => {
    const url = `/api/db?q=${encodeURIComponent(JSON.stringify(props.query))}`;
    fetch(url)
      .then((res) => res.json())
      .then((data) => {
        console.log(data);
        setData(data);
      })
      .catch((e) => {
        console.log(e);
      });
  }, []);
  return <div>hi</div>;
}
