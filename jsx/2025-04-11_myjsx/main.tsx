const Title = (props: { children: string }) => <h1>{props.children}</h1>;

const element = (
  <div>
    <Title>Hello, world!</Title>
    <p>Goodbye, world!</p>
  </div>
);

console.log(element);
