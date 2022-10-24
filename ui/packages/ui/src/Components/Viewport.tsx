type Props = {
  children: React.ReactNode;
  className?: string;
};

export function Viewport(props: Props) {
  return (
    <main
      className={
        "m-2 w-full xs:w-[300px] sm:w-[640px] hd:w-full wide:w-[1160px] ml-auto mr-auto " +
        props.className
      }
    >
      {props.children}
    </main>
  );
}
