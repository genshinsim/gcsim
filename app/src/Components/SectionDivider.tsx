type Props = {
  children: React.ReactNode;
  fontClass?: string;
};

export function SectionDivider({
  children,
  fontClass = "font-bold text-xl",
}: Props) {
  return (
    <div className="flex flex-row place-items-center mt-2 mb-2">
      <div className="ml-1 flex-grow border-t h-0 mr-1 mt-1"></div>
      <span className={fontClass}>{children}</span>
      <div className="ml-1 flex-grow border-t h-0 mr-1 mt-1"></div>
    </div>
  );
}
