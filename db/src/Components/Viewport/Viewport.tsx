type Props = {
  children: React.ReactNode;
  className?: string;
};

export function Viewport(props: Props) {
  return (
    <main
      className={
        "m-2 w-full xs:w-[300px] sm:w-[640px] md:w-[750px] md:ml-2 md:mr-2 wide:w-[1160px] ml-auto mr-auto " +
        props.className
      }
    >
      {props.children}
    </main>
  );
}
