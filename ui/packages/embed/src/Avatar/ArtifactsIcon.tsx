const ArtifactsIcon = ({
  sets,
  half = false,
}: {
  sets: string[];
  half: boolean;
}): JSX.Element => {
  const artifacts: JSX.Element[] = [
    <image
      key={0}
      filter="url(#outlinew)"
      href={`/api/assets/artifacts/${sets[0]}_flower.png`}
      height="35"
      width={sets.length > 1 || half ? "17.5" : "35"}
      x="20"
      y="33"
      preserveAspectRatio={
        sets.length > 1 || half ? "xMinYMid slice" : undefined
      }
    ></image>,
  ];
  if (sets.length > 1) {
    artifacts.push(
      <image
        key={1}
        filter="url(#outlinew)"
        href={`/api/assets/artifacts/${sets[1]}_flower.png`}
        height="35"
        width="17.5"
        x="37.5"
        y="33"
        preserveAspectRatio="xMaxYMid slice"
      ></image>
    );
  }

  return <g filter="url(#outlineb)">{artifacts}</g>;
};

export default ArtifactsIcon;
