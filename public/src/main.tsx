import { render } from "preact"

import { App } from "./components/app";

function main() {
  const container = document.getElementById("root");

  if (container == null) {
    return;
  }

  render(<App />, container);
}

main();
