interface HasChildren {
  children?: string | string[];
  id?: string;
}

declare namespace JSX {
  interface IntrinsicElements {
    div: HasChildren;
    h1: HasChildren;
    p: HasChildren;
  }
}
