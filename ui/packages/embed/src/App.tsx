import { parsed } from ".";
import { ErrorBoundary } from "react-error-boundary";
import Layout from "./Layout";


function fallbackRender({ error  }) {
  return (
    <div id="card" role="alert">
      <p>Something went wrong:</p>
      <pre style={{ color: "red" }}>{error.message}</pre>
    </div>
  );
}

const App = ({ }) => {
  if ("err" in parsed && parsed.err != "") {
    return <div id="card">{parsed.err}</div>;
  }

  return (
    <ErrorBoundary fallbackRender={fallbackRender}>
      <Layout data={parsed} />
    </ErrorBoundary>
  );
};

export default App;
