type Props = {
  children: React.ReactNode;
  className?: string;
};

export function Main(props: Props) {
  return (
    <main
      className={
        "m-2 xs:w-[25rem] md:w-[46rem] wide:w-[72rem] ml-auto mr-auto " +
        props.className
      }
    >
      {props.children}
    </main>
  );
}
