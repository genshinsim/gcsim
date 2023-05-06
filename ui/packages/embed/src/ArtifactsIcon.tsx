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
      height="43"
      width={sets.length > 1 || half ? "20.5" : "43"}
      x="30"
      y="52"
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
        height="43"
        width="20.5"
        x="52.5"
        y="52"
        preserveAspectRatio="xMaxYMid slice"
      ></image>
    );
  }

  return <g filter="url(#outlineb)">{artifacts}</g>;
};

export default ArtifactsIcon;
