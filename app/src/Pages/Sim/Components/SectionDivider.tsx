type Props = {
  children: React.ReactNode;
};

export function SectionDivider(props: Props) {
  return (
    <div className="flex flex-row place-items-center">
      <div className="ml-1 flex-grow border-t h-0 mr-1 mt-1"></div>
      <span className="font-bold text-xl">{props.children}</span>
      <div className="ml-1 flex-grow border-t h-0 mr-1 mt-1"></div>
    </div>
  );
}
